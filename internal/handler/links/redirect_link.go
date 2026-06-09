package links

import (
	"errors"
	"net/http"

	"github.com/HOangAG2207/GoBeK03Echo/internal/service/links"
	"github.com/HOangAG2207/GoBeK03Echo/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// RedirectURL godoc
// @Summary      Redirect short URL
// @Description  Redirect to original URL by short code
// @Tags         Links
// @Param        code   path      string  true  "Short URL code"
// @Success      302    {string}  string  "Found - Redirect via Location header"
// @Header       302    {string}  Location  "Original URL"
// @Failure      404    {object}  utils.ErrorResponse "URL not found"
// @Failure      500    {object}  utils.ErrorResponse "Internal server error"
// @Router       /v1/links/redirect/{code} [get]
func (h *handler) RedirectLink(ctx echo.Context) error {
	code := ctx.Param("code")

	url, err := h.service.GetLink(ctx.Request().Context(), code)
	if err != nil {
		if errors.Is(err, links.ErrCodeNotFound) {
			return utils.Fail(ctx, http.StatusNotFound, "URL not found", nil)
		}
		log.Error().Err(err).Str("From", "links.handler.RedirectLink").Msg("Fail to get url from code")
		return utils.Fail500(ctx, err)
	}

	return ctx.Redirect(http.StatusFound, url)
}
