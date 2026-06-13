package user

import (
	"context"
	"time"

	"sandbox-api-gin/internal/application/command"
	"sandbox-api-gin/internal/application/dto"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/model"
	"sandbox-api-gin/internal/domain/repository"
)

type RegisterUserUseCase struct {
	userRepo repository.UserRepository
}

func NewRegisterUserUseCase(userRepo repository.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{userRepo: userRepo}
}

func (u *RegisterUserUseCase) Execute(ctx context.Context, cmd *command.RegisterUserCommand) (*dto.UserDto, error) {
	exists, err := u.userRepo.ExistsByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.NewDuplicateError("登録済みのユーザです")
	}

	now := time.Now()
	user := &model.User{
		UserID:       cmd.UserID,
		EmailAddress: cmd.Email,
		NickName:     cmd.NickName,
		Approved:     false,
		ApprovedAt:   nil,
		Admin:        false,
		Blocked:      false,
		Deleted:      false,
		CreatedAt:    now,
		CreatedBy:    cmd.UserID,
		UpdatedAt:    now,
		UpdatedBy:    cmd.UserID,
	}

	if err := u.userRepo.InsertUser(ctx, user); err != nil {
		return nil, err
	}

	return dto.UserDtoFrom(user), nil
}
