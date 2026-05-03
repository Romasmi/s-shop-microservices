package grpc

import (
	api "github.com/Romasmi/s-shop-microservices/billing-service/internal/api"
	"github.com/Romasmi/s-shop-microservices/billing-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(app interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}) *grpc.Server {
	grpcServer := grpc.NewServer()

	billingHandler := NewBillingHandler(app)
	api.RegisterBillingServiceServer(grpcServer, billingHandler)

	reflection.Register(grpcServer)

	return grpcServer
}
