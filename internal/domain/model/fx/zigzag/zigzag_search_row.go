package zigzag

import "time"

type ZigZagSearchRow struct {
	Symbol string
	Depth  int

	Target4h *WaveWithSma
	Current  *WaveWithSma
	Previous *WaveInfo
	Next     *WaveInfo
	Next2    *WaveInfo

	WaveDxy4h float64
	WaveDxy1h float64

	FractalWaveList []FractalWaveInfo
}

type WaveInfo struct {
	WaveStart  time.Time
	WaveEnd    time.Time
	Wave       int
	Resistance float64
	Support    float64
}

type WaveWithSma struct {
	WaveStart  time.Time
	WaveEnd    time.Time
	Wave       int
	Resistance float64
	Support    float64
	Sma4h200s  SmaPrice
	Sma4h75s   SmaPrice
	Sma4h20s   SmaPrice
	Sma1h200s  SmaPrice
	Sma15m200s SmaPrice
}

type SmaPrice struct {
	PriceS float64
	PriceE float64
}

type FractalWaveInfo struct {
	WaveStart time.Time
	Wave      int
}
