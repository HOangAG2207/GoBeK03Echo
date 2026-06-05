package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

func Success(c echo.Context, status int, message string, data any) error {
	return c.JSON(status, SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Fail(c echo.Context, status int, message string, err any) error {
	return c.JSON(status, ErrorResponse{
		Status:  "fail",
		Message: message,
		Error:   err,
	})
}
func Fail500(c echo.Context, err error) error {
	return Fail(c, http.StatusInternalServerError, "Internal Server Error", err)
}
