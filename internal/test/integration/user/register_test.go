package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HOangAG2207/GoBeK03Echo/internal/api"
	fixtures_test "github.com/HOangAG2207/GoBeK03Echo/internal/test/fixtures"
	pkgutils "github.com/HOangAG2207/GoBeK03Echo/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIntegration_RegisterUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB

		setupTestHTTP func(e api.Engine) *httptest.ResponseRecorder

		expectedStatusCode      int
		expectedMessageResponse string
	}{
		{
			name: "success",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures_test.NewFixtureDB(t, &fixtures_test.UserCommonTestDB{})
			},

			setupTestHTTP: func(e api.Engine) *httptest.ResponseRecorder {
				body := `{
					"username":"newuser",
					"password":"123456Aa@",
					"display_name":"New User",
					"email":"newuser@example.com"
				}`

				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				return rec
			},

			expectedStatusCode:      http.StatusOK,
			expectedMessageResponse: "Register user successfully!",
		},

		{
			name: "duplicate username",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures_test.NewFixtureDB(t, &fixtures_test.UserCommonTestDB{})
			},

			setupTestHTTP: func(e api.Engine) *httptest.ResponseRecorder {
				body := `{
					"username":"hoang01",
					"password":"123456Aa@",
					"display_name":"Dup",
					"email":"dup@example.com"
				}`

				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				return rec
			},

			expectedStatusCode:      http.StatusInternalServerError,
			expectedMessageResponse: "Internal Server Error",
		},

		{
			name: "invalid json",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures_test.NewFixtureDB(t, &fixtures_test.UserCommonTestDB{})
			},

			setupTestHTTP: func(e api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(`{"username":`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				return rec
			},

			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: "invalid request body",
		},

		{
			name: "validation error",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures_test.NewFixtureDB(t, &fixtures_test.UserCommonTestDB{})
			},

			setupTestHTTP: func(e api.Engine) *httptest.ResponseRecorder {
				body := `{
					"username":"",
					"password":"123",
					"display_name":"Test",
					"email":"invalid"
				}`

				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", strings.NewReader(body))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				return rec
			},

			expectedStatusCode:      http.StatusBadRequest,
			expectedMessageResponse: "validation error",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// ===== DB =====
			db := tc.setupDB(t)

			// ===== ENGINE =====
			e := api.NewEngine(&api.EngineOpts{
				Cfg: &api.Config{
					Port: "8080",
				},
				DB:            db,
				Redis:         nil, // ❌ bỏ Redis đúng yêu cầu
				RandomCodeGen: nil,
				PassHashing:   pkgutils.NewPasswordHashing(),
			})

			// ===== CALL API =====
			rec := tc.setupTestHTTP(e)

			// ===== ASSERT =====
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedMessageResponse)
		})
	}
}
