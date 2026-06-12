package fxrepository

import (
	"context"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type SummerTimeRepository interface {
	Count(ctx context.Context) (int, error)
	Search(ctx context.Context, page, size int) ([]fxmodel.SummerTime, error)
	Get(ctx context.Context, targetYear int16) (*fxmodel.SummerTime, error)
	Exists(ctx context.Context, targetYear int16) (bool, error)
	Add(ctx context.Context, summerTime fxmodel.SummerTime) error
	Update(ctx context.Context, summerTime fxmodel.SummerTime) error
	UpdateYear(ctx context.Context, summerTime fxmodel.SummerTime, baseYear int16) error
}
