package fxusecase

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type SearchCountryUseCase struct {
	repo fxrepository.CountryRepository
}

func NewSearchCountryUseCase(repo fxrepository.CountryRepository) *SearchCountryUseCase {
	return &SearchCountryUseCase{repo: repo}
}

type CountrySearchResult struct {
	TotalCount  int
	CountryList []fxdto.CountryDto
	Page        int
	Size        int
}

func (r CountrySearchResult) TotalPage() int {
	if r.TotalCount == 0 {
		return 0
	}
	return (r.TotalCount + r.Size - 1) / r.Size
}

func (uc *SearchCountryUseCase) Execute(ctx context.Context, page, size int) (CountrySearchResult, error) {
	count, err := uc.repo.Count(ctx)
	if err != nil {
		return CountrySearchResult{}, err
	}

	list := make([]fxdto.CountryDto, 0)
	if count > 0 {
		countries, err := uc.repo.Search(ctx, page, size)
		if err != nil {
			return CountrySearchResult{}, err
		}
		for _, c := range countries {
			list = append(list, fxdto.CountryDtoFromDomain(c))
		}
	}

	return CountrySearchResult{
		TotalCount:  count,
		CountryList: list,
		Page:        page,
		Size:        size,
	}, nil
}
