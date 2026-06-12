package fxrequest

import fxdto "sandbox-api-gin/internal/application/dto/fx"

type CountryRequest struct {
	Country fxdto.CountryDto `json:"country" binding:"required"`
}
