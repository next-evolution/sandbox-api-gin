package fxrepository

import (
	"context"

	"sandbox-api-gin/internal/domain/model"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type CountryRepository interface {
	GetList(ctx context.Context) ([]model.KeyValue, error)
	Count(ctx context.Context) (int, error)
	Search(ctx context.Context, page, size int) ([]fxmodel.Country, error)
	Get(ctx context.Context, code string) (*fxmodel.Country, error)
	Exists(ctx context.Context, code string) (bool, error)
	Add(ctx context.Context, country fxmodel.Country) error
	Update(ctx context.Context, country fxmodel.Country) error
	UpdateCode(ctx context.Context, country fxmodel.Country, baseCode string) error
	RefreshCache(ctx context.Context) error
}
