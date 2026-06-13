package zigzag

import "time"

type ZigZag struct {
	Symbol string
	Depth  int
	BarDateTime time.Time

	Resistance         float64
	ResistanceFractal  float64
	Support            float64
	SupportFractal     float64
	PriceHigh          float64
	PriceLow           float64
	BackStepHigh       float64
	BackStepLow        float64
	FractalHigh        float64
	FractalLow         float64

	ResistanceBarDateTime         time.Time
	ResistanceFractalBarDateTime  time.Time
	SupportBarDateTime            time.Time
	SupportFractalBarDateTime     time.Time
	PriceHighBarDateTime          time.Time
	PriceLowBarDateTime           time.Time
	BackStepHighBarDateTime       time.Time
	BackStepLowBarDateTime        time.Time

	Wave             int
	BreakResistance  bool
	BreakSupport     bool
	WaveFractal      int
	BreakResistanceFractal bool
	BreakSupportFractal    bool
	UpTrend          bool
	BackStepUp       int
	BackStepDown     int

	// DB非保存 — JOINで取得するバーデータ
	BarHighPrice  float64
	BarLowPrice   float64
	BarClosePrice float64
	ExistsZigzag  bool
}
