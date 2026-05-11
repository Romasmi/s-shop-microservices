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
	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/infrastructure/db/postgres"
	grpcint "github.com/Romasmi/s-shop-microservices/delivery-service/internal/interface/grpc"
	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/usecase/delivery"
	api "github.com/Romasmi/s-shop/gen/go/delivery"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type App struct {
	cfg  *config.Config
	pool *pgxpool.Pool
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		a.cfg.Db.User, a.cfg.Db.Password, a.cfg.Db.Host, a.cfg.Db.Port, a.cfg.Db.Name)

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	a.pool = pool

	repo := postgres.NewCourierRepository(pool)
	uc := delivery.NewUseCase(repo)
	handler := grpcint.NewDeliveryHandler(uc)

	grpcServer := grpc.NewServer()
	api.RegisterDeliveryServiceServer(grpcServer, handler)

	addr := fmt.Sprintf(":%d", a.cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		slog.Info("Delivery service starting", "addr", addr)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
	a.pool.Close()
	return nil
}
