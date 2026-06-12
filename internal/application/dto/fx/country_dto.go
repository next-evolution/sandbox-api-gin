package fxdto

import (
	"time"

	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type CountryDto struct {
	Code         string `json:"code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	CurrencyCode string `json:"currencyCode" binding:"required"`
	NameEn       string `json:"nameEn" binding:"required"`
	NameShort    string `json:"nameShort" binding:"required"`
	SortOrder    int16  `json:"sortOrder"`
}

func CountryDtoFromDomain(c fxmodel.Country) CountryDto {
	return CountryDto{
		Code:         c.Code,
		Name:         c.Name,
		CurrencyCode: c.CurrencyCode,
		NameEn:       c.NameEn,
		NameShort:    c.NameShort,
		SortOrder:    c.SortOrder,
	}
}

func (d CountryDto) ToDomain(author string) fxmodel.Country {
	now := time.Now()
	return fxmodel.Country{
		Code:         d.Code,
		Name:         d.Name,
		CurrencyCode: d.CurrencyCode,
		NameEn:       d.NameEn,
		NameShort:    d.NameShort,
		SortOrder:    d.SortOrder,
		Deleted:      false,
		CreatedAt:    now,
		CreatedBy:    author,
		UpdatedAt:    now,
		UpdatedBy:    author,
	}
}
