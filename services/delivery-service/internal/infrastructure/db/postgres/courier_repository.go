package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/domain/courier"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourierRepository struct {
	pool *pgxpool.Pool
}

func NewCourierRepository(pool *pgxpool.Pool) *CourierRepository {
	return &CourierRepository{pool: pool}
}

func (r *CourierRepository) AddCourier(ctx context.Context, c *courier.Courier) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO couriers (id, name) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name", c.ID, c.Name)
	if err != nil {
		return fmt.Errorf("failed to add courier: %w", err)
	}
	return nil
}

func (r *CourierRepository) GetAvailableCourier(ctx context.Context, from, to int64) (uuid.UUID, error) {
	var courierID uuid.UUID
	err := r.pool.QueryRow(ctx, `
		SELECT id 
		FROM couriers 
		WHERE id NOT IN (
			SELECT courier_id 
			FROM slots 
			WHERE to_date > $1 AND from_date < $2
		) 
		LIMIT 1`, from, to).Scan(&courierID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, courier.ErrNoAvailableCourier
		}
		return uuid.Nil, fmt.Errorf("failed to get available courier: %w", err)
	}
	return courierID, nil
}

func (r *CourierRepository) CreateSlot(ctx context.Context, s *courier.Slot) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO slots (order_id, courier_id, from_date, to_date) VALUES ($1, $2, $3, $4)", s.OrderID, s.CourierID, s.FromDate, s.ToDate)
	if err != nil {
		return fmt.Errorf("failed to create slot: %w", err)
	}
	return nil
}

func (r *CourierRepository) DeleteSlotByOrder(ctx context.Context, orderID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM slots WHERE order_id = $1", orderID)
	if err != nil {
		return fmt.Errorf("failed to delete slot: %w", err)
	}
	return nil
}
