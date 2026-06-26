package helpers

import (
	"encoding/json"
	"maps"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BaseResponse struct {
	Status string `json:"status" example:"success"`
	Info   string `json:"info"`
}

type ErrorResponse struct {
	Status string `json:"status" example:"false"`
	Info   string `json:"info"`
	Error  any    `json:"error,omitempty"`
}

// SUCCESS KHÔNG CÓ DATA WRAPPER
func SuccessFlat(c echo.Context, status int, info string, data any) error {
	resp := map[string]any{
		"status": "success",
		"info":   info,
	}

	// merge struct vào root
	if m, ok := structToMap(data); ok {
		maps.Copy(resp, m)
	}

	return c.JSON(status, resp)
}
func SuccessWrapData(c echo.Context, status int, info string, data any) error {
	resp := map[string]any{
		"status": "success",
		"info":   info,
		"data":   data,
	}

	return c.JSON(status, resp)
}

func Fail(c echo.Context, status int, info string, err any) error {
	return c.JSON(status, ErrorResponse{
		Status: "error",
		Info:   info,
		Error:  err,
	})
}

func Fail500(c echo.Context, err error) error {
	var msg any = "Internal Server Error"

	if err != nil {
		msg = err.Error()
	}

	return Fail(c, http.StatusInternalServerError, "Internal Server Error", msg)
}
func structToMap(data any) (map[string]any, bool) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, false
	}

	var result map[string]any
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, false
	}

	return result, true
}
