package order

import (
	"context"
	"fmt"
	"time"

	"github.com/Romasmi/s-shop-microservices/order-service/internal/domain/order"
	billingapi "github.com/Romasmi/s-shop/gen/go/billing"
	deliveryapi "github.com/Romasmi/s-shop/gen/go/delivery"
	api "github.com/Romasmi/s-shop/gen/go/order"
	userapi "github.com/Romasmi/s-shop/gen/go/user"
	warehouseapi "github.com/Romasmi/s-shop/gen/go/warehouse"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PlaceOrderInput struct {
	UserID    string
	Price     int64
	ProductID string
	Count     int32
}

type PlaceOrderUseCase struct {
	repo            Repository
	userClient      userapi.UserServiceClient
	billingClient   billingapi.BillingServiceClient
	warehouseClient warehouseapi.WarehouseServiceClient
	deliveryClient  deliveryapi.DeliveryServiceClient
	producer        EventProducer
}

func NewPlaceOrderUseCase(
	repo Repository,
	userClient userapi.UserServiceClient,
	billingClient billingapi.BillingServiceClient,
	warehouseClient warehouseapi.WarehouseServiceClient,
	deliveryClient deliveryapi.DeliveryServiceClient,
	producer EventProducer,
) *PlaceOrderUseCase {
	return &PlaceOrderUseCase{
		repo:            repo,
		userClient:      userClient,
		billingClient:   billingClient,
		warehouseClient: warehouseClient,
		deliveryClient:  deliveryClient,
		producer:        producer,
	}
}

func (uc *PlaceOrderUseCase) Do(ctx context.Context, input PlaceOrderInput) (*order.Order, error) {
	// 1. Fetch user details
	userResp, err := uc.userClient.GetUser(ctx, &userapi.GetUserRequest{UserId: input.UserID})
	if err != nil {
		return nil, err
	}

	orderID := uuid.New().String()
	o := &order.Order{
		ID:        orderID,
		UserID:    input.UserID,
		Price:     input.Price,
		Status:    "PENDING",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// SAGA START

	// 1. Warehouse Reserve
	warehouseResp, err := uc.warehouseClient.ReserveProduct(ctx, &warehouseapi.ReserveProductRequest{
		ProductId: input.ProductID,
		Count:     input.Count,
	})
	if err != nil {
		return nil, fmt.Errorf("warehouse service error: %w", err)
	}
	if !warehouseResp.Success {
		o.Status = "FAILED"
		_ = uc.repo.CreateOrder(ctx, o)
		return nil, fmt.Errorf("warehouse reservation failed: %s", warehouseResp.Error)
	}
	reservationID := warehouseResp.ReservationId

	// 2. Delivery Reserve
	fromDate := time.Now().Add(24 * time.Hour)
	toDate := fromDate.Add(2 * time.Hour)
	deliveryResp, err := uc.deliveryClient.ReserveCourier(ctx, &deliveryapi.ReserveCourierRequest{
		OrderId:  orderID,
		FromDate: timestamppb.New(fromDate),
		ToDate:   timestamppb.New(toDate),
	})

	if err != nil || !deliveryResp.Success {
		// ROLLBACK Warehouse
		_, _ = uc.warehouseClient.CancelReservation(ctx, &warehouseapi.CancelReservationRequest{
			ReservationId: reservationID,
		})
		o.Status = "FAILED"
		_ = uc.repo.CreateOrder(ctx, o)
		if err != nil {
			return nil, fmt.Errorf("delivery service error: %w", err)
		}
		return nil, fmt.Errorf("delivery reservation failed: %s", deliveryResp.Error)
	}

	// 3. Billing Withdraw
	withdrawResp, err := uc.billingClient.Withdraw(ctx, &billingapi.WithdrawRequest{
		UserId:         input.UserID,
		Amount:         input.Price,
		IdempotencyKey: orderID,
	})

	if err != nil || !withdrawResp.Success {
		// ROLLBACK Delivery
		_, _ = uc.deliveryClient.CancelCourier(ctx, &deliveryapi.CancelCourierRequest{
			OrderId: orderID,
		})
		// ROLLBACK Warehouse
		_, _ = uc.warehouseClient.CancelReservation(ctx, &warehouseapi.CancelReservationRequest{
			ReservationId: reservationID,
		})

		o.Status = "FAILED"
		_ = uc.repo.CreateOrder(ctx, o)

		if err != nil {
			return nil, fmt.Errorf("billing service error: %w", err)
		}
		return nil, fmt.Errorf("billing failed: %s", withdrawResp.Reason)
	}

	// SAGA SUCCESS - Finalize

	// Release Warehouse Product (decrease count)
	_, _ = uc.warehouseClient.ReleaseProduct(ctx, &warehouseapi.ReleaseProductRequest{
		ProductId:     input.ProductID,
		ReservationId: reservationID,
	})

	o.Status = "SUCCESS"

	// Persist order
	if err := uc.repo.CreateOrder(ctx, o); err != nil {
		return nil, fmt.Errorf("failed to persist order: %w", err)
	}

	// Emit order.placed
	event := &api.OrderPlaced{
		EventId:    uuid.New().String(),
		OccurredAt: timestamppb.New(time.Now()),
		Order: &api.OrderPlaced_Order{
			Id:     o.ID,
			UserId: o.UserID,
			Price:  o.Price,
			Status: o.Status,
		},
		PaymentResult: &api.OrderPlaced_PaymentResult{
			Success:             withdrawResp.Success,
			Reason:              withdrawResp.Reason,
			AccountBalanceAfter: withdrawResp.UpdatedAccount.Balance,
		},
		User: &api.OrderPlaced_UserInfo{
			Id:    userResp.Id,
			Email: userResp.Email,
		},
	}

	if uc.producer != nil {
		if err := uc.producer.EmitOrderPlaced(ctx, event); err != nil {
			fmt.Printf("failed to emit order.placed event: %v\n", err)
		}
	}

	return o, nil
}
