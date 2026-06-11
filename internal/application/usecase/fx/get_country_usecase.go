package fxusecase

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type GetCountryUseCase struct {
	repo fxrepository.CountryRepository
}

func NewGetCountryUseCase(repo fxrepository.CountryRepository) *GetCountryUseCase {
	return &GetCountryUseCase{repo: repo}
}

func (uc *GetCountryUseCase) Get(ctx context.Context, code string) (fxdto.CountryDto, error) {
	country, err := uc.repo.Get(ctx, code)
	if err != nil {
		return fxdto.CountryDto{}, err
	}
	if country == nil {
		return fxdto.CountryDto{}, apperror.NewNotFoundError(code)
	}
	return fxdto.CountryDtoFromDomain(*country), nil
}
