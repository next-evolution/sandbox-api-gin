package fxservice

import (
	"math"

	"sandbox-api-gin/internal/domain/model/fx/zigzag"
)

type ZigZagDomainService struct{}

func NewZigZagDomainService() *ZigZagDomainService {
	return &ZigZagDomainService{}
}

// CalculateFibonacci はwave・resistance・supportからFibonacci水準を計算する。
// uptrend = (wave > 0) == (wave % 2 != 0)
func (s *ZigZagDomainService) CalculateFibonacci(wave int, resistance, support float64, scale int) zigzag.ZigZagFibonacci {
	uptrend := (wave > 0) == (wave%2 != 0)
	var f1, f0 float64
	if uptrend {
		f1 = resistance
		f0 = support
	} else {
		f1 = support
		f0 = resistance
	}
	r := resistance - support
	pow := math.Pow(10, float64(scale))

	return zigzag.ZigZagFibonacci{
		F1:         math.Round(f1*pow) / pow,
		F7:         fibPrice(f1, r, zigzag.FibF7, uptrend, scale),
		F6:         fibPrice(f1, r, zigzag.FibF6, uptrend, scale),
		F5:         fibPrice(f1, r, zigzag.FibF5, uptrend, scale),
		F3:         fibPrice(f1, r, zigzag.FibF3, uptrend, scale),
		F2:         fibPrice(f1, r, zigzag.FibF2, uptrend, scale),
		F0:         math.Round(f0*pow) / pow,
		PriceRange: r,
	}
}

// GetFibonacciRate はFibonacci水準に対するtargetPriceの戻り率（%）を返す。
func (s *ZigZagDomainService) GetFibonacciRate(fib zigzag.ZigZagFibonacci, targetPrice float64) float64 {
	var numerator float64
	if fib.F1 > fib.F0 {
		numerator = targetPrice - fib.F0
	} else {
		numerator = fib.F0 - targetPrice
	}
	if fib.PriceRange == 0 {
		return 0
	}
	return math.Round(numerator/fib.PriceRange*100*1000) / 1000
}

// CalculatePrevious はpreviousListの高値・安値・日時を集約して直前ZigZagを返す。
func (s *ZigZagDomainService) CalculatePrevious(previousList []*zigzag.ZigZag) *zigzag.ZigZag {
	result := *previousList[0]
	for _, target := range previousList {
		result.BarDateTime = target.BarDateTime
		if target.BarHighPrice > result.BarHighPrice {
			result.BarHighPrice = target.BarHighPrice
		}
		if target.BarLowPrice < result.BarLowPrice {
			result.BarLowPrice = target.BarLowPrice
		}
	}
	return &result
}

// PreviousNotExists はpreviousListにExistsZigzag=trueのレコードが1件も存在しない場合にtrueを返す。
func (s *ZigZagDomainService) PreviousNotExists(previousList []*zigzag.ZigZag) bool {
	for _, z := range previousList {
		if z.ExistsZigzag {
			return false
		}
	}
	return true
}

func fibPrice(f1, rangeVal, ratio float64, uptrend bool, scale int) float64 {
	offset := (1 - ratio) * rangeVal
	pow := math.Pow(10, float64(scale))
	if uptrend {
		return math.Round((f1-offset)*pow) / pow
	}
	return math.Round((f1+offset)*pow) / pow
}
