package fxdto

import fxmodel "sandbox-api-gin/internal/domain/model/fx"

type EconomicIndicatorDto struct {
	ID               *int64 `json:"id"`
	CountryCode      string `json:"countryCode" binding:"required"`
	Name             string `json:"name" binding:"required"`
	Importance       string `json:"importance" binding:"required"`
	Description      string `json:"description"`
	UnitOfValue      string `json:"unitOfValue"`
	CountryName      string `json:"countryName"`
	CountryNameShort string `json:"countryNameShort"`
}

func EconomicIndicatorDtoFromDomain(m fxmodel.EconomicIndicator) EconomicIndicatorDto {
	return EconomicIndicatorDto{
		ID:               &m.ID,
		CountryCode:      m.CountryCode,
		Name:             m.Name,
		Importance:       m.Importance,
		Description:      m.Description,
		UnitOfValue:      m.UnitOfValue,
		CountryName:      m.CountryName,
		CountryNameShort: m.CountryNameShort,
	}
}
