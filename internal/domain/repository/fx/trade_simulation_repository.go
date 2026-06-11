package fxrepository

import (
	"context"
	"time"

	"sandbox-api-gin/internal/domain/model/fx"
)

type TradeSimulationRepository interface {
	GetPrice(ctx context.Context, symbol string, contractAt time.Time) (*fx.PriceInfo, error)
}
