package user

import (
	"context"
	"time"

	admincommand "sandbox-api-gin/internal/application/command/admin"
	"sandbox-api-gin/internal/application/dto"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/repository"
)

type BlockUserUseCase struct {
	userRepo repository.UserRepository
}

func NewBlockUserUseCase(userRepo repository.UserRepository) *BlockUserUseCase {
	return &BlockUserUseCase{userRepo: userRepo}
}

func (uc *BlockUserUseCase) Execute(ctx context.Context, cmd *admincommand.BlockUserCommand) (*dto.UserDto, error) {
	user, err := uc.userRepo.FindByUserID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("ユーザが存在しません")
	}

	if err := user.CheckBlockDuplicate(cmd.Blocked); err != nil {
		return nil, err
	}

	user.Blocked = cmd.Blocked
	user.UpdatedAt = time.Now()
	user.UpdatedBy = cmd.UpdatedBy

	if err := uc.userRepo.UpdateBlocked(ctx, user); err != nil {
		return nil, err
	}

	return dto.UserDtoFrom(user), nil
}
