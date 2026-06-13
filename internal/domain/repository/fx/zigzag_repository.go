package fxrepository

import (
	"context"
	"time"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	"sandbox-api-gin/internal/domain/model/fx/zigzag"
)

type ZigZagRepository interface {
	// 検索
	SearchCount(ctx context.Context, barType fxmodel.BarType, symbol string, depth int,
		barDateTimeMin, barDateTimeMax time.Time,
		wave, previousWave, nextWave, next2Wave, wave4h int) (int, error)

	Search(ctx context.Context, barType fxmodel.BarType, symbol string, depth int,
		barDateTimeMin, barDateTimeMax time.Time,
		wave, previousWave, nextWave, next2Wave, wave4h int,
		page, size int) ([]*zigzag.ZigZagSearchRow, error)

	// ステータス
	GetStatusList(ctx context.Context, symbolType string, barType fxmodel.BarType, depth int) ([]*zigzag.ZigZagStatus, error)
	GetStatus(ctx context.Context, barType fxmodel.BarType, symbol string, depth int) (*zigzag.ZigZagStatus, error)

	// バーデータ
	GetBarDataList(ctx context.Context, barType fxmodel.BarType, symbol string, depth int, waveStart time.Time) ([]*zigzag.ZigZagBarDataRow, error)

	// 生成
	TargetBarCount(ctx context.Context, barType fxmodel.BarType, symbol string, barDateTime time.Time) (int, error)
	PreviousList(ctx context.Context, barType fxmodel.BarType, symbol string, depth int, barDateTime time.Time, limit int) ([]*zigzag.ZigZag, error)
	TargetList(ctx context.Context, barType fxmodel.BarType, symbol string, depth int, barDateTime time.Time, limit int) ([]*zigzag.ZigZag, error)
	Insert(ctx context.Context, barType fxmodel.BarType, z *zigzag.ZigZag) (int64, error)
	Update(ctx context.Context, barType fxmodel.BarType, z *zigzag.ZigZag) (int64, error)

	// Wave
	DeleteWave(ctx context.Context, barType fxmodel.BarType, symbol string, depth int, barDateTime time.Time) error
	GetLastWave(ctx context.Context, barType fxmodel.BarType, symbol string, depth int) (*zigzag.ZigZagWave, error)
	InsertWaveBulk(ctx context.Context, barType fxmodel.BarType, symbol string, depth int, waveList []*zigzag.ZigZagWave) error
}
