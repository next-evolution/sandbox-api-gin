package fxrequest

type BarDataSearchRequest struct {
	BarType     string `json:"barType" binding:"required"`
	Symbol      string `json:"symbol" binding:"required"`
	BarDateFrom string `json:"barDateFrom"`
	BarDateTo   string `json:"barDateTo"`
	SortAsc     bool   `json:"sortAsc"`
	Page        int    `json:"page" binding:"required,min=1"`
	Size        int    `json:"size" binding:"required,min=1"`
}
