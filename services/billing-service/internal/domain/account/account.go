package account

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	UserID    uuid.UUID `json:"user_id"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
