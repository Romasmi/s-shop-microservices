package warehouse

import (
	"context"
	"fmt"

	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/domain/product"
	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
}

func NewUseCase(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) AddProduct(ctx context.Context, productID uuid.UUID, count int32) (*product.Product, error) {
	p, err := uc.repo.GetProduct(ctx, productID)
	if err != nil && err != product.ErrProductNotFound {
		return nil, err
	}
	if p == nil {
		p = &product.Product{ID: productID, Count: count}
	} else {
		p.Count += count
	}
	if err := uc.repo.UpsertProduct(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *UseCase) ReserveProduct(ctx context.Context, productID uuid.UUID, count int32) (uuid.UUID, error) {
	p, err := uc.repo.GetProduct(ctx, productID)
	if err != nil {
		return uuid.Nil, err
	}
	if p.Count < count {
		return uuid.Nil, product.ErrInsufficientStock
	}

	reservationID := uuid.New()
	r := &product.Reservation{
		ID:        reservationID,
		ProductID: productID,
		Count:     count,
	}

	if err := uc.repo.CreateReservation(ctx, r); err != nil {
		return uuid.Nil, err
	}

	return reservationID, nil
}

func (uc *UseCase) ReleaseProduct(ctx context.Context, productID uuid.UUID, reservationID uuid.UUID) error {
	r, err := uc.repo.GetReservation(ctx, reservationID)
	if err != nil {
		return err
	}
	if r.ProductID != productID {
		return fmt.Errorf("reservation %s is for product %s, not %s", reservationID, r.ProductID, productID)
	}

	// Instruction: "releaseProduct(productid, reservationId) decrease product count"
	if err := uc.repo.UpdateProductCount(ctx, productID, -r.Count); err != nil {
		return err
	}

	if err := uc.repo.DeleteReservation(ctx, reservationID); err != nil {
		// Rollback count if deletion fails
		_ = uc.repo.UpdateProductCount(ctx, productID, r.Count)
		return err
	}

	return nil
}

func (uc *UseCase) CancelReservation(ctx context.Context, reservationID uuid.UUID) error {
	if err := uc.repo.DeleteReservation(ctx, reservationID); err != nil {
		return err
	}
	return nil
}
