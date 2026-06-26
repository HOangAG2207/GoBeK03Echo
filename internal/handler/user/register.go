package user

import (
	"net/http"

	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers"
	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers/validator"
	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type RegisterUserRequest struct {
	Username    string `json:"username" validate:"required,gt=0"`
	Displayname string `json:"display_name" validate:"required,gt=0"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,gte=8"`
}
type RegisterUserSwaggerResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    model.User `json:"data"`
}

// New User generate.
// @Summary      Create a new user
// @Description  Create a new user with the provided information
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      RegisterUserRequest  true  "User to create"
// @Success      201   {object}  RegisterUserSwaggerResponse
// @Failure      400   {object}  helpers.ErrorResponse
// @Failure      500   {object}  helpers.ErrorResponse
// @Router       /v1/users/register [post]
func (h *handler) Register(ctx echo.Context) error {
	req := &RegisterUserRequest{}

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
	return helpers.SuccessWrapData(
		ctx,
		http.StatusOK,
		"Register user successfully!",
		user,
	)
}
