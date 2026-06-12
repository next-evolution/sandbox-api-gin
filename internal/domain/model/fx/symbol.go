package fx

import "time"

type Symbol struct {
	Symbol           string
	SymbolType       string
	Name             string
	ValidScale       int16
	TargetVolatility float64
	SortOrder        int
	Deleted          bool
	CreatedAt        time.Time
	CreatedBy        string
	UpdatedAt        time.Time
	UpdatedBy        string
}
