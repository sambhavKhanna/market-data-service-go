package kafka

import (
	"github.com/segmentio/kafka-go"
)

const (
	KafkaBroker      = "kafka:29092"
	RawPriceEventsTopic = "raw-price-events"
)

// NewWriter creates a new Kafka writer for publishing messages.
func NewWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(KafkaBroker),
		Topic:    RawPriceEventsTopic,
		Balancer: &kafka.LeastBytes{},
	}
}

// NewReader creates a new Kafka reader for consuming messages.
func NewReader(groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{KafkaBroker},
		Topic:    RawPriceEventsTopic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}
