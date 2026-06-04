package main

import (
	"log"

	"github.com/HOangAG2207/GoBeK03Echo/internal/api"
)

// @title           GoBe K03 Project API
// @version         1.0.0
// @description     REST API cho hệ thống GoBe K03 (health check, URL shortener, ...)
// @description     Sử dụng Echo framework + Redis + Clean Architecture

// @termsOfService  https://example.com/terms

// @host      localhost:8081
// @BasePath
func main() {

	cfg, err := api.NewConfig()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// redisClient, err := pkgredis.NewRedisClient("")
	// if err != nil {
	// 	panic(err)
	// }

	app := api.NewEngine(&api.EngineOpts{
		Cfg: cfg,
	})

	// Start server (chạy HTTP server)
	if err := app.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
