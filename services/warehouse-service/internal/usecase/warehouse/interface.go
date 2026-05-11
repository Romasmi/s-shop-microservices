package warehouse

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/domain/product"
)

type Repository interface {
	GetProduct(ctx context.Context, productID string) (*product.Product, error)
	UpsertProduct(ctx context.Context, p *product.Product) error
	CreateReservation(ctx context.Context, r *product.Reservation) error
	GetReservation(ctx context.Context, reservationID string) (*product.Reservation, error)
	DeleteReservation(ctx context.Context, reservationID string) error
	UpdateProductCount(ctx context.Context, productID string, delta int32) error
}
