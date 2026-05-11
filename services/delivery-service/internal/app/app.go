package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/config"
	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/infrastructure/db/memory"
	grpcint "github.com/Romasmi/s-shop-microservices/delivery-service/internal/interface/grpc"
	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/usecase/delivery"
	api "github.com/Romasmi/s-shop/gen/go/delivery"
	"google.golang.org/grpc"
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run() error {
	repo := memory.NewRepository()
	uc := delivery.NewUseCase(repo)
	handler := grpcint.NewDeliveryHandler(uc)

	grpcServer := grpc.NewServer()
	api.RegisterDeliveryServiceServer(grpcServer, handler)

	addr := fmt.Sprintf(":%d", a.cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("Delivery service starting", "addr", addr)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
