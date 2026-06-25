package pkgjwt

import (
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJwtGenerator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		keyPath         string
		expectErrString string
	}{
		{
			name:            "success",
			keyPath:         filepath.FromSlash("./test.private.key"),
			expectErrString: "",
		},
		{
			name:            "file not found",
			keyPath:         filepath.FromSlash("./none.key"),
			expectErrString: "open",
		},
		{
			name:            "invalid private key",
			keyPath:         filepath.FromSlash("./test.public.key"),
			expectErrString: "structure error",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewGeneratorJWT(tc.keyPath)

			if err != nil {
				assert.ErrorContains(t, err, tc.expectErrString)
			}
		})
	}
}

func TestJwtGenerator_GenerateToken(t *testing.T) {
	gen, err := NewGeneratorJWT(filepath.FromSlash("./test.private.key"))
	assert.NoError(t, err)

	claims := map[string]any{
		"sub":  "1234",
		"name": "test",
		"test": true,
		"iat":  1516239022,
	}

	token, err := gen.GenerateJWTToken(claims)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3)

	// ===== HEADER =====
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	assert.NoError(t, err)

	var header map[string]any
	err = json.Unmarshal(headerBytes, &header)
	assert.NoError(t, err)

	assert.Equal(t, "RS256", header["alg"])
	assert.Equal(t, "JWT", header["typ"])

	// ===== PAYLOAD =====
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	assert.NoError(t, err)

	var decoded map[string]any
	err = json.Unmarshal(payloadBytes, &decoded)
	assert.NoError(t, err)

	for k, v := range claims {
		assert.EqualValues(t, v, decoded[k]) // 👈 xử lý int vs float64
	}
}
