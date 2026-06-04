package pkgredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewClient(envPrefix string) (*Redis, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.redis.Address,
		Password: cfg.redis.Password,
		DB:       cfg.redis.DB,
	})

	return &Redis{client: client}, nil
}
func toDuration(expire int64) time.Duration {
	if expire <= 0 {
		return 0 // no expiration
	}
	return time.Duration(expire) * time.Second
}
func (r *Redis) SetValue(ctx context.Context, key string, value string, expire int64) error {
	return r.client.Set(ctx, key, value, toDuration(expire)).Err()
}
func (r *Redis) GetValue(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
