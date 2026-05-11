package product

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type Product struct {
	ID    string
	Count int32
}

type Reservation struct {
	ID        string
	ProductID string
	Count     int32
}
