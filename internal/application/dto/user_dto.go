package dto

import (
	"encoding/json"
	"sandbox-api-gin/internal/domain/model"
	"time"
)

// JavaのLocalDateTime形式 "yyyy-MM-dd HH:mm:ss" に対応するカスタム時刻型
type DateTime struct {
	time.Time
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Format("2006-01-02 15:04:05"))
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

type UserDto struct {
	ID           int64     `json:"id"`
	UserID       string    `json:"userId"`
	EmailAddress string    `json:"emailAddress"`
	NickName     string    `json:"nickName"`
	Approved     bool      `json:"approved"`
	ApprovedAt   *DateTime `json:"approvedAt"`
	Admin        bool      `json:"admin"`
	Blocked      bool      `json:"blocked"`
	CreatedAt    DateTime  `json:"createdAt"`
	UpdatedAt    DateTime  `json:"updatedAt"`
}

func UserDtoFrom(user *model.User) *UserDto {
	dto := &UserDto{
		ID:           user.ID,
		UserID:       user.UserID,
		EmailAddress: user.EmailAddress,
		NickName:     user.NickName,
		Approved:     user.Approved,
		Admin:        user.Admin,
		Blocked:      user.Blocked,
		CreatedAt:    DateTime{user.CreatedAt},
		UpdatedAt:    DateTime{user.UpdatedAt},
	}
	if user.ApprovedAt != nil {
		dt := DateTime{*user.ApprovedAt}
		dto.ApprovedAt = &dt
	}
	return dto
}
