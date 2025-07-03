package market_data

import (
	"time"

	"gorm.io/datatypes"
)

type RawMarketData struct {
	ID        uint      `gorm:"primaryKey"`
	Symbol    string    `gorm:"size:10;not null;index"`
	Provider  string    `gorm:"size:50;not null"`
	Response  datatypes.JSON `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

func (RawMarketData) TableName() string {
	return "raw_market_data"
}

type PricePoint struct {
	ID        uint      `gorm:"primaryKey"`
	Symbol    string    `gorm:"size:10;not null;index"`
	Price     float64   `gorm:"type:decimal(10,4);not null"`
	Timestamp time.Time `gorm:"not null;index"`
	Provider  string    `gorm:"size:50;not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

func (PricePoint) TableName() string {
	return "price_points"
}

type MovingAverage struct {
	ID            uint      `gorm:"primaryKey"`
	Symbol        string    `gorm:"size:10;not null;index"`
	MovingAverage float64   `gorm:"type:decimal(10,4);not null"`
	WindowSize    int       `gorm:"not null"`
	Timestamp     time.Time `gorm:"not null;index"`
	CreatedAt     time.Time `gorm:"not null;default:now()"`
}

func (MovingAverage) TableName() string {
	return "moving_averages"
}

type PollingJob struct {
	ID              uint      `gorm:"primaryKey"`
	Symbol          string    `gorm:"size:10;not null;index"`
	IntervalSeconds int       `gorm:"not null"`
	IsActive        bool      `gorm:"not null;default:true"`
	CreatedAt       time.Time `gorm:"not null;default:now()"`
	UpdatedAt       time.Time `gorm:"not null;default:now()"`
}

func (PollingJob) TableName() string {
	return "polling_jobs"
}
