package usecase

import (
	"context"
	"encoding/base64"
	"log/slog"

	"sandbox-api-gin/internal/application/command"
	"sandbox-api-gin/internal/domain/repository"
)

type LogoutUseCase struct {
	sessionRepo repository.SessionRepository
}

func NewLogoutUseCase(sessionRepo repository.SessionRepository) *LogoutUseCase {
	return &LogoutUseCase{sessionRepo: sessionRepo}
}

func (u *LogoutUseCase) Execute(ctx context.Context, cmd *command.LogoutCommand) {
	authUser := cmd.AuthUser
	if authUser == nil {
		return
	}

	userIDBytes, err := base64.StdEncoding.DecodeString(cmd.EncodedUserID)
	if err != nil {
		userIDBytes, err = base64.RawStdEncoding.DecodeString(cmd.EncodedUserID)
		if err != nil {
			slog.Error("logout decode failed", "error", err)
			return
		}
	}
	userIDValue := string(userIDBytes)

	if authUser.Sub == userIDValue {
		if err := u.sessionRepo.DeleteBySub(ctx, authUser.Sub); err != nil {
			slog.Error("logout session delete failed", "error", err)
		}
	} else {
		slog.Error("logout failed", "reqUserId", userIDValue, "tokenUserId", authUser.Sub)
	}
}
