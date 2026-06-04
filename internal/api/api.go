package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app    *echo.Echo
	config *Config
}

type EngineOpts struct {
	Cfg *Config
}

func NewEngine(opts *EngineOpts) Engine {
	e := &engine{
		app:    echo.New(),
		config: opts.Cfg,
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
