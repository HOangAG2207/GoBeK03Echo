package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app    *echo.Echo
	config *Config
	redis  *redis.Client
}

type EngineOpts struct {
	Cfg   *Config
	Redis *redis.Client
}

func NewEngine(opts *EngineOpts) Engine {
	e := &engine{
		app:    echo.New(),
		config: opts.Cfg,
		redis:  opts.Redis,
	}

	e.initMiddleware()

	// tách ra
	e.initRoutes()

	return e
}

func (e *engine) initMiddleware() {
	e.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.app.Use(middleware.RequestLogger())
	e.app.Use(middleware.Recover())
}

func (e *engine) Start() error {
	port := e.config.App.Port
	if port[0] != ':' {
		port = fmt.Sprintf(":%s", port)
	}

	log.Printf("Server running at %s\n", port)
	return e.app.Start(port)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}
