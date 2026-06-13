package zigzag

import "time"

var minDateTime = time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)

const backStep3 = 3

type price struct {
	BarDateTime time.Time
	Price       float64
}

func newPrice(p float64, dt time.Time) price {
	return price{Price: p, BarDateTime: dt}
}

func copyPrice(p price) price {
	return price{Price: p.Price, BarDateTime: p.BarDateTime}
}

// ZigZagCalculation はJavaのZigZagCalculationに相当する。
// JavaのBigDecimalはfloat64でポーティングしている。
type ZigZagCalculation struct {
	Symbol   string
	Depth    int
	BarDateTime time.Time

	Resistance         price
	ResistanceFractal  price
	Support            price
	SupportFractal     price
	PriceHigh          price
	PriceLow           price
	BackStepHigh       price
	BackStepLow        price

	FractalHigh float64
	FractalLow  float64

	Wave             int
	WaveStart        bool
	BreakResistance  bool
	BreakSupport     bool
	WaveFractal      int
	BreakResistanceFractal bool
	BreakSupportFractal    bool
	UpTrend          bool
	BackStepUp       int
	BackStepDown     int

	Message       string
	ID            int
	WaveList      []ZigZagWave
	WaveFractalList []ZigZagWave
}

func NewZigZagCalculation(z *ZigZag) *ZigZagCalculation {
	return &ZigZagCalculation{
		Symbol:      z.Symbol,
		Depth:       z.Depth,
		BarDateTime: z.BarDateTime,

		Resistance:         newPrice(z.Resistance, z.ResistanceBarDateTime),
		ResistanceFractal:  newPrice(z.ResistanceFractal, z.ResistanceFractalBarDateTime),
		Support:            newPrice(z.Support, z.SupportBarDateTime),
		SupportFractal:     newPrice(z.SupportFractal, z.SupportFractalBarDateTime),
		PriceHigh:          newPrice(z.PriceHigh, z.PriceHighBarDateTime),
		PriceLow:           newPrice(z.PriceLow, z.PriceLowBarDateTime),
		BackStepHigh:       newPrice(z.BackStepHigh, z.BackStepHighBarDateTime),
		BackStepLow:        newPrice(z.BackStepLow, z.BackStepLowBarDateTime),
		FractalHigh:        z.FractalHigh,
		FractalLow:         z.FractalLow,

		Wave:             z.Wave,
		BreakResistance:  z.BreakResistance,
		BreakSupport:     z.BreakSupport,
		WaveFractal:      z.Wave,
		BreakResistanceFractal: z.BreakResistance,
		BreakSupportFractal:    z.BreakSupport,
		UpTrend:          z.UpTrend,
		BackStepUp:       z.BackStepUp,
		BackStepDown:     z.BackStepDown,

		WaveList:      []ZigZagWave{},
		WaveFractalList: []ZigZagWave{},
	}
}

func (c *ZigZagCalculation) Snapshot() *ZigZagCalculation {
	return &ZigZagCalculation{
		Symbol:      c.Symbol,
		Depth:       c.Depth,
		BarDateTime: c.BarDateTime,

		Resistance:         copyPrice(c.Resistance),
		ResistanceFractal:  copyPrice(c.ResistanceFractal),
		Support:            copyPrice(c.Support),
		SupportFractal:     copyPrice(c.SupportFractal),
		PriceHigh:          copyPrice(c.PriceHigh),
		PriceLow:           copyPrice(c.PriceLow),
		BackStepHigh:       copyPrice(c.BackStepHigh),
		BackStepLow:        copyPrice(c.BackStepLow),
		FractalHigh:        c.FractalHigh,
		FractalLow:         c.FractalLow,

		Wave:             c.Wave,
		WaveStart:        false,
		BreakResistance:  c.BreakResistance,
		BreakSupport:     c.BreakSupport,
		WaveFractal:      c.WaveFractal,
		BreakResistanceFractal: c.BreakResistanceFractal,
		BreakSupportFractal:    c.BreakSupportFractal,
		UpTrend:          c.UpTrend,
		BackStepUp:       c.BackStepUp,
		BackStepDown:     c.BackStepDown,
	}
}

func (c *ZigZagCalculation) ToEntity() *ZigZag {
	return &ZigZag{
		Symbol:      c.Symbol,
		Depth:       c.Depth,
		BarDateTime: c.BarDateTime,

		Resistance:        c.Resistance.Price,
		ResistanceFractal: c.ResistanceFractal.Price,
		Support:           c.Support.Price,
		SupportFractal:    c.SupportFractal.Price,
		PriceHigh:         c.PriceHigh.Price,
		PriceLow:          c.PriceLow.Price,
		BackStepHigh:      c.BackStepHigh.Price,
		BackStepLow:       c.BackStepLow.Price,
		FractalHigh:       c.FractalHigh,
		FractalLow:        c.FractalLow,

		ResistanceBarDateTime:        c.Resistance.BarDateTime,
		ResistanceFractalBarDateTime: c.ResistanceFractal.BarDateTime,
		SupportBarDateTime:           c.Support.BarDateTime,
		SupportFractalBarDateTime:    c.SupportFractal.BarDateTime,
		PriceHighBarDateTime:         c.PriceHigh.BarDateTime,
		PriceLowBarDateTime:          c.PriceLow.BarDateTime,
		BackStepHighBarDateTime:      c.BackStepHigh.BarDateTime,
		BackStepLowBarDateTime:       c.BackStepLow.BarDateTime,

		Wave:             c.Wave,
		BreakResistance:  c.BreakResistance,
		BreakSupport:     c.BreakSupport,
		WaveFractal:      c.WaveFractal,
		BreakResistanceFractal: c.BreakResistanceFractal,
		BreakSupportFractal:    c.BreakSupportFractal,
		UpTrend:          c.UpTrend,
		BackStepUp:       c.BackStepUp,
		BackStepDown:     c.BackStepDown,
	}
}

func (c *ZigZagCalculation) Calculate(r *ZigZagCalculation, target, previous *ZigZag) {
	c.BarDateTime = target.BarDateTime
	c.Message = ""

	if target.BarLowPrice < r.Support.Price && target.BarHighPrice > r.Resistance.Price {
		if r.UpTrend {
			c.breakResistanceSupport(r, target)
		} else {
			c.breakSupportResistance(r, target)
		}
	} else {
		c.calculateBackStepUp(r, target, previous)
		c.calculateBackStepDown(r, target, previous)
		c.calculateLatest(r, target)
		c.calculateFractal(r, target)
		c.calculateBreakResistance(r, target)
		c.calculateBreakSupport(r, target)
	}
}

func (c *ZigZagCalculation) breakResistanceSupport(r *ZigZagCalculation, target *ZigZag) {
	c.UpTrend = false
	c.BreakResistance = false
	c.BreakSupport = true

	c.Resistance = newPrice(target.BarHighPrice, target.BarDateTime)
	c.ResistanceFractal = newPrice(target.BarHighPrice, target.BarDateTime)
	c.PriceHigh = newPrice(target.BarHighPrice, target.BarDateTime)
	c.Support = copyPrice(r.SupportFractal)

	if r.Wave > 0 {
		c.addWave(r.Resistance.BarDateTime, c.ResistanceFractal.BarDateTime,
			r.Resistance.Price, r.Support.Price, r.Wave, "BRS_RV1")
		c.Wave = -1
	} else {
		if r.Wave%2 == 0 {
			c.addWave(r.Support.BarDateTime, c.ResistanceFractal.BarDateTime,
				r.FractalHigh, r.Support.Price, r.Wave, "BRS_DW2")
		} else {
			c.addWave(r.Resistance.BarDateTime, r.Support.BarDateTime,
				r.Resistance.Price, r.Support.Price, r.Wave, "BRS_DW1")
		}
		c.Wave = c.Wave - 1
	}

	c.Support = newPrice(target.BarLowPrice, target.BarDateTime)
	c.SupportFractal = newPrice(target.BarLowPrice, target.BarDateTime)
	c.PriceLow = newPrice(target.BarLowPrice, target.BarDateTime)
	c.FractalHigh = target.BarLowPrice
}

func (c *ZigZagCalculation) breakSupportResistance(r *ZigZagCalculation, target *ZigZag) {
	c.UpTrend = true
	c.BreakResistance = true
	c.BreakSupport = false

	c.Support = newPrice(target.BarLowPrice, target.BarDateTime)
	c.SupportFractal = newPrice(target.BarLowPrice, target.BarDateTime)
	c.PriceLow = newPrice(target.BarLowPrice, target.BarDateTime)
	c.Resistance = copyPrice(r.ResistanceFractal)

	if r.Wave < 0 {
		c.addWave(r.Resistance.BarDateTime, c.SupportFractal.BarDateTime,
			r.Resistance.Price, c.Support.Price, r.Wave, "BSR_RV1")
		c.Wave = 1
	} else {
		if r.Wave%2 == 0 {
			c.addWave(r.Resistance.BarDateTime, c.SupportFractal.BarDateTime,
				r.Resistance.Price, r.FractalLow, r.Wave, "BSR_UP2")
		} else {
			c.addWave(r.Support.BarDateTime, r.Resistance.BarDateTime,
				r.Resistance.Price, r.Support.Price, r.Wave, "BSR_UP1")
		}
		c.Wave = c.Wave + 1
	}

	c.Resistance = newPrice(target.BarHighPrice, target.BarDateTime)
	c.ResistanceFractal = newPrice(target.BarHighPrice, target.BarDateTime)
	c.PriceHigh = newPrice(target.BarHighPrice, target.BarDateTime)
	c.FractalLow = target.BarHighPrice
}

func (c *ZigZagCalculation) calculateBackStepUp(r *ZigZagCalculation, target, previous *ZigZag) {
	c.updateBackStepHigh(r, target)

	if r.BackStepUp == backStep3 {
		c.commitBackStepUp(r, target.BarDateTime)
	}

	if c.BackStepUp < backStep3 && c.BackStepUp > 0 && r.BackStepUp > 0 {
		c.BackStepUp++
		if target.BarHighPrice > r.BackStepHigh.Price || (r.BreakSupport && c.BackStepUp == 2) {
			if target.BarLowPrice > r.Support.Price {
				c.commitBackStepUp(r, target.BarDateTime)
			}
		}
		if target.BarLowPrice < previous.BarLowPrice {
			c.BackStepUp = 0
		}
	}

	if target.BarHighPrice >= previous.BarHighPrice &&
		r.BackStepUp == 0 &&
		!r.UpTrend &&
		target.BarLowPrice > r.SupportFractal.Price {
		c.BackStepUp = 1
		c.BackStepHigh = newPrice(target.BarHighPrice, target.BarDateTime)
	}
}

func (c *ZigZagCalculation) calculateBackStepDown(r *ZigZagCalculation, target, previous *ZigZag) {
	c.updateBackStepLow(r, target)

	if r.BackStepDown == backStep3 {
		c.commitBackStepDown(r, target.BarDateTime)
	}

	if c.BackStepDown < backStep3 && c.BackStepDown > 0 && r.BackStepDown > 0 {
		c.BackStepDown++
		if target.BarLowPrice < r.BackStepLow.Price || (r.BreakResistance && c.BackStepDown == 2) {
			if target.BarHighPrice < r.Resistance.Price {
				c.commitBackStepDown(r, target.BarDateTime)
			}
		}
		if target.BarHighPrice > previous.BarHighPrice {
			c.BackStepDown = 0
		}
	}

	if target.BarLowPrice <= previous.BarLowPrice &&
		r.BackStepDown == 0 &&
		r.UpTrend &&
		target.BarHighPrice < r.Resistance.Price {
		c.BackStepDown = 1
		c.BackStepLow = newPrice(target.BarLowPrice, target.BarDateTime)
	}
}

func (c *ZigZagCalculation) calculateLatest(r *ZigZagCalculation, target *ZigZag) {
	if target.BarHighPrice > r.PriceHigh.Price {
		c.PriceHigh = newPrice(target.BarHighPrice, target.BarDateTime)
	}
	if target.BarLowPrice < r.PriceLow.Price {
		c.PriceLow = newPrice(target.BarLowPrice, target.BarDateTime)
	}
	if target.BarHighPrice > r.FractalHigh {
		c.FractalHigh = target.BarHighPrice
	}
	if target.BarLowPrice < r.FractalLow {
		c.FractalLow = target.BarLowPrice
	}
}

func (c *ZigZagCalculation) calculateFractal(r *ZigZagCalculation, target *ZigZag) {
	if target.BarHighPrice > r.ResistanceFractal.Price {
		c.ResistanceFractal = newPrice(target.BarHighPrice, target.BarDateTime)
		c.PriceLow = copyPrice(c.ResistanceFractal)
		c.UpTrend = true
		c.BackStepUp = 0
	}
	if target.BarLowPrice < r.SupportFractal.Price {
		c.SupportFractal = newPrice(target.BarLowPrice, target.BarDateTime)
		c.PriceHigh = copyPrice(c.SupportFractal)
		c.UpTrend = false
		c.BackStepDown = 0
	}
}

func (c *ZigZagCalculation) calculateBreakResistance(r *ZigZagCalculation, target *ZigZag) {
	if target.BarHighPrice > r.Resistance.Price {
		c.UpTrend = true
		if target.BarLowPrice >= c.Support.Price {
			c.BreakSupport = false
		}
		c.Resistance = newPrice(target.BarHighPrice, target.BarDateTime)
		c.ResistanceFractal = newPrice(target.BarHighPrice, target.BarDateTime)
		c.PriceHigh = newPrice(target.BarHighPrice, target.BarDateTime)

		if !r.BreakResistance {
			if r.BreakSupport && !c.BreakSupport {
				c.addWave(r.Resistance.BarDateTime, c.Support.BarDateTime,
					r.Resistance.Price, c.Support.Price, r.Wave, "BR1_RV1x")
			}
			c.BreakResistance = true
			c.Support = copyPrice(c.SupportFractal)
			if r.Wave < 0 {
				c.addWave(r.Support.BarDateTime, c.SupportFractal.BarDateTime,
					r.Resistance.Price, r.Support.Price, r.Wave, "BR1_RV1")
				c.Wave = 1
			} else {
				if r.Wave%2 == 0 {
					c.addWave(r.Resistance.BarDateTime, c.SupportFractal.BarDateTime,
						r.Resistance.Price, r.FractalLow, r.Wave, "BR1_UP2")
				} else {
					c.addWave(r.Support.BarDateTime, c.Resistance.BarDateTime,
						c.Resistance.Price, r.Support.Price, r.Wave, "BR1_UP1")
				}
				c.Wave = c.Wave + 1
			}
		}

		if r.BreakResistance && r.UpTrend && r.BackStepDown == 2 && c.BackStepDown > 0 {
			c.SupportFractal = copyPrice(c.BackStepLow)
			c.Support = copyPrice(c.SupportFractal)
			if r.Wave < 0 {
				c.addWave(r.Support.BarDateTime, c.SupportFractal.BarDateTime,
					r.Resistance.Price, r.Support.Price, r.Wave, "BR2_RV1")
				c.Wave = 1
			} else {
				if r.Wave%2 == 0 {
					c.addWave(r.Resistance.BarDateTime, c.SupportFractal.BarDateTime,
						r.Resistance.Price, r.FractalLow, r.Wave, "BR2_UP2")
				} else {
					c.addWave(r.Support.BarDateTime, c.Resistance.BarDateTime,
						c.Resistance.Price, r.Support.Price, r.Wave, "BR2_UP1")
				}
				c.Wave = c.Wave + 1
			}
		}
		c.FractalLow = target.BarHighPrice
	}
}

func (c *ZigZagCalculation) calculateBreakSupport(r *ZigZagCalculation, target *ZigZag) {
	if target.BarLowPrice < r.Support.Price {
		c.UpTrend = false
		if target.BarHighPrice <= c.Resistance.Price {
			c.BreakResistance = false
		}
		c.Support = newPrice(target.BarLowPrice, target.BarDateTime)
		c.SupportFractal = newPrice(target.BarLowPrice, target.BarDateTime)
		c.PriceLow = newPrice(target.BarLowPrice, target.BarDateTime)

		if !r.BreakSupport {
			if r.BreakResistance && !c.BreakResistance {
				c.addWave(r.Support.BarDateTime, c.Resistance.BarDateTime,
					c.Resistance.Price, r.Support.Price, r.Wave, "BS1_RV1x")
			}
			c.BreakSupport = true
			c.Resistance = copyPrice(c.ResistanceFractal)
			if r.Wave > 0 {
				c.addWave(r.Resistance.BarDateTime, r.ResistanceFractal.BarDateTime,
					r.Resistance.Price, r.Support.Price, r.Wave, "BS1_RV1")
				c.Wave = -1
			} else {
				if r.Wave%2 == 0 {
					c.addWave(r.Support.BarDateTime, r.ResistanceFractal.BarDateTime,
						r.FractalHigh, r.Support.Price, r.Wave, "BS1_DW2")
				} else {
					c.addWave(r.Support.BarDateTime, c.ResistanceFractal.BarDateTime,
						r.Resistance.Price, r.Support.Price, r.Wave, "BS1_DW1")
				}
				c.Wave = c.Wave - 1
			}
		}

		if r.BreakSupport && !r.UpTrend && r.BackStepUp == 2 && c.BackStepUp > 0 {
			c.ResistanceFractal = copyPrice(c.BackStepHigh)
			c.Resistance = copyPrice(c.ResistanceFractal)
			if r.Wave > 0 {
				c.addWave(r.Resistance.BarDateTime, c.ResistanceFractal.BarDateTime,
					r.Resistance.Price, r.Support.Price, r.Wave, "BS2_RV1")
				c.Wave = -1
			} else {
				if r.Wave%2 == 0 {
					c.addWave(r.Support.BarDateTime, c.ResistanceFractal.BarDateTime,
						r.FractalHigh, r.Support.Price, r.Wave, "BS2_DW2")
				} else {
					c.addWave(r.Resistance.BarDateTime, r.Support.BarDateTime,
						r.Resistance.Price, r.Support.Price, r.Wave, "BS2_DW1")
				}
				c.Wave = c.Wave - 1
			}
		}
		c.FractalHigh = target.BarLowPrice
	}
}

func (c *ZigZagCalculation) updateBackStepHigh(r *ZigZagCalculation, target *ZigZag) {
	if r.BackStepUp > 0 && target.BarHighPrice > r.BackStepHigh.Price {
		c.BackStepHigh = newPrice(target.BarHighPrice, target.BarDateTime)
	}
}

func (c *ZigZagCalculation) updateBackStepLow(r *ZigZagCalculation, target *ZigZag) {
	if r.BackStepDown > 0 && target.BarLowPrice < r.BackStepLow.Price {
		c.BackStepLow = newPrice(target.BarLowPrice, target.BarDateTime)
	}
}

func (c *ZigZagCalculation) commitBackStepUp(r *ZigZagCalculation, targetDT time.Time) {
	c.BackStepUp = 0
	c.UpTrend = true

	if r.BreakSupport {
		c.BreakSupport = false
		if r.Wave%2 == 0 {
			c.addWave(r.Resistance.BarDateTime, r.Support.BarDateTime,
				r.Resistance.Price, r.Support.Price, r.Wave, "BSU_DW2")
		} else {
			c.addWave(r.Resistance.BarDateTime, r.Support.BarDateTime,
				r.Resistance.Price, r.Support.Price, r.Wave, "BSU_DW1")
		}
		c.Wave--
	}

	if r.BreakResistance {
		c.BreakResistance = false
		if r.Wave%2 == 0 {
			c.addWave(r.Resistance.BarDateTime, c.SupportFractal.BarDateTime,
				r.Resistance.Price, r.FractalLow, r.Wave, "BSU_UP2")
		} else {
			c.addWave(r.Support.BarDateTime, c.Resistance.BarDateTime,
				c.Resistance.Price, r.Support.Price, r.Wave, "BSU_UP1")
		}
		c.Wave++
	}

	c.ResistanceFractal = copyPrice(c.BackStepHigh)
	c.SupportFractal = copyPrice(r.PriceLow)
	c.PriceHigh = copyPrice(c.BackStepHigh)
	c.PriceLow = copyPrice(c.BackStepHigh)

	if r.BreakResistance && r.UpTrend {
		c.BreakResistance = false
	}
}

func (c *ZigZagCalculation) commitBackStepDown(r *ZigZagCalculation, targetDT time.Time) {
	c.BackStepDown = 0
	c.UpTrend = false

	if r.BreakResistance {
		c.BreakResistance = false
		if r.Wave%2 == 0 {
			c.addWave(r.Resistance.BarDateTime, c.SupportFractal.BarDateTime,
				r.Resistance.Price, r.FractalLow, r.Wave, "BSD_UP2")
		} else {
			c.addWave(r.Support.BarDateTime, c.Resistance.BarDateTime,
				c.Resistance.Price, r.Support.Price, r.Wave, "BSD_UP1")
		}
		c.Wave++
	}

	if r.BreakSupport {
		c.BreakSupport = false
		if r.Wave%2 == 0 {
			c.addWave(r.Support.BarDateTime, c.ResistanceFractal.BarDateTime,
				r.FractalHigh, r.Support.Price, r.Wave, "BSD_DW2")
		} else {
			c.addWave(r.Resistance.BarDateTime, r.Support.BarDateTime,
				r.Resistance.Price, r.Support.Price, r.Wave, "BSD_DW1")
		}
		c.Wave--
	}

	c.ResistanceFractal = copyPrice(r.PriceHigh)
	c.SupportFractal = copyPrice(c.BackStepLow)
	c.PriceHigh = copyPrice(c.BackStepLow)
	c.PriceLow = copyPrice(c.BackStepLow)

	if r.BreakSupport && !r.UpTrend {
		c.BreakSupport = false
	}
}

func (c *ZigZagCalculation) addWave(from, to time.Time, waveResistance, waveSupport float64, waveNo int, memo string) {
	if from.Equal(to) {
		return
	}

	previousWaveStart := minDateTime
	previousWaveNo := 0

	if len(c.WaveList) > 0 {
		last := c.WaveList[len(c.WaveList)-1]
		if last.WaveStart.Equal(from) && last.WaveEnd.Equal(to) {
			return
		}
		previousWaveStart = last.WaveStart
		previousWaveNo = last.Wave
	}

	c.WaveList = append(c.WaveList, ZigZagWave{
		WaveStart:         from,
		WaveEnd:           to,
		Resistance:        waveResistance,
		Support:           waveSupport,
		Wave:              waveNo,
		PreviousWaveStart: previousWaveStart,
		PreviousWave:      previousWaveNo,
		WaveMemo:          memo,
	})
}
