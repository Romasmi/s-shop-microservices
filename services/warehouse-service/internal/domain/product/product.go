package product

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type Product struct {
	ID    uuid.UUID
	Count int32
}

type Reservation struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Count     int32
}
