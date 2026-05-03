package grpc

import (
	"github.com/Romasmi/s-shop-microservices/order-service/internal/usecase"
	api "github.com/Romasmi/s-shop-microservices/order-service/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(app interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}) *grpc.Server {
	grpcServer := grpc.NewServer()

	orderHandler := NewOrderHandler(app)
	api.RegisterOrderServiceServer(grpcServer, orderHandler)

	reflection.Register(grpcServer)

	return grpcServer
}
