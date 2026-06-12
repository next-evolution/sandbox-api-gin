package fxcommand

import "sandbox-api-gin/internal/domain/model/fx"

type TradeSimulationCommand struct {
	RiskAmount    float64
	FirstLotRatio float64
	Entry         *fx.TradeEntry
	PositionList  []*fx.TradePosition
}
