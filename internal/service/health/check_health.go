package health

import (
	"context"
	"errors"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
)

var (
	ErrRedisConnection = errors.New("Redis: kết nối thất bại")
)

func (s *service) CheckHealth(ctx context.Context) (*model.HealthCheckResponse, error) {
	if err := s.repo.PingRedis(ctx); err != nil {
		return nil, ErrRedisConnection
	}
	// id := s.instanceID
	// if id == "" {
	// 	id = utils.UuidGenerator()
	// }
	return &model.HealthCheckResponse{
		Message:     "OK",
		ServiceName: s.serviceName,
		InstanceID:  s.instanceID,
	}, nil
}
