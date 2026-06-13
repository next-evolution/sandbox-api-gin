package fxrequest

import (
	"time"

	"sandbox-api-gin/internal/api/dto/request"
)

type ZigZagSearchRequest struct {
	request.ApiSearchRequest
	BarType              string    `json:"barType" binding:"required"`
	Symbol               string    `json:"symbol" binding:"required"`
	Depth                int16     `json:"depth" binding:"required,min=1"`
	BarDateTimeMin       time.Time `json:"barDateTimeMin" binding:"required"`
	BarDateTimeMax       time.Time `json:"barDateTimeMax" binding:"required"`
	Wave                 int       `json:"wave"`
	PreviousWave         int       `json:"previousWave"`
	NextWave             int       `json:"nextWave"`
	Next2Wave            int       `json:"next2Wave"`
	Direction4h200       int       `json:"direction4h200"`
	Direction4h75        int       `json:"direction4h75"`
	Direction4h20        int       `json:"direction4h20"`
	Direction1h200       int       `json:"direction1h200"`
	Direction15m200      int       `json:"direction15m200"`
	Wave4h               int       `json:"wave4h"`
	DirectionTarget4h200 int       `json:"directionTarget4h200"`
}

type ZigZagStatusRequest struct {
	SymbolType string `json:"symbolType" binding:"required"`
	BarType    string `json:"barType" binding:"required"`
	Depth      int16  `json:"depth" binding:"required,min=1"`
}

type ZigZagGenerateRequest struct {
	Symbol      string    `json:"symbol" binding:"required"`
	BarType     string    `json:"barType" binding:"required"`
	Depth       int16     `json:"depth" binding:"required,min=1"`
	BarDateTime time.Time `json:"barDateTime" binding:"required"`
	LoadSize    int       `json:"loadSize" binding:"required,min=1"`
}

type ZigZagBarDataRequest struct {
	BarType   string    `json:"barType" binding:"required"`
	Symbol    string    `json:"symbol" binding:"required"`
	Depth     int16     `json:"depth" binding:"required,min=1"`
	WaveStart time.Time `json:"waveStart" binding:"required"`
	Wave      int       `json:"wave"`
}
