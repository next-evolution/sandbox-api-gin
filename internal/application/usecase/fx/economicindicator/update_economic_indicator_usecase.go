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

type UpdateEconomicIndicatorUseCase struct {
	repo fxrepository.EconomicIndicatorRepository
}

func NewUpdateEconomicIndicatorUseCase(repo fxrepository.EconomicIndicatorRepository) *UpdateEconomicIndicatorUseCase {
	return &UpdateEconomicIndicatorUseCase{repo: repo}
}

func (uc *UpdateEconomicIndicatorUseCase) Execute(ctx context.Context, countryCode, code string, dto fxdto.EconomicIndicatorDto) error {
	existing, err := uc.repo.Get(ctx, countryCode, code)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperror.NewNotFoundError(fmt.Sprintf("(%s) %s", countryCode, code))
	}

	newCountryCode := dto.CountryCode
	if countryCode != newCountryCode {
		dup, err := uc.repo.Exists(ctx, newCountryCode, dto.Name)
		if err != nil {
			return err
		}
		if dup {
			return apperror.NewDuplicateError(fmt.Sprintf("(%s) %s", newCountryCode, dto.Name))
		}
	}

	toUpdate := fxmodel.EconomicIndicator{
		Code:        code,
		CountryCode: newCountryCode,
		Name:        dto.Name,
		Importance:  dto.Importance,
		Description: dto.Description,
		UnitOfValue: dto.UnitOfValue,
		UpdatedAt:   time.Now(),
		UpdatedBy:   "UpdateEconomicIndicatorUseCase",
	}

	if err := uc.repo.Update(ctx, toUpdate, countryCode); err != nil {
		return err
	}

	if err := uc.repo.RefreshCache(ctx, countryCode); err != nil {
		return err
	}
	if countryCode != newCountryCode {
		return uc.repo.RefreshCache(ctx, newCountryCode)
	}
	return nil
}
