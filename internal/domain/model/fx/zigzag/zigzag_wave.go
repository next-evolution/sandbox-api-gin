package zigzag

import "time"

type ZigZagWave struct {
	Symbol            string
	Depth             int
	WaveStart         time.Time
	WaveEnd           time.Time
	Wave              int
	Resistance        float64
	Support           float64
	PreviousWaveStart time.Time
	PreviousWave      int
	WaveMemo          string
}
