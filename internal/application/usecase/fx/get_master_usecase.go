package fxusecase

import (
	"context"

	"sandbox-api-gin/internal/domain/model"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type GetMasterUseCase struct {
	symbolRepo             fxrepository.SymbolRepository
	countryRepo            fxrepository.CountryRepository
	economicIndicatorRepo  fxrepository.EconomicIndicatorRepository
}

func NewGetMasterUseCase(
	symbolRepo fxrepository.SymbolRepository,
	countryRepo fxrepository.CountryRepository,
	economicIndicatorRepo fxrepository.EconomicIndicatorRepository,
) *GetMasterUseCase {
	return &GetMasterUseCase{
		symbolRepo:            symbolRepo,
		countryRepo:           countryRepo,
		economicIndicatorRepo: economicIndicatorRepo,
	}
}

func (uc *GetMasterUseCase) Symbol(ctx context.Context, symbolType string) ([]model.KeyValue, error) {
	return uc.symbolRepo.GetList(ctx, symbolType)
}

func (uc *GetMasterUseCase) Country(ctx context.Context) ([]model.KeyValue, error) {
	return uc.countryRepo.GetList(ctx)
}

func (uc *GetMasterUseCase) CurrencyPair(ctx context.Context) ([]model.KeyValue, error) {
	return uc.symbolRepo.GetList(ctx, "Trade")
}

func (uc *GetMasterUseCase) CurrencyIndex(ctx context.Context) ([]model.KeyValue, error) {
	return uc.symbolRepo.GetList(ctx, "Analyze")
}

func (uc *GetMasterUseCase) EconomicIndicator(ctx context.Context, countryCode string) ([]model.KeyValue, error) {
	return uc.economicIndicatorRepo.GetList(ctx, countryCode)
}
