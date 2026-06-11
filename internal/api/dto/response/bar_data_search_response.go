package response

import fxmodel "sandbox-api-gin/internal/domain/model/fx"

type BarDataSearchResponse struct {
	ReturnCode  ReturnCode       `json:"returnCode"`
	TotalCount  int              `json:"totalCount"`
	SearchCount int              `json:"searchCount"`
	TotalPage   int              `json:"totalPage"`
	List        []fxmodel.BarData `json:"list"`
}
