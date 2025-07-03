package market_data

import (
	"testing"
)

func TestAlphaVantageProvider_GetLatestPrice(t *testing.T) {
	provider := NewAlphaVantageProvider()
	price, err := provider.GetLatestPrice("TSLA")
	if err != nil {
		t.Fatalf("Error getting latest price: %v", err)
	}

	if price <= 0 {
		t.Fatalf("Latest price is not greater than 0: %f", price)
	}
}