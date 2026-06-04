package links

import (
	"context"

	pkgredis "github.com/HOangAG2207/GoBeK03Echo/pkg/redis"
)

type Repository interface {
	StoreURL(ctx context.Context, code, url string, exptime int64) error
	GetURL(ctx context.Context, code string) (string, error)
}

type repository struct {
	redis *pkgredis.Redis
}

func NewRepository(r *pkgredis.Redis) Repository {
	return &repository{
		redis: r,
	}
}
