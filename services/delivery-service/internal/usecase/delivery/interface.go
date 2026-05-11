package delivery

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/domain/courier"
)

type Repository interface {
	AddCourier(ctx context.Context, c *courier.Courier) error
	GetAvailableCourier(ctx context.Context, from, to int64) (string, error)
	CreateSlot(ctx context.Context, s *courier.Slot) error
	DeleteSlotByOrder(ctx context.Context, orderID string) error
}
