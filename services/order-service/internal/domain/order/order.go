package order

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Price     int64     `json:"price"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
