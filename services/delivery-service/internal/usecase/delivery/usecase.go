package delivery

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/domain/courier"
	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) AddCourier(ctx context.Context, name string) (*courier.Courier, error) {
	c := &courier.Courier{
		ID:   uuid.New().String(),
		Name: name,
	}
	if err := uc.repo.AddCourier(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *UseCase) ReserveCourier(ctx context.Context, orderID string, from, to int64) (string, error) {
	courierID, err := uc.repo.GetAvailableCourier(ctx, from, to)
	if err != nil {
		return "", err
	}

	slot := &courier.Slot{
		OrderID:   orderID,
		CourierID: courierID,
		FromDate:  from,
		ToDate:    to,
	}

	if err := uc.repo.CreateSlot(ctx, slot); err != nil {
		return "", err
	}

	return courierID, nil
}

func (uc *UseCase) CancelCourier(ctx context.Context, orderID string) error {
	return uc.repo.DeleteSlotByOrder(ctx, orderID)
}
