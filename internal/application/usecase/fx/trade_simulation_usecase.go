package fxusecase

import (
	"context"

	fxcommand "sandbox-api-gin/internal/application/command/fx"
	"sandbox-api-gin/internal/domain/model/fx"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
	fxservice "sandbox-api-gin/internal/domain/service/fx"
)

type SimulationResult struct {
	Entry        *fx.TradeEntry
	PositionList []*fx.TradePosition
}

type TradeSimulationUseCase struct {
	tradeSimulationRepo fxrepository.TradeSimulationRepository
	calculator          *fxservice.FxTradeCalculator
}

func NewTradeSimulationUseCase(
	tradeSimulationRepo fxrepository.TradeSimulationRepository,
	calculator *fxservice.FxTradeCalculator,
) *TradeSimulationUseCase {
	return &TradeSimulationUseCase{
		tradeSimulationRepo: tradeSimulationRepo,
		calculator:          calculator,
	}
}

func (uc *TradeSimulationUseCase) Execute(ctx context.Context, cmd *fxcommand.TradeSimulationCommand) (*SimulationResult, error) {
	entry := cmd.Entry
	positionList := cmd.PositionList

	priceInfo, err := uc.tradeSimulationRepo.GetPrice(ctx, entry.Symbol, entry.ContractAt.Time)
	if err != nil {
		return nil, err
	}

	filteredPositions, err := uc.initialize(entry, positionList, priceInfo)
	if err != nil {
		return nil, err
	}

	uc.calculator.Calculate(entry, filteredPositions, cmd.RiskAmount, cmd.FirstLotRatio)

	return &SimulationResult{Entry: entry, PositionList: filteredPositions}, nil
}

func (uc *TradeSimulationUseCase) initialize(
	entry *fx.TradeEntry,
	positionList []*fx.TradePosition,
	priceInfo *fx.PriceInfo,
) ([]*fx.TradePosition, error) {
	entry.ApplyPrice(priceInfo)

	var filteredPositions []*fx.TradePosition

	if positionList[0].SettlementPrice == 0 {
		positionList[0].SettlementPrice = entry.ComputeDefaultSettlementPrice(0.6)
		if len(positionList) > 1 {
			positionList[1].SettlementPrice = entry.ComputeDefaultSettlementPrice(0.9)
			if len(positionList) == 3 {
				positionList[2].SettlementPrice = entry.ComputeDefaultSettlementPrice(1.2)
			}
		}
		filteredPositions = make([]*fx.TradePosition, len(positionList))
		copy(filteredPositions, positionList)
	} else {
		var positionNumber int16 = 1
		for _, pos := range positionList {
			if pos.SettlementPrice > 0 {
				pos.PositionNumber = positionNumber
				filteredPositions = append(filteredPositions, pos)
				positionNumber++
			}
		}
	}

	entry.ApplyDefaultLossPrice()

	return filteredPositions, nil
}
