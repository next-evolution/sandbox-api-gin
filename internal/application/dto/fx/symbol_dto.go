package fxdto

import (
	"time"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type SymbolDto struct {
	Symbol           string  `json:"symbol" binding:"required"`
	SymbolType       string  `json:"symbolType" binding:"required"`
	Name             string  `json:"name" binding:"required"`
	ValidScale       int16   `json:"validScale"`
	TargetVolatility float64 `json:"targetVolatility"`
	SortOrder        int     `json:"sortOrder"`
}

func SymbolDtoFromDomain(s fxmodel.Symbol) SymbolDto {
	return SymbolDto{
		Symbol:           s.Symbol,
		SymbolType:       s.SymbolType,
		Name:             s.Name,
		ValidScale:       s.ValidScale,
		TargetVolatility: s.TargetVolatility,
		SortOrder:        s.SortOrder,
	}
}

func (d SymbolDto) ToDomain(author string) fxmodel.Symbol {
	now := time.Now()
	return fxmodel.Symbol{
		Symbol:           d.Symbol,
		SymbolType:       d.SymbolType,
		Name:             d.Name,
		ValidScale:       d.ValidScale,
		TargetVolatility: d.TargetVolatility,
		SortOrder:        d.SortOrder,
		Deleted:          false,
		CreatedAt:        now,
		CreatedBy:        author,
		UpdatedAt:        now,
		UpdatedBy:        author,
	}
}
