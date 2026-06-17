package economicindicatordata

import (
	"context"
	"fmt"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type AddEconomicIndicatorDataUseCase struct {
	repo fxrepository.EconomicIndicatorDataRepository
}

func NewAddEconomicIndicatorDataUseCase(repo fxrepository.EconomicIndicatorDataRepository) *AddEconomicIndicatorDataUseCase {
	return &AddEconomicIndicatorDataUseCase{repo: repo}
}

func (uc *AddEconomicIndicatorDataUseCase) Execute(ctx context.Context, dto fxdto.EconomicIndicatorDataDto) error {
	exists, err := uc.repo.Exists(ctx, dto.Code, dto.CountryCode, dto.Publication.Time)
	if err != nil {
		return err
	}
	if exists {
		return apperror.NewDuplicateError(fmt.Sprintf("(%s) %s / %s", dto.CountryCode, dto.Code, dto.Publication.Format("2006-01-02 15:04:05")))
	}

	data := fxmodel.EconomicIndicatorData{
		Code:          dto.Code,
		CountryCode:   dto.CountryCode,
		Publication:   dto.Publication.Time,
		SubTitle:      dto.SubTitle,
		ResultValue:   dto.ResultValue,
		ForecastValue: dto.ForecastValue,
		PreviousValue: dto.PreviousValue,
		Memo:          dto.Memo,
	}

	return uc.repo.Add(ctx, data)
}
