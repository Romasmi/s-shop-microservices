package grpc

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/order-service/internal/domain/order"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/usecase"
	orderuc "github.com/Romasmi/s-shop-microservices/order-service/internal/usecase/order"
	api "github.com/Romasmi/s-shop-microservices/order-service/pkg/api"
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
		UserID: req.UserId,
		Price:  req.Price,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to place order: %v", err)
	}

	o := resp.(*order.Order)
	return &api.Order{
		Id:     o.ID,
		UserId: o.UserID,
		Price:  o.Price,
		Status: o.Status,
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *api.GetOrderRequest) (*api.Order, error) {
	handler := h.app.GetHandler(usecase.UseCaseGetOrder)
	resp, err := handler.Do(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	o := resp.(*order.Order)
	return &api.Order{
		Id:     o.ID,
		UserId: o.UserID,
		Price:  o.Price,
		Status: o.Status,
	}, nil
}
