package fxusecase

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type SearchSummerTimeUseCase struct {
	repo fxrepository.SummerTimeRepository
}

func NewSearchSummerTimeUseCase(repo fxrepository.SummerTimeRepository) *SearchSummerTimeUseCase {
	return &SearchSummerTimeUseCase{repo: repo}
}

type SummerTimeSearchResult struct {
	TotalCount     int
	SummerTimeList []fxdto.SummerTimeDto
	Page           int
	Size           int
}

func (r SummerTimeSearchResult) TotalPage() int {
	if r.TotalCount == 0 {
		return 0
	}
	return (r.TotalCount + r.Size - 1) / r.Size
}

func (uc *SearchSummerTimeUseCase) Execute(ctx context.Context, page, size int) (SummerTimeSearchResult, error) {
	count, err := uc.repo.Count(ctx)
	if err != nil {
		return SummerTimeSearchResult{}, err
	}

	list := make([]fxdto.SummerTimeDto, 0)
	if count > 0 {
		items, err := uc.repo.Search(ctx, page, size)
		if err != nil {
			return SummerTimeSearchResult{}, err
		}
		for _, s := range items {
			list = append(list, fxdto.SummerTimeDtoFromDomain(s))
		}
	}

	return SummerTimeSearchResult{
		TotalCount:     count,
		SummerTimeList: list,
		Page:           page,
		Size:           size,
	}, nil
}
