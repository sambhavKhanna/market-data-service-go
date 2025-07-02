package market_data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Provider interface {
	GetProviderName() string
	GetLatestPrice(symbol string) (float64, error)
}

type AlphaVantageProvider struct {
	apiKey string
	baseURL string
	function string
}

type GlobalQuoteResponse struct {
	GlobalQuote GlobalQuote `json:"Global Quote"`
}

type GlobalQuote struct {
	Symbol           string `json:"01. symbol"`
	Open             string `json:"02. open"`
	High             string `json:"03. high"`
	Low              string `json:"04. low"`
	Price            string `json:"05. price"`
	Volume           string `json:"06. volume"`
	LatestTradingDay string `json:"07. latest trading day"`
	PreviousClose    string `json:"08. previous close"`
	Change           string `json:"09. change"`
	ChangePercent    string `json:"10. change percent"`
}

func NewAlphaVantageProvider(apiKey string) *AlphaVantageProvider {
	return &AlphaVantageProvider{
		apiKey: apiKey,
		baseURL: "https://www.alphavantage.co/query",
		function: "GLOBAL_QUOTE",
	}
}

func (p *AlphaVantageProvider) GetProviderName() string {
	return "alpha_vantage"
}

func (p *AlphaVantageProvider) GetLatestPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("%s?function=%s&symbol=%s&apikey=%s", p.baseURL, p.function, symbol, p.apiKey)
	resp, err := http.Get(url)

	if err != nil {
		return 0, fmt.Errorf("error fetching latest price: %w", err)
	}
	defer resp.Body.Close()

	var quote GlobalQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return 0, fmt.Errorf("decode json: %w", err)
	}

	price, err := strconv.ParseFloat(quote.GlobalQuote.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing latest price: %w", err)
	}

	return price, nil
}