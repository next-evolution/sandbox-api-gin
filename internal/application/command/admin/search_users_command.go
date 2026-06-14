package admin

type SearchUsersCommand struct {
	EmailAddress string
	Approved     *bool
	Page         int
	Size         int
}
