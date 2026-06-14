package infraredis

import (
	"context"
	"fmt"
	"sort"
	"strings"

	goredis "github.com/redis/go-redis/v9"

	"sandbox-api-gin/internal/domain/repository"
)

type RedisMasterCacheRepository struct {
	client *goredis.Client
}

func NewRedisMasterCacheRepository(client *goredis.Client) repository.MasterCacheRepository {
	return &RedisMasterCacheRepository{client: client}
}

func (r *RedisMasterCacheRepository) GetStatus(ctx context.Context) (string, error) {
	keys, err := r.client.Keys(ctx, "master*").Result()
	if err != nil {
		return "", err
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, key := range keys {
		keyType, err := r.client.Type(ctx, key).Result()
		if err != nil {
			return "", err
		}
		count := 0
		if keyType == "list" {
			n, err := r.client.LLen(ctx, key).Result()
			if err != nil {
				return "", err
			}
			count = int(n)
		}
		sb.WriteString(fmt.Sprintf("%s=%d\n", key, count))
	}
	return sb.String(), nil
}

func (r *RedisMasterCacheRepository) DeleteByPattern(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}
