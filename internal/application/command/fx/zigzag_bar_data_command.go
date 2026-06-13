package fxcommand

import (
	"time"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type ZigZagBarDataCommand struct {
	BarType   fxmodel.BarType
	Symbol    string
	Depth     int16
	WaveStart time.Time
	Wave      int
}
