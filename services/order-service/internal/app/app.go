package app

import (
	"context"
	"fmt"

	billingapi "github.com/Romasmi/s-shop-microservices/billing-service/pkg/api"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/config"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/infrastructure/db/postgres"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/infrastructure/kafka"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/usecase"
	orderuc "github.com/Romasmi/s-shop-microservices/order-service/internal/usecase/order"
	api "github.com/Romasmi/s-shop-microservices/order-service/pkg/api"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	Cfg           *config.Config
	Pool          *pgxpool.Pool
	OrderRepo     *postgres.OrderRepository
	OrderProducer *kafka.OrderProducer
	Handlers      map[usecase.UseCaseID]usecase.Handler
	userConn      *grpc.ClientConn
	billingConn   *grpc.ClientConn
}

func NewApp(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Name)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	orderRepo := postgres.NewOrderRepository(pool)
	orderProducer := kafka.NewOrderProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)

	// gRPC clients
	userConn, err := grpc.Dial(cfg.UserServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial user service: %w", err)
	}
	userClient := api.NewUserServiceClient(userConn)

	billingConn, err := grpc.Dial(cfg.BillingServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial billing service: %w", err)
	}
	billingClient := billingapi.NewBillingServiceClient(billingConn)

	app := &App{
		Cfg:           cfg,
		Pool:          pool,
		OrderRepo:     orderRepo,
		OrderProducer: orderProducer,
		Handlers:      make(map[usecase.UseCaseID]usecase.Handler),
		userConn:      userConn,
		billingConn:   billingConn,
	}

	app.registerHandlers(userClient, billingClient)

	return app, nil
}

func (a *App) registerHandlers(userClient api.UserServiceClient, billingClient billingapi.BillingServiceClient) {
	a.Handlers[usecase.UseCasePlaceOrder] = usecase.NewHandler(orderuc.NewPlaceOrderUseCase(a.OrderRepo, userClient, billingClient, a.OrderProducer))
	a.Handlers[usecase.UseCaseGetOrder] = usecase.NewHandler(orderuc.NewGetOrderUseCase(a.OrderRepo))
}

func (a *App) GetHandler(id usecase.UseCaseID) usecase.Handler {
	return a.Handlers[id]
}

func (a *App) Close() {
	if a.Pool != nil {
		a.Pool.Close()
	}
	if a.OrderProducer != nil {
		a.OrderProducer.Close()
	}
	if a.userConn != nil {
		a.userConn.Close()
	}
	if a.billingConn != nil {
		a.billingConn.Close()
	}
}

func (a *App) GetConfig() *config.Config {
	return a.Cfg
}

func (a *App) Ping(ctx context.Context) error {
	if a.Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	return a.Pool.Ping(ctx)
}
