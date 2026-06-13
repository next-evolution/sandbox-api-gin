package fxdto

import "time"

type ZigZagSearchItem struct {
	Symbol          string
	Depth           int
	Target4h        ZigZagInfoSmaFibonacci
	Current         ZigZagInfoSmaFibonacci
	Previous        ZigZagInfo
	Next            ZigZagInfo
	Next2           ZigZagInfo
	NextRsRate      float64
	Next2MaxRate    float64
	WaveDxy4h       float64
	WaveDxy1h       float64
	FractalWaveList []ZigZagFractalWave
}

type ZigZagInfo struct {
	WaveStart  time.Time
	WaveEnd    time.Time
	Wave       int
	Resistance float64
	Support    float64
}

type ZigZagInfoSmaFibonacci struct {
	WaveStart  time.Time
	WaveEnd    time.Time
	Wave       int
	Resistance float64
	Support    float64
	Fibonacci  *ZigZagFibonacciItem
	Sma4h200s  ZigZagSma
	Sma4h75s   ZigZagSma
	Sma4h20s   ZigZagSma
	Sma1h200s  ZigZagSma
	Sma15m200s ZigZagSma
}

type ZigZagFibonacciItem struct {
	F1         float64
	F7         float64
	F6         float64
	F5         float64
	F3         float64
	F2         float64
	F0         float64
	PriceRange float64
}

type ZigZagSma struct {
	PriceS    float64
	PriceE    float64
	Deviation float64
	Fibonacci float64
	Direction int
	Position  int
}

type ZigZagFractalWave struct {
	WaveStart time.Time
	Wave      int
}
