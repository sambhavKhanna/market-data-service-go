package market_data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"unicode"

	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, db)

	var handler http.Handler = mux

	return handler
}

func addRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.Handle("/prices/latest", GetLatestPrice(db))
	mux.Handle("/prices/poll", PollPrices(db))
}

func GetLatestPrice(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		symbol := query.Get("symbol")
		if symbol == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		runes := []rune(symbol)
		runes[0] = unicode.ToUpper(runes[0])
		symbol = string(runes)

		var pricePoint PricePoint
		err := db.Where("symbol = ?", symbol).Order("timestamp desc").First(&pricePoint).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Fallback to Alpha Vantage provider
				provider := NewAlphaVantageProvider()
				price, err := provider.GetLatestPrice(symbol)
				if err != nil {
					fmt.Print(err)
					w.WriteHeader(http.StatusBadGateway) // Or another appropriate error
					return
				}
				resp := struct {
					Symbol    string    `json:"symbol"`
					Price     float64   `json:"price"`
					Timestamp time.Time `json:"timestamp"`
					Provider  string    `json:"provider"`
				}{
					Symbol:    symbol,
					Price:     price,
					Timestamp: time.Now().UTC(),
					Provider:  "alpha_vantage",
				}
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
				return
			} else {
				// Other database error
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		resp := struct {
			Symbol    string    `json:"symbol"`
			Price     float64   `json:"price"`
			Timestamp time.Time `json:"timestamp"`
			Provider  string    `json:"provider"`
		}{
			Symbol:    pricePoint.Symbol,
			Price:     pricePoint.Price,
			Timestamp: pricePoint.Timestamp,
			Provider:  pricePoint.Provider,
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}

func PollPrices(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var poll struct {
			Symbols  []string `json:"symbols"`
			Interval int32    `json:"interval"`
		}

		if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, symbol := range poll.Symbols {
			job := PollingJob{
				Symbol:          symbol,
				IntervalSeconds: int(poll.Interval),
				IsActive:        true,
			}
			if err := db.Create(&job).Error; err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return 
			}
		}

		w.WriteHeader(http.StatusAccepted)
	})
}

