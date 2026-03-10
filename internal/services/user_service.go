package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/Romasmi/s-shop-microservices/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func CreateUserService(
	userRepo *repository.UserRepository,
) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(ctx context.Context, payload *user.User) (*user.User, error) {
	userId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
