package user

import (
	"context"

	"github.com/HOangAG2207/GoBeK03Echo/internal/model"
)

func (s *service) CreateUser(ctx context.Context, displayName, userName, passWord, email string) (*model.User, error) {
	hashedPassword, err := s.passwordHashing.Hash(passWord)
	if err != nil {
		return nil, err
	}
	newUser := &model.User{
		Username:    userName,
		Password:    hashedPassword,
		DisplayName: displayName,
		Email:       email,
	}
	createdUser, err := s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
