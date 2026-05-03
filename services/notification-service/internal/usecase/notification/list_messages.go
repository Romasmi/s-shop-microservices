package notification

import (
	"context"
	"github.com/Romasmi/s-shop-microservices/notification-service/internal/domain/message"
)

type ListMessagesUseCase struct {
	repo Repository
}

func NewListMessagesUseCase(repo Repository) *ListMessagesUseCase {
	return &ListMessagesUseCase{repo: repo}
}

func (uc *ListMessagesUseCase) Do(ctx context.Context, userID string) ([]*message.Message, error) {
	return uc.repo.ListMessages(ctx, userID)
}
