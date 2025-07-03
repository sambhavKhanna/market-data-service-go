package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/sambhavKhanna/market_data/infra/database"
	"github.com/sambhavKhanna/market_data/infra/kafka"
	"github.com/sambhavKhanna/market_data/internal/market_data"
)

const (
	WindowSize = 5
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	kafkaReader := kafka.NewReader("price-processor")

	for {
		m, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}

		var priceEvent struct {
			Symbol    string    `json:"symbol"`
			Price     float64   `json:"price"`
			Timestamp time.Time `json:"timestamp"`
		}

		if err := json.Unmarshal(m.Value, &priceEvent); err != nil {
			log.Printf("error unmarshalling message: %v", err)
			continue
		}

		// Store the raw price point
		pricePoint := market_data.PricePoint{
			Symbol:    priceEvent.Symbol,
			Price:     priceEvent.Price,
			Timestamp: priceEvent.Timestamp,
			Provider:  "alpha_vantage", // Or get from event
		}
		if err := db.Create(&pricePoint).Error; err != nil {
			log.Printf("error saving price point: %v", err)
			continue
		}

		// Calculate and store the moving average
		var pricePoints []market_data.PricePoint
		if err := db.Where("symbol = ?", priceEvent.Symbol).Order("timestamp desc").Limit(WindowSize).Find(&pricePoints).Error; err != nil {
			log.Printf("error fetching price points for moving average: %v", err)
			continue
		}

		if len(pricePoints) == WindowSize {
			var sum float64
			for _, p := range pricePoints {
				sum += p.Price
			}
			movingAverage := sum / float64(WindowSize)

			ma := market_data.MovingAverage{
				Symbol:        priceEvent.Symbol,
				MovingAverage: movingAverage,
				WindowSize:    WindowSize,
				Timestamp:     priceEvent.Timestamp,
			}

			if err := db.Create(&ma).Error; err != nil {
				log.Printf("error saving moving average: %v", err)
			}
		}

		fmt.Printf("processed price for %s: %f\n", priceEvent.Symbol, priceEvent.Price)
	}
}
