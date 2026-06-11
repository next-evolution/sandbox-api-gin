package infradbfx

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	goredis "github.com/redis/go-redis/v9"

	"sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
	"sandbox-api-gin/internal/infrastructure/external"
)

const usdJpy = "USDJPY"

type MySQLTradeSimulationRepository struct {
	db          *sqlx.DB
	redisClient *goredis.Client
	gaitame     *external.GaitameRateService
}

func NewMySQLTradeSimulationRepository(
	db *sqlx.DB,
	redisClient *goredis.Client,
	gaitame *external.GaitameRateService,
) fxrepository.TradeSimulationRepository {
	return &MySQLTradeSimulationRepository{
		db:          db,
		redisClient: redisClient,
		gaitame:     gaitame,
	}
}

func (r *MySQLTradeSimulationRepository) GetPrice(ctx context.Context, symbol string, contractAt time.Time) (*fx.PriceInfo, error) {
	barDateTime := interval15m(contractAt)
	// JavaのHM_FORMATTER DateTimeFormatter.ofPattern("yyyyMMddHHmm")に相当
	contractHm := barDateTime.Format("200601021504")

	priceUsdJpy, err := r.getPriceWithCache(ctx, usdJpy, contractHm)
	if err != nil {
		return nil, err
	}

	price := priceUsdJpy
	if symbol != usdJpy {
		price, err = r.getPriceWithCache(ctx, symbol, contractHm)
		if err != nil {
			return nil, err
		}
	}

	return &fx.PriceInfo{
		Symbol:      symbol,
		BarDateTime: barDateTime,
		Price:       price,
		PriceUsdJpy: priceUsdJpy,
	}, nil
}

func (r *MySQLTradeSimulationRepository) getPriceWithCache(ctx context.Context, symbol, contractHm string) (float64, error) {
	redisKey := fmt.Sprintf("price:%s_%s", symbol, contractHm)

	if cached, err := r.redisClient.Get(ctx, redisKey).Result(); err == nil {
		return strconv.ParseFloat(cached, 64)
	}

	price, err := r.getOpenPrice(ctx, symbol, contractHm)
	if err != nil {
		return 0, err
	}

	if price == 0 {
		if err := r.gaitame.RefreshRate(ctx, contractHm); err != nil {
			return 0, err
		}
		if refreshed, err := r.redisClient.Get(ctx, redisKey).Result(); err == nil {
			price, _ = strconv.ParseFloat(refreshed, 64)
		}
	}

	if price > 0 {
		if err := r.redisClient.Set(ctx, redisKey, strconv.FormatFloat(price, 'f', -1, 64), 60*time.Minute).Err(); err != nil {
			slog.Warn("価格キャッシュ保存失敗", "key", redisKey, "error", err)
		}
	}

	return price, nil
}

func (r *MySQLTradeSimulationRepository) getOpenPrice(ctx context.Context, symbol, contractHm string) (float64, error) {
	// JavaのTradeSimulationMapper.getOpenPrice SQLに相当。
	// 15分足 → 1時間足 → 0 の優先順でCOALESCE。
	query := `
		SELECT COALESCE(m15.open_price, h1.open_price, 0) AS open_price
		FROM (SELECT ? AS symbol) AS tmp
		LEFT JOIN fx_bar_15m m15
		  ON m15.symbol = tmp.symbol
		 AND m15.bar_date_time = STR_TO_DATE(?, '%Y%m%d%H%i')
		LEFT JOIN fx_bar_1h h1
		  ON h1.symbol = tmp.symbol
		 AND h1.bar_date_time = STR_TO_DATE(?, '%Y%m%d%H')`

	var price float64
	if err := r.db.GetContext(ctx, &price, query, symbol, contractHm, contractHm); err != nil {
		return 0, err
	}
	return price, nil
}

// interval15m は15分足バーの開始時刻に切り捨てる。
// JavaのTradeSimulationRepositoryImpl.interval15m()に相当。
func interval15m(dt time.Time) time.Time {
	minute := dt.Minute()
	var aligned int
	switch {
	case minute >= 45:
		aligned = 45
	case minute >= 30:
		aligned = 30
	case minute >= 15:
		aligned = 15
	default:
		aligned = 0
	}
	return time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), aligned, 0, 0, dt.Location())
}
