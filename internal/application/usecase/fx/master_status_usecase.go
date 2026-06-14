package fxusecase

import (
	"context"

	"sandbox-api-gin/internal/domain/repository"
)

type MasterStatusUseCase struct {
	masterCacheRepo repository.MasterCacheRepository
}

func NewMasterStatusUseCase(masterCacheRepo repository.MasterCacheRepository) *MasterStatusUseCase {
	return &MasterStatusUseCase{masterCacheRepo: masterCacheRepo}
}

func (uc *MasterStatusUseCase) Execute(ctx context.Context) (string, error) {
	return uc.masterCacheRepo.GetStatus(ctx)
}
