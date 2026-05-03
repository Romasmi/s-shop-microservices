package kafka

import (
	"context"
	"fmt"

	api "github.com/Romasmi/s-shop-microservices/order-service/pkg/api"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/encoding/protojson"
)

type OrderProducer struct {
	writer *kafka.Writer
}

func NewOrderProducer(brokers []string, topic string) *OrderProducer {
	return &OrderProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *OrderProducer) EmitOrderPlaced(ctx context.Context, event *api.OrderPlaced) error {
	payload, err := protojson.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal order.placed event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Order.Id),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("failed to write order.placed message: %w", err)
	}

	return nil
}

func (p *OrderProducer) Close() error {
	return p.writer.Close()
}
