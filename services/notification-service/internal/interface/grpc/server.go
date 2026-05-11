package grpc

import (
	"github.com/Romasmi/s-shop-microservices/notification-service/internal/usecase"
	api "github.com/Romasmi/s-shop/gen/go/notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(app interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}) *grpc.Server {
	grpcServer := grpc.NewServer()

	notificationHandler := NewNotificationHandler(app)
	api.RegisterNotificationServiceServer(grpcServer, notificationHandler)

	reflection.Register(grpcServer)

	return grpcServer
}
