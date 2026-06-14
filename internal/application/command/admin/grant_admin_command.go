package admin

type GrantAdminCommand struct {
	UserID    string
	Admin     bool
	UpdatedBy string
}
