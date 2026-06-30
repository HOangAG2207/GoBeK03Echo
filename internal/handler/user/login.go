package user

import (
	"errors"
	"net/http"

	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers"
	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers/validator"
	"github.com/HOangAG2207/GoBeK03Echo/internal/service/user"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Pasword  string `json:"password" validate:"required"`
}
type LoginSwaggerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// Login. 		 Authentication
// @Summary      Return JWT Token if Correct Incredentials
// @Description  Return JWT Token if Correct Incredentials
// @Tags         User
// @Accept       application/json
// @Produce      application/json
// @Param        user  body      LoginRequest  true  "Input Required"
// @Success      200   {object}  LoginSwaggerResponse
// @Failure      400   {object}  helpers.ErrorResponse
// @Failure      500   {object}  helpers.ErrorResponse
// @Router       /v1/users/login [post]
func (h *handler) Login(ctx echo.Context) error {
	req := &LoginRequest{}

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
	//gọi service
	tokenStr, err := h.userService.Login(ctx.Request().Context(), req.Username, req.Pasword)
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredential) {
			return helpers.Fail(
				ctx,
				http.StatusBadRequest,
				"Invalid username or password",
				err.Error(),
			)
		}
		log.Error().Err(err).Str("From", "links.handler.Login").Msg("Fail to Login")
		return helpers.Fail500(ctx, nil)
	}

	// 6. Response
	return helpers.SuccessWithMode(
		ctx,
		http.StatusOK,
		"Logged in successfully!",
		tokenStr, // token,
		helpers.ModeWrap,
	)
}
