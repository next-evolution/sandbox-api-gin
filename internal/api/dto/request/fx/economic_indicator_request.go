package fxrequest

import fxdto "sandbox-api-gin/internal/application/dto/fx"

type EconomicIndicatorRequest struct {
	Indicator fxdto.EconomicIndicatorDto `json:"indicator" binding:"required"`
}

type EconomicIndicatorSearchRequest struct {
	Page        int    `json:"page" binding:"required,min=1"`
	Size        int    `json:"size" binding:"required,min=1"`
	CountryCode string `json:"countryCode"`
	Importance  string `json:"importance"`
	Name        string `json:"name"`
}
