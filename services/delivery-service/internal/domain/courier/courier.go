package courier

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNoAvailableCourier = errors.New("no available courier")
)

type Courier struct {
	ID   uuid.UUID
	Name string
}

type Slot struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
	FromDate  int64
	ToDate    int64
}
