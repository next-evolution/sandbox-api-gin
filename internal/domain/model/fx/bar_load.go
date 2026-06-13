package fx

import "time"

type BarLoadData struct {
	Symbol      string
	BarDateTime time.Time
	OpenPrice   float64
	HighPrice   float64
	LowPrice    float64
	ClosePrice  float64
	Volume      int
}

type BarLoadSma struct {
	Symbol      string
	BarDateTime time.Time
	SmaRange    int
	SmaPrice    *float64
	SmaCross    bool
}

func NewBarLoadSma(symbol string, barDateTime time.Time, smaRange int, smaPrice *float64, highPrice, lowPrice float64) BarLoadSma {
	smaCross := smaPrice != nil && highPrice >= *smaPrice && lowPrice <= *smaPrice
	return BarLoadSma{
		Symbol:      symbol,
		BarDateTime: barDateTime,
		SmaRange:    smaRange,
		SmaPrice:    smaPrice,
		SmaCross:    smaCross,
	}
}

type BarLoadRsi struct {
	Symbol      string
	BarDateTime time.Time
	RsiRange    int
	RsiValue    *float64
	RsiMa       *float64
}

type BarCsvImportCheck struct {
	Symbol      string
	ExistsCount int
	DiffCount   int
}
