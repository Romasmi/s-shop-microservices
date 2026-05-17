package delivery

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/domain/courier"
	"github.com/google/uuid"
)

type Repository interface {
	AddCourier(ctx context.Context, c *courier.Courier) error
	GetAvailableCourier(ctx context.Context, from, to int64) (uuid.UUID, error)
	CreateSlot(ctx context.Context, s *courier.Slot) error
	DeleteSlotByOrder(ctx context.Context, orderID uuid.UUID) error
}
