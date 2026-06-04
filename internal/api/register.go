package api

import (
	healthHandler "github.com/HOangAG2207/GoBeK03Echo/internal/handler/health"
	healthRepo "github.com/HOangAG2207/GoBeK03Echo/internal/repository/health"
	healthService "github.com/HOangAG2207/GoBeK03Echo/internal/service/health"
)

type Handlers struct {
	Health healthHandler.Handler
}

func (e *engine) registerHandlers() *Handlers {

	// ===== HEALTH =====
	healthRepo := healthRepo.NewRepository(e.redis)

	healthService := healthService.NewService(
		healthRepo,
		e.config.App.ServiceName,
		e.config.App.InstanceID,
	)

	healthHandler := healthHandler.NewHandler(healthService)

	return &Handlers{
		Health: healthHandler,
	}
}
