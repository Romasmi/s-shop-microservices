package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/domain/product"
)

type Repository struct {
	mu           sync.RWMutex
	products     map[string]*product.Product
	reservations map[string]*product.Reservation
}

func NewRepository() *Repository {
	return &Repository{
		products:     make(map[string]*product.Product),
		reservations: make(map[string]*product.Reservation),
	}
}

func (r *Repository) GetProduct(ctx context.Context, productID string) (*product.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.products[productID]
	if !ok {
		return nil, product.ErrProductNotFound
	}
	return p, nil
}

func (r *Repository) UpsertProduct(ctx context.Context, p *product.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.products[p.ID] = p
	return nil
}

func (r *Repository) CreateReservation(ctx context.Context, res *product.Reservation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reservations[res.ID] = res
	return nil
}

func (r *Repository) GetReservation(ctx context.Context, reservationID string) (*product.Reservation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res, ok := r.reservations[reservationID]
	if !ok {
		return nil, fmt.Errorf("reservation not found")
	}
	return res, nil
}

func (r *Repository) DeleteReservation(ctx context.Context, reservationID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.reservations, reservationID)
	return nil
}

func (r *Repository) UpdateProductCount(ctx context.Context, productID string, delta int32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.products[productID]
	if !ok {
		return product.ErrProductNotFound
	}
	p.Count += delta
	return nil
}
