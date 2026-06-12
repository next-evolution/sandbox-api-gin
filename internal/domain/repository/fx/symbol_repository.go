package fxrepository

import (
	"context"

	"sandbox-api-gin/internal/domain/model"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type SymbolRepository interface {
	GetList(ctx context.Context, symbolType string) ([]model.KeyValue, error)
	GetTradingSymbols(ctx context.Context) ([]string, error)
	Count(ctx context.Context, symbolType string) (int, error)
	Search(ctx context.Context, symbolType string, page, size int) ([]fxmodel.Symbol, error)
	Get(ctx context.Context, symbol string) (*fxmodel.Symbol, error)
	Exists(ctx context.Context, symbol string) (bool, error)
	Add(ctx context.Context, symbol fxmodel.Symbol) error
	Update(ctx context.Context, symbol fxmodel.Symbol) error
	UpdateSymbol(ctx context.Context, symbol fxmodel.Symbol, baseSymbol string) error
	RefreshCache(ctx context.Context, symbolType string) error
}
