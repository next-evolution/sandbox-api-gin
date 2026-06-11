package fxusecase

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type SearchSymbolUseCase struct {
	repo fxrepository.SymbolRepository
}

func NewSearchSymbolUseCase(repo fxrepository.SymbolRepository) *SearchSymbolUseCase {
	return &SearchSymbolUseCase{repo: repo}
}

type SymbolSearchResult struct {
	TotalCount int
	SymbolList []fxdto.SymbolDto
	Page       int
	Size       int
}

func (r SymbolSearchResult) TotalPage() int {
	if r.TotalCount == 0 {
		return 0
	}
	return (r.TotalCount + r.Size - 1) / r.Size
}

func (uc *SearchSymbolUseCase) Execute(ctx context.Context, symbolType string, page, size int) (SymbolSearchResult, error) {
	count, err := uc.repo.Count(ctx, symbolType)
	if err != nil {
		return SymbolSearchResult{}, err
	}

	list := make([]fxdto.SymbolDto, 0)
	if count > 0 {
		symbols, err := uc.repo.Search(ctx, symbolType, page, size)
		if err != nil {
			return SymbolSearchResult{}, err
		}
		for _, s := range symbols {
			list = append(list, fxdto.SymbolDtoFromDomain(s))
		}
	}

	return SymbolSearchResult{
		TotalCount: count,
		SymbolList: list,
		Page:       page,
		Size:       size,
	}, nil
}
