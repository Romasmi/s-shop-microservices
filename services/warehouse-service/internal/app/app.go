package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/config"
	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/infrastructure/db/postgres"
	grpcint "github.com/Romasmi/s-shop-microservices/warehouse-service/internal/interface/grpc"
	httpint "github.com/Romasmi/s-shop-microservices/warehouse-service/internal/interface/http"
	"github.com/Romasmi/s-shop-microservices/warehouse-service/internal/usecase/warehouse"
	api "github.com/Romasmi/s-shop/gen/go/warehouse"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	pool       *pgxpool.Pool
	grpcServer *grpc.Server
	gwServer   *http.Server
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

	repo := postgres.NewProductRepository(pool)
	uc := warehouse.NewUseCase(repo)
	handler := grpcint.NewWarehouseHandler(uc)

	a.grpcServer = grpc.NewServer()
	api.RegisterWarehouseServiceServer(a.grpcServer, handler)

	grpcAddr := fmt.Sprintf(":%d", a.cfg.Server.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	go func() {
		slog.Info("Warehouse service gRPC starting", "addr", grpcAddr)
		if err := a.grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	var errGw error
	a.gwServer, errGw = httpint.NewGatewayServer(grpcAddr, a.cfg.Server.Port)
	if errGw != nil {
		return fmt.Errorf("failed to create gateway server: %w", errGw)
	}

	go func() {
		slog.Info("Warehouse service HTTP gateway starting", "addr", a.gwServer.Addr)
		if err := a.gwServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP gateway error", "error", err)
		}
	}()

	<-ctx.Done()
	return a.Shutdown(context.Background())
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.grpcServer != nil {
		slog.Info("Shutting down gRPC server...")
		a.grpcServer.GracefulStop()
	}
	if a.gwServer != nil {
		slog.Info("Shutting down HTTP gateway...")
		if err := a.gwServer.Shutdown(ctx); err != nil {
			slog.Error("HTTP gateway shutdown error", "error", err)
		}
	}
	if a.pool != nil {
		a.pool.Close()
	}
	return nil
}
