package response

import fxdto "sandbox-api-gin/internal/application/dto/fx"

type SummerTimeSearchResponse struct {
	ReturnCode  ReturnCode           `json:"returnCode"`
	TotalCount  int                  `json:"totalCount"`
	SearchCount int                  `json:"searchCount"`
	TotalPage   int                  `json:"totalPage"`
	List        []fxdto.SummerTimeDto `json:"list"`
}
