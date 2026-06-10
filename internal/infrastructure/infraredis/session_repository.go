package infraredis

import (
	"context"
	"encoding/json"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"sandbox-api-gin/internal/domain/model"
	"sandbox-api-gin/internal/domain/repository"
)

const keyPrefix = "session:"

type RedisSessionRepository struct {
	client     *goredis.Client
	sessionTTL time.Duration
}

func NewRedisSessionRepository(client *goredis.Client, sessionTTL int) repository.SessionRepository {
	return &RedisSessionRepository{
		client:     client,
		sessionTTL: time.Duration(sessionTTL) * time.Second,
	}
}

func (r *RedisSessionRepository) Save(ctx context.Context, authUser *model.AuthUser) error {
	key := keyPrefix + authUser.Sub
	data, err := json.Marshal(authUser)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.sessionTTL).Err()
}

func (r *RedisSessionRepository) FindBySub(ctx context.Context, sub string) (*model.AuthUser, error) {
	key := keyPrefix + sub
	data, err := r.client.Get(ctx, key).Bytes()
	if err == goredis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var authUser model.AuthUser
	if err := json.Unmarshal(data, &authUser); err != nil {
		return nil, err
	}
	return &authUser, nil
}

func (r *RedisSessionRepository) DeleteBySub(ctx context.Context, sub string) error {
	key := keyPrefix + sub
	return r.client.Del(ctx, key).Err()
}

func (r *RedisSessionRepository) Update(ctx context.Context, authUser *model.AuthUser) error {
	key := keyPrefix + authUser.Sub
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return r.client.Expire(ctx, key, r.sessionTTL).Err()
	}
	return nil
}
