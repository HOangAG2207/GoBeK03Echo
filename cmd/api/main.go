package main

import (
	"log"

	"github.com/HOangAG2207/GoBeK03Echo/internal/infrastructure"
)

// @title           GoBe K03 Project API
// @version         1.0.3
// @description     REST API cho hệ thống GoBe K03 (health check, URL shortener, ...)
// @description     Sử dụng Echo framework + Redis + Clean Architecture

// @termsOfService  https://example.com/terms

// @host      localhost:8080
// @BasePath /
// @schemes http

// Tags phân loại API (sử dụng trong Swagger UI)
// @tag.name        Health
// @tag.name        Links
// @tag.name		User
func main() {
	app := infrastructure.CreateAPI()

	// Start server (chạy HTTP server)
	if err := app.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
