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
		ID:   uuid.New(),
		Name: name,
	}
	if err := uc.repo.AddCourier(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (uc *UseCase) ReserveCourier(ctx context.Context, orderID uuid.UUID, from, to int64) (uuid.UUID, error) {
	courierID, err := uc.repo.GetAvailableCourier(ctx, from, to)
	if err != nil {
		return uuid.Nil, err
	}

	slot := &courier.Slot{
		OrderID:   orderID,
		CourierID: courierID,
		FromDate:  from,
		ToDate:    to,
	}

	if err := uc.repo.CreateSlot(ctx, slot); err != nil {
		return uuid.Nil, err
	}

	return courierID, nil
}

func (uc *UseCase) CancelCourier(ctx context.Context, orderID uuid.UUID) error {
	return uc.repo.DeleteSlotByOrder(ctx, orderID)
}
