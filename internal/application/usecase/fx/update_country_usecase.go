package fxusecase

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type UpdateCountryUseCase struct {
	repo fxrepository.CountryRepository
}

func NewUpdateCountryUseCase(repo fxrepository.CountryRepository) *UpdateCountryUseCase {
	return &UpdateCountryUseCase{repo: repo}
}

func (uc *UpdateCountryUseCase) Execute(ctx context.Context, baseCode string, dto fxdto.CountryDto) error {
	if baseCode == dto.Code {
		exists, err := uc.repo.Exists(ctx, baseCode)
		if err != nil {
			return err
		}
		if !exists {
			return apperror.NewUpdateError(baseCode)
		}
		if err := uc.repo.Update(ctx, dto.ToDomain("UpdateCountryUseCase")); err != nil {
			return err
		}
	} else {
		exists, err := uc.repo.Exists(ctx, dto.Code)
		if err != nil {
			return err
		}
		if exists {
			return apperror.NewDuplicateError(dto.Code)
		}
		if err := uc.repo.UpdateCode(ctx, dto.ToDomain("UpdateCountryUseCase"), baseCode); err != nil {
			return err
		}
	}

	return uc.repo.RefreshCache(ctx)
}
