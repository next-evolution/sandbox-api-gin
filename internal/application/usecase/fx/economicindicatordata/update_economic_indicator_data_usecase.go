package economicindicatordata

import (
	"context"
	"fmt"
	"time"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type UpdateEconomicIndicatorDataUseCase struct {
	repo fxrepository.EconomicIndicatorDataRepository
}

func NewUpdateEconomicIndicatorDataUseCase(repo fxrepository.EconomicIndicatorDataRepository) *UpdateEconomicIndicatorDataUseCase {
	return &UpdateEconomicIndicatorDataUseCase{repo: repo}
}

func (uc *UpdateEconomicIndicatorDataUseCase) Execute(ctx context.Context, economicIndicatorID int64, publication time.Time, dto fxdto.EconomicIndicatorDataDto) error {
	existing, err := uc.repo.Get(ctx, economicIndicatorID, publication)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperror.NewNotFoundError(publication.Format("2006/01/02 15:04"))
	}

	isIDDiff := economicIndicatorID != dto.ID
	isPublicationDiff := !publication.Equal(dto.Publication.Time)

	displayName := buildDisplayName(existing, publication, dto.Publication.Time, isPublicationDiff)

	toUpdate := fxmodel.EconomicIndicatorData{
		ID:            dto.ID,
		Publication:   dto.Publication.Time,
		SubTitle:      dto.SubTitle,
		ResultValue:   dto.ResultValue,
		ForecastValue: dto.ForecastValue,
		PreviousValue: dto.PreviousValue,
		Memo:          dto.Memo,
	}

	if isIDDiff {
		dup, err := uc.repo.Exists(ctx, dto.ID, dto.Publication.Time)
		if err != nil {
			return err
		}
		if dup {
			return apperror.NewDuplicateError(displayName)
		}
		return uc.repo.UpdateID(ctx, toUpdate, economicIndicatorID, publication)
	}

	if isPublicationDiff {
		dup, err := uc.repo.Exists(ctx, dto.ID, dto.Publication.Time)
		if err != nil {
			return err
		}
		if dup {
			return apperror.NewDuplicateError(displayName)
		}
	}
	return uc.repo.Update(ctx, toUpdate, publication)
}

func buildDisplayName(existing *fxmodel.EconomicIndicatorData, publication, newPublication time.Time, isPublicationDiff bool) string {
	const dtfLayout = "2006/01/02 15:04"
	if isPublicationDiff {
		return fmt.Sprintf("[%s] -> [%s]", publication.Format(dtfLayout), newPublication.Format(dtfLayout))
	}
	return fmt.Sprintf("[%s] (%s) %s", publication.Format(dtfLayout), existing.CountryNameShort, existing.Name)
}
