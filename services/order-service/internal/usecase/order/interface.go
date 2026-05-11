package order

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/order-service/internal/domain/order"
	api "github.com/Romasmi/s-shop/gen/go/order"
	"github.com/google/uuid"
)

type Repository interface {
	CreateOrder(ctx context.Context, o *order.Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*order.Order, error)
}

type EventProducer interface {
	EmitOrderPlaced(ctx context.Context, event *api.OrderPlaced) error
}
