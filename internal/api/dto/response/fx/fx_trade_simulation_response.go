package fxresponse

import (
	"sandbox-api-gin/internal/api/dto/response"
	"sandbox-api-gin/internal/domain/model/fx"
)

type TradeSimulationResponse struct {
	response.ApiResponse
	Entry        *fx.TradeEntry      `json:"entry"`
	PositionList []*fx.TradePosition `json:"positionList"`
}
