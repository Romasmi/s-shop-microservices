package notification

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/notification-service/internal/domain/message"
	"github.com/google/uuid"
)

type Repository interface {
	ListMessages(ctx context.Context, userID uuid.UUID) ([]*message.Message, error)
}
