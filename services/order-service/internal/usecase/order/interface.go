package order

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/order-service/internal/domain/order"
	api "github.com/Romasmi/s-shop-microservices/order-service/pkg/api"
)

type Repository interface {
	CreateOrder(ctx context.Context, o *order.Order) error
	GetOrder(ctx context.Context, id string) (*order.Order, error)
}

type EventProducer interface {
	EmitOrderPlaced(ctx context.Context, event *api.OrderPlaced) error
}
