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

		expectedOutput *model.User
		expectedError  bool

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
				Displayname: "New User 01",
				Username:    "New User 01",
				Password:    "$2a$10$xxx",
				Email:       "newuser01@example.com",
			},
			expectedOutput: &model.User{
				Base: model.Base{
					ID: "de305d54-75b4-431b-adb2-eb6b9e546090",
				},
				Displayname: "New User 01",
				Username:    "New User 01",
				Email:       "newuser01@example.com",
				// ⚠️ Password thường không assert trực tiếp (hash có thể khác)
			},
			expectedError: false,
			verifyFunc: func(db *gorm.DB, user *model.User) {
				checkUser := &model.User{}

				err := db.Where("id = ?", user.ID).First(checkUser).Error
				assert.NoError(t, err)

				assert.Equal(t, user.ID, checkUser.ID)
				assert.Equal(t, user.Username, checkUser.Username)
				assert.Equal(t, user.Email, checkUser.Email)
				assert.Equal(t, user.Displayname, checkUser.Displayname)
			},
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
				Displayname: "New User 01",
				Username:    "New User 02",
				Password:    "$2a$10$xxx",
				Email:       "hoang01@gmail.com",
			},
			expectedOutput: nil,
			expectedError:  true,
		},
		{
			name: "Create user fail - duplicate username",
			setupMockDB: func(t *testing.T) *gorm.DB {
				return fixtures_test.NewFixtureDB(t, &fixtures_test.UserCommonTestDB{})
			},
			inputUser: &model.User{
				Base: model.Base{
					ID: "de305d54-75b4-431b-adb2-eb6b9e546091",
				},
				Displayname: "Another User",
				Username:    "hoang01",
				Password:    "$2a$10$xxx",
				Email:       "another@example.com",
			},
			expectedOutput: nil,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		tc := tc // fix race khi dùng t.Parallel

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockDB := tc.setupMockDB(t)
			repo := NewRepository(mockDB)

			result, err := repo.CreateUser(ctx, tc.inputUser)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)

			// ✅ So sánh output
			assert.Equal(t, tc.expectedOutput.ID, result.ID)
			assert.Equal(t, tc.expectedOutput.Username, result.Username)
			assert.Equal(t, tc.expectedOutput.Email, result.Email)
			assert.Equal(t, tc.expectedOutput.Displayname, result.Displayname)

			// ⚠️ Password: không nên assert equality nếu có hash
			// assert.Equal(t, tc.expectedOutput.Password, result.Password)

			if tc.verifyFunc != nil {
				tc.verifyFunc(mockDB, result)
			}
		})
	}
}
