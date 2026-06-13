package fxcommand

import fxmodel "sandbox-api-gin/internal/domain/model/fx"

type ZigZagStatusCommand struct {
	SymbolType string
	BarType    fxmodel.BarType
	Depth      int
}
