package kafka_client

import (
	"fmt"
	"log"

	"github.com/Romasmi/s-shop-microservices/internal/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Connection struct {
	Consumer *kafka.Consumer
	Producer *kafka.Producer
	Config   *config.Kafka
}

func CreateKafkaConnection(cfg *config.Kafka) (*Connection, error) {
	connection := &Connection{
		Config: cfg,
	}
	err := connection.ConnectProducer()
	if err != nil {
		return nil, fmt.Errorf("error while connection Kafka producer: %v\n", err)
	}
	err = connection.ConnectConsumer()
	if err != nil {
		return nil, fmt.Errorf("error while connection Kafka consumer: %v\n", err)
	}
	return connection, nil
}

func (k *Connection) ConnectProducer() error {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": k.Config.Brokers,
		"acks":              "all",
		"client.id":         "socialNetwork",
	}

	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	k.Producer = producer
	go k.handleDeliveryReports()

	return nil
}

func (k *Connection) Produce(topic string, key, value []byte) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   key,
		Value: value,
	}
	err := k.Producer.Produce(message, nil)
	if err != nil {
		return fmt.Errorf("failed to produce message: %v", err)
	}
	return nil
}

func (k *Connection) ConnectConsumer() error {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": k.Config.Brokers,
		"group.id":          "socialNetworkGroup",
		"auto.offset.reset": "smallest",
	}
	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer: %v", err)
	}
	k.Consumer = consumer
	return nil
}

func (k *Connection) Close() {
	if k.Producer != nil {
		k.Producer.Flush(15 * 1000)
		k.Producer.Close()
	}

	if k.Consumer != nil {
		err := k.Consumer.Close()
		if err != nil {
			log.Printf("error while closing Kafka consumer: %v", err)
		}
	}
}

func (k *Connection) handleDeliveryReports() {
	for e := range k.Producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}
