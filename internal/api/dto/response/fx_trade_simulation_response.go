package response

import "sandbox-api-gin/internal/domain/model/fx"

type TradeSimulationResponse struct {
	ApiResponse
	Entry        *fx.TradeEntry    `json:"entry"`
	PositionList []*fx.TradePosition `json:"positionList"`
}
