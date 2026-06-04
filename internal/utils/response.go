package utils

import (
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
	Error   string `json:"error,omitempty"`
}

func Success(c echo.Context, status int, message string, data any) error {
	return c.JSON(status, SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Fail(c echo.Context, status int, message string, err error) error {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	return c.JSON(status, ErrorResponse{
		Status:  "fail",
		Message: message,
		Error:   errMsg,
	})
}
func Fail500(c echo.Context, err error) error {
	return Fail(c, 500, "Internal Server Error", err)
}
