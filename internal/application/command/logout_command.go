package command

import "sandbox-api-gin/internal/domain/model"

type LogoutCommand struct {
	AuthUser      *model.AuthUser
	EncodedUserID string
}
