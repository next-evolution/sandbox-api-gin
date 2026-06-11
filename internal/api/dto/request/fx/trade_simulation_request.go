package fxrequest

import (
	"time"

	"sandbox-api-gin/internal/domain/model/fx"
)

type TradeSimulationRequest struct {
	RiskAmount    float64        `json:"riskAmount"`
	FirstLotRatio float64        `json:"firstLotRatio" binding:"gt=0"`
	Entry         EntryParam     `json:"entry" binding:"required"`
	PositionList  []PositionParam `json:"positionList" binding:"required,min=1"`
}

type EntryParam struct {
	ID            *int64     `json:"id"`
	TradeVersion  string     `json:"tradeVersion" binding:"required"`
	EntryType     fx.EntryType `json:"entryType"`
	Symbol        string     `json:"symbol" binding:"required"`
	TradeType     fx.TradeType `json:"tradeType"`
	ContractAt    time.Time  `json:"contractAt"`
	FibonacciType string     `json:"fibonacciType" binding:"required"`
	FibonacciBar  string     `json:"fibonacciBar" binding:"required"`
	ContractPrice float64    `json:"contractPrice"`
	LossPrice     float64    `json:"lossPrice"`
	PositionRatio int        `json:"positionRatio"`
	PriceJpy      float64    `json:"priceJpy"`
	Lot           *float64   `json:"lot"`
	SettlementAmount int     `json:"settlementAmount"`
	LossPips      int        `json:"lossPips"`
	SettlementRatio *float64 `json:"settlementRatio"`
	Comment       *string    `json:"comment"`
	ImagePath     *string    `json:"imagePath"`
}

func (ep *EntryParam) ToDomain() *fx.TradeEntry {
	lot := 0.0
	if ep.Lot != nil {
		lot = *ep.Lot
	}
	settlementRatio := 0.0
	if ep.SettlementRatio != nil {
		settlementRatio = *ep.SettlementRatio
	}
	return &fx.TradeEntry{
		ID:               ep.ID,
		TradeVersion:     ep.TradeVersion,
		EntryType:        ep.EntryType,
		Symbol:           ep.Symbol,
		TradeType:        ep.TradeType,
		ContractAt:       fx.LocalDateTime{Time: ep.ContractAt},
		FibonacciType:    ep.FibonacciType,
		FibonacciBar:     ep.FibonacciBar,
		ContractPrice:    ep.ContractPrice,
		LossPrice:        ep.LossPrice,
		PositionRatio:    ep.PositionRatio,
		PriceJpy:         ep.PriceJpy,
		Lot:              lot,
		SettlementAmount: ep.SettlementAmount,
		LossPips:         ep.LossPips,
		SettlementRatio:  settlementRatio,
		Comment:          ep.Comment,
		ImagePath:        ep.ImagePath,
	}
}

type PositionParam struct {
	ID              *int64   `json:"id"`
	PositionNumber  int16    `json:"positionNumber" binding:"gt=0"`
	SettlementPrice float64  `json:"settlementPrice"`
	SettlementPips  int      `json:"settlementPips"`
	SettlementRatio *float64 `json:"settlementRatio"`
	Lot             *float64 `json:"lot"`
	ProfitAmount    int      `json:"profitAmount"`
	LossAmount      int      `json:"lossAmount"`
}

func (pp *PositionParam) ToDomain() *fx.TradePosition {
	lot := 0.0
	if pp.Lot != nil {
		lot = *pp.Lot
	}
	settlementRatio := 0.0
	if pp.SettlementRatio != nil {
		settlementRatio = *pp.SettlementRatio
	}
	return &fx.TradePosition{
		ID:              pp.ID,
		PositionNumber:  pp.PositionNumber,
		SettlementPrice: pp.SettlementPrice,
		SettlementPips:  pp.SettlementPips,
		SettlementRatio: settlementRatio,
		Lot:             lot,
		ProfitAmount:    pp.ProfitAmount,
		LossAmount:      pp.LossAmount,
	}
}
