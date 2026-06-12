package fx

import "time"

type EconomicIndicator struct {
	ID               int64
	CountryCode      string
	Name             string
	Importance       string
	Description      string
	UnitOfValue      string
	CountryName      string
	CountryNameShort string
	Deleted          bool
	CreatedAt        time.Time
	CreatedBy        string
	UpdatedAt        time.Time
	UpdatedBy        string
}
