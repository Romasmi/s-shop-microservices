package memory

import (
	"context"
	"sync"

	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/domain/courier"
)

type Repository struct {
	mu       sync.RWMutex
	couriers map[string]*courier.Courier
	slots    []*courier.Slot
}

func NewRepository() *Repository {
	return &Repository{
		couriers: make(map[string]*courier.Courier),
		slots:    make([]*courier.Slot, 0),
	}
}

func (r *Repository) AddCourier(ctx context.Context, c *courier.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.couriers[c.ID] = c
	return nil
}

func (r *Repository) GetAvailableCourier(ctx context.Context, from, to int64) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for id := range r.couriers {
		available := true
		for _, s := range r.slots {
			if s.CourierID == id {
				// Simple overlap check
				if (from >= s.FromDate && from < s.ToDate) || (to > s.FromDate && to <= s.ToDate) || (from <= s.FromDate && to >= s.ToDate) {
					available = false
					break
				}
			}
		}
		if available {
			return id, nil
		}
	}

	return "", courier.ErrNoAvailableCourier
}

func (r *Repository) CreateSlot(ctx context.Context, s *courier.Slot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.slots = append(r.slots, s)
	return nil
}

func (r *Repository) DeleteSlotByOrder(ctx context.Context, orderID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, s := range r.slots {
		if s.OrderID == orderID {
			r.slots = append(r.slots[:i], r.slots[i+1:]...)
			return nil
		}
	}
	return nil
}
