package public

import (
	"encoding/json"
	"net/http"
	"testing"
	"tourGo/coinmate"
)

func TestGetTransactionsSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": [
				{
					"timestamp": 1640995200,
					"transactionId": "tx123456",
					"price": 50000.0,
					"amount": 1.5,
					"currencyPair": "BTC_EUR",
					"tradeType": "BUY"
				},
				{
					"timestamp": 1640995260,
					"transactionId": "tx123457",
					"price": 50010.0,
					"amount": 0.5,
					"currencyPair": "BTC_EUR",
					"tradeType": "SELL"
				}
			]
		}`),
	}

	mockClient := &MockClient{response: mockResponse}
	transactions := &Transactions{Client: mockClient}

	response, err := transactions.GetTransactions("BTC_EUR", 60)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(response.Data))
	}

	// Check first transaction
	tx1 := response.Data[0]
	if tx1.Timestamp != 1640995200 {
		t.Errorf("Expected timestamp to be 1640995200, got %d", tx1.Timestamp)
	}

	if tx1.TransactionId != "tx123456" {
		t.Errorf("Expected transaction ID to be 'tx123456', got '%s'", tx1.TransactionId)
	}

	if tx1.Price != 50000.0 {
		t.Errorf("Expected price to be 50000.0, got %f", tx1.Price)
	}

	if tx1.Amount != 1.5 {
		t.Errorf("Expected amount to be 1.5, got %f", tx1.Amount)
	}

	if tx1.CurrencyPair != "BTC_EUR" {
		t.Errorf("Expected currency pair to be 'BTC_EUR', got '%s'", tx1.CurrencyPair)
	}

	if tx1.TradeType != "BUY" {
		t.Errorf("Expected trade type to be 'BUY', got '%s'", tx1.TradeType)
	}

	// Check second transaction
	tx2 := response.Data[1]
	if tx2.Timestamp != 1640995260 {
		t.Errorf("Expected timestamp to be 1640995260, got %d", tx2.Timestamp)
	}

	if tx2.TransactionId != "tx123457" {
		t.Errorf("Expected transaction ID to be 'tx123457', got '%s'", tx2.TransactionId)
	}

	if tx2.Price != 50010.0 {
		t.Errorf("Expected price to be 50010.0, got %f", tx2.Price)
	}

	if tx2.Amount != 0.5 {
		t.Errorf("Expected amount to be 0.5, got %f", tx2.Amount)
	}

	if tx2.CurrencyPair != "BTC_EUR" {
		t.Errorf("Expected currency pair to be 'BTC_EUR', got '%s'", tx2.CurrencyPair)
	}

	if tx2.TradeType != "SELL" {
		t.Errorf("Expected trade type to be 'SELL', got '%s'", tx2.TradeType)
	}
}

func TestGetTransactionsEmptyResponse(t *testing.T) {
	// Create mock response with empty data
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": []
		}`),
	}

	mockClient := &MockClient{response: mockResponse}
	transactions := &Transactions{Client: mockClient}

	response, err := transactions.GetTransactions("BTC_EUR", 60)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if len(response.Data) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(response.Data))
	}
}

func TestGetTransactionsErrorResponse(t *testing.T) {
	// Create mock error response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": true,
			"errorMessage": "Invalid currency pair",
			"data": []
		}`),
	}

	mockClient := &MockClient{response: mockResponse}
	transactions := &Transactions{Client: mockClient}

	response, err := transactions.GetTransactions("INVALID_PAIR", 60)

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

func TestGetTransactionsHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte("Bad Request"),
	}

	mockClient := &MockClient{response: mockResponse}
	transactions := &Transactions{Client: mockClient}

	response, err := transactions.GetTransactions("BTC_EUR", 60)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should return empty response when HTTP status is not OK
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestGetTransactionsNetworkError(t *testing.T) {
	// Create mock network error
	mockClient := &MockClient{err: &http.ProtocolError{}}
	transactions := &Transactions{Client: mockClient}

	response, err := transactions.GetTransactions("BTC_EUR", 60)

	if err == nil {
		t.Error("Expected error for network failure")
	}

	// Should return empty response on network error
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestGetTransactionsInvalidJSON(t *testing.T) {
	// Create mock response with invalid JSON
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       []byte("invalid json"),
	}

	mockClient := &MockClient{response: mockResponse}
	transactions := &Transactions{Client: mockClient}

	response, err := transactions.GetTransactions("BTC_EUR", 60)

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	// Should return empty response on JSON error
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestGetTransactionsDifferentTimeRanges(t *testing.T) {
	// Test different time ranges
	testCases := []struct {
		currencyPair string
		minutes      uint64
		description  string
	}{
		{"BTC_EUR", 5, "5 minutes"},
		{"ETH_EUR", 15, "15 minutes"},
		{"LTC_EUR", 60, "60 minutes"},
		{"BTC_EUR", 1440, "24 hours"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mockResponse := &coinmate.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body: []byte(`{
					"error": false,
					"errorMessage": "",
					"data": []
				}`),
			}

			mockClient := &MockClient{response: mockResponse}
			transactions := &Transactions{Client: mockClient}

			response, err := transactions.GetTransactions(tc.currencyPair, tc.minutes)

			if err != nil {
				t.Errorf("Expected no error for %s, got %v", tc.description, err)
			}

			if response.Error {
				t.Errorf("Expected no error in response for %s", tc.description)
			}
		})
	}
}

func TestGetTransactionsEmptyCurrencyPair(t *testing.T) {
	mockClient := &MockClient{response: &coinmate.Response{}}
	transactions := &Transactions{Client: mockClient}

	response, err := transactions.GetTransactions("", 60)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should handle empty currency pair gracefully
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestTransactionsDataStructure(t *testing.T) {
	// Test that TransactionsData struct can be marshaled/unmarshaled correctly
	data := TransactionsData{
		Timestamp:     1640995200,
		TransactionId: "tx123456",
		Price:         50000.0,
		Amount:        1.5,
		CurrencyPair:  "BTC_EUR",
		TradeType:     "BUY",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal TransactionsData: %v", err)
	}

	var unmarshaledData TransactionsData
	err = json.Unmarshal(jsonData, &unmarshaledData)
	if err != nil {
		t.Errorf("Failed to unmarshal TransactionsData: %v", err)
	}

	if unmarshaledData.Timestamp != data.Timestamp {
		t.Errorf("Expected timestamp to be %d, got %d", data.Timestamp, unmarshaledData.Timestamp)
	}

	if unmarshaledData.TransactionId != data.TransactionId {
		t.Errorf("Expected transaction ID to be %s, got %s", data.TransactionId, unmarshaledData.TransactionId)
	}

	if unmarshaledData.Price != data.Price {
		t.Errorf("Expected price to be %f, got %f", data.Price, unmarshaledData.Price)
	}

	if unmarshaledData.Amount != data.Amount {
		t.Errorf("Expected amount to be %f, got %f", data.Amount, unmarshaledData.Amount)
	}

	if unmarshaledData.CurrencyPair != data.CurrencyPair {
		t.Errorf("Expected currency pair to be %s, got %s", data.CurrencyPair, unmarshaledData.CurrencyPair)
	}

	if unmarshaledData.TradeType != data.TradeType {
		t.Errorf("Expected trade type to be %s, got %s", data.TradeType, unmarshaledData.TradeType)
	}
}

