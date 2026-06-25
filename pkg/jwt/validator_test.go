package pkgjwt

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJwtValidator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		keyPath string

		expectedErrString string
	}{
		{
			name:              "success",
			keyPath:           filepath.FromSlash("./test.public.key"),
			expectedErrString: "",
		},
		{
			name:              "file not found",
			keyPath:           filepath.FromSlash("./none.key"),
			expectedErrString: "open",
		},
		{
			name:              "file not found",
			keyPath:           filepath.FromSlash("./test.private.key"),
			expectedErrString: "structure error",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewValidatorJWT(tc.keyPath)
			if err != nil {
				assert.ErrorContains(t, err, tc.expectedErrString)
			}
		})
	}
}
func TestJwtValidator_ValidateToken(t *testing.T) {
	// init generator (tạo token)
	gen, err := NewGeneratorJWT(filepath.FromSlash("./test.private.key"))
	assert.NoError(t, err)

	// init validator
	validator, err := NewValidatorJWT(filepath.FromSlash("./test.public.key"))
	assert.NoError(t, err)

	claims := map[string]any{
		"sub":  "1234",
		"name": "test",
	}

	// ===== SUCCESS =====
	token, err := gen.GenerateJWTToken(claims)
	assert.NoError(t, err)

	parsedClaims, err := validator.ValidateJWTToken(token)
	assert.NoError(t, err)

	for k, v := range claims {
		assert.EqualValues(t, v, parsedClaims[k])
	}

	// ===== INVALID TOKEN STRING =====
	_, err = validator.ValidateJWTToken("invalid.token.here")
	assert.ErrorIs(t, err, ErrInvalidToken)

	// ===== TAMPERED TOKEN =====
	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3)

	// sửa payload → làm sai signature
	parts[1] = parts[1] + "tamper"
	tamperedToken := strings.Join(parts, ".")

	_, err = validator.ValidateJWTToken(tamperedToken)
	assert.ErrorIs(t, err, ErrInvalidToken)
}
