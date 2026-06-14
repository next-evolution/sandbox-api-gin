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

func (u *User) CheckAlreadyApproved() error {
	if u.Approved {
		return apperror.NewDuplicateError("承認済みです")
	}
	return nil
}

func (u *User) CheckBlockDuplicate(newBlocked bool) error {
	if u.Blocked == newBlocked {
		if newBlocked {
			return apperror.NewDuplicateError("Block済みです")
		}
		return apperror.NewDuplicateError("Block解除済みです")
	}
	return nil
}

func (u *User) CheckAdminDuplicate(newAdmin bool) error {
	if u.Admin == newAdmin {
		if newAdmin {
			return apperror.NewDuplicateError("admin権限設定済みです")
		}
		return apperror.NewDuplicateError("admin権限設定剥奪済みです")
	}
	return nil
}
