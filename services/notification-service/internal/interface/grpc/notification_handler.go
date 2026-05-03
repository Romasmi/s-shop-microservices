package grpc

import (
	"context"

	api "github.com/Romasmi/s-shop-microservices/notification-service/internal/api"
	"github.com/Romasmi/s-shop-microservices/notification-service/internal/domain/message"
	"github.com/Romasmi/s-shop-microservices/notification-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationHandler struct {
	api.UnimplementedNotificationServiceServer
	app interface {
		GetHandler(id usecase.UseCaseID) usecase.Handler
	}
}

func NewNotificationHandler(app interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}) *NotificationHandler {
	return &NotificationHandler{app: app}
}

func (h *NotificationHandler) ListMessages(ctx context.Context, req *api.ListMessagesRequest) (*api.ListMessagesResponse, error) {
	handler := h.app.GetHandler(usecase.UseCaseListMessages)
	resp, err := handler.Do(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list messages: %v", err)
	}

	msgs := resp.([]*message.Message)
	var apiMsgs []*api.Message
	for _, m := range msgs {
		apiMsgs = append(apiMsgs, &api.Message{
			Id:        m.ID,
			UserId:    m.UserID,
			OrderId:   m.OrderID,
			Type:      m.Type,
			Timestamp: m.CreatedAt.String(),
		})
	}

	return &api.ListMessagesResponse{
		Messages: apiMsgs,
	}, nil
}
