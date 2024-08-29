package Clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ExchangeRateResponse struct {
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

type APIClient interface {
	GetCurrencyRates(currency string) (map[string]float64, error)
}

type httpClient struct {
	apiKey  string
	baseURL string
}

func NewHttpClient(key string) APIClient {
	return &httpClient{
		apiKey:  key,
		baseURL: fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/", key),
	}
}

func (c *httpClient) GetCurrencyRates(currency string) (map[string]float64, error) {
	response, err := http.Get(c.baseURL + "/" + currency)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var exchangeRateResponse ExchangeRateResponse
	err = json.Unmarshal(body, &exchangeRateResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return exchangeRateResponse.ConversionRates, nil
}
