package country

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type AddCountryUseCase struct {
	repo fxrepository.CountryRepository
}

func NewAddCountryUseCase(repo fxrepository.CountryRepository) *AddCountryUseCase {
	return &AddCountryUseCase{repo: repo}
}

func (uc *AddCountryUseCase) Execute(ctx context.Context, dto fxdto.CountryDto) error {
	exists, err := uc.repo.Exists(ctx, dto.Code)
	if err != nil {
		return err
	}
	if exists {
		return apperror.NewDuplicateError(dto.Code)
	}

	if err := uc.repo.Add(ctx, dto.ToDomain("AddCountryUseCase")); err != nil {
		return err
	}

	return uc.repo.RefreshCache(ctx)
}
