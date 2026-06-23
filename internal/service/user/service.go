package user

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"github.com/HOangAG2207/GoBeK03Echo/internal/repository/user"
	pkgutils "github.com/HOangAG2207/GoBeK03Echo/pkg/utils"
)

//go:generate mockery --name Service --filename user_service_mock.go --output ./mocks

type Service interface {
	CreateUser(ctx context.Context, displayName, userName, passWord, email string) (*model.User, error)
}
type service struct {
	userRepo        user.Repository
	passwordHashing pkgutils.PasswordHashing
}

func NewService(repo user.Repository, passHashing pkgutils.PasswordHashing) Service {
	return &service{
		userRepo:        repo,
		passwordHashing: passHashing,
	}
}
