package fx

import "time"

type PriceInfo struct {
	Symbol      string
	BarDateTime time.Time
	Price       float64
	PriceUsdJpy float64
}
