package billing

import (
	"context"
	"github.com/Romasmi/s-shop-microservices/billing-service/internal/domain/account"
	"github.com/google/uuid"
)

type TopUpInput struct {
	UserID uuid.UUID
	Amount int64
}

type TopUpUseCase struct {
	repo Repository
}

func NewTopUpUseCase(repo Repository) *TopUpUseCase {
	return &TopUpUseCase{repo: repo}
}

func (uc *TopUpUseCase) Do(ctx context.Context, input TopUpInput) (*account.Account, error) {
	return uc.repo.UpdateBalance(ctx, input.UserID, input.Amount)
}
