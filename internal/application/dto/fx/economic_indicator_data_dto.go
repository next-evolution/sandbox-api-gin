package fxdto

import (
	"sandbox-api-gin/internal/application/dto"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type EconomicIndicatorDataDto struct {
	ID               int64        `json:"id"`
	CountryCode      string       `json:"countryCode"`
	Name             string       `json:"name"`
	Importance       string       `json:"importance"`
	Description      string       `json:"description"`
	Publication      dto.DateTime `json:"publication"`
	PublicationDate  string       `json:"publicationDate"`
	PublicationTime  string       `json:"publicationTime"`
	DayOfWeek        int16        `json:"dayOfWeek"`
	SubTitle         string       `json:"subTitle"`
	ResultValue      string       `json:"resultValue"`
	ForecastValue    string       `json:"forecastValue"`
	PreviousValue    string       `json:"previousValue"`
	UnitOfValue      string       `json:"unitOfValue"`
	Memo             string       `json:"memo"`
	CountryName      string       `json:"countryName"`
	CountryNameShort string       `json:"countryNameShort"`
}

func EconomicIndicatorDataDtoFromDomain(m fxmodel.EconomicIndicatorData) EconomicIndicatorDataDto {
	return EconomicIndicatorDataDto{
		ID:               m.ID,
		CountryCode:      m.CountryCode,
		Name:             m.Name,
		Importance:       m.Importance,
		Description:      m.Description,
		Publication:      dto.DateTime{Time: m.Publication},
		PublicationDate:  m.PublicationDate,
		PublicationTime:  m.PublicationTime,
		DayOfWeek:        m.DayOfWeek,
		SubTitle:         m.SubTitle,
		ResultValue:      m.ResultValue,
		ForecastValue:    m.ForecastValue,
		PreviousValue:    m.PreviousValue,
		UnitOfValue:      m.UnitOfValue,
		Memo:             m.Memo,
		CountryName:      m.CountryName,
		CountryNameShort: m.CountryNameShort,
	}
}
