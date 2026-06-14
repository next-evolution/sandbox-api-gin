package adminrequest

type UserAdminRequest struct {
	Admin *bool `json:"admin" binding:"required"`
}
