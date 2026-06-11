package infradbfx

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type MySQLBarDataRepository struct {
	db *sqlx.DB
}

func NewMySQLBarDataRepository(db *sqlx.DB) fxrepository.BarDataRepository {
	return &MySQLBarDataRepository{db: db}
}

type fxBarDataStatusRecord struct {
	Symbol          string  `db:"symbol"`
	BarDateTimeMinS *string `db:"barDateTimeMinS"`
	BarDateTimeMaxS *string `db:"barDateTimeMaxS"`
	Count           int     `db:"count"`
}

type fxBarDataRecord struct {
	Symbol      string    `db:"symbol"`
	BarDateTime time.Time `db:"barDateTime"`
	OpenPrice   float64   `db:"openPrice"`
	HighPrice   float64   `db:"highPrice"`
	LowPrice    float64   `db:"lowPrice"`
	ClosePrice  float64   `db:"closePrice"`
	Volume      int       `db:"volume"`
	HighProfit  float64   `db:"highProfit"`
	LowProfit   float64   `db:"lowProfit"`
	CloseProfit float64   `db:"closeProfit"`
	RangeProfit float64   `db:"rangeProfit"`
	RsiValue    float64   `db:"rsiValue"`
	RsiMa       float64   `db:"rsiMa"`
}

func (r *MySQLBarDataRepository) StatusList(ctx context.Context, symbolType string, barType fxmodel.BarType) ([]fxmodel.BarDataStatus, error) {
	query := fmt.Sprintf(`
		SELECT
			c.symbol                                              AS symbol,
			DATE_FORMAT(MIN(b.bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMinS,
			DATE_FORMAT(MAX(b.bar_date_time), '%%Y-%%m-%%d %%H:%%i') AS barDateTimeMaxS,
			COUNT(b.symbol)                                       AS `+"`count`"+`
		FROM fx_symbol c
		LEFT JOIN fx_bar_%s b ON b.symbol = c.symbol
		WHERE (c.deleted+0) = 0 AND c.symbol_type = ?
		GROUP BY c.symbol
		ORDER BY c.sort_order`, barType.Suffix())

	var recs []fxBarDataStatusRecord
	if err := r.db.SelectContext(ctx, &recs, query, symbolType); err != nil {
		return nil, err
	}
	result := make([]fxmodel.BarDataStatus, len(recs))
	for i, rec := range recs {
		result[i] = fxmodel.BarDataStatus{
			Symbol:          rec.Symbol,
			BarDateTimeMinS: rec.BarDateTimeMinS,
			BarDateTimeMaxS: rec.BarDateTimeMaxS,
			Count:           rec.Count,
		}
	}
	return result, nil
}

func (r *MySQLBarDataRepository) SearchCount(ctx context.Context, symbol string, barType fxmodel.BarType, barDateFrom, barDateTo string) (int, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM fx_bar_%s B
		INNER JOIN fx_symbol C ON C.symbol = B.symbol
		WHERE B.symbol = ?`, suffix)
	args := []any{symbol}
	query, args = appendDateRange(query, args, barDateFrom, barDateTo)
	query += ` AND (C.deleted+0) = 0`

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *MySQLBarDataRepository) Search(ctx context.Context, symbol string, barType fxmodel.BarType, barDateFrom, barDateTo string, sortAsc bool, page, size int) ([]fxmodel.BarData, error) {
	suffix := barType.Suffix()
	query := fmt.Sprintf(`
		SELECT
			B.symbol                                   AS symbol,
			B.bar_date_time                            AS barDateTime,
			ROUND(B.open_price,   C.valid_scale)       AS openPrice,
			ROUND(B.high_price,   C.valid_scale)       AS highPrice,
			ROUND(B.low_price,    C.valid_scale)       AS lowPrice,
			ROUND(B.close_price,  C.valid_scale)       AS closePrice,
			B.volume                                   AS volume,
			ROUND(B.high_profit,  C.valid_scale)       AS highProfit,
			ROUND(B.low_profit,   C.valid_scale)       AS lowProfit,
			ROUND(B.close_profit, C.valid_scale)       AS closeProfit,
			ROUND(B.range_profit, C.valid_scale)       AS rangeProfit,
			COALESCE(R.rsi_value, 0)                   AS rsiValue,
			COALESCE(R.rsi_ma,    0)                   AS rsiMa
		FROM fx_bar_%s B
		INNER JOIN fx_symbol C ON C.symbol = B.symbol
		LEFT JOIN fx_bar_%s_rsi R
			ON R.symbol = B.symbol AND R.bar_date_time = B.bar_date_time AND R.rsi_range = 14
		WHERE B.symbol = ?`, suffix, suffix)
	args := []any{symbol}
	query, args = appendDateRange(query, args, barDateFrom, barDateTo)
	query += ` AND (C.deleted+0) = 0`

	if sortAsc {
		query += ` ORDER BY B.bar_date_time`
	} else {
		query += ` ORDER BY B.bar_date_time DESC`
	}

	offset := (page - 1) * size
	if offset > 0 {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, size, offset)
	} else {
		query += ` LIMIT ?`
		args = append(args, size)
	}

	var recs []fxBarDataRecord
	if err := r.db.SelectContext(ctx, &recs, query, args...); err != nil {
		return nil, err
	}
	result := make([]fxmodel.BarData, len(recs))
	for i, rec := range recs {
		result[i] = fxmodel.BarData{
			Symbol:      rec.Symbol,
			BarDateTime: fxmodel.LocalDateTime{Time: rec.BarDateTime},
			OpenPrice:   rec.OpenPrice,
			HighPrice:   rec.HighPrice,
			LowPrice:    rec.LowPrice,
			ClosePrice:  rec.ClosePrice,
			Volume:      rec.Volume,
			HighProfit:  rec.HighProfit,
			LowProfit:   rec.LowProfit,
			CloseProfit: rec.CloseProfit,
			RangeProfit: rec.RangeProfit,
			RsiValue:    rec.RsiValue,
			RsiMa:       rec.RsiMa,
		}
	}
	return result, nil
}

// appendDateRange はyyyyMMdd形式の日付範囲条件をクエリに追加する。
func appendDateRange(query string, args []any, barDateFrom, barDateTo string) (string, []any) {
	if barDateFrom != "" && barDateTo != "" {
		query += ` AND B.bar_date_time BETWEEN STR_TO_DATE(CONCAT(?, '000000'), '%Y%m%d%H%i%s') AND STR_TO_DATE(CONCAT(?, '235959'), '%Y%m%d%H%i%s')`
		args = append(args, barDateFrom, barDateTo)
	} else if barDateFrom != "" {
		query += ` AND B.bar_date_time >= STR_TO_DATE(CONCAT(?, '000000'), '%Y%m%d%H%i%s')`
		args = append(args, barDateFrom)
	} else if barDateTo != "" {
		query += ` AND B.bar_date_time <= STR_TO_DATE(CONCAT(?, '235959'), '%Y%m%d%H%i%s')`
		args = append(args, barDateTo)
	}
	return query, args
}
