package api

type Handlers struct {
	// Health handler.HealthHandler
}

func (e *engine) registerHandlers() *Handlers {

	// ===== HEALTH =====
	// healthRepo := repository.NewHealthRepository(e.redisClient)

	// healthService := service.NewHealthService(
	// 	healthRepo,
	// 	e.config.ServiceName,
	// 	e.config.InstanceID,
	// )

	// healthHandler := handler.NewHealthHandler(healthService)

	return &Handlers{
		// Health: healthHandler,
	}
}
