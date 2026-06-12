package fx

import "time"

type SummerTime struct {
	TargetYear int16
	ApplyStart time.Time
	ApplyEnd   time.Time
}
