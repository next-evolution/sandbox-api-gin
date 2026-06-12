package fx

import (
	"math"
	"strings"
)

type TradeEntry struct {
	ID               *int64        `json:"id"`
	TradeVersion     string        `json:"tradeVersion"`
	EntryType        EntryType     `json:"entryType"`
	Symbol           string        `json:"symbol"`
	TradeType        TradeType     `json:"tradeType"`
	ContractAt       LocalDateTime `json:"contractAt"`
	FibonacciType    string        `json:"fibonacciType"`
	FibonacciBar     string        `json:"fibonacciBar"`
	ContractPrice    float64       `json:"contractPrice"`
	LossPrice        float64       `json:"lossPrice"`
	PositionRatio    int           `json:"positionRatio"`
	PriceJpy         float64       `json:"priceJpy"`
	Lot              float64       `json:"lot"`
	SettlementAmount int           `json:"settlementAmount"`
	LossPips         int           `json:"lossPips"`
	SettlementRatio  float64       `json:"settlementRatio"`
	Comment          *string       `json:"comment"`
	ImagePath        *string       `json:"imagePath"`
}

func (e *TradeEntry) IsDollarCurrency() bool {
	return len(e.Symbol) > 0 && !strings.HasSuffix(e.Symbol, "JPY")
}

func (e *TradeEntry) ApplyPrice(priceInfo *PriceInfo) {
	if e.ID == nil || *e.ID == 0 {
		id := int64(-1)
		e.ID = &id
	}
	e.ContractAt = LocalDateTime{priceInfo.BarDateTime}
	if e.ContractPrice == 0 {
		e.ContractPrice = priceInfo.Price
	}
	if e.PriceJpy == 0 {
		e.PriceJpy = priceInfo.PriceUsdJpy
	}
}

// ComputeDefaultSettlementPrice はデフォルト決済価格を計算する。
// JavaのTradeEntry.computeDefaultSettlementPrice()に相当。
func (e *TradeEntry) ComputeDefaultSettlementPrice(plusPips float64) float64 {
	var offset float64
	if e.IsDollarCurrency() {
		// JavaのBigDecimal.divide(100, 5, HALF_UP)に相当
		offset = math.Round(plusPips/100*1e5) / 1e5
	} else {
		offset = plusPips
	}
	if e.TradeType == TradeTypeL {
		return e.ContractPrice + offset
	}
	return e.ContractPrice - offset
}

// ApplyDefaultLossPrice はlossPrice未設定時にデフォルト値を適用する。
func (e *TradeEntry) ApplyDefaultLossPrice() {
	if e.LossPrice != 0 {
		return
	}
	var rangeVal float64
	if e.IsDollarCurrency() {
		rangeVal = 0.003
	} else {
		rangeVal = 0.3
	}
	if e.TradeType == TradeTypeL {
		e.LossPrice = e.ContractPrice - rangeVal
	} else {
		e.LossPrice = e.ContractPrice + rangeVal
	}
}
