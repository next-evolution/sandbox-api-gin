package user

import (
	"context"
	"time"

	admincommand "sandbox-api-gin/internal/application/command/admin"
	"sandbox-api-gin/internal/application/dto"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/repository"
)

type GrantAdminUseCase struct {
	userRepo repository.UserRepository
}

func NewGrantAdminUseCase(userRepo repository.UserRepository) *GrantAdminUseCase {
	return &GrantAdminUseCase{userRepo: userRepo}
}

func (uc *GrantAdminUseCase) Execute(ctx context.Context, cmd *admincommand.GrantAdminCommand) (*dto.UserDto, error) {
	user, err := uc.userRepo.FindByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("ユーザが存在しません")
	}

	if err := user.CheckAdminDuplicate(cmd.Admin); err != nil {
		return nil, err
	}

	user.Admin = cmd.Admin
	user.UpdatedAt = time.Now()
	user.UpdatedBy = cmd.UpdatedBy

	if err := uc.userRepo.UpdateAdmin(ctx, user); err != nil {
		return nil, err
	}

	return dto.UserDtoFrom(user), nil
}
