package postgres

import (
	"context"
	"fmt"

	"github.com/Romasmi/s-shop-microservices/order-service/internal/domain/order"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, o *order.Order) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO orders (id, user_id, price, status, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())", o.ID, o.UserID, o.Price, o.Status)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, id string) (*order.Order, error) {
	o := &order.Order{}
	err := r.pool.QueryRow(ctx, "SELECT id, user_id, price, status, created_at, updated_at FROM orders WHERE id = $1", id).Scan(&o.ID, &o.UserID, &o.Price, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return o, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := r.pool.Exec(ctx, "UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2", status, id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
