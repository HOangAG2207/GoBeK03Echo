package pkgredis

import (
	"github.com/redis/go-redis/v9"
)

func NewClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.redis.Address,
		Password: cfg.redis.Password,
		DB:       cfg.redis.DB,
	})

	return client, nil
}
