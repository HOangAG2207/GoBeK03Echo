package health

import (
	"net/http"

	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers"
	"github.com/labstack/echo/v4"
)

// CheckHealth godoc
// @Summary      Check service health
// @Description  Kiểm tra trạng thái hoạt động của service và kết nối Redis
// @Tags         Health
// @Produce      json
// @Success      200 {object} model.HealthCheckSwaggerResponse
// @Failure      500 {object} helpers.ErrorResponse
// @Router       /v1/health-check [get]
func (h *handler) CheckHealth(c echo.Context) error {
	response, err := h.service.CheckHealth(c.Request().Context())
	if err != nil {
		return helpers.Fail500(c, err)
	}

	return helpers.SuccessWithMode(c, http.StatusOK, "Health check passed", response, helpers.ModeFlat)
}
