package zigzagusecase

import (
	"context"
	"math"

	fxcommand "sandbox-api-gin/internal/application/command/fx"
	fxdto "sandbox-api-gin/internal/application/dto/fx"
	"sandbox-api-gin/internal/domain/model/fx/zigzag"
	fxrepository "sandbox-api-gin/internal/domain/repository/fx"
)

const (
	fibF7 = 0.786
	fibF6 = 0.618
	fibF5 = 0.5
	fibF3 = 0.382
	fibF2 = 0.236
)

type SearchZigZagUseCase struct {
	repo fxrepository.ZigZagRepository
}

func NewSearchZigZagUseCase(repo fxrepository.ZigZagRepository) *SearchZigZagUseCase {
	return &SearchZigZagUseCase{repo: repo}
}

type SearchZigZagResult struct {
	TotalCount int
	List       []*fxdto.ZigZagSearchItem
	Page       int
	Size       int
}

func (r *SearchZigZagResult) TotalPage() int {
	if r.TotalCount == 0 {
		return 0
	}
	return (r.TotalCount + r.Size - 1) / r.Size
}

func (uc *SearchZigZagUseCase) Execute(ctx context.Context, cmd fxcommand.ZigZagSearchCommand) (*SearchZigZagResult, error) {
	count, err := uc.repo.SearchCount(ctx, cmd.BarType, cmd.Symbol, cmd.Depth,
		cmd.BarDateTimeMin, cmd.BarDateTimeMax,
		cmd.Wave, cmd.PreviousWave, cmd.NextWave, cmd.Next2Wave, cmd.Wave4h)
	if err != nil {
		return nil, err
	}

	var rows []*zigzag.ZigZagSearchRow
	if count > 0 {
		rows, err = uc.repo.Search(ctx, cmd.BarType, cmd.Symbol, cmd.Depth,
			cmd.BarDateTimeMin, cmd.BarDateTimeMax,
			cmd.Wave, cmd.PreviousWave, cmd.NextWave, cmd.Next2Wave, cmd.Wave4h,
			cmd.Page, cmd.Size)
		if err != nil {
			return nil, err
		}
	}

	list := uc.toItemList(rows, cmd)

	if len(rows) != len(list) {
		count -= (len(rows) - len(list))
	}

	return &SearchZigZagResult{
		TotalCount: count,
		List:       list,
		Page:       cmd.Page,
		Size:       cmd.Size,
	}, nil
}

func (uc *SearchZigZagUseCase) toItemList(rows []*zigzag.ZigZagSearchRow, cmd fxcommand.ZigZagSearchCommand) []*fxdto.ZigZagSearchItem {
	list := make([]*fxdto.ZigZagSearchItem, 0, len(rows))

	for _, row := range rows {
		item := toItem(row)
		scale := 3
		if len(row.Symbol) < 3 || row.Symbol[len(row.Symbol)-3:] != "JPY" {
			scale = 5
		}

		cur := &item.Current
		if cur.Resistance > 0 && cur.Support > 0 {
			fib := calcFibonacci(cur, scale)
			cur.Fibonacci = fib
			calcSma(&cur.Sma4h200s, fib, cur.Wave)
			calcSma(&cur.Sma4h75s, fib, cur.Wave)
			calcSma(&cur.Sma4h20s, fib, cur.Wave)
			calcSma(&cur.Sma1h200s, fib, cur.Wave)
			calcSma(&cur.Sma15m200s, fib, cur.Wave)

			item.NextRsRate = 0
			if (cur.Wave > 0 && item.Next2.Wave > cur.Wave) || (cur.Wave < 0 && item.Next2.Wave < cur.Wave) {
				var targetPrice float64
				if cur.Wave > 0 {
					targetPrice = item.Next2.Support
				} else {
					targetPrice = item.Next2.Resistance
				}
				item.NextRsRate = fibonacciValue(fib, targetPrice)
			}

			item.Next2MaxRate = 0
			if cur.Wave > 0 && item.Next2.Wave > cur.Wave {
				if fib.PriceRange != 0 {
					item.Next2MaxRate = roundScale((item.Next2.Resistance-cur.Support)/fib.PriceRange, 3)
				}
			}
			if cur.Wave < 0 && item.Next2.Wave < cur.Wave {
				if fib.PriceRange != 0 {
					item.Next2MaxRate = roundScale((cur.Resistance-item.Next2.Support)/fib.PriceRange, 3)
				}
			}
		}

		if isExclude(cmd.Direction4h200, cur.Sma4h200s.Direction) {
			continue
		}
		if isExclude(cmd.Direction4h75, cur.Sma4h75s.Direction) {
			continue
		}
		if isExclude(cmd.Direction4h20, cur.Sma4h20s.Direction) {
			continue
		}
		if isExclude(cmd.Direction1h200, cur.Sma1h200s.Direction) {
			continue
		}
		if isExclude(cmd.Direction15m200, cur.Sma15m200s.Direction) {
			continue
		}

		t4h := &item.Target4h
		if t4h.Resistance > 0 && t4h.Support > 0 {
			fibTarget := calcFibonacci(t4h, scale)
			t4h.Fibonacci = fibTarget
			calcSma(&t4h.Sma4h200s, fibTarget, t4h.Wave)
			calcSma(&t4h.Sma4h75s, fibTarget, t4h.Wave)
			calcSma(&t4h.Sma4h20s, fibTarget, t4h.Wave)
			calcSma(&t4h.Sma1h200s, fibTarget, t4h.Wave)
			calcSma(&t4h.Sma15m200s, fibTarget, t4h.Wave)

			if isExclude(cmd.DirectionTarget4h200, t4h.Sma4h200s.Direction) {
				continue
			}
		}

		list = append(list, item)
	}

	return list
}

func toItem(row *zigzag.ZigZagSearchRow) *fxdto.ZigZagSearchItem {
	item := &fxdto.ZigZagSearchItem{
		Symbol:    row.Symbol,
		Depth:     row.Depth,
		Target4h:  toInfoSmaFibonacci(row.Target4h),
		Current:   toInfoSmaFibonacci(row.Current),
		Previous:  toInfo(row.Previous),
		Next:      toInfo(row.Next),
		Next2:     toInfo(row.Next2),
		WaveDxy4h: row.WaveDxy4h,
		WaveDxy1h: row.WaveDxy1h,
	}
	if row.FractalWaveList != nil {
		item.FractalWaveList = make([]fxdto.ZigZagFractalWave, len(row.FractalWaveList))
		for i, fw := range row.FractalWaveList {
			item.FractalWaveList[i] = fxdto.ZigZagFractalWave{WaveStart: fw.WaveStart, Wave: fw.Wave}
		}
	}
	return item
}

func toInfoSmaFibonacci(src *zigzag.WaveWithSma) fxdto.ZigZagInfoSmaFibonacci {
	if src == nil {
		return fxdto.ZigZagInfoSmaFibonacci{
			Sma4h200s:  emptySma(),
			Sma4h75s:   emptySma(),
			Sma4h20s:   emptySma(),
			Sma1h200s:  emptySma(),
			Sma15m200s: emptySma(),
		}
	}
	return fxdto.ZigZagInfoSmaFibonacci{
		WaveStart:  src.WaveStart,
		WaveEnd:    src.WaveEnd,
		Wave:       src.Wave,
		Resistance: src.Resistance,
		Support:    src.Support,
		Sma4h200s:  toSma(src.Sma4h200s),
		Sma4h75s:   toSma(src.Sma4h75s),
		Sma4h20s:   toSma(src.Sma4h20s),
		Sma1h200s:  toSma(src.Sma1h200s),
		Sma15m200s: toSma(src.Sma15m200s),
	}
}

func toInfo(src *zigzag.WaveInfo) fxdto.ZigZagInfo {
	if src == nil {
		return fxdto.ZigZagInfo{}
	}
	return fxdto.ZigZagInfo{
		WaveStart:  src.WaveStart,
		WaveEnd:    src.WaveEnd,
		Wave:       src.Wave,
		Resistance: src.Resistance,
		Support:    src.Support,
	}
}

func toSma(src zigzag.SmaPrice) fxdto.ZigZagSma {
	return fxdto.ZigZagSma{
		PriceS: src.PriceS,
		PriceE: src.PriceE,
	}
}

func emptySma() fxdto.ZigZagSma {
	return fxdto.ZigZagSma{}
}

func calcFibonacci(info *fxdto.ZigZagInfoSmaFibonacci, scale int) *fxdto.ZigZagFibonacciItem {
	uptrend := (info.Wave > 0) == (info.Wave%2 != 0)
	var f1, f0 float64
	if uptrend {
		f1 = info.Resistance
		f0 = info.Support
	} else {
		f1 = info.Support
		f0 = info.Resistance
	}
	r := info.Resistance - info.Support

	return &fxdto.ZigZagFibonacciItem{
		F1:         roundScale(f1, scale),
		F7:         fibPrice(f1, r, fibF7, uptrend, scale),
		F6:         fibPrice(f1, r, fibF6, uptrend, scale),
		F5:         fibPrice(f1, r, fibF5, uptrend, scale),
		F3:         fibPrice(f1, r, fibF3, uptrend, scale),
		F2:         fibPrice(f1, r, fibF2, uptrend, scale),
		F0:         roundScale(f0, scale),
		PriceRange: r,
	}
}

func fibPrice(f1, rangeVal, ratio float64, uptrend bool, scale int) float64 {
	offset := (1 - ratio) * rangeVal
	if uptrend {
		return roundScale(f1-offset, scale)
	}
	return roundScale(f1+offset, scale)
}

func fibonacciValue(fib *fxdto.ZigZagFibonacciItem, targetPrice float64) float64 {
	var rate float64
	if fib.F1 > fib.F0 {
		rate = targetPrice - fib.F0
	} else {
		rate = fib.F0 - targetPrice
	}
	if fib.PriceRange == 0 {
		return 0
	}
	return roundScale(rate/fib.PriceRange*100, 3)
}

func calcSma(sma *fxdto.ZigZagSma, fib *fxdto.ZigZagFibonacciItem, wave int) {
	if sma.PriceS <= 0 || sma.PriceE <= 0 {
		return
	}
	if fib.PriceRange == 0 {
		return
	}

	sma.Deviation = roundScale((sma.PriceE-sma.PriceS)/fib.PriceRange*100, 3)

	var fibRange float64
	if wave > 0 {
		fibRange = fib.F1 - sma.PriceE
	} else {
		fibRange = -(fib.F1 - sma.PriceE)
	}
	sma.Fibonacci = roundScale((1-fibRange/fib.PriceRange)*100, 2)

	dev := int(sma.Deviation)
	direction := 0
	switch {
	case dev > 15:
		direction = 2
	case dev > 5:
		direction = 1
	case dev < -15:
		direction = -2
	case dev < -5:
		direction = -1
	}
	sma.Direction = direction

	fibInt := int(sma.Fibonacci)
	switch {
	case fibInt > 100:
		sma.Position = 1
	case fibInt < 0:
		sma.Position = -1
	default:
		sma.Position = 0
	}
}

func isExclude(searchDir, targetDir int) bool {
	if searchDir == 999 {
		return false
	}
	if searchDir == 0 {
		return targetDir != 0
	}
	if searchDir > 0 {
		return targetDir < searchDir
	}
	return targetDir > searchDir
}

func roundScale(v float64, scale int) float64 {
	pow := math.Pow(10, float64(scale))
	return math.Round(v*pow) / pow
}
