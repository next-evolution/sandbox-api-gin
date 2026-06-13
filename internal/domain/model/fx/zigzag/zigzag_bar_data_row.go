package zigzag

import "time"

type ZigZagBarDataRow struct {
	BarDateTime time.Time
	OpenPrice   float64
	HighPrice   float64
	LowPrice    float64
	ClosePrice  float64
	Sma200      float64
	Sma75       float64
	Sma20       float64
}
