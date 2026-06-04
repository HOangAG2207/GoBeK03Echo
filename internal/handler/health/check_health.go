package health

import (
	"github.com/HOangAG2207/GoBeK03Echo/internal/utils"
	"github.com/labstack/echo/v4"
)

// CheckHealth godoc
// @Summary      Check service health
// @Description  Kiểm tra trạng thái hoạt động của service và kết nối Redis
// @Tags         Health
// @Produce      json
// @Success      200 {object} utils.SuccessResponse{data=model.HealthCheckResponse}
// @Failure      500 {object} utils.ErrorResponse
// @Router       /v1/health-check [get]
func (h *handler) CheckHealth(c echo.Context) error {
	response, err := h.service.CheckHealth(c.Request().Context())
	if err != nil {
		return utils.Fail500(c, err)
	}
	return utils.Success(c, 200, "Health check passed", response)
}
