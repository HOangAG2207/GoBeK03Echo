package health

import (
	"context"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name Repository --filename health_repository_mock.go --output ./mocks
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
