package fxservice

import (
	"math"

	"sandbox-api-gin/internal/domain/model/fx"
)

const (
	unitValue        = 100000.0
	hundred          = 100.0
	thousand         = 1000.0
	unitValueHundred = unitValue * hundred

	defaultRiskAmount    = 10000.0
	defaultFirstLotRatio = 0.3
)

type FxTradeCalculator struct{}

func NewFxTradeCalculator() *FxTradeCalculator {
	return &FxTradeCalculator{}
}

func (c *FxTradeCalculator) Calculate(
	entry *fx.TradeEntry,
	positionList []*fx.TradePosition,
	riskAmount float64,
	firstLotRatio float64,
) {
	effectiveRisk := defaultRiskAmount
	if riskAmount != 0 {
		effectiveRisk = riskAmount
	}

	// JavaのfirstLotRatio.divide(HUNDRED, 2, HALF_UP)に相当
	effectiveRatio := defaultFirstLotRatio
	if firstLotRatio != 0 {
		effectiveRatio = roundScale(firstLotRatio/hundred, 2)
	}

	c.calculateLot(entry, positionList, effectiveRatio, effectiveRisk)

	if entry.Lot > 0 {
		c.calculateAmount(entry, positionList)
		c.calculateEntrySettlementRatio(entry)
	}
}

func (c *FxTradeCalculator) PriceRange(base, value float64) float64 {
	if base > value {
		return base - value
	}
	return value - base
}

func (c *FxTradeCalculator) Pips(base, value float64, isDollarCurrency bool) int {
	r := c.PriceRange(base, value)
	if isDollarCurrency {
		return int(r * thousand * hundred)
	}
	return int(r * thousand)
}

func (c *FxTradeCalculator) ProfitAmount(entry *fx.TradeEntry, position *fx.TradePosition) int {
	var isProfit bool
	if entry.TradeType == fx.TradeTypeL {
		isProfit = position.SettlementPrice > entry.ContractPrice
	} else {
		isProfit = position.SettlementPrice < entry.ContractPrice
	}

	pip := c.PriceRange(entry.ContractPrice, position.SettlementPrice)
	var amountJpy float64
	if entry.IsDollarCurrency() {
		amountJpy = pip * position.Lot * unitValueHundred
	} else {
		amountJpy = pip * position.Lot * unitValue
	}

	var amount int
	if entry.IsDollarCurrency() {
		// JavaのpriceJpy.divide(HUNDRED, 2, HALF_UP)に相当
		amount = int(amountJpy * roundScale(entry.PriceJpy/hundred, 2))
	} else {
		amount = int(amountJpy)
	}

	if isProfit {
		return amount
	}
	return amount * -1
}

func (c *FxTradeCalculator) LossAmount(entry *fx.TradeEntry, lot float64) int {
	pip := c.PriceRange(entry.ContractPrice, entry.LossPrice)
	var amountJpy float64
	if entry.IsDollarCurrency() {
		amountJpy = pip * lot * unitValueHundred
	} else {
		amountJpy = pip * lot * unitValue
	}

	if entry.IsDollarCurrency() {
		return int(amountJpy * roundScale(entry.PriceJpy/hundred, 2))
	}
	return int(amountJpy)
}

func (c *FxTradeCalculator) TotalLot(entry *fx.TradeEntry, riskAmount float64) float64 {
	lossValue := c.PriceRange(entry.ContractPrice, entry.LossPrice)
	if entry.IsDollarCurrency() {
		// scale=3, 最後はscale-1=2
		lotTemp := roundScale(riskAmount/lossValue, 3)
		lotTemp = roundScale(lotTemp/unitValue, 3)
		return roundScale(lotTemp/entry.PriceJpy, 2)
	}
	// scale=2
	lotTemp := roundScale(riskAmount/lossValue, 2)
	return roundScale(lotTemp/unitValue, 2)
}

func (c *FxTradeCalculator) calculateLot(
	entry *fx.TradeEntry,
	positionList []*fx.TradePosition,
	firstLotRatio float64,
	riskAmount float64,
) {
	entry.Lot = c.TotalLot(entry, riskAmount)

	switch len(positionList) {
	case 1:
		positionList[0].Lot = entry.Lot
	case 2:
		// JavaのMathContext(2, HALF_UP): 有効数字2桁に丸める
		positionList[0].Lot = roundToSigFigs(entry.Lot*firstLotRatio, 2)
		positionList[1].Lot = entry.Lot - positionList[0].Lot
	case 3:
		positionList[0].Lot = roundToSigFigs(entry.Lot*firstLotRatio, 2)
		// JavaのBigDecimal.divide(TWO, 2, HALF_UP)に相当
		positionList[1].Lot = roundScale((entry.Lot-positionList[0].Lot)/2, 2)
		positionList[2].Lot = entry.Lot - positionList[0].Lot - positionList[1].Lot
	}
}

func (c *FxTradeCalculator) calculateAmount(entry *fx.TradeEntry, positionList []*fx.TradePosition) {
	profitAmountTotal := 0

	for _, position := range positionList {
		position.ProfitAmount = c.ProfitAmount(entry, position)
		position.LossAmount = c.LossAmount(entry, position.Lot)
		profitAmountTotal += position.ProfitAmount

		position.SettlementRatio = roundScale(
			float64(position.ProfitAmount)/float64(position.LossAmount), 2)

		position.SettlementPips = c.Pips(
			position.SettlementPrice, entry.ContractPrice, entry.IsDollarCurrency())
	}

	entry.SettlementAmount = profitAmountTotal
	entry.LossPips = c.Pips(entry.ContractPrice, entry.LossPrice, entry.IsDollarCurrency())
}

func (c *FxTradeCalculator) calculateEntrySettlementRatio(entry *fx.TradeEntry) {
	entry.SettlementRatio = roundScale(
		float64(entry.SettlementAmount)/float64(c.LossAmount(entry, entry.Lot)), 2)
}

// roundScale は小数点以下scale桁でHALF_UP丸めを行う。
// JavaのBigDecimal.divide(x, scale, HALF_UP)やsetScale(scale, HALF_UP)に相当。
func roundScale(v float64, scale int) float64 {
	factor := math.Pow(10, float64(scale))
	return math.Round(v*factor) / factor
}

// roundToSigFigs は有効数字n桁でHALF_UP丸めを行う。
// JavaのMathContext(n, HALF_UP)を使ったBigDecimal.multiply()に相当。
func roundToSigFigs(v float64, n int) float64 {
	if v == 0 {
		return 0
	}
	d := math.Pow(10, float64(n)-1-math.Floor(math.Log10(math.Abs(v))))
	return math.Round(v*d) / d
}
