package repository

import "context"

type MasterCacheRepository interface {
	GetStatus(ctx context.Context) (string, error)
	DeleteByPattern(ctx context.Context, pattern string) error
}
