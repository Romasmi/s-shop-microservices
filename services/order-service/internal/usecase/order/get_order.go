package order

import (
	"context"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/domain/order"
)

type GetOrderUseCase struct {
	repo Repository
}

func NewGetOrderUseCase(repo Repository) *GetOrderUseCase {
	return &GetOrderUseCase{repo: repo}
}

func (uc *GetOrderUseCase) Do(ctx context.Context, id string) (*order.Order, error) {
	return uc.repo.GetOrder(ctx, id)
}
