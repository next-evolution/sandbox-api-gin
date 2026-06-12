package economicindicator

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type SearchEconomicIndicatorUseCase struct {
	repo fxrepository.EconomicIndicatorRepository
}

func NewSearchEconomicIndicatorUseCase(repo fxrepository.EconomicIndicatorRepository) *SearchEconomicIndicatorUseCase {
	return &SearchEconomicIndicatorUseCase{repo: repo}
}

type EconomicIndicatorSearchResult struct {
	TotalCount int
	List       []fxdto.EconomicIndicatorDto
	Page       int
	Size       int
}

func (r EconomicIndicatorSearchResult) TotalPage() int {
	if r.TotalCount == 0 {
		return 0
	}
	return (r.TotalCount + r.Size - 1) / r.Size
}

func (uc *SearchEconomicIndicatorUseCase) Execute(ctx context.Context, page, size int, countryCode, importance, name string) (EconomicIndicatorSearchResult, error) {
	count, err := uc.repo.Count(ctx, countryCode, importance, name)
	if err != nil {
		return EconomicIndicatorSearchResult{}, err
	}

	list := make([]fxdto.EconomicIndicatorDto, 0)
	if count > 0 {
		indicators, err := uc.repo.Search(ctx, page, size, countryCode, importance, name)
		if err != nil {
			return EconomicIndicatorSearchResult{}, err
		}
		for _, m := range indicators {
			list = append(list, fxdto.EconomicIndicatorDtoFromDomain(m))
		}
	}

	return EconomicIndicatorSearchResult{
		TotalCount: count,
		List:       list,
		Page:       page,
		Size:       size,
	}, nil
}
