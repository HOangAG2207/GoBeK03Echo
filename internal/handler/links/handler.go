package links

import (
	"github.com/HOangAG2207/GoBeK03Echo/internal/service/links"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	ShortenLink(ctx echo.Context) error
	RedirectLink(ctx echo.Context) error
}

type handler struct {
	service links.Service
}

func NewHandler(svc links.Service) Handler {
	return &handler{
		service: svc,
	}
}
