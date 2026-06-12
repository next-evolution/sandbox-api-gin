package fxcommand

import fxmodel "sandbox-api-gin/internal/domain/model/fx"

type SearchBarDataCommand struct {
	Symbol      string
	BarType     fxmodel.BarType
	BarDateFrom string
	BarDateTo   string
	SortAsc     bool
	Page        int
	Size        int
}
