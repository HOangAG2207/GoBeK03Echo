package user

import (
	"net/http"

	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers"
	"github.com/HOangAG2207/GoBeK03Echo/internal/helpers/validator"
	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/labstack/echo/v4"
)

func (h *handler) Login(ctx echo.Context) error {
	req := &model.LoginRequest{}

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

	// 6. Response
	return helpers.Success(
		ctx,
		http.StatusOK,
		"Login successfully!",
		nil, // token,
	)
}
