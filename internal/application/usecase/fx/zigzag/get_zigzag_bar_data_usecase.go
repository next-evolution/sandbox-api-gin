package zigzagusecase

import (
	"context"

	fxcommand "sandbox-api-gin/internal/application/command/fx"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	"sandbox-api-gin/internal/domain/model/fx/zigzag"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

type GetZigZagBarDataUseCase struct {
	repo fxrepository.ZigZagRepository
}

func NewGetZigZagBarDataUseCase(repo fxrepository.ZigZagRepository) *GetZigZagBarDataUseCase {
	return &GetZigZagBarDataUseCase{repo: repo}
}

type BarDataResult struct {
	BarType fxmodel.BarType
	Symbol  string
	Depth   int16
	Wave    int
	List    []*zigzag.ZigZagBarDataRow
}

func (uc *GetZigZagBarDataUseCase) Execute(ctx context.Context, cmd fxcommand.ZigZagBarDataCommand) (*BarDataResult, error) {
	list, err := uc.repo.GetBarDataList(ctx, cmd.BarType, cmd.Symbol, int(cmd.Depth), cmd.WaveStart)
	if err != nil {
		return nil, err
	}
	return &BarDataResult{
		BarType: cmd.BarType,
		Symbol:  cmd.Symbol,
		Depth:   cmd.Depth,
		Wave:    cmd.Wave,
		List:    list,
	}, nil
}
