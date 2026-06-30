package user

import (
	"context"
	"testing"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	repoMocks "github.com/HOangAG2207/GoBeK03Echo/internal/repository/user/mocks"
	jwtMocks "github.com/HOangAG2207/GoBeK03Echo/pkg/jwt/mocks"
	pkgutils "github.com/HOangAG2207/GoBeK03Echo/pkg/utils"
	passHashingMocks "github.com/HOangAG2207/GoBeK03Echo/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockPasswordHashing func(t *testing.T) *passHashingMocks.PasswordHashing
		setupMockRepo            func(ctx context.Context) *repoMocks.Repository
		setupMockJwt             func(t *testing.T) *jwtMocks.JwtGenerator

		inputUsername    string
		inputPassword    string
		inputDisplayName string
		inputEmail       string

		expectedOutput *model.User
		expectedError  error
	}{
		{
			name: "Create user successfully",

			setupMockPasswordHashing: func(t *testing.T) *passHashingMocks.PasswordHashing {
				hashingMock := passHashingMocks.NewPasswordHashing(t)
				hashingMock.On("Hash", "password123").
					Return("$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", nil)
				return hashingMock
			},

			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				repoMock := repoMocks.NewRepository(t)
				repoMock.On("CreateUser", ctx, &model.User{
					Username:    "testuser",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					Displayname: "Test User",
					Email:       "testuser@example.com",
				}).Return(&model.User{
					Base: model.Base{
						ID: "de305d54-75b4-431b-adb2-eb6b9e546099",
					},
					Username:    "testuser",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					Displayname: "Test User",
					Email:       "testuser@example.com",
				}, nil)
				return repoMock
			},

			setupMockJwt: func(t *testing.T) *jwtMocks.JwtGenerator {
				jwtMock := jwtMocks.NewJwtGenerator(t)
				jwtMock.On("GenerateToken", mock.Anything).
					Maybe().
					Return("dummy-token", nil)
				return jwtMock
			},

			inputUsername:    "testuser",
			inputPassword:    "password123",
			inputDisplayName: "Test User",
			inputEmail:       "testuser@example.com",

			expectedOutput: &model.User{
				Base: model.Base{
					ID: "de305d54-75b4-431b-adb2-eb6b9e546099",
				},
				Username:    "testuser",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Displayname: "Test User",
				Email:       "testuser@example.com",
			},
		},

		{
			name: "Fail to hash password",

			setupMockPasswordHashing: func(t *testing.T) *passHashingMocks.PasswordHashing {
				hashingMock := passHashingMocks.NewPasswordHashing(t)
				hashingMock.On("Hash", "badpassword").
					Return("", pkgutils.ErrCannotGenerateHash)
				return hashingMock
			},

			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				return repoMocks.NewRepository(t)
			},

			setupMockJwt: func(t *testing.T) *jwtMocks.JwtGenerator {
				return jwtMocks.NewJwtGenerator(t)
			},

			inputUsername:    "testuser2",
			inputPassword:    "badpassword",
			inputDisplayName: "Test User 2",
			inputEmail:       "testuser2@example.com",

			expectedError: pkgutils.ErrCannotGenerateHash,
		},

		{
			name: "Fail to create user in repository",

			setupMockPasswordHashing: func(t *testing.T) *passHashingMocks.PasswordHashing {
				hashingMock := passHashingMocks.NewPasswordHashing(t)
				hashingMock.On("Hash", "password123").
					Return("$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e", nil)
				return hashingMock
			},

			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				repoMock := repoMocks.NewRepository(t)
				repoMock.On("CreateUser", ctx, &model.User{
					Username:    "testuser3",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					Displayname: "Test User 3",
					Email:       "testuser3@example.com",
				}).Return(nil, assert.AnError)
				return repoMock
			},

			setupMockJwt: func(t *testing.T) *jwtMocks.JwtGenerator {
				return jwtMocks.NewJwtGenerator(t)
			},

			inputUsername:    "testuser3",
			inputPassword:    "password123",
			inputDisplayName: "Test User 3",
			inputEmail:       "testuser3@example.com",

			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			passwordHashingMock := tc.setupMockPasswordHashing(t)
			userRepoMock := tc.setupMockRepo(ctx)
			jwtMock := tc.setupMockJwt(t)

			userService := NewService(
				userRepoMock,
				passwordHashingMock,
				jwtMock,
			)

			res, err := userService.CreateUser(
				ctx,
				tc.inputDisplayName,
				tc.inputUsername,
				tc.inputPassword,
				tc.inputEmail,
			)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
