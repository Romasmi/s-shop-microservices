package grpc

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/usecase/warehouse"
	api "github.com/Romasmi/s-shop/gen/go/warehouse"
	"github.com/google/uuid"
)

type WarehouseHandler struct {
	api.UnimplementedWarehouseServiceServer
	uc *warehouse.UseCase
}

func NewWarehouseHandler(uc *warehouse.UseCase) *WarehouseHandler {
	return &WarehouseHandler{uc: uc}
}

func (h *WarehouseHandler) AddProduct(ctx context.Context, req *api.AddProductRequest) (*api.Product, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, err
	}
	p, err := h.uc.AddProduct(ctx, productID, req.Count)
	if err != nil {
		return nil, err
	}
	return &api.Product{
		ProductId: p.ID.String(),
		Count:     p.Count,
	}, nil
}

func (h *WarehouseHandler) ReserveProduct(ctx context.Context, req *api.ReserveProductRequest) (*api.ReserveProductResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return &api.ReserveProductResponse{Success: false, Error: err.Error()}, nil
	}
	reservationID, err := h.uc.ReserveProduct(ctx, productID, req.Count)
	if err != nil {
		return &api.ReserveProductResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}
	return &api.ReserveProductResponse{
		Success:       true,
		ReservationId: reservationID.String(),
	}, nil
}

func (h *WarehouseHandler) ReleaseProduct(ctx context.Context, req *api.ReleaseProductRequest) (*api.ReleaseProductResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return &api.ReleaseProductResponse{Success: false, Error: err.Error()}, nil
	}
	reservationID, err := uuid.Parse(req.ReservationId)
	if err != nil {
		return &api.ReleaseProductResponse{Success: false, Error: err.Error()}, nil
	}
	err = h.uc.ReleaseProduct(ctx, productID, reservationID)
	if err != nil {
		return &api.ReleaseProductResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}
	return &api.ReleaseProductResponse{
		Success: true,
	}, nil
}

func (h *WarehouseHandler) CancelReservation(ctx context.Context, req *api.CancelReservationRequest) (*api.CancelReservationResponse, error) {
	reservationID, err := uuid.Parse(req.ReservationId)
	if err != nil {
		return &api.CancelReservationResponse{Success: false, Error: err.Error()}, nil
	}
	err = h.uc.CancelReservation(ctx, reservationID)
	if err != nil {
		return &api.CancelReservationResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}
	return &api.CancelReservationResponse{
		Success: true,
	}, nil
}
