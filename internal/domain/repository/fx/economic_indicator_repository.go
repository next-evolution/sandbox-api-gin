package fxrepository

import (
	"context"

	"sandbox-api-gin/internal/domain/model"
)

type EconomicIndicatorRepository interface {
	GetList(ctx context.Context, countryCode string) ([]model.KeyValue, error)
}
