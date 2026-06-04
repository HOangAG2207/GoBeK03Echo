package health

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	PingRedis(ctx context.Context) error
}
type repository struct {
	redis *redis.Client
}

func NewRepository(r *redis.Client) Repository {
	return &repository{
		redis: r,
	}
}
