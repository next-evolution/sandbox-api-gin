package fxusecase

import (
	"context"
	"fmt"

	fxcommand "sandbox-api-gin/internal/application/command/fx"
	fxdto "sandbox-api-gin/internal/application/dto/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type StatusBarDataUseCase struct {
	repo fxrepository.BarDataRepository
}

func NewStatusBarDataUseCase(repo fxrepository.BarDataRepository) *StatusBarDataUseCase {
	return &StatusBarDataUseCase{repo: repo}
}

func (uc *StatusBarDataUseCase) Execute(ctx context.Context, cmd fxcommand.StatusBarDataCommand) ([]fxdto.BarDataImportResult, error) {
	statuses, err := uc.repo.StatusList(ctx, cmd.SymbolType, cmd.BarType)
	if err != nil {
		return nil, err
	}

	results := make([]fxdto.BarDataImportResult, len(statuses))
	for i, s := range statuses {
		results[i] = fxdto.BarDataImportResult{
			Symbol:      s.Symbol,
			ExistsCount: s.Count,
			Message:     buildStatusMessage(s.BarDateTimeMinS, s.BarDateTimeMaxS),
		}
	}
	return results, nil
}

func buildStatusMessage(minS, maxS *string) *string {
	if minS == nil && maxS == nil {
		return nil
	}
	msg := fmt.Sprintf("%s~%s",
		strOrEmpty(minS),
		strOrEmpty(maxS),
	)
	return &msg
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
