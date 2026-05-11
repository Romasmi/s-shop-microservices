package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/config"
	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/infrastructure/db/memory"
	grpcint "github.com/Romasmi/s-shop-microservices/warehouse-service/internal/interface/grpc"
	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/usecase/warehouse"
	api "github.com/Romasmi/s-shop/gen/go/warehouse"
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
	uc := warehouse.NewUseCase(repo)
	handler := grpcint.NewWarehouseHandler(uc)

	grpcServer := grpc.NewServer()
	api.RegisterWarehouseServiceServer(grpcServer, handler)

	addr := fmt.Sprintf(":%d", a.cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("Warehouse service starting", "addr", addr)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
