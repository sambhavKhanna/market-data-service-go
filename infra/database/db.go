package database

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DbConnection = "host=db user=user password=password dbname=market_data port=5432 sslmode=disable"
)

// New creates a new database connection and configures the connection pool.
func New() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DbConnection), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
