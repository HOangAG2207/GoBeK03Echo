package pkgutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing_Hash(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		inputPassword string

		expectedError error
	}{
		{
			name:          "Hash successfully",
			inputPassword: "123456",
			expectedError: nil,
		},
		{
			name:          "Hash empty password",
			inputPassword: "",
			expectedError: nil, // bcrypt vẫn hash được chuỗi rỗng
		},
	}

	hasher := NewPasswordHashing()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hashedPassword, err := hasher.Hash(tc.inputPassword)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, hashedPassword)

			// verify hash đúng với password
			ok := hasher.CompareHash(hashedPassword, tc.inputPassword)
			assert.True(t, ok)
		})
	}
}
func TestPasswordHashing_CompareHash(t *testing.T) {
	t.Parallel()

	hasher := NewPasswordHashing()

	hashedPassword, err := hasher.Hash("123456")
	assert.NoError(t, err)

	testCases := []struct {
		name string

		hashedPassword string
		inputPassword  string

		expectedResult bool
	}{
		{
			name:           "Compare success",
			hashedPassword: hashedPassword,
			inputPassword:  "123456",
			expectedResult: true,
		},
		{
			name:           "Compare wrong password",
			hashedPassword: hashedPassword,
			inputPassword:  "wrong",
			expectedResult: false,
		},
		{
			name:           "Compare invalid hash",
			hashedPassword: "invalid-hash",
			inputPassword:  "123456",
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := hasher.CompareHash(tc.hashedPassword, tc.inputPassword)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
