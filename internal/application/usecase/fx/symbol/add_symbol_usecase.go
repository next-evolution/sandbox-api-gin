package symbol

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type AddSymbolUseCase struct {
	repo fxrepository.SymbolRepository
}

func NewAddSymbolUseCase(repo fxrepository.SymbolRepository) *AddSymbolUseCase {
	return &AddSymbolUseCase{repo: repo}
}

func (uc *AddSymbolUseCase) Execute(ctx context.Context, dto fxdto.SymbolDto) error {
	exists, err := uc.repo.Exists(ctx, dto.Symbol)
	if err != nil {
		return err
	}
	if exists {
		return apperror.NewDuplicateError(dto.Symbol)
	}

	symbol := dto.ToDomain("AddSymbolUseCase")
	if err := uc.repo.Add(ctx, symbol); err != nil {
		return err
	}

	return uc.repo.RefreshCache(ctx, dto.SymbolType)
}
