package public

import (
	"encoding/json"
	"net/http"
	"testing"
	"tourGo/coinmate"
)

func TestGetOrderBookSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": {
				"asks": [
					{"price": 50050.0, "amount": 1.5},
					{"price": 50060.0, "amount": 2.0}
				],
				"bids": [
					{"price": 49950.0, "amount": 1.0},
					{"price": 49940.0, "amount": 2.5}
				]
			}
		}`),
	}

	mockClient := &MockClient{response: mockResponse}
	orderBook := &OrderBook{Client: mockClient}

	response, err := orderBook.GetOrderBook("BTC_EUR", false)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	// Check asks
	if len(response.Data.Asks) != 2 {
		t.Errorf("Expected 2 asks, got %d", len(response.Data.Asks))
	}

	if response.Data.Asks[0].Price != 50050.0 {
		t.Errorf("Expected first ask price to be 50050.0, got %f", response.Data.Asks[0].Price)
	}

	if response.Data.Asks[0].Amount != 1.5 {
		t.Errorf("Expected first ask amount to be 1.5, got %f", response.Data.Asks[0].Amount)
	}

	if response.Data.Asks[1].Price != 50060.0 {
		t.Errorf("Expected second ask price to be 50060.0, got %f", response.Data.Asks[1].Price)
	}

	if response.Data.Asks[1].Amount != 2.0 {
		t.Errorf("Expected second ask amount to be 2.0, got %f", response.Data.Asks[1].Amount)
	}

	// Check bids
	if len(response.Data.Bids) != 2 {
		t.Errorf("Expected 2 bids, got %d", len(response.Data.Bids))
	}

	if response.Data.Bids[0].Price != 49950.0 {
		t.Errorf("Expected first bid price to be 49950.0, got %f", response.Data.Bids[0].Price)
	}

	if response.Data.Bids[0].Amount != 1.0 {
		t.Errorf("Expected first bid amount to be 1.0, got %f", response.Data.Bids[0].Amount)
	}

	if response.Data.Bids[1].Price != 49940.0 {
		t.Errorf("Expected second bid price to be 49940.0, got %f", response.Data.Bids[1].Price)
	}

	if response.Data.Bids[1].Amount != 2.5 {
		t.Errorf("Expected second bid amount to be 2.5, got %f", response.Data.Bids[1].Amount)
	}
}

func TestGetOrderBookWithGrouping(t *testing.T) {
	// Create mock response for grouped order book
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": {
				"asks": [
					{"price": 50000.0, "amount": 5.5}
				],
				"bids": [
					{"price": 49900.0, "amount": 3.5}
				]
			}
		}`),
	}

	mockClient := &MockClient{response: mockResponse}
	orderBook := &OrderBook{Client: mockClient}

	response, err := orderBook.GetOrderBook("BTC_EUR", true)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	// Check grouped asks
	if len(response.Data.Asks) != 1 {
		t.Errorf("Expected 1 grouped ask, got %d", len(response.Data.Asks))
	}

	if response.Data.Asks[0].Price != 50000.0 {
		t.Errorf("Expected grouped ask price to be 50000.0, got %f", response.Data.Asks[0].Price)
	}

	if response.Data.Asks[0].Amount != 5.5 {
		t.Errorf("Expected grouped ask amount to be 5.5, got %f", response.Data.Asks[0].Amount)
	}

	// Check grouped bids
	if len(response.Data.Bids) != 1 {
		t.Errorf("Expected 1 grouped bid, got %d", len(response.Data.Bids))
	}

	if response.Data.Bids[0].Price != 49900.0 {
		t.Errorf("Expected grouped bid price to be 49900.0, got %f", response.Data.Bids[0].Price)
	}

	if response.Data.Bids[0].Amount != 3.5 {
		t.Errorf("Expected grouped bid amount to be 3.5, got %f", response.Data.Bids[0].Amount)
	}
}

func TestGetOrderBookErrorResponse(t *testing.T) {
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
	orderBook := &OrderBook{Client: mockClient}

	response, err := orderBook.GetOrderBook("INVALID_PAIR", false)

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

func TestGetOrderBookHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte("Bad Request"),
	}

	mockClient := &MockClient{response: mockResponse}
	orderBook := &OrderBook{Client: mockClient}

	_, err := orderBook.GetOrderBook("BTC_EUR", false)

	if err == nil {
		t.Errorf("Expected error for non-200 response")
	}
}

func TestGetOrderBookNetworkError(t *testing.T) {
	// Create mock network error
	mockClient := &MockClient{err: &http.ProtocolError{}}
	orderBook := &OrderBook{Client: mockClient}

	_, err := orderBook.GetOrderBook("BTC_EUR", false)

	if err == nil {
		t.Error("Expected error for network failure")
	}
}

func TestGetOrderBookInvalidJSON(t *testing.T) {
	// Create mock response with invalid JSON
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       []byte("invalid json"),
	}

	mockClient := &MockClient{response: mockResponse}
	orderBook := &OrderBook{Client: mockClient}

	_, err := orderBook.GetOrderBook("BTC_EUR", false)

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestGetOrderBookEmptyCurrencyPair(t *testing.T) {
	mockClient := &MockClient{response: &coinmate.Response{}}
	orderBook := &OrderBook{Client: mockClient}

	_, err := orderBook.GetOrderBook("", false)

	if err == nil {
		t.Error("Expected error for empty currency pair")
	}
}

func TestOrderBookDataStructure(t *testing.T) {
	// Test that OrderBookData struct can be marshaled/unmarshaled correctly
	data := OrderBookData{
		Asks: []OrderBookAsksBids{
			{Price: 50050.0, Amount: 1.5},
			{Price: 50060.0, Amount: 2.0},
		},
		Bids: []OrderBookAsksBids{
			{Price: 49950.0, Amount: 1.0},
			{Price: 49940.0, Amount: 2.5},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal OrderBookData: %v", err)
	}

	var unmarshaledData OrderBookData
	err = json.Unmarshal(jsonData, &unmarshaledData)
	if err != nil {
		t.Errorf("Failed to unmarshal OrderBookData: %v", err)
	}

	// Check asks
	if len(unmarshaledData.Asks) != 2 {
		t.Errorf("Expected 2 asks, got %d", len(unmarshaledData.Asks))
	}

	if unmarshaledData.Asks[0].Price != data.Asks[0].Price {
		t.Errorf("Expected first ask price to be %f, got %f", data.Asks[0].Price, unmarshaledData.Asks[0].Price)
	}

	if unmarshaledData.Asks[0].Amount != data.Asks[0].Amount {
		t.Errorf("Expected first ask amount to be %f, got %f", data.Asks[0].Amount, unmarshaledData.Asks[0].Amount)
	}

	// Check bids
	if len(unmarshaledData.Bids) != 2 {
		t.Errorf("Expected 2 bids, got %d", len(unmarshaledData.Bids))
	}

	if unmarshaledData.Bids[0].Price != data.Bids[0].Price {
		t.Errorf("Expected first bid price to be %f, got %f", data.Bids[0].Price, unmarshaledData.Bids[0].Price)
	}

	if unmarshaledData.Bids[0].Amount != data.Bids[0].Amount {
		t.Errorf("Expected first bid amount to be %f, got %f", data.Bids[0].Amount, unmarshaledData.Bids[0].Amount)
	}
}

func TestOrderBookAsksBidsStructure(t *testing.T) {
	// Test that OrderBookAsksBids struct can be marshaled/unmarshaled correctly
	data := OrderBookAsksBids{
		Price:  50050.0,
		Amount: 1.5,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal OrderBookAsksBids: %v", err)
	}

	var unmarshaledData OrderBookAsksBids
	err = json.Unmarshal(jsonData, &unmarshaledData)
	if err != nil {
		t.Errorf("Failed to unmarshal OrderBookAsksBids: %v", err)
	}

	if unmarshaledData.Price != data.Price {
		t.Errorf("Expected price to be %f, got %f", data.Price, unmarshaledData.Price)
	}

	if unmarshaledData.Amount != data.Amount {
		t.Errorf("Expected amount to be %f, got %f", data.Amount, unmarshaledData.Amount)
	}
}
