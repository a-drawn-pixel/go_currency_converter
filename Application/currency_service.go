package Application

import (
	"errors"
	"fmt"
	"go_currency_converter/Application/Clients"
)

type CurrencyService interface {
	GetRates(currency string) (map[string]float64, error)
	Convert(amount float64, fromCurrency string, toCurrency string) (float64, error)
}

type currencyService struct {
	apiClient Clients.APIClient
}

func NewCurrencyService(client Clients.APIClient) CurrencyService {
	return &currencyService{
		apiClient: client,
	}
}

func (s *currencyService) GetRates(currency string) (map[string]float64, error) {
	if currency == "" {
		return nil, errors.New("currency cannot be empty")
	}

	response, err := s.apiClient.GetCurrencyRates(currency)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func (s *currencyService) Convert(amount float64, fromCurrency string, toCurrency string) (float64, error) {
	if fromCurrency == "" || toCurrency == "" || amount < 0 {
		return 0, errors.New("invalid parameters")
	}

	rates, err := s.GetRates(fromCurrency)
	if err != nil {
		return 0, err
	}

	rate, exists := rates[toCurrency]
	if !exists {
		return 0, fmt.Errorf("cannot find currency: %s", toCurrency)
	}
	return amount * rate, nil
}
