package links

import "context"

func (r *repository) GetURL(ctx context.Context, code string) (string, error) {
	return r.redis.GetValue(ctx, code)
}
