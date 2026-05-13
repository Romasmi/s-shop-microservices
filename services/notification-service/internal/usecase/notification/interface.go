package notification

import (
	"context"
	"github.com/Romasmi/s-shop-microservices/notification-service/internal/domain/message"
)

type Repository interface {
	ListMessages(ctx context.Context, userID string) ([]*message.Message, error)
}
