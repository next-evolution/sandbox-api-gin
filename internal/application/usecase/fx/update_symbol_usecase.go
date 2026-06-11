package fxusecase

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type UpdateSymbolUseCase struct {
	repo fxrepository.SymbolRepository
}

func NewUpdateSymbolUseCase(repo fxrepository.SymbolRepository) *UpdateSymbolUseCase {
	return &UpdateSymbolUseCase{repo: repo}
}

func (uc *UpdateSymbolUseCase) Execute(ctx context.Context, baseSymbol string, dto fxdto.SymbolDto) error {
	domainSymbol := dto.ToDomain("UpdateSymbolUseCase")

	if baseSymbol == dto.Symbol {
		exists, err := uc.repo.Exists(ctx, baseSymbol)
		if err != nil {
			return err
		}
		if !exists {
			return apperror.NewUpdateError(baseSymbol)
		}
		if err := uc.repo.Update(ctx, domainSymbol); err != nil {
			return err
		}
	} else {
		exists, err := uc.repo.Exists(ctx, dto.Symbol)
		if err != nil {
			return err
		}
		if exists {
			return apperror.NewDuplicateError(dto.Symbol)
		}
		if err := uc.repo.UpdateSymbol(ctx, domainSymbol, baseSymbol); err != nil {
			return err
		}
	}

	return uc.repo.RefreshCache(ctx, dto.SymbolType)
}
