package adminrequest

type UserBlockRequest struct {
	Blocked *bool `json:"blocked" binding:"required"`
}
