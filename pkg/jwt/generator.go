package pkgjwt

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockery --name JwtGenerator --filename generator.go --output ./mocks
type JwtGenerator interface {
	GenerateJWTToken(claims jwt.MapClaims) (string, error)
}
type jwtGenerator struct {
	privateKey *rsa.PrivateKey
}

func NewGeneratorJWT(privateKeyPath string) (JwtGenerator, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtGenerator{
		privateKey: privateKey,
	}, nil
}
func (j *jwtGenerator) GenerateJWTToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
