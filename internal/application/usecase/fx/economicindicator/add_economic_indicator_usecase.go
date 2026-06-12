package economicindicator

import (
	"context"
	"fmt"
	"time"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type AddEconomicIndicatorUseCase struct {
	repo fxrepository.EconomicIndicatorRepository
}

func NewAddEconomicIndicatorUseCase(repo fxrepository.EconomicIndicatorRepository) *AddEconomicIndicatorUseCase {
	return &AddEconomicIndicatorUseCase{repo: repo}
}

func (uc *AddEconomicIndicatorUseCase) Execute(ctx context.Context, dto fxdto.EconomicIndicatorDto) error {
	exists, err := uc.repo.Exists(ctx, dto.CountryCode, dto.Name)
	if err != nil {
		return err
	}
	if exists {
		return apperror.NewDuplicateError(fmt.Sprintf("(%s) %s", dto.CountryCode, dto.Name))
	}

	now := time.Now()
	indicator := fxmodel.EconomicIndicator{
		CountryCode: dto.CountryCode,
		Name:        dto.Name,
		Importance:  dto.Importance,
		Description: dto.Description,
		UnitOfValue: dto.UnitOfValue,
		Deleted:     false,
		CreatedAt:   now,
		CreatedBy:   "AddEconomicIndicatorUseCase",
		UpdatedAt:   now,
		UpdatedBy:   "AddEconomicIndicatorUseCase",
	}

	if err := uc.repo.Add(ctx, indicator); err != nil {
		return err
	}

	return uc.repo.RefreshCache(ctx, dto.CountryCode)
}
