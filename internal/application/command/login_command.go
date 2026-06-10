package command

import "sandbox-api-gin/internal/domain/model"

type LoginCommand struct {
	AuthUser     *model.AuthUser
	EncodedEmail string
}
