package fxcommand

import (
	"io"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type ImportCsvBarDataCommand struct {
	Symbol           string
	BarType          fxmodel.BarType
	SkipLatest       bool
	FileReader       io.Reader
	OriginalFileName string
	FileSize         int64
	UserSub          string
}
