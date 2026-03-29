package services

import (
	"context"

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

func (s *UserService) CreateUser(ctx context.Context, payload *user.User) (*user.User, error) {
	return s.userRepo.CreateUser(ctx, payload)
}

func (s *UserService) GetUserById(ctx context.Context, userId uuid.UUID) (*user.User, error) {
	return s.userRepo.GetUserById(ctx, userId)
}

func (s *UserService) UpdateUser(ctx context.Context, payload *user.User) (*user.User, error) {
	return s.userRepo.UpdateUser(ctx, payload)
}

func (s *UserService) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	return s.userRepo.DeleteUser(ctx, userId)
}
