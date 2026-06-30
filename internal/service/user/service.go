package user

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/HOangAG2207/GoBeK03Echo/internal/repository/user"
	pkgjwt "github.com/HOangAG2207/GoBeK03Echo/pkg/jwt"
	pkgutils "github.com/HOangAG2207/GoBeK03Echo/pkg/utils"
)

//go:generate mockery --name Service --filename user_service_mock.go --output ./mocks
type Service interface {
	CreateUser(ctx context.Context, displayName, userName, passWord, email string) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, error)
}
type service struct {
	userRepo        user.Repository
	passwordHashing pkgutils.PasswordHashing
	jwtGenerator    pkgjwt.JwtGenerator
}

func NewService(repo user.Repository, passHashing pkgutils.PasswordHashing, jwtGen pkgjwt.JwtGenerator) Service {
	return &service{
		userRepo:        repo,
		passwordHashing: passHashing,
		jwtGenerator:    jwtGen,
	}
}
