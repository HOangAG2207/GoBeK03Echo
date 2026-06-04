package health

import (
	"github.com/HOangAG2207/GoBeK03Echo/internal/service/health"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	CheckHealth(c echo.Context) error
}
type handler struct {
	service health.Service
}

func NewHandler(s health.Service) Handler {
	return &handler{
		service: s,
	}
}
