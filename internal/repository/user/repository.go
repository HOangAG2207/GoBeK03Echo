package user

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
