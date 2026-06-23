package links

import (
	"net/http"

	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers"
	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers/validator"
	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// ShortenURL godoc
// @Summary      Shorten URL
// @Description  Generate a short code for a given URL
// @Tags         Links
// @Accept       json
// @Produce      json
// @Param		 request body model.ShortenURLRequest true "Shorten URL request (exp must be > 0)"
// @Success      200  {object}  model.ShortenURLSwaggerResponse
// @Failure      400  {object}  helpers.ErrorResponse
// @Failure      500  {object}  helpers.ErrorResponse
// @Router       /v1/links/shorten [post]
func (h *handler) ShortenLink(ctx echo.Context) error {
	req := new(model.ShortenURLRequest)

	// 1. Bind
	if err := ctx.Bind(req); err != nil {
		return helpers.Fail(ctx, http.StatusBadRequest, "invalid request body", nil)
	}
	// 2. Validate URL bằng validator
	if err := validator.ValidateStruct(req); err != nil {
		return helpers.Fail(
			ctx,
			http.StatusBadRequest,
			"validation error",
			validator.FormatValidationErrors(err),
		)
	}

	// 5. Call service
	code, err := h.service.ShortenLink(
		ctx.Request().Context(),
		req.URL,
		7,
		req.Exp,
	)
	if err != nil {
		log.Error().Err(err).Str("From", "links.handler.ShortenLink").Msg("Fail to short url")
		return helpers.Fail500(ctx, nil)
	}

	// 6. Response
	return helpers.Success(
		ctx,
		http.StatusOK,
		"Shorten URL generated successfully!",
		model.ShortenURLResponse{
			Code: code,
		},
	)
}
