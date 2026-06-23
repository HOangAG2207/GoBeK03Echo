package user

import (
	"testing"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	fixtures_test "github.com/HOangAG2207/GoBeK03Echo/internal/test/fixtures"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRepository_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockDB func(t *testing.T) *gorm.DB
		inputUser   *model.User

		expectedError bool

		verifyFunc func(db *gorm.DB, user *model.User)
	}{
		{
			name: "Create user successfully",
			setupMockDB: func(t *testing.T) *gorm.DB {
				return fixtures_test.NewFixtureDB(t, &fixtures_test.UserCommonTestDB{})
			},
			inputUser: &model.User{
				Base: model.Base{
					ID: "de305d54-75b4-431b-adb2-eb6b9e546090",
				},
				DisplayName: "New User 01",
				Username:    "New User 01",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Email:       "newuser01@example.com",
			},
			verifyFunc: func(db *gorm.DB, user *model.User) {
				checkUser := &model.User{}

				err := db.Where("id = ?", user.ID).First(checkUser).Error
				assert.NoError(t, err)

				assert.Equal(t, user.ID, checkUser.ID)
				assert.Equal(t, user.Username, checkUser.Username)
				assert.Equal(t, user.Email, checkUser.Email)
				assert.Equal(t, user.DisplayName, checkUser.DisplayName)
			},
			expectedError: false,
		},
		{
			name: "Create user fail - duplicate email",
			setupMockDB: func(t *testing.T) *gorm.DB {
				return fixtures_test.NewFixtureDB(t, &fixtures_test.UserCommonTestDB{})
			},
			inputUser: &model.User{
				Base: model.Base{
					ID: "de305d54-75b4-431b-adb2-eb6b9e546090",
				},
				DisplayName: "New User 01",
				Username:    "New User 02",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				Email:       "hoang01@gmail.com",
			},
			verifyFunc: func(db *gorm.DB, user *model.User) {

			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockDB := tc.setupMockDB(t)

			repo := NewRepository(mockDB)

			result, err := repo.CreateUser(ctx, tc.inputUser)
			if tc.expectedError {
				assert.NotNil(t, err)
				return
			}

			assert.NoError(t, err)

			tc.verifyFunc(mockDB, result)
		})
	}
}
