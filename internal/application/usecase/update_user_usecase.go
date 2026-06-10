package usecase

import (
	"context"
	"time"

	"sandbox-api-gin/internal/application/command"
	"sandbox-api-gin/internal/application/dto"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/repository"
)

type UpdateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewUpdateUserUseCase(userRepo repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{userRepo: userRepo}
}

func (u *UpdateUserUseCase) Execute(ctx context.Context, cmd *command.UpdateUserCommand) (*dto.UserDto, error) {
	user, err := u.userRepo.FindByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("ユーザが存在しません")
	}

	user.NickName = cmd.NickName
	user.UpdatedAt = time.Now()
	user.UpdatedBy = cmd.UpdatedBy

	if err := u.userRepo.UpdateNickName(ctx, user); err != nil {
		return nil, err
	}

	return dto.UserDtoFrom(user), nil
}
