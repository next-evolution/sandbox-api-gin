package user

import (
	"context"
	"encoding/base64"

	"sandbox-api-gin/internal/application/command"
	"sandbox-api-gin/internal/application/dto"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/model"
	"sandbox-api-gin/internal/domain/repository"
)

type LoginUseCase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewLoginUseCase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *LoginUseCase {
	return &LoginUseCase{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (u *LoginUseCase) Execute(ctx context.Context, cmd *command.LoginCommand) (*dto.UserDto, error) {
	authUser := cmd.AuthUser
	if authUser == nil {
		return nil, apperror.NewAuthenticationError("login failed.")
	}

	// BASE64デコード
	decodedEmail, err := decodeBase64(cmd.EncodedEmail)
	if err != nil {
		return nil, apperror.NewAuthenticationError("login failed.")
	}

	// JWT内emailとリクエストbodyのemailを照合
	if authUser.Email != decodedEmail {
		return nil, apperror.NewAuthenticationError("login failed.")
	}

	// ユーザーテーブルから取得（adminフラグ含む）
	user, err := u.userRepo.Login(ctx, authUser.Sub, decodedEmail)
	if err != nil {
		return nil, err
	}

	// blockedチェック
	if user != nil {
		if err := user.CheckBlocked(); err != nil {
			return nil, err
		}
	}

	// adminフラグ付きAuthUserをRedisに保存
	adminFlag := false
	if user != nil {
		adminFlag = user.Admin
	}
	authUserWithAdmin := &model.AuthUser{
		Sub:           authUser.Sub,
		Email:         authUser.Email,
		EmailVerified: authUser.EmailVerified,
		Admin:         adminFlag,
	}
	if err := u.sessionRepo.Save(ctx, authUserWithAdmin); err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}
	return dto.UserDtoFrom(user), nil
}

func decodeBase64(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(encoded)
		if err != nil {
			return "", err
		}
	}
	return string(decoded), nil
}
