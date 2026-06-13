package zigzagusecase

import (
	"context"
	"fmt"
	"log/slog"

	fxcommand "sandbox-api-gin/internal/application/command/fx"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/model/fx/zigzag"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
	fxservice "sandbox-api-gin/internal/domain/service/fx"
)

type GenerateZigZagUseCase struct {
	repo          fxrepository.ZigZagRepository
	domainService *fxservice.ZigZagDomainService
}

func NewGenerateZigZagUseCase(
	repo fxrepository.ZigZagRepository,
	domainService *fxservice.ZigZagDomainService,
) *GenerateZigZagUseCase {
	return &GenerateZigZagUseCase{repo: repo, domainService: domainService}
}

type GenerateZigZagResult struct {
	Status *zigzag.ZigZagStatus
	Warn   bool
}

func (uc *GenerateZigZagUseCase) Execute(ctx context.Context, cmd fxcommand.ZigZagGenerateCommand) (*GenerateZigZagResult, error) {
	slog.Debug("GenerateZigZag", "cmd", cmd)

	status, err := uc.repo.GetStatus(ctx, cmd.BarType, cmd.Symbol, cmd.Depth)
	if err != nil {
		return nil, err
	}

	maxLimit, err := uc.repo.TargetBarCount(ctx, cmd.BarType, cmd.Symbol, cmd.BarDateTime)
	if err != nil {
		return nil, err
	}
	if maxLimit > cmd.LoadSize {
		maxLimit = cmd.LoadSize
	}

	previousList, err := uc.repo.PreviousList(ctx, cmd.BarType, cmd.Symbol, cmd.Depth, cmd.BarDateTime, cmd.Depth-1)
	if err != nil {
		return nil, err
	}

	if len(previousList) == 0 || uc.domainService.PreviousNotExists(previousList) {
		if len(previousList) == 0 {
			status.Message = "previousList empty."
		} else {
			status.Message = "previous not exists."
		}
		return &GenerateZigZagResult{Status: status, Warn: true}, nil
	}

	targetList, err := uc.repo.TargetList(ctx, cmd.BarType, cmd.Symbol, cmd.Depth, cmd.BarDateTime, maxLimit)
	if err != nil {
		return nil, err
	}

	if len(targetList) == 0 {
		status.Message = "targetList empty."
		return &GenerateZigZagResult{Status: status, Warn: true}, nil
	}

	status.ZigzagCount = len(targetList)
	slog.Debug("GenerateZigZag targets", "symbol", cmd.Symbol, "previous", len(previousList), "target", len(targetList))

	calc := zigzag.NewZigZagCalculation(previousList[len(previousList)-1])
	calc.ID = 0

	for _, target := range targetList {
		previous := uc.domainService.CalculatePrevious(previousList)
		snapshot := calc.Snapshot()
		calc.ID++

		calc.Calculate(snapshot, target, previous)

		entity := calc.ToEntity()
		entity.Symbol = cmd.Symbol
		entity.Depth = cmd.Depth

		if target.ExistsZigzag {
			rows, err := uc.repo.Update(ctx, cmd.BarType, entity)
			if err != nil {
				return nil, err
			}
			if rows != 1 {
				return nil, apperror.NewUpdateError(fmt.Sprintf("%s:%s", cmd.Symbol, calc.BarDateTime))
			}
		} else {
			rows, err := uc.repo.Insert(ctx, cmd.BarType, entity)
			if err != nil {
				return nil, err
			}
			if rows != 1 {
				return nil, apperror.NewInsertError(fmt.Sprintf("%s:%s", cmd.Symbol, calc.BarDateTime))
			}
		}

		if len(previousList) == cmd.Depth-1 {
			previousList = previousList[1:]
		}
		previousList = append(previousList, target)
	}

	slog.Debug("GenerateZigZag done", "symbol", cmd.Symbol, "count", len(targetList))

	if err := uc.processWaveList(ctx, cmd, calc); err != nil {
		return nil, err
	}

	resultStatus, err := uc.repo.GetStatus(ctx, cmd.BarType, cmd.Symbol, cmd.Depth)
	if err != nil {
		return nil, err
	}
	waveCount := 0
	if calc.WaveList != nil {
		waveCount = len(calc.WaveList)
	}
	resultStatus.Message = fmt.Sprintf("target=%d wave=%d", len(targetList), waveCount)

	return &GenerateZigZagResult{Status: resultStatus, Warn: false}, nil
}

func (uc *GenerateZigZagUseCase) processWaveList(ctx context.Context, cmd fxcommand.ZigZagGenerateCommand, calc *zigzag.ZigZagCalculation) error {
	if len(calc.WaveList) == 0 {
		return nil
	}

	if err := uc.repo.DeleteWave(ctx, cmd.BarType, cmd.Symbol, cmd.Depth, cmd.BarDateTime); err != nil {
		return err
	}
	slog.Debug("deleteWave", "from", cmd.BarDateTime)

	lastWave, err := uc.repo.GetLastWave(ctx, cmd.BarType, cmd.Symbol, cmd.Depth)
	if err != nil {
		return err
	}
	slog.Debug("lastWave", "wave", lastWave)

	if lastWave != nil {
		calc.WaveList[0].PreviousWaveStart = lastWave.WaveStart
		calc.WaveList[0].PreviousWave = lastWave.Wave
	}

	// 前後のwaveが連続していないものを除外する
	enableList := make([]*zigzag.ZigZagWave, 0)
	var prev *zigzag.ZigZagWave
	for i := range calc.WaveList {
		w := &calc.WaveList[i]
		if prev != nil && w.WaveStart.Equal(prev.WaveEnd) {
			enableList = append(enableList, w)
		}
		prev = w
	}

	if len(enableList) > 0 {
		if err := uc.repo.InsertWaveBulk(ctx, cmd.BarType, cmd.Symbol, cmd.Depth, enableList); err != nil {
			slog.Error("insertWaveBulk error", "symbol", cmd.Symbol, "barDateTime", calc.BarDateTime, "error", err)
			return apperror.NewInsertError(fmt.Sprintf("%s:%s", cmd.Symbol, calc.BarDateTime))
		}
	}

	return nil
}
