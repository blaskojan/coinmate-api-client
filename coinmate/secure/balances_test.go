package secure

import (
	"encoding/json"
	"net/http"
	"testing"
	"tourGo/coinmate"
)

// Mock client for testing
type MockSecureClient struct {
	coinmate.ClientInterface
	response *coinmate.Response
	err      error
}

func (m *MockSecureClient) GetBaseUrl() string {
	return "https://coinmate.io/api"
}

func (m *MockSecureClient) MakePublicRequest(r coinmate.Request) (coinmate.Response, error) {
	return coinmate.Response{}, nil
}

func (m *MockSecureClient) MakeSecureRequest(r coinmate.Request) (coinmate.Response, error) {
	if m.err != nil {
		return coinmate.Response{}, m.err
	}
	return *m.response, nil
}

func (m *MockSecureClient) GetNonce() string {
	return "1234567890"
}

func (m *MockSecureClient) GetSignature(clientId, apiKey, nonce, privateKey string) string {
	return "test-signature"
}

func (m *MockSecureClient) GetRequestBody(map[string]string) []byte {
	return []byte("test-body")
}

func TestGetBalancesSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": {
				"BTC": {
					"currency": "BTC",
					"balance": 1.5,
					"reserved": 0.1,
					"available": 1.4
				},
				"EUR": {
					"currency": "EUR",
					"balance": 1000.0,
					"reserved": 50.0,
					"available": 950.0
				}
			}
		}`),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	balances := &Balances{Client: mockClient}

	response, err := balances.GetBalances()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 currencies, got %d", len(response.Data))
	}

	// Check BTC balance
	btcBalance := response.Data["BTC"]
	if btcBalance.Currency != "BTC" {
		t.Errorf("Expected BTC currency, got %s", btcBalance.Currency)
	}
	if btcBalance.Balance != 1.5 {
		t.Errorf("Expected BTC balance to be 1.5, got %f", btcBalance.Balance)
	}
	if btcBalance.Reserved != 0.1 {
		t.Errorf("Expected BTC reserved to be 0.1, got %f", btcBalance.Reserved)
	}
	if btcBalance.Available != 1.4 {
		t.Errorf("Expected BTC available to be 1.4, got %f", btcBalance.Available)
	}

	// Check EUR balance
	eurBalance := response.Data["EUR"]
	if eurBalance.Currency != "EUR" {
		t.Errorf("Expected EUR currency, got %s", eurBalance.Currency)
	}
	if eurBalance.Balance != 1000.0 {
		t.Errorf("Expected EUR balance to be 1000.0, got %f", eurBalance.Balance)
	}
	if eurBalance.Reserved != 50.0 {
		t.Errorf("Expected EUR reserved to be 50.0, got %f", eurBalance.Reserved)
	}
	if eurBalance.Available != 950.0 {
		t.Errorf("Expected EUR available to be 950.0, got %f", eurBalance.Available)
	}
}

func TestGetBalancesErrorResponse(t *testing.T) {
	// Create mock error response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": true,
			"errorMessage": "Authentication failed",
			"data": {}
		}`),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	balances := &Balances{Client: mockClient}

	response, err := balances.GetBalances()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !response.Error {
		t.Error("Expected error in response")
	}

	if response.ErrorMessage != "Authentication failed" {
		t.Errorf("Expected error message 'Authentication failed', got '%s'", response.ErrorMessage)
	}
}

func TestGetBalancesHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusUnauthorized,
		Status:     "401 Unauthorized",
		Body:       []byte("Unauthorized"),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	balances := &Balances{Client: mockClient}

	_, err := balances.GetBalances()

	if err == nil {
		t.Errorf("Expected error for non-200 response")
	}
}

func TestGetBalancesNetworkError(t *testing.T) {
	// Create mock network error
	mockClient := &MockSecureClient{err: &http.ProtocolError{}}
	balances := &Balances{Client: mockClient}

	_, err := balances.GetBalances()

	if err == nil {
		t.Error("Expected error for network failure")
	}
}

func TestGetBalancesInvalidJSON(t *testing.T) {
	// Create mock response with invalid JSON
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       []byte("invalid json"),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	balances := &Balances{Client: mockClient}

	_, err := balances.GetBalances()

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestBalancesDataStructure(t *testing.T) {
	// Test that BalanceCurrency struct can be marshaled/unmarshaled correctly
	data := BalanceCurrency{
		Currency:  "BTC",
		Balance:   1.5,
		Reserved:  0.1,
		Available: 1.4,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal BalanceCurrency: %v", err)
	}

	var unmarshaledData BalanceCurrency
	err = json.Unmarshal(jsonData, &unmarshaledData)
	if err != nil {
		t.Errorf("Failed to unmarshal BalanceCurrency: %v", err)
	}

	if unmarshaledData.Currency != data.Currency {
		t.Errorf("Expected currency to be %s, got %s", data.Currency, unmarshaledData.Currency)
	}

	if unmarshaledData.Balance != data.Balance {
		t.Errorf("Expected balance to be %f, got %f", data.Balance, unmarshaledData.Balance)
	}

	if unmarshaledData.Reserved != data.Reserved {
		t.Errorf("Expected reserved to be %f, got %f", data.Reserved, unmarshaledData.Reserved)
	}

	if unmarshaledData.Available != data.Available {
		t.Errorf("Expected available to be %f, got %f", data.Available, unmarshaledData.Available)
	}
}
