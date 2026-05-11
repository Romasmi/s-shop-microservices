package grpc

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/order-service/internal/domain/order"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/usecase"
	orderuc "github.com/Romasmi/s-shop-microservices/order-service/internal/usecase/order"
	api "github.com/Romasmi/s-shop/gen/go/order"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	api.UnimplementedOrderServiceServer
	app interface {
		GetHandler(id usecase.UseCaseID) usecase.Handler
	}
}

func NewOrderHandler(app interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}) *OrderHandler {
	return &OrderHandler{app: app}
}

func (h *OrderHandler) PlaceOrder(ctx context.Context, req *api.PlaceOrderRequest) (*api.Order, error) {
	handler := h.app.GetHandler(usecase.UseCasePlaceOrder)
	resp, err := handler.Do(ctx, orderuc.PlaceOrderInput{
		UserID:    req.UserId,
		Price:     req.Price,
		ProductID: req.ProductId,
		Count:     req.Count,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Errorf(codes.Internal, "failed to place order: %v", err)
	}

	o := resp.(*order.Order)
	return &api.Order{
		Id:     o.ID.String(),
		UserId: o.UserID.String(),
		Price:  o.Price,
		Status: o.Status,
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *api.GetOrderRequest) (*api.Order, error) {
	handler := h.app.GetHandler(usecase.UseCaseGetOrder)
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order id: %v", err)
	}
	resp, err := handler.Do(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	o := resp.(*order.Order)
	return &api.Order{
		Id:     o.ID.String(),
		UserId: o.UserID.String(),
		Price:  o.Price,
		Status: o.Status,
	}, nil
}
