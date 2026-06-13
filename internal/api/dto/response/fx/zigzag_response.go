package fxresponse

import (
	"time"

	"sandbox-api-gin/internal/api/dto/response"
)

// ZigZagSearchResponse は POST /v1/fx/zigzag のレスポンス
type ZigZagSearchResponse struct {
	ReturnCode  response.ReturnCode `json:"returnCode"`
	TotalCount  int                 `json:"totalCount"`
	SearchCount int                 `json:"searchCount"`
	TotalPage   int                 `json:"totalPage"`
	List        []ZigZagResult      `json:"list"`
}

type ZigZagResult struct {
	Symbol          string              `json:"symbol"`
	Depth           int                 `json:"depth"`
	Target4h        ZigZagInfoSmaFib    `json:"target4h"`
	Current         ZigZagInfoSmaFib    `json:"current"`
	Previous        ZigZagInfo          `json:"previous"`
	Next            ZigZagInfo          `json:"next"`
	Next2           ZigZagInfo          `json:"next2"`
	NextRsRate      float64             `json:"nextRsRate"`
	Next2MaxRate    float64             `json:"next2MaxRate"`
	WaveDxy4h       float64             `json:"waveDxy4h"`
	WaveDxy1h       float64             `json:"waveDxy1h"`
	FractalWaveList []ZigZagFractalWave `json:"fractalWaveList"`
}

type ZigZagInfo struct {
	WaveStart  *time.Time `json:"waveStart"`
	WaveEnd    *time.Time `json:"waveEnd"`
	Wave       int        `json:"wave"`
	Resistance float64    `json:"resistance"`
	Support    float64    `json:"support"`
}

type ZigZagInfoSmaFib struct {
	WaveStart  *time.Time       `json:"waveStart"`
	WaveEnd    *time.Time       `json:"waveEnd"`
	Wave       int              `json:"wave"`
	Resistance float64          `json:"resistance"`
	Support    float64          `json:"support"`
	Fibonacci  *ZigZagFibonacci `json:"fibonacci"`
	Sma4h200s  ZigZagSma        `json:"sma4h200s"`
	Sma4h75s   ZigZagSma        `json:"sma4h75s"`
	Sma4h20s   ZigZagSma        `json:"sma4h20s"`
	Sma1h200s  ZigZagSma        `json:"sma1h200s"`
	Sma15m200s ZigZagSma        `json:"sma15m200s"`
}

type ZigZagFibonacci struct {
	F1         float64 `json:"f1"`
	F7         float64 `json:"f7"`
	F6         float64 `json:"f6"`
	F5         float64 `json:"f5"`
	F3         float64 `json:"f3"`
	F2         float64 `json:"f2"`
	F0         float64 `json:"f0"`
	PriceRange float64 `json:"priceRange"`
}

type ZigZagSma struct {
	PriceS    float64 `json:"priceS"`
	PriceE    float64 `json:"priceE"`
	Deviation float64 `json:"deviation"`
	Fibonacci float64 `json:"fibonacci"`
	Direction int     `json:"direction"`
	Position  int     `json:"position"`
}

type ZigZagFractalWave struct {
	WaveStart *time.Time `json:"waveStart"`
	Wave      int        `json:"wave"`
}

// ZigZagStatusResponse は POST /v1/fx/zigzag/status のレスポンス
type ZigZagStatusResponse struct {
	ReturnCode  response.ReturnCode `json:"returnCode"`
	TotalCount  int                 `json:"totalCount"`
	SearchCount int                 `json:"searchCount"`
	TotalPage   int                 `json:"totalPage"`
	List        []ZigZagStatusItem  `json:"list"`
}

type ZigZagStatusItem struct {
	SymbolType           string `json:"symbolType"`
	BarType              string `json:"barType"`
	Symbol               string `json:"symbol"`
	Depth                int16  `json:"depth"`
	BarDateTimeMin       string `json:"barDateTimeMin"`
	BarDateTimeMax       string `json:"barDateTimeMax"`
	BarCount             int    `json:"barCount"`
	BarDateTimeMinZigZag string `json:"barDateTimeMinZigZag"`
	BarDateTimeMaxZigZag string `json:"barDateTimeMaxZigZag"`
	ZigzagCount          int    `json:"zigzagCount"`
	BreakResistanceCount int    `json:"breakResistanceCount"`
	BreakSupportCount    int    `json:"breakSupportCount"`
	Message              string `json:"message,omitempty"`
}

// ZigZagGenerateResponse は POST /v1/fx/zigzag/generate のレスポンス
type ZigZagGenerateResponse struct {
	ReturnCode response.ReturnCode `json:"returnCode"`
	Message    string              `json:"message,omitempty"`
	Status     ZigZagStatusItem    `json:"status"`
}

// ZigZagBarDataResponse は POST /v1/fx/zigzag/bar-data のレスポンス
type ZigZagBarDataResponse struct {
	ReturnCode       response.ReturnCode `json:"returnCode"`
	BarType          string              `json:"barType"`
	Symbol           string              `json:"symbol"`
	Depth            int16               `json:"depth"`
	Wave             int                 `json:"wave"`
	ZigZagBarDataList []ZigZagBarData    `json:"zigZagBarDataList"`
}

type ZigZagBarData struct {
	BarDateTime time.Time `json:"barDateTime"`
	OpenPrice   float64   `json:"openPrice"`
	HighPrice   float64   `json:"highPrice"`
	LowPrice    float64   `json:"lowPrice"`
	ClosePrice  float64   `json:"closePrice"`
	Sma200      float64   `json:"sma200"`
	Sma75       float64   `json:"sma75"`
	Sma20       float64   `json:"sma20"`
}
