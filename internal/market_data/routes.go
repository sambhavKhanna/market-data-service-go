package market_data

import (
	"net/http"
	"time"
	"encoding/json"
	"unicode"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux)

	var handler http.Handler = mux

	return handler
}

func addRoutes(mux *http.ServeMux) {
	mux.Handle("/prices/latest", GetLatestPrice())
}

func GetLatestPrice() http.Handler {
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

		provider := NewAlphaVantageProvider("")
		price, err := provider.GetLatestPrice(symbol)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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
	})
}

func PollPrices() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var poll struct {
			Symbols   []string    `json:"symbols"`
			Interval     int32   	`json:"interval"`
		}

		if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}

