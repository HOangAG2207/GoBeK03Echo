package user

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TokenExpirationDuration = 1 * time.Hour

var ErrInvalidCredential = errors.New("Invalid Credentials")

type TokenInfo struct {
	ID       string
	Username string
}

func (t *TokenInfo) ToMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"sub":      t.ID,
		"username": t.Username,
		"exp":      time.Now().Add(TokenExpirationDuration).Unix(),
		"iat":      time.Now().Unix(),
	}
}

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	// get user to username
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", ErrInvalidCredential
	}
	// check password -> compare hash password
	check := s.passwordHashing.CompareHash(user.Password, password)
	if !check {
		return "", ErrInvalidCredential
	}

	// generate token
	tokenInfo := &TokenInfo{
		ID:       user.ID,
		Username: user.Username,
	}
	token, err := s.jwtGenerator.GenerateJWTToken(tokenInfo.ToMapClaims())
	if err != nil {
		return "", err
	}

	// return token
	return token, nil
}
