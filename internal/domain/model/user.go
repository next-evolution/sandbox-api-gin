package model

import (
	"sandbox-api-gin/internal/domain/apperror"
	"time"
)

type User struct {
	ID           int64
	UserID       string
	EmailAddress string
	NickName     string
	Approved     bool
	ApprovedAt   *time.Time
	Admin        bool
	Blocked      bool
	Deleted      bool
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
}

func (u *User) CheckBlocked() error {
	if u.Blocked {
		return apperror.NewAuthenticationError("blocked.")
	}
	return nil
}
