package fx

type BarData struct {
	Symbol      string        `json:"symbol"`
	BarDateTime LocalDateTime `json:"barDateTime"`
	OpenPrice   float64       `json:"openPrice"`
	HighPrice   float64       `json:"highPrice"`
	LowPrice    float64       `json:"lowPrice"`
	ClosePrice  float64       `json:"closePrice"`
	Volume      int           `json:"volume"`
	HighProfit  float64       `json:"highProfit"`
	LowProfit   float64       `json:"lowProfit"`
	CloseProfit float64       `json:"closeProfit"`
	RangeProfit float64       `json:"rangeProfit"`
	RsiValue    float64       `json:"rsiValue"`
	RsiMa       float64       `json:"rsiMa"`
}

type BarDataStatus struct {
	Symbol             string
	BarDateTimeMinS    *string
	BarDateTimeMaxS    *string
	Count              int
}
