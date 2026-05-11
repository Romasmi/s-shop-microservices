package app

import (
	"context"
	"fmt"

	"github.com/Romasmi/s-shop-microservices/order-service/internal/config"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/infrastructure/db/postgres"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/infrastructure/kafka"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/usecase"
	orderuc "github.com/Romasmi/s-shop-microservices/order-service/internal/usecase/order"
	billingapi "github.com/Romasmi/s-shop/gen/go/billing"
	deliveryapi "github.com/Romasmi/s-shop/gen/go/delivery"
	userapi "github.com/Romasmi/s-shop/gen/go/user"
	warehouseapi "github.com/Romasmi/s-shop/gen/go/warehouse"
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
	warehouseConn *grpc.ClientConn
	deliveryConn  *grpc.ClientConn
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
	userConn, err := grpc.NewClient(cfg.UserServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial user service: %w", err)
	}
	userClient := userapi.NewUserServiceClient(userConn)

	billingConn, err := grpc.NewClient(cfg.BillingServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial billing service: %w", err)
	}
	billingClient := billingapi.NewBillingServiceClient(billingConn)

	warehouseConn, err := grpc.NewClient(cfg.WarehouseServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial warehouse service: %w", err)
	}
	warehouseClient := warehouseapi.NewWarehouseServiceClient(warehouseConn)

	deliveryConn, err := grpc.NewClient(cfg.DeliveryServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial delivery service: %w", err)
	}
	deliveryClient := deliveryapi.NewDeliveryServiceClient(deliveryConn)

	app := &App{
		Cfg:           cfg,
		Pool:          pool,
		OrderRepo:     orderRepo,
		OrderProducer: orderProducer,
		Handlers:      make(map[usecase.UseCaseID]usecase.Handler),
		userConn:      userConn,
		billingConn:   billingConn,
		warehouseConn: warehouseConn,
		deliveryConn:  deliveryConn,
	}

	app.registerHandlers(userClient, billingClient, warehouseClient, deliveryClient)

	return app, nil
}

func (a *App) registerHandlers(userClient userapi.UserServiceClient, billingClient billingapi.BillingServiceClient, warehouseClient warehouseapi.WarehouseServiceClient, deliveryClient deliveryapi.DeliveryServiceClient) {
	a.Handlers[usecase.UseCasePlaceOrder] = usecase.NewHandler(orderuc.NewPlaceOrderUseCase(a.OrderRepo, userClient, billingClient, warehouseClient, deliveryClient, a.OrderProducer))
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
	if a.warehouseConn != nil {
		a.warehouseConn.Close()
	}
	if a.deliveryConn != nil {
		a.deliveryConn.Close()
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
