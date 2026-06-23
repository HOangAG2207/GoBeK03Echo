package user

import (
	"github.com/HOangAG2207/GoBeK03Echo/internal/service/user"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Register(ctx echo.Context) error
}

type handler struct {
	userService user.Service
}

func NewHandler(userSvc user.Service) Handler {
	return &handler{
		userService: userSvc,
	}
}
