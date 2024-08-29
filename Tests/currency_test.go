package Tests

import (
	"errors"
	"go_currency_converter/Application"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAPIClient is a mock for the Clients.APIClient
type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) GetCurrencyRates(currency string) (map[string]float64, error) {
	args := m.Called(currency)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func TestCurrencyService_GetRates_Success(t *testing.T) {
	// Arrange
	mockClient := new(MockAPIClient)
	mockClient.On("GetCurrencyRates", "USD").Return(map[string]float64{
		"EUR": 0.85,
		"GBP": 0.75,
	}, nil)

	service := Application.NewCurrencyService(mockClient)

	// Act
	rates, err := service.GetRates("USD")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, rates)
	assert.Equal(t, 0.85, rates["EUR"])
	assert.Equal(t, 0.75, rates["GBP"])

	mockClient.AssertExpectations(t)
}

func TestCurrencyService_GetRates_EmptyCurrency(t *testing.T) {
	// Arrange
	mockClient := new(MockAPIClient)
	service := Application.NewCurrencyService(mockClient)

	// Act
	rates, err := service.GetRates("")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, rates)
	assert.EqualError(t, err, "currency cannot be empty")
}

func TestCurrencyService_GetRates_ErrorFromAPI(t *testing.T) {
	// Arrange
	mockClient := new(MockAPIClient)
	mockClient.On("GetCurrencyRates", "USD").Return((map[string]float64)(nil), errors.New("API error"))

	service := Application.NewCurrencyService(mockClient)

	// Act
	rates, err := service.GetRates("USD")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, rates)
	assert.EqualError(t, err, "API error")

	mockClient.AssertExpectations(t)
}

func TestCurrencyService_Convert_Success(t *testing.T) {
	// Arrange
	mockClient := new(MockAPIClient)
	mockClient.On("GetCurrencyRates", "USD").Return(map[string]float64{
		"EUR": 0.85,
	}, nil)

	service := Application.NewCurrencyService(mockClient)

	// Act
	converted, err := service.Convert(100, "USD", "EUR")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 85.0, converted)

	mockClient.AssertExpectations(t)
}

func TestCurrencyService_Convert_InvalidParameters(t *testing.T) {
	// Arrange
	mockClient := new(MockAPIClient)
	service := Application.NewCurrencyService(mockClient)

	// Act & Assert
	_, err := service.Convert(-100, "USD", "EUR")
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid parameters")

	_, err = service.Convert(100, "", "EUR")
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid parameters")

	_, err = service.Convert(100, "USD", "")
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid parameters")
}

func TestCurrencyService_Convert_CurrencyNotFound(t *testing.T) {
	// Arrange
	mockClient := new(MockAPIClient)
	mockClient.On("GetCurrencyRates", "USD").Return(map[string]float64{
		"EUR": 0.85,
	}, nil)

	service := Application.NewCurrencyService(mockClient)

	// Act
	_, err := service.Convert(100, "USD", "GBP")

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "cannot find currency: GBP")

	mockClient.AssertExpectations(t)
}

func TestCurrencyService_Convert_GetRatesFails(t *testing.T) {
	// Arrange
	mockClient := new(MockAPIClient)
	mockClient.On("GetCurrencyRates", "USD").Return((map[string]float64)(nil), errors.New("API error"))

	service := Application.NewCurrencyService(mockClient)

	// Act
	_, err := service.Convert(100, "USD", "EUR")

	// Assert
	assert.Error(t, err)
	mockClient.AssertExpectations(t)
}
