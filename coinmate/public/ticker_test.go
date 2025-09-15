package public

import (
	"encoding/json"
	"net/http"
	"testing"
	"tourGo/coinmate"
)

// Mock client for testing
type MockClient struct {
	coinmate.ClientInterface
	response *coinmate.Response
	err      error
}

func (m *MockClient) GetBaseUrl() string {
	return "https://coinmate.io/api"
}

func (m *MockClient) MakePublicRequest(r coinmate.Request) (coinmate.Response, error) {
	if m.err != nil {
		return coinmate.Response{}, m.err
	}
	return *m.response, nil
}

func (m *MockClient) MakeSecureRequest(r coinmate.Request) (coinmate.Response, error) {
	return coinmate.Response{}, nil
}

func (m *MockClient) GetNonce() string {
	return "1234567890"
}

func (m *MockClient) GetSignature(clientId, apiKey, nonce, privateKey string) string {
	return "test-signature"
}

func (m *MockClient) GetRequestBody(map[string]string) []byte {
	return []byte("test-body")
}

func TestGetTickerSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": {
				"last": 50000.0,
				"high": 51000.0,
				"low": 49000.0,
				"amount": 100.5,
				"bid": 49950.0,
				"ask": 50050.0,
				"change": 500.0,
				"open": 49500.0,
				"timestamp": 1640995200
			}
		}`),
	}

	mockClient := &MockClient{response: mockResponse}
	ticker := &Ticker{Client: mockClient}

	response, err := ticker.GetTicker("BTC_EUR")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if response.Data.Last != 50000.0 {
		t.Errorf("Expected Last to be 50000.0, got %f", response.Data.Last)
	}

	if response.Data.High != 51000.0 {
		t.Errorf("Expected High to be 51000.0, got %f", response.Data.High)
	}

	if response.Data.Low != 49000.0 {
		t.Errorf("Expected Low to be 49000.0, got %f", response.Data.Low)
	}

	if response.Data.Amount != 100.5 {
		t.Errorf("Expected Amount to be 100.5, got %f", response.Data.Amount)
	}

	if response.Data.Bid != 49950.0 {
		t.Errorf("Expected Bid to be 49950.0, got %f", response.Data.Bid)
	}

	if response.Data.Ask != 50050.0 {
		t.Errorf("Expected Ask to be 50050.0, got %f", response.Data.Ask)
	}

	if response.Data.Change != 500.0 {
		t.Errorf("Expected Change to be 500.0, got %f", response.Data.Change)
	}

	if response.Data.Open != 49500.0 {
		t.Errorf("Expected Open to be 49500.0, got %f", response.Data.Open)
	}

	if response.Data.Timestamp != 1640995200 {
		t.Errorf("Expected Timestamp to be 1640995200, got %d", response.Data.Timestamp)
	}
}

func TestGetTickerErrorResponse(t *testing.T) {
	// Create mock error response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": true,
			"errorMessage": "Invalid currency pair",
			"data": {}
		}`),
	}

	mockClient := &MockClient{response: mockResponse}
	ticker := &Ticker{Client: mockClient}

	response, err := ticker.GetTicker("INVALID_PAIR")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !response.Error {
		t.Error("Expected error in response")
	}

	if response.ErrorMessage != "Invalid currency pair" {
		t.Errorf("Expected error message 'Invalid currency pair', got '%s'", response.ErrorMessage)
	}
}

func TestGetTickerHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte("Bad Request"),
	}

	mockClient := &MockClient{response: mockResponse}
	ticker := &Ticker{Client: mockClient}

	response, err := ticker.GetTicker("BTC_EUR")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should return empty response when HTTP status is not OK
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestGetTickerNetworkError(t *testing.T) {
	// Create mock network error
	mockClient := &MockClient{err: &http.ProtocolError{}}
	ticker := &Ticker{Client: mockClient}

	response, err := ticker.GetTicker("BTC_EUR")

	if err == nil {
		t.Error("Expected error for network failure")
	}

	// Should return empty response on network error
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestGetTickerInvalidJSON(t *testing.T) {
	// Create mock response with invalid JSON
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       []byte("invalid json"),
	}

	mockClient := &MockClient{response: mockResponse}
	ticker := &Ticker{Client: mockClient}

	response, err := ticker.GetTicker("BTC_EUR")

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	// Should return empty response on JSON error
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestGetTickerEmptyCurrencyPair(t *testing.T) {
	mockClient := &MockClient{response: &coinmate.Response{}}
	ticker := &Ticker{Client: mockClient}

	response, err := ticker.GetTicker("")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should handle empty currency pair gracefully
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestTickerDataStructure(t *testing.T) {
	// Test that TickerData struct can be marshaled/unmarshaled correctly
	data := TickerData{
		Last:      50000.0,
		High:      51000.0,
		Low:       49000.0,
		Amount:    100.5,
		Bid:       49950.0,
		Ask:       50050.0,
		Change:    500.0,
		Open:      49500.0,
		Timestamp: 1640995200,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal TickerData: %v", err)
	}

	var unmarshaledData TickerData
	err = json.Unmarshal(jsonData, &unmarshaledData)
	if err != nil {
		t.Errorf("Failed to unmarshal TickerData: %v", err)
	}

	if unmarshaledData.Last != data.Last {
		t.Errorf("Expected Last to be %f, got %f", data.Last, unmarshaledData.Last)
	}

	if unmarshaledData.High != data.High {
		t.Errorf("Expected High to be %f, got %f", data.High, unmarshaledData.High)
	}

	if unmarshaledData.Low != data.Low {
		t.Errorf("Expected Low to be %f, got %f", data.Low, unmarshaledData.Low)
	}

	if unmarshaledData.Amount != data.Amount {
		t.Errorf("Expected Amount to be %f, got %f", data.Amount, unmarshaledData.Amount)
	}

	if unmarshaledData.Bid != data.Bid {
		t.Errorf("Expected Bid to be %f, got %f", data.Bid, unmarshaledData.Bid)
	}

	if unmarshaledData.Ask != data.Ask {
		t.Errorf("Expected Ask to be %f, got %f", data.Ask, unmarshaledData.Ask)
	}

	if unmarshaledData.Change != data.Change {
		t.Errorf("Expected Change to be %f, got %f", data.Change, unmarshaledData.Change)
	}

	if unmarshaledData.Open != data.Open {
		t.Errorf("Expected Open to be %f, got %f", data.Open, unmarshaledData.Open)
	}

	if unmarshaledData.Timestamp != data.Timestamp {
		t.Errorf("Expected Timestamp to be %d, got %d", data.Timestamp, unmarshaledData.Timestamp)
	}
}

