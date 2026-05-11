package courier

import "errors"

var (
	ErrNoAvailableCourier = errors.New("no available courier")
)

type Courier struct {
	ID   string
	Name string
}

type Slot struct {
	OrderID   string
	CourierID string
	FromDate  int64
	ToDate    int64
}
