package fxcommand

import fxmodel "sandbox-api-gin/internal/domain/model/fx"

type StatusBarDataCommand struct {
	SymbolType string
	BarType    fxmodel.BarType
}
