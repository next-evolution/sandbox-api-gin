package fx

import "time"

type EconomicIndicator struct {
	Code             string
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
