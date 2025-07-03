package market_data

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

func PublishPriceEvent(writer *kafka.Writer, symbol string, price float64) error {
	msg := map[string]interface{}{
		"symbol":    symbol,
		"price":     price,
		"timestamp": time.Now().UTC(),
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(symbol),
			Value: jsonMsg,
		},
	)
}