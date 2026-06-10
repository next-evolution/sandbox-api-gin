package usecase

import (
	"context"

	"sandbox-api-gin/internal/application/command"
	"sandbox-api-gin/internal/application/dto"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/repository"
)

type GetProfileUseCase struct {
	userRepo repository.UserRepository
}

func NewGetProfileUseCase(userRepo repository.UserRepository) *GetProfileUseCase {
	return &GetProfileUseCase{userRepo: userRepo}
}

// Execute はユーザープロフィールを返す。未承認の場合は (nil, nil) を返す。
func (u *GetProfileUseCase) Execute(ctx context.Context, cmd *command.GetProfileCommand) (*dto.UserDto, error) {
	user, err := u.userRepo.FindByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("ユーザが存在しません")
	}
	if !user.Approved {
		return nil, nil
	}
	return dto.UserDtoFrom(user), nil
}
