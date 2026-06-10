package request

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
}

type LogoutRequest struct {
	UserID string `json:"userId" binding:"required"`
}
