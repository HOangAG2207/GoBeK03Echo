package health

import "context"

func (h *repository) PingRedis(ctx context.Context) error {
	return h.redis.Ping(ctx).Err()
}
