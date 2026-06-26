package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/HOangAG2207/GoBeK03Echo/internal/service/user/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_RegisterUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		inputRequest *RegisterUserRequest

		setupRequest func(e *echo.Echo, inputRequest *RegisterUserRequest) echo.Context

		setupMockSvc func(inputRequest *RegisterUserRequest) *mocks.Service

		expectedCode int
	}{
		{
			name: "success",
			inputRequest: &RegisterUserRequest{
				Displayname: "hoang",
				Username:    "hoang01",
				Password:    "12345678",
				Email:       "hoang@gmail.com",
			},
			setupRequest: func(e *echo.Echo, input *RegisterUserRequest) echo.Context {
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				return e.NewContext(req, rec)
			},
			setupMockSvc: func(input *RegisterUserRequest) *mocks.Service {
				mockSvc := new(mocks.Service)

				mockSvc.
					On("CreateUser",
						mock.Anything,
						input.Displayname,
						input.Username,
						input.Password,
						input.Email,
					).
					Return(&model.User{
						Username: input.Username,
						Email:    input.Email,
					}, nil)

				return mockSvc
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "bind error - invalid json",
			inputRequest: nil,
			setupRequest: func(e *echo.Echo, _ *RegisterUserRequest) echo.Context {
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer([]byte("{invalid-json")))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				return e.NewContext(req, rec)
			},
			setupMockSvc: func(_ *RegisterUserRequest) *mocks.Service {
				return new(mocks.Service)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validate error",
			inputRequest: &RegisterUserRequest{
				Username: "hoang01", // thiếu field
			},
			setupRequest: func(e *echo.Echo, input *RegisterUserRequest) echo.Context {
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				return e.NewContext(req, rec)
			},
			setupMockSvc: func(_ *RegisterUserRequest) *mocks.Service {
				return new(mocks.Service)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			inputRequest: &RegisterUserRequest{
				Displayname: "hoang",
				Username:    "hoang01",
				Password:    "12345678",
				Email:       "hoang@gmail.com",
			},
			setupRequest: func(e *echo.Echo, input *RegisterUserRequest) echo.Context {
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				return e.NewContext(req, rec)
			},
			setupMockSvc: func(input *RegisterUserRequest) *mocks.Service {
				mockSvc := new(mocks.Service)

				mockSvc.
					On("CreateUser",
						mock.Anything,
						input.Displayname,
						input.Username,
						input.Password,
						input.Email,
					).
					Return(nil, assert.AnError)

				return mockSvc
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			ctx := tc.setupRequest(e, tc.inputRequest)

			mockSvc := tc.setupMockSvc(tc.inputRequest)

			h := NewHandler(mockSvc)

			err := h.Register(ctx)

			// echo handler thường return nil
			assert.NoError(t, err)

			rec := ctx.Response().Writer.(*httptest.ResponseRecorder)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
