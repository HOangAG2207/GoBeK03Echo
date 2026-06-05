package links

import (
	"context"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name Repository --filename links_repository_mock.go --output ./mocks
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
