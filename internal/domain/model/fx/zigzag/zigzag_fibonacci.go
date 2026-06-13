package zigzag

const (
	FibF7 = 0.786
	FibF6 = 0.618
	FibF5 = 0.5
	FibF3 = 0.382
	FibF2 = 0.236
)

type ZigZagFibonacci struct {
	F1         float64
	F7         float64
	F6         float64
	F5         float64
	F3         float64
	F2         float64
	F0         float64
	PriceRange float64
}
