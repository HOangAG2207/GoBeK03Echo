package links

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers"
)

func (r *repository) StoreURL(ctx context.Context, code, url string, exptime int64) error {
	return r.redis.Set(ctx, code, url, helpers.Int64ToDuration(exptime)).Err()
}
