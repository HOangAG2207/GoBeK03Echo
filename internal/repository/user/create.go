package user

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
)

func (r *repository) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
