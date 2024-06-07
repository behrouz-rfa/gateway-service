package ports

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository interface {
	SetTTL(ctx context.Context, key string, value interface{}, ttl time.Duration)
	Set(ctx context.Context, key string, value interface{})
	Get(ctx context.Context, key string, dest interface{}) error
	SetCachedResponse(query string, response []byte)
	GetCachedResponse(query string) ([]byte, bool)
	GetClient() *redis.Client
}
