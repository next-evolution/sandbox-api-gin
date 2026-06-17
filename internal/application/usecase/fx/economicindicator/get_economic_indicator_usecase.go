package economicindicator

import (
	"context"
	"fmt"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type GetEconomicIndicatorUseCase struct {
	repo fxrepository.EconomicIndicatorRepository
}

func NewGetEconomicIndicatorUseCase(repo fxrepository.EconomicIndicatorRepository) *GetEconomicIndicatorUseCase {
	return &GetEconomicIndicatorUseCase{repo: repo}
}

func (uc *GetEconomicIndicatorUseCase) Execute(ctx context.Context, countryCode, code string) (fxdto.EconomicIndicatorDto, error) {
	indicator, err := uc.repo.Get(ctx, countryCode, code)
	if err != nil {
		return fxdto.EconomicIndicatorDto{}, err
	}
	if indicator == nil {
		return fxdto.EconomicIndicatorDto{}, apperror.NewNotFoundError(fmt.Sprintf("(%s) %s", countryCode, code))
	}
	return fxdto.EconomicIndicatorDtoFromDomain(*indicator), nil
}
