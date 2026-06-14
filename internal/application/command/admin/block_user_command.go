package admin

type BlockUserCommand struct {
	UserID    string
	Blocked   bool
	UpdatedBy string
}
