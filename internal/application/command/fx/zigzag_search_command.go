package fxcommand

import (
	"time"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type ZigZagSearchCommand struct {
	BarType        fxmodel.BarType
	Symbol         string
	Depth          int
	BarDateTimeMin time.Time
	BarDateTimeMax time.Time
	Wave           int
	PreviousWave   int
	NextWave       int
	Next2Wave      int
	Direction4h200 int
	Direction4h75  int
	Direction4h20  int
	Direction1h200 int
	Direction15m200 int
	Wave4h         int
	DirectionTarget4h200 int
	Page           int
	Size           int
}
