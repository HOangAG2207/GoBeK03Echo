package links

import (
	"errors"
	"net/http"

	"github.com/HOangAG2207/GoBeK03Echo/internal/service/links"
	"github.com/HOangAG2207/GoBeK03Echo/internal/utils"
	"github.com/labstack/echo/v4"
)

// RedirectURL godoc
// @Summary      Redirect short URL
// @Description  Redirect to original URL by short code
// @Tags         Links
// @Param        code   path      string  true "Short URL code"
// @Success      302    {string}  string  "Redirect to original URL"
// @Failure      400    {object}  utils.ErrorResponse
// @Failure      404    {object}  utils.ErrorResponse
// @Failure      500    {object}  utils.ErrorResponse
// @Router       /v1/links/redirect/{code} [get]
func (h *handler) RedirectURL(ctx echo.Context) error {
	code := ctx.Param("code")

	// 1. Validate
	if code == "" {
		return utils.Fail(ctx, http.StatusBadRequest, "code is required", nil)
	}

	// 2. Get URL
	url, err := h.service.GetURL(ctx.Request().Context(), code)
	if err != nil {
		if errors.Is(err, links.ErrCodeNotFound) {
			return utils.Fail(ctx, http.StatusNotFound, "URL not found", nil)
		}
		return utils.Fail500(ctx, nil)
	}

	// 3. Redirect
	return ctx.Redirect(http.StatusFound, url)
}
