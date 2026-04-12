package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/Romasmi/s-shop-microservices/internal/repository"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo       *repository.UserRepository
	authServiceURL string
}

func CreateUserService(
	userRepo *repository.UserRepository,
	authServiceURL string,
) *UserService {
	return &UserService{
		userRepo:       userRepo,
		authServiceURL: authServiceURL,
	}
}

func (s *UserService) CreateUser(ctx context.Context, payload *user.User) (*user.User, error) {
	// Generate random password because user doesn't provide it yet
	password := gofakeit.Password(true, true, true, true, false, 12)
	payload.Password = password

	newUser, err := s.userRepo.CreateUser(ctx, payload)
	if err != nil {
		return nil, err
	}

	if err := s.registerInAuth(ctx, newUser.ID.String(), newUser.Username, password); err != nil {
		// Ideally we should rollback user creation or mark it as partially created
		// TODO add saga later
		return nil, fmt.Errorf("error registering user in auth service: %w", err)
	}

	newUser.Password = password
	return newUser, nil
}

func (s *UserService) registerInAuth(ctx context.Context, userID, login, password string) error {
	reqBody, _ := json.Marshal(map[string]string{
		"user_id":  userID,
		"login":    login,
		"password": password,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.authServiceURL+"/register", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	return nil
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
