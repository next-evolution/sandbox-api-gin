package symbol

import (
	"context"

	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/apperror"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type GetSymbolUseCase struct {
	repo fxrepository.SymbolRepository
}

func NewGetSymbolUseCase(repo fxrepository.SymbolRepository) *GetSymbolUseCase {
	return &GetSymbolUseCase{repo: repo}
}

func (uc *GetSymbolUseCase) Get(ctx context.Context, symbol string) (fxdto.SymbolDto, error) {
	s, err := uc.repo.Get(ctx, symbol)
	if err != nil {
		return fxdto.SymbolDto{}, err
	}
	if s == nil {
		return fxdto.SymbolDto{}, apperror.NewNotFoundError(symbol)
	}
	return fxdto.SymbolDtoFromDomain(*s), nil
}
