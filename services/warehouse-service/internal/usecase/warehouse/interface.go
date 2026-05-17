package warehouse

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/domain/product"
	"github.com/google/uuid"
)

type Repository interface {
	GetProduct(ctx context.Context, productID uuid.UUID) (*product.Product, error)
	UpsertProduct(ctx context.Context, p *product.Product) error
	CreateReservation(ctx context.Context, r *product.Reservation) error
	GetReservation(ctx context.Context, reservationID uuid.UUID) (*product.Reservation, error)
	DeleteReservation(ctx context.Context, reservationID uuid.UUID) error
	UpdateProductCount(ctx context.Context, productID uuid.UUID, delta int32) error
}
