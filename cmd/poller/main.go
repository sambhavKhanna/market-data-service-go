package main

import (
	"log"
	"time"

	"github.com/sambhavKhanna/market_data/infra/database"
	"github.com/sambhavKhanna/market_data/infra/kafka"
	"github.com/sambhavKhanna/market_data/internal/market_data"
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	kafkaWriter := kafka.NewWriter()
	provider := market_data.NewAlphaVantageProvider()

	for {
		var jobs []market_data.PollingJob
		if err := db.Where("is_active = ?", true).Find(&jobs).Error; err != nil {
			log.Printf("error fetching polling jobs: %v", err)
			time.Sleep(10 * time.Second) // Wait before retrying
			continue
		}

		for _, job := range jobs {
			go func(job market_data.PollingJob) {
				ticker := time.NewTicker(time.Duration(job.IntervalSeconds) * time.Second)
				for range ticker.C {
					price, err := provider.GetLatestPrice(job.Symbol)
					if err != nil {
						log.Printf("error fetching price for %s: %v", job.Symbol, err)
						continue
					}

					if err := market_data.PublishPriceEvent(kafkaWriter, job.Symbol, price); err != nil {
						log.Printf("error publishing price event for %s: %v", job.Symbol, err)
					}
				}
			}(job)
		}

		time.Sleep(1 * time.Minute) // Re-query for new jobs every minute
	}
}
