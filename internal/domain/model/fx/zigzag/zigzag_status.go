package zigzag

type ZigZagStatus struct {
	SymbolType string
	BarType    string
	Symbol     string
	Depth      int16

	BarDateTimeMin string
	BarDateTimeMax string
	BarCount       int

	BarDateTimeMinZigZag string
	BarDateTimeMaxZigZag string
	ZigzagCount          int
	BreakResistanceCount int
	BreakSupportCount    int

	Message string
}
