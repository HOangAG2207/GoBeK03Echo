package user

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/repository/user"
)

type Service interface {
	CreateUser(ctx context.Context, displayName, userName, passWord, email string) error
}
type service struct {
	repo user.Repository
}

func NewService(repo user.Repository) Service {
	return &service{
		repo: repo,
	}
}
