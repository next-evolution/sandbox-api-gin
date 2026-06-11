package fxrequest

import fxdto "sandbox-api-gin/internal/application/dto/fx"

type SymbolRequest struct {
	Symbol fxdto.SymbolDto `json:"symbol" binding:"required"`
}

type SymbolSearchRequest struct {
	SymbolType string `json:"symbolType" binding:"required"`
	Page       int    `json:"page" binding:"required,min=1"`
	Size       int    `json:"size" binding:"required,min=1"`
}
