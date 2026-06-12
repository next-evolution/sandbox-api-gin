package summertime

import (
	"context"
	"fmt"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type AddSummerTimeUseCase struct {
	repo fxrepository.SummerTimeRepository
}

func NewAddSummerTimeUseCase(repo fxrepository.SummerTimeRepository) *AddSummerTimeUseCase {
	return &AddSummerTimeUseCase{repo: repo}
}

func (uc *AddSummerTimeUseCase) Execute(ctx context.Context, dto fxdto.SummerTimeDto) error {
	exists, err := uc.repo.Exists(ctx, dto.TargetYear)
	if err != nil {
		return err
	}
	if exists {
		return apperror.NewDuplicateError(fmt.Sprintf("%d", dto.TargetYear))
	}

	return uc.repo.Add(ctx, dto.ToDomain())
}
