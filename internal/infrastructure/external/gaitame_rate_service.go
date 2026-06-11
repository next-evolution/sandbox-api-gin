package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	goredis "github.com/redis/go-redis/v9"

	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type GaitameRateService struct {
	symbolRepo  fxrepository.SymbolRepository
	httpClient  *http.Client
	baseURL     string
	redisClient *goredis.Client
}

func NewGaitameRateService(
	symbolRepo fxrepository.SymbolRepository,
	baseURL string,
	redisClient *goredis.Client,
) *GaitameRateService {
	return &GaitameRateService{
		symbolRepo:  symbolRepo,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		baseURL:     baseURL,
		redisClient: redisClient,
	}
}

func (s *GaitameRateService) RefreshRate(ctx context.Context, contractHm string) error {
	symbols, err := s.symbolRepo.GetTradingSymbols(ctx)
	if err != nil {
		return fmt.Errorf("シンボルリスト取得エラー: %w", err)
	}
	if len(symbols) == 0 {
		slog.Warn("シンボルリストが空のため、Gaitameレート取得をスキップします", "target", contractHm)
		return nil
	}

	symbolSet := make(map[string]struct{}, len(symbols))
	for _, s := range symbols {
		symbolSet[s] = struct{}{}
	}

	resp, err := s.httpClient.Get(s.baseURL + "/v3/info/prices/rate")
	if err != nil {
		return fmt.Errorf("Gaitameレート取得エラー。target=%s: %w", contractHm, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Gaitameレスポンスcloseエラー", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Gaitameレート取得エラー。target=%s status=%d", contractHm, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Gaitameレスポンス読み込みエラー。target=%s: %w", contractHm, err)
	}

	var dto gaitameRateDto
	if err := json.Unmarshal(body, &dto); err != nil {
		return fmt.Errorf("Gaitameレスポンスパースエラー。target=%s: %w", contractHm, err)
	}

	if len(dto.Data) == 0 {
		slog.Info("Gaitameレート取得: レスポンスが空", "target", contractHm)
		return nil
	}

	slog.Info("Gaitameレート取得", "target", contractHm, "size", len(dto.Data))

	for _, rate := range dto.Data {
		if _, ok := symbolSet[rate.Pair]; !ok {
			continue
		}
		key := fmt.Sprintf("price:%s_%s", rate.Pair, contractHm)
		if err := s.redisClient.Set(ctx, key, fmt.Sprintf("%g", rate.Open), 60*time.Minute).Err(); err != nil {
			slog.Warn("Gaitameレートキャッシュ保存失敗", "key", key, "error", err)
		}
	}

	return nil
}

type gaitameRateDto struct {
	Status int           `json:"status"`
	Data   []gaitameRate `json:"data"`
}

type gaitameRate struct {
	Pair string  `json:"pair"`
	Open float64 `json:"open"`
}
