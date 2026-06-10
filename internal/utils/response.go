package utils

import (
	"encoding/json"
	"maps"
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
	// return c.JSON(status, SuccessResponse{
	// 	Status:  "success",
	// 	Message: message,
	// 	Data:    data,
	// })
	resp := map[string]any{
		"status":  "success",
		"message": message,
	}

	if data != nil {
		if m, ok := toMap(data); ok {
			maps.Copy(resp, m)
		}
	}

	return c.JSON(status, resp)
}

func Fail(c echo.Context, status int, message string, err any) error {
	return c.JSON(status, ErrorResponse{
		// Status:  "fail",
		Message: message,
		Error:   err,
	})
}
func Fail500(c echo.Context, err error) error {
	return Fail(c, http.StatusInternalServerError, "Internal Server Error", err)
}
func toMap(data any) (map[string]any, bool) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, false
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, false
	}

	return m, true
}
