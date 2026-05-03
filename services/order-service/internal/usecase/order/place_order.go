package order

import (
	"context"
	"fmt"
	"time"

	billingapi "github.com/Romasmi/s-shop-microservices/billing-service/pkg/api"
	"github.com/Romasmi/s-shop-microservices/order-service/internal/domain/order"
	api "github.com/Romasmi/s-shop-microservices/order-service/pkg/api"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PlaceOrderInput struct {
	UserID string
	Price  int64
}

type PlaceOrderUseCase struct {
	repo          Repository
	userClient    api.UserServiceClient
	billingClient billingapi.BillingServiceClient
	producer      EventProducer
}

func NewPlaceOrderUseCase(repo Repository, userClient api.UserServiceClient, billingClient billingapi.BillingServiceClient, producer EventProducer) *PlaceOrderUseCase {
	return &PlaceOrderUseCase{
		repo:          repo,
		userClient:    userClient,
		billingClient: billingClient,
		producer:      producer,
	}
}

func (uc *PlaceOrderUseCase) Do(ctx context.Context, input PlaceOrderInput) (*order.Order, error) {
	// 1. Fetch user details
	userResp, err := uc.userClient.GetUser(ctx, &api.GetUserRequest{Id: input.UserID})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
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

	// 2. Call billing service Withdraw
	withdrawResp, err := uc.billingClient.Withdraw(ctx, &billingapi.WithdrawRequest{
		UserId:         input.UserID,
		Amount:         input.Price,
		IdempotencyKey: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("billing service call failed: %w", err)
	}

	if withdrawResp.Success {
		o.Status = "SUCCESS"
	} else {
		o.Status = "FAILED"
	}

	// 3. Persist order
	if err := uc.repo.CreateOrder(ctx, o); err != nil {
		return nil, fmt.Errorf("failed to persist order: %w", err)
	}

	// 4. Emit order.placed
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
			// Log error but don't fail the order if it was successful?
			// Actually instructions say "Abort with error if this call fails" only for GetUser.
			fmt.Printf("failed to emit order.placed event: %v\n", err)
		}
	}

	return o, nil
}
