package fxusecase

import (
	"context"
	"fmt"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type GetSummerTimeUseCase struct {
	repo fxrepository.SummerTimeRepository
}

func NewGetSummerTimeUseCase(repo fxrepository.SummerTimeRepository) *GetSummerTimeUseCase {
	return &GetSummerTimeUseCase{repo: repo}
}

func (uc *GetSummerTimeUseCase) Get(ctx context.Context, targetYear int16) (fxdto.SummerTimeDto, error) {
	item, err := uc.repo.Get(ctx, targetYear)
	if err != nil {
		return fxdto.SummerTimeDto{}, err
	}
	if item == nil {
		return fxdto.SummerTimeDto{}, apperror.NewNotFoundError(fmt.Sprintf("%d", targetYear))
	}
	return fxdto.SummerTimeDtoFromDomain(*item), nil
}
