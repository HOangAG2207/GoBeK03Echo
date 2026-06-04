package links

import "context"

func (r *repository) StoreURL(ctx context.Context, code, url string, exptime int64) error {
	return r.redis.SetValue(ctx, code, url, exptime)
}
