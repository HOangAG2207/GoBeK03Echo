package user

import (
	"net/http"

	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers"
	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers/validator"
	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// New User generate.
// @Summary      Create a new user
// @Description  Create a new user with the provided information
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      model.RegisterUserRequest  true  "User to create"
// @Success      201   {object}  model.RegisterUserSwaggerResponse
// @Failure      400   {object}  helpers.ErrorResponse
// @Failure      500   {object}  helpers.ErrorResponse
// @Router       /v1/users/register [post]
func (h *handler) Register(ctx echo.Context) error {
	req := &model.RegisterUserRequest{}

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
	// 3. Call service
	user, err := h.userService.CreateUser(ctx.Request().Context(), req.Displayname, req.Username, req.Password, req.Email)
	if err != nil {
		log.Error().Err(err).Str("From", "links.handler.Register").Msg("Fail to register user")
		return helpers.Fail500(ctx, nil)
	}
	// 6. Response
	return helpers.Success(
		ctx,
		http.StatusOK,
		"Register user successfully!",
		user,
	)
}
