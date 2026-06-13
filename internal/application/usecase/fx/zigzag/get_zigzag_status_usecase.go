package zigzagusecase

import (
	"context"

	fxcommand "sandbox-api-gin/internal/application/command/fx"
	"sandbox-api-gin/internal/domain/model/fx/zigzag"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type GetZigZagStatusUseCase struct {
	repo fxrepository.ZigZagRepository
}

func NewGetZigZagStatusUseCase(repo fxrepository.ZigZagRepository) *GetZigZagStatusUseCase {
	return &GetZigZagStatusUseCase{repo: repo}
}

func (uc *GetZigZagStatusUseCase) Execute(ctx context.Context, cmd fxcommand.ZigZagStatusCommand) ([]*zigzag.ZigZagStatus, error) {
	return uc.repo.GetStatusList(ctx, cmd.SymbolType, cmd.BarType, cmd.Depth)
}
