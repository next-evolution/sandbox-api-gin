package fxusecase

import (
	"context"

	"sandbox-api-gin/internal/domain/repository"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

var symbolTypes = []string{"Trade", "Analyze"}

type MasterRefreshUseCase struct {
	countryRepo           fxrepository.CountryRepository
	symbolRepo            fxrepository.SymbolRepository
	economicIndicatorRepo fxrepository.EconomicIndicatorRepository
	masterCacheRepo       repository.MasterCacheRepository
}

func NewMasterRefreshUseCase(
	countryRepo fxrepository.CountryRepository,
	symbolRepo fxrepository.SymbolRepository,
	economicIndicatorRepo fxrepository.EconomicIndicatorRepository,
	masterCacheRepo repository.MasterCacheRepository,
) *MasterRefreshUseCase {
	return &MasterRefreshUseCase{
		countryRepo:           countryRepo,
		symbolRepo:            symbolRepo,
		economicIndicatorRepo: economicIndicatorRepo,
		masterCacheRepo:       masterCacheRepo,
	}
}

func (uc *MasterRefreshUseCase) Execute(ctx context.Context) (string, error) {
	if err := uc.countryRepo.RefreshCache(ctx); err != nil {
		return "", err
	}

	for _, st := range symbolTypes {
		if err := uc.symbolRepo.RefreshCache(ctx, st); err != nil {
			return "", err
		}
	}

	countries, err := uc.countryRepo.GetList(ctx)
	if err != nil {
		return "", err
	}
	for _, c := range countries {
		if err := uc.economicIndicatorRepo.RefreshCache(ctx, c.Key); err != nil {
			return "", err
		}
	}

	if err := uc.masterCacheRepo.DeleteByPattern(ctx, "price*"); err != nil {
		return "", err
	}

	return uc.masterCacheRepo.GetStatus(ctx)
}
