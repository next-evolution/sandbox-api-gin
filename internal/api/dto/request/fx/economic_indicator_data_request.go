package fxrequest

import fxdto "sandbox-api-gin/internal/application/dto/fx"

type EconomicIndicatorDataRequest struct {
	Data fxdto.EconomicIndicatorDataDto `json:"data" binding:"required"`
}

type EconomicIndicatorDataSearchRequest struct {
	Page                int    `json:"page" binding:"required,min=1"`
	Size                int    `json:"size" binding:"required,min=1"`
	Code                string `json:"code"`
	Importance          string `json:"importance"`
	CountryCode         string `json:"countryCode"`
	PublicationBaseDate string `json:"publicationBaseDate"`
	SortAsc             bool   `json:"sortAsc"`
}
