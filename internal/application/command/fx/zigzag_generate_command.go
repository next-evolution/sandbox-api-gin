package fxcommand

import (
	"time"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type ZigZagGenerateCommand struct {
	Symbol      string
	BarType     fxmodel.BarType
	Depth       int
	BarDateTime time.Time
	LoadSize    int
}
