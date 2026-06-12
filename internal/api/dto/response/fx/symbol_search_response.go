package fxresponse

import (
	"sandbox-api-gin/internal/api/dto/response"
	fxdto "sandbox-api-gin/internal/application/dto/fx"
)

type SymbolSearchResponse struct {
	ReturnCode  response.ReturnCode `json:"returnCode"`
	Message     string              `json:"message,omitempty"`
	TotalCount  int                 `json:"totalCount"`
	SearchCount int                 `json:"searchCount"`
	TotalPage   int                 `json:"totalPage"`
	List        []fxdto.SymbolDto   `json:"list"`
}
