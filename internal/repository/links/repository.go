package links

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	StoreURL(ctx context.Context, code, url string, exptime int64) error
	GetURL(ctx context.Context, code string) (string, error)
}

type repository struct {
	redis *redis.Client
}

func NewRepository(r *redis.Client) Repository {
	return &repository{
		redis: r,
	}
}
