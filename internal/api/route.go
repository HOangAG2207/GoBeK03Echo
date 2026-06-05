package api

import (
	_ "github.com/HOangAG2207/GoBeK03Echo/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (e *engine) initRoutes() {

	h := e.registerHandlers()

	api := e.app.Group("/v1")

	// Swagger
	api.GET("/docs/*", echoSwagger.WrapHandler)

	// Health
	api.GET("/health-check", h.Health.CheckHealth)

	// Links
	// - Shorten link
	api.POST("/links/shorten", h.Links.ShortenURL)
	// - Redirect link from code
	api.GET("/links/redirect/:code", h.Links.RedirectURL)
}
