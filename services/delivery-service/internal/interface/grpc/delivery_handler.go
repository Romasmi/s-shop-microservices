package grpc

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/usecase/delivery"
	api "github.com/Romasmi/s-shop/gen/go/delivery"
	"github.com/google/uuid"
)

type DeliveryHandler struct {
	api.UnimplementedDeliveryServiceServer
	uc *delivery.UseCase
}

func NewDeliveryHandler(uc *delivery.UseCase) *DeliveryHandler {
	return &DeliveryHandler{uc: uc}
}

func (h *DeliveryHandler) AddCourier(ctx context.Context, req *api.AddCourierRequest) (*api.Courier, error) {
	c, err := h.uc.AddCourier(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &api.Courier{
		CourierId: c.ID.String(),
		Name:      c.Name,
	}, nil
}

func (h *DeliveryHandler) ReserveCourier(ctx context.Context, req *api.ReserveCourierRequest) (*api.ReserveCourierResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return &api.ReserveCourierResponse{Success: false, Error: err.Error()}, nil
	}
	courierID, err := h.uc.ReserveCourier(ctx, orderID, req.FromDate.Seconds, req.ToDate.Seconds)
	if err != nil {
		return &api.ReserveCourierResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}
	return &api.ReserveCourierResponse{
		Success:   true,
		CourierId: courierID.String(),
	}, nil
}

func (h *DeliveryHandler) CancelCourier(ctx context.Context, req *api.CancelCourierRequest) (*api.CancelCourierResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return &api.CancelCourierResponse{Success: false, Error: err.Error()}, nil
	}
	err = h.uc.CancelCourier(ctx, orderID)
	if err != nil {
		return &api.CancelCourierResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}
	return &api.CancelCourierResponse{
		Success: true,
	}, nil
}
