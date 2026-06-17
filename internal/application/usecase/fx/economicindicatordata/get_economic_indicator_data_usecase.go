package economicindicatordata

import (
	"context"
	"fmt"
	"time"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type GetEconomicIndicatorDataUseCase struct {
	repo fxrepository.EconomicIndicatorDataRepository
}

func NewGetEconomicIndicatorDataUseCase(repo fxrepository.EconomicIndicatorDataRepository) *GetEconomicIndicatorDataUseCase {
	return &GetEconomicIndicatorDataUseCase{repo: repo}
}

func (uc *GetEconomicIndicatorDataUseCase) Execute(ctx context.Context, code, countryCode string, publication time.Time) (fxdto.EconomicIndicatorDataDto, error) {
	data, err := uc.repo.Get(ctx, code, countryCode, publication)
	if err != nil {
		return fxdto.EconomicIndicatorDataDto{}, err
	}
	if data == nil {
		return fxdto.EconomicIndicatorDataDto{}, apperror.NewNotFoundError(fmt.Sprintf("(%s) %s / %s", countryCode, code, publication.Format("2006-01-02 15:04:05")))
	}
	return fxdto.EconomicIndicatorDataDtoFromDomain(*data), nil
}
