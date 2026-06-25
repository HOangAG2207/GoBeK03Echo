package pkgjwt

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("Invalid Token")
	ErrExtractToken = errors.New("Unable to Extract Token")
)

//go:generate mockery --name JwtValidator--filename validator.go --output ./mocks
type JwtValidator interface {
	ValidateJWTToken(tokenStr string) (jwt.MapClaims, error)
}

type jwtValidator struct {
	publicKey *rsa.PublicKey
}

func NewValidatorJWT(publicKeyPath string) (JwtValidator, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}
	return &jwtValidator{
		publicKey: publicKey,
	}, nil
}
func (j *jwtValidator) ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return j.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, ErrExtractToken
}
