package fxrepository

import (
	"context"

	"sandbox-api-gin/internal/domain/model"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type EconomicIndicatorRepository interface {
	GetList(ctx context.Context, countryCode string) ([]model.KeyValue, error)
	Count(ctx context.Context, countryCode, importance, name string) (int, error)
	Search(ctx context.Context, page, size int, countryCode, importance, name string) ([]fxmodel.EconomicIndicator, error)
	Get(ctx context.Context, id int64) (*fxmodel.EconomicIndicator, error)
	Exists(ctx context.Context, countryCode, name string) (bool, error)
	Add(ctx context.Context, indicator fxmodel.EconomicIndicator) error
	Update(ctx context.Context, indicator fxmodel.EconomicIndicator, countryCode string) error
	GetEconomicIndicatorList(ctx context.Context, countryCode string) ([]fxmodel.EconomicIndicator, error)
	RefreshCache(ctx context.Context, countryCode string) error
}
