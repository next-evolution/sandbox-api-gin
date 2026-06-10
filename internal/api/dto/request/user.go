package request

type UserRegistrationRequest struct {
	NickName string `json:"nickName" binding:"required,max=50"`
}

type UserUpdateRequest struct {
	NickName string `json:"nickName" binding:"required,max=50"`
}
