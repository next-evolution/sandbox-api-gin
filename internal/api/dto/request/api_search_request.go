package request

type ApiSearchRequest struct {
	Page int `json:"page" binding:"required,min=1"`
	Size int `json:"size" binding:"required,min=1"`
}
