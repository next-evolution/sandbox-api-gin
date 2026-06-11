package fxdto

import (
	"encoding/json"
	"time"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

// LocalDate はレスポンス用に "yyyy-MM-dd" 形式でシリアライズされる日付型。
type LocalDate struct {
	time.Time
}

func (d LocalDate) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Format("2006-01-02"))
}

func (d *LocalDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

type SummerTimeDto struct {
	TargetYear int16     `json:"targetYear" binding:"required"`
	ApplyStart LocalDate `json:"applyStart"`
	ApplyEnd   LocalDate `json:"applyEnd"`
}

func SummerTimeDtoFromDomain(s fxmodel.SummerTime) SummerTimeDto {
	return SummerTimeDto{
		TargetYear: s.TargetYear,
		ApplyStart: LocalDate{s.ApplyStart},
		ApplyEnd:   LocalDate{s.ApplyEnd},
	}
}

func (d SummerTimeDto) ToDomain() fxmodel.SummerTime {
	return fxmodel.SummerTime{
		TargetYear: d.TargetYear,
		ApplyStart: d.ApplyStart.Time,
		ApplyEnd:   d.ApplyEnd.Time,
	}
}
