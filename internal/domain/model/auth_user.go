package model

type AuthUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	Admin         bool   `json:"admin"`
}

func (a *AuthUser) IsAdmin() bool {
	return a.Admin
}
