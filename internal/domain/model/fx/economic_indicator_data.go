package fx

import "time"

type EconomicIndicatorData struct {
	Code             string
	CountryCode      string
	Name             string
	Importance       string
	Description      string
	Publication      time.Time
	PublicationDate  string
	PublicationTime  string
	DayOfWeek        int16
	SubTitle         string
	ResultValue      string
	ForecastValue    string
	PreviousValue    string
	UnitOfValue      string
	Memo             string
	CountryName      string
	CountryNameShort string
}
