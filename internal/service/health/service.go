package health

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/HOangAG2207/GoBeK03Echo/internal/repository/health"
)

type Service interface {
	CheckHealth(ctx context.Context) (*model.HealthCheckResponse, error)
}

type service struct {
	repo        health.Repository
	serviceName string
	instanceID  string
}

func NewService(r health.Repository, serviceName, instanceID string) Service {
	return &service{
		repo:        r,
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}
