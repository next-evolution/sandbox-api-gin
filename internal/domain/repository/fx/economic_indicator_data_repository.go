package fxrepository

import (
	"context"
	"time"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type EconomicIndicatorDataRepository interface {
	Count(ctx context.Context, code, importance, countryCode, publicationBaseDate string) (int, error)
	Search(ctx context.Context, code, importance, countryCode, publicationBaseDate string, page, size int, sortAsc bool) ([]fxmodel.EconomicIndicatorData, error)
	Get(ctx context.Context, code, countryCode string, publication time.Time) (*fxmodel.EconomicIndicatorData, error)
	Exists(ctx context.Context, code, countryCode string, publication time.Time) (bool, error)
	Add(ctx context.Context, data fxmodel.EconomicIndicatorData) error
	Update(ctx context.Context, data fxmodel.EconomicIndicatorData, publication time.Time) error
	UpdateCode(ctx context.Context, data fxmodel.EconomicIndicatorData, code, countryCode string, publication time.Time) error
	DeleteLoad(ctx context.Context) error
	InsertLoad(ctx context.Context, data fxmodel.EconomicIndicatorData) error
	LoadDiff(ctx context.Context) ([]fxmodel.EconomicIndicatorData, error)
	InsertFromLoad(ctx context.Context) error
}
