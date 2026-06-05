package links

import (
	"context"
	"errors"

	"github.com/HOangAG2207/GoBeK03Echo/internal/utils"
	"github.com/redis/go-redis/v9"
)

const maxRetyAttempts = 10

var ErrMaxRetriesExceeded = errors.New("maximum retry attempts exceeded for generating unique URL code")

func (s *service) ShortenURL(ctx context.Context, url string, codeLength int, exptime int64) (string, error) {
	for range maxRetyAttempts {

		code, err := utils.GenerateCode(codeLength)
		if err != nil {
			return "", err
		}

		// check code exists
		_, err = s.repo.GetURL(ctx, code)

		// exists → retry
		if err == nil {
			continue
		}

		// real error
		if err != redis.Nil {
			return "", err
		}

		// not exists → store
		err = s.repo.StoreURL(ctx, code, url, exptime)
		if err != nil {
			return "", err
		}

		return code, nil
	}

	return "", ErrMaxRetriesExceeded
}
