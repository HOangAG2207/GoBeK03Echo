package user

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"gorm.io/gorm"
)

//go:generate mockery --name Repository --filename links_repository_mock.go --output ./mocks
type Repository interface {
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
