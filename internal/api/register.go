package api

import (
	healthHandler "github.com/HOangAG2207/GoBeK03Echo/internal/handler/health"
	linksHandler "github.com/HOangAG2207/GoBeK03Echo/internal/handler/links"
	healthRepo "github.com/HOangAG2207/GoBeK03Echo/internal/repository/health"
	linksRepo "github.com/HOangAG2207/GoBeK03Echo/internal/repository/links"
	healthService "github.com/HOangAG2207/GoBeK03Echo/internal/service/health"
	linksService "github.com/HOangAG2207/GoBeK03Echo/internal/service/links"
)

type Handlers struct {
	Health healthHandler.Handler
	Links  linksHandler.Handler
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

	// ===== LINKS =====
	linksRepo := linksRepo.NewRepository(e.redis)
	linksService := linksService.NewService(linksRepo)
	linksHandler := linksHandler.NewHandler(linksService)

	return &Handlers{
		Health: healthHandler,
		Links:  linksHandler,
	}
}
