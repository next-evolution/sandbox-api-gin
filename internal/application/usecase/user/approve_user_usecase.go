package user

import (
	"context"
	"time"

	admincommand "sandbox-api-gin/internal/application/command/admin"
	"sandbox-api-gin/internal/application/dto"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/repository"
)

type ApproveUserUseCase struct {
	userRepo repository.UserRepository
}

func NewApproveUserUseCase(userRepo repository.UserRepository) *ApproveUserUseCase {
	return &ApproveUserUseCase{userRepo: userRepo}
}

func (uc *ApproveUserUseCase) Execute(ctx context.Context, cmd *admincommand.ApproveUserCommand) (*dto.UserDto, error) {
	user, err := uc.userRepo.FindByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("ユーザが存在しません")
	}

	if err := user.CheckAlreadyApproved(); err != nil {
		return nil, err
	}

	now := time.Now()
	user.Approved = true
	user.ApprovedAt = &now
	user.UpdatedAt = now
	user.UpdatedBy = cmd.UpdatedBy

	if err := uc.userRepo.UpdateApproved(ctx, user); err != nil {
		return nil, err
	}

	return dto.UserDtoFrom(user), nil
}
