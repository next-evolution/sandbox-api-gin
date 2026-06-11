package fxusecase

import (
	"context"

	fxcommand "sandbox-api-gin/internal/application/command/fx"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type SearchBarDataUseCase struct {
	repo fxrepository.BarDataRepository
}

func NewSearchBarDataUseCase(repo fxrepository.BarDataRepository) *SearchBarDataUseCase {
	return &SearchBarDataUseCase{repo: repo}
}

type BarDataSearchResult struct {
	TotalCount  int
	BarDataList []fxmodel.BarData
	Page        int
	Size        int
}

func (r BarDataSearchResult) TotalPage() int {
	if r.TotalCount == 0 {
		return 0
	}
	return (r.TotalCount + r.Size - 1) / r.Size
}

func (uc *SearchBarDataUseCase) Execute(ctx context.Context, cmd fxcommand.SearchBarDataCommand) (BarDataSearchResult, error) {
	count, err := uc.repo.SearchCount(ctx, cmd.Symbol, cmd.BarType, cmd.BarDateFrom, cmd.BarDateTo)
	if err != nil {
		return BarDataSearchResult{}, err
	}

	list := make([]fxmodel.BarData, 0)
	if count > 0 {
		list, err = uc.repo.Search(ctx, cmd.Symbol, cmd.BarType, cmd.BarDateFrom, cmd.BarDateTo, cmd.SortAsc, cmd.Page, cmd.Size)
		if err != nil {
			return BarDataSearchResult{}, err
		}
	}

	return BarDataSearchResult{
		TotalCount:  count,
		BarDataList: list,
		Page:        cmd.Page,
		Size:        cmd.Size,
	}, nil
}
