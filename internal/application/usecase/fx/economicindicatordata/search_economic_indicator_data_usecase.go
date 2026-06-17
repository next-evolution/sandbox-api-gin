package economicindicatordata

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type SearchEconomicIndicatorDataUseCase struct {
	repo fxrepository.EconomicIndicatorDataRepository
}

func NewSearchEconomicIndicatorDataUseCase(repo fxrepository.EconomicIndicatorDataRepository) *SearchEconomicIndicatorDataUseCase {
	return &SearchEconomicIndicatorDataUseCase{repo: repo}
}

type EconomicIndicatorDataSearchResult struct {
	TotalCount int
	List       []fxdto.EconomicIndicatorDataDto
	Page       int
	Size       int
}

func (r EconomicIndicatorDataSearchResult) TotalPage() int {
	if r.TotalCount == 0 {
		return 0
	}
	return (r.TotalCount + r.Size - 1) / r.Size
}

func (uc *SearchEconomicIndicatorDataUseCase) Execute(ctx context.Context, code, importance, countryCode, publicationBaseDate string, page, size int, sortAsc bool) (EconomicIndicatorDataSearchResult, error) {
	count, err := uc.repo.Count(ctx, code, importance, countryCode, publicationBaseDate)
	if err != nil {
		return EconomicIndicatorDataSearchResult{}, err
	}

	list := make([]fxdto.EconomicIndicatorDataDto, 0)
	if count > 0 {
		items, err := uc.repo.Search(ctx, code, importance, countryCode, publicationBaseDate, page, size, sortAsc)
		if err != nil {
			return EconomicIndicatorDataSearchResult{}, err
		}
		for _, m := range items {
			list = append(list, fxdto.EconomicIndicatorDataDtoFromDomain(m))
		}
	}

	return EconomicIndicatorDataSearchResult{
		TotalCount: count,
		List:       list,
		Page:       page,
		Size:       size,
	}, nil
}
