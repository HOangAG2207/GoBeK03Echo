package links

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrCodeNotFound = errors.New("Code not found")

func (s *service) GetURL(ctx context.Context, code string) (string, error) {
	url, err := s.repo.GetURL(ctx, code)
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeNotFound
	}

	if err != nil {
		return "", err
	}

	return url, nil
}
