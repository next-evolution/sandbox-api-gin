package response

import fxdto "sandbox-api-gin/internal/application/dto/fx"

type CountrySearchResponse struct {
	ReturnCode  ReturnCode        `json:"returnCode"`
	TotalCount  int               `json:"totalCount"`
	SearchCount int               `json:"searchCount"`
	TotalPage   int               `json:"totalPage"`
	List        []fxdto.CountryDto `json:"list"`
}
