package links

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/utils"
)

func (r *repository) StoreURL(ctx context.Context, code, url string, exptime int64) error {
	return r.redis.Set(ctx, code, url, utils.Int64ToDuration(exptime)).Err()
}
