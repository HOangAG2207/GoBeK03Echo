package user

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
)

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}

	err := r.db.WithContext(ctx).Where("username = ? ", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
