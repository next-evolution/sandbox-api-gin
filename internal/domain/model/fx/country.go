package fx

import "time"

type Country struct {
	Code         string
	Name         string
	CurrencyCode string
	NameEn       string
	NameShort    string
	SortOrder    int16
	Deleted      bool
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
}
