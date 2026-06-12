package fxrepository

import (
	"context"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type BarDataRepository interface {
	StatusList(ctx context.Context, symbolType string, barType fxmodel.BarType) ([]fxmodel.BarDataStatus, error)
	SearchCount(ctx context.Context, symbol string, barType fxmodel.BarType, barDateFrom, barDateTo string) (int, error)
	Search(ctx context.Context, symbol string, barType fxmodel.BarType, barDateFrom, barDateTo string, sortAsc bool, page, size int) ([]fxmodel.BarData, error)
}
