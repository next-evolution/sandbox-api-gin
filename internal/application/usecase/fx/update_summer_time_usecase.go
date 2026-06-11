package fxusecase

import (
	"context"
	"fmt"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type UpdateSummerTimeUseCase struct {
	repo fxrepository.SummerTimeRepository
}

func NewUpdateSummerTimeUseCase(repo fxrepository.SummerTimeRepository) *UpdateSummerTimeUseCase {
	return &UpdateSummerTimeUseCase{repo: repo}
}

func (uc *UpdateSummerTimeUseCase) Execute(ctx context.Context, baseYear int16, dto fxdto.SummerTimeDto) error {
	if baseYear == dto.TargetYear {
		exists, err := uc.repo.Exists(ctx, baseYear)
		if err != nil {
			return err
		}
		if !exists {
			return apperror.NewUpdateError(fmt.Sprintf("%d", baseYear))
		}
		return uc.repo.Update(ctx, dto.ToDomain())
	}

	// 年変更の場合: 新しい年の重複チェック（Country パターンに合わせて修正版）
	exists, err := uc.repo.Exists(ctx, dto.TargetYear)
	if err != nil {
		return err
	}
	if exists {
		return apperror.NewDuplicateError(fmt.Sprintf("%d", dto.TargetYear))
	}
	return uc.repo.UpdateYear(ctx, dto.ToDomain(), baseYear)
}
