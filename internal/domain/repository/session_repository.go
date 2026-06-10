package repository

import (
	"context"

	"sandbox-api-gin/internal/domain/model"
)

type SessionRepository interface {
	Save(ctx context.Context, authUser *model.AuthUser) error
	FindBySub(ctx context.Context, sub string) (*model.AuthUser, error)
	DeleteBySub(ctx context.Context, sub string) error
	Update(ctx context.Context, authUser *model.AuthUser) error
}
