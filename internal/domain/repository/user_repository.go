package repository

import (
	"context"

	"sandbox-api-gin/internal/domain/model"
)

type UserRepository interface {
	Login(ctx context.Context, userID, email string) (*model.User, error)
	FindByUserID(ctx context.Context, userID string) (*model.User, error)
	ExistsByUserID(ctx context.Context, userID string) (bool, error)
	InsertUser(ctx context.Context, user *model.User) error
	UpdateNickName(ctx context.Context, user *model.User) error
}
