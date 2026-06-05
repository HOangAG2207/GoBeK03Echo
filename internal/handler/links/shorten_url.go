package links

import (
	"net/http"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/HOangAG2207/GoBeK03Echo/internal/utils"
	"github.com/labstack/echo/v4"
)

// ShortenURL godoc
// @Summary      Shorten URL
// @Description  Generate a short code for a given URL
// @Tags         Links
// @Accept       json
// @Produce      json
// @Param        request  body      model.ShortenURLRequest  true  "Shorten URL request"
// @Success      200  {object}  utils.SuccessResponse{data=model.ShortenURLResponse}
// @Failure      400  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /v1/links/shorten [post]
func (h *handler) ShortenURL(ctx echo.Context) error {
	req := new(model.ShortenURLRequest)

	// 1. Bind
	if err := ctx.Bind(req); err != nil {
		return utils.Fail(ctx, http.StatusBadRequest, "invalid request body", nil)
	}

	// 2. Validate URL
	if req.URL == "" {
		return utils.Fail(ctx, http.StatusBadRequest, "url is required", nil)
	}

	// 3. Validate exp
	if req.Exp <= 0 {
		return utils.Fail(ctx, http.StatusBadRequest, "exp must be > 0", nil)
	}

	// 5. Call service
	code, err := h.service.ShortenURL(
		ctx.Request().Context(),
		req.URL,
		7,
		req.Exp,
	)
	if err != nil {
		return utils.Fail500(ctx, nil)
	}

	// 6. Response
	return utils.Success(ctx, http.StatusOK, "Shorten URL generated successfully!", model.ShortenURLResponse{
		Code: code,
	})
}
