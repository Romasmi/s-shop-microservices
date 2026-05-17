package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/domain/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) GetProduct(ctx context.Context, productID uuid.UUID) (*product.Product, error) {
	p := &product.Product{}
	err := r.pool.QueryRow(ctx, "SELECT id, count FROM products WHERE id = $1", productID).Scan(&p.ID, &p.Count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return p, nil
}

func (r *ProductRepository) UpsertProduct(ctx context.Context, p *product.Product) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO products (id, count) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET count = EXCLUDED.count", p.ID, p.Count)
	if err != nil {
		return fmt.Errorf("failed to upsert product: %w", err)
	}
	return nil
}

func (r *ProductRepository) CreateReservation(ctx context.Context, res *product.Reservation) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO reservations (id, product_id, count) VALUES ($1, $2, $3)", res.ID, res.ProductID, res.Count)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetReservation(ctx context.Context, reservationID uuid.UUID) (*product.Reservation, error) {
	res := &product.Reservation{}
	err := r.pool.QueryRow(ctx, "SELECT id, product_id, count FROM reservations WHERE id = $1", reservationID).Scan(&res.ID, &res.ProductID, &res.Count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("reservation not found")
		}
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	return res, nil
}

func (r *ProductRepository) DeleteReservation(ctx context.Context, reservationID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM reservations WHERE id = $1", reservationID)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}
	return nil
}

func (r *ProductRepository) UpdateProductCount(ctx context.Context, productID uuid.UUID, delta int32) error {
	_, err := r.pool.Exec(ctx, "UPDATE products SET count = count + $1 WHERE id = $2", delta, productID)
	if err != nil {
		return fmt.Errorf("failed to update product count: %w", err)
	}
	return nil
}
