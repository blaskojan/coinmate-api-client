package secure

import (
	"encoding/json"
	"net/http"
	"testing"
	"tourGo/coinmate"
)

func TestGetHistorySuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": [
				{
					"id": 12345,
					"timestamp": 1640995200,
					"type": "BUY",
					"price": 50000.0,
					"remainingAmount": 0.0,
					"originalAmount": 1.0,
					"status": "FILLED",
					"stopPrice": 0.0,
					"orderTradeType": "LIMIT",
					"hidden": false
				}
			]
		}`),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	response, err := order.GetHistory("BTC_EUR", 10)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 order, got %d", len(response.Data))
	}

	orderData := response.Data[0]
	if orderData.Id != 12345 {
		t.Errorf("Expected order ID to be 12345, got %d", orderData.Id)
	}
	if orderData.Type != "BUY" {
		t.Errorf("Expected order type to be 'BUY', got '%s'", orderData.Type)
	}
	if orderData.Price != 50000.0 {
		t.Errorf("Expected order price to be 50000.0, got %f", orderData.Price)
	}
}

func TestGetHistoryHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusUnauthorized,
		Status:     "401 Unauthorized",
		Body:       []byte("Unauthorized"),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	_, err := order.GetHistory("BTC_EUR", 10)

	if err == nil {
		t.Errorf("Expected error for non-200 response")
	}
}

func TestGetHistoryNetworkError(t *testing.T) {
	// Create mock network error
	mockClient := &MockSecureClient{err: &http.ProtocolError{}}
	order := &Order{Client: mockClient}

	_, err := order.GetHistory("BTC_EUR", 10)

	if err == nil {
		t.Error("Expected error for network failure")
	}
}

func TestGetOpenOrdersSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": [
				{
					"id": 12346,
					"timestamp": 1640995200,
					"type": "SELL",
					"currencyPair": "BTC_EUR",
					"price": 51000.0,
					"amount": 0.5,
					"orderTradeType": "LIMIT",
					"stopPrice": 0.0,
					"hidden": false
				}
			]
		}`),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	response, err := order.GetOpenOrders("BTC_EUR")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 open order, got %d", len(response.Data))
	}

	orderData := response.Data[0]
	if orderData.Id != 12346 {
		t.Errorf("Expected order ID to be 12346, got %d", orderData.Id)
	}
	if orderData.Type != "SELL" {
		t.Errorf("Expected order type to be 'SELL', got '%s'", orderData.Type)
	}
	if orderData.CurrencyPair != "BTC_EUR" {
		t.Errorf("Expected currency pair to be 'BTC_EUR', got '%s'", orderData.CurrencyPair)
	}
}

func TestGetOpenOrdersHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusUnauthorized,
		Status:     "401 Unauthorized",
		Body:       []byte("Unauthorized"),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	_, err := order.GetOpenOrders("BTC_EUR")

	if err == nil {
		t.Errorf("Expected error for non-200 response")
	}
}

func TestBuyLimitSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": 12347
		}`),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	response, err := order.BuyLimit(1.0, 50000.0, 0.0, "BTC_EUR", false, false, 0)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if response.OrderId != 12347 {
		t.Errorf("Expected order ID to be 12347, got %d", response.OrderId)
	}
}

func TestBuyLimitHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte("Bad Request"),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	_, err := order.BuyLimit(1.0, 50000.0, 0.0, "BTC_EUR", false, false, 0)

	if err == nil {
		t.Errorf("Expected error for non-200 response")
	}
}

func TestSellLimitSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": 12348
		}`),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	response, err := order.SellLimit(1.0, 51000.0, 0.0, "BTC_EUR", false, false, 0)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if response.OrderId != 12348 {
		t.Errorf("Expected order ID to be 12348, got %d", response.OrderId)
	}
}

func TestSellLimitHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte("Bad Request"),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	_, err := order.SellLimit(1.0, 51000.0, 0.0, "BTC_EUR", false, false, 0)

	if err == nil {
		t.Errorf("Expected error for non-200 response")
	}
}

func TestCancelOrderSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": true
		}`),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	response, err := order.CancelOrder(12345)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if !response.Data {
		t.Error("Expected cancel order to return true")
	}
}

func TestCancelOrderHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte("Bad Request"),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	_, err := order.CancelOrder(12345)

	if err == nil {
		t.Errorf("Expected error for non-200 response")
	}
}

func TestCancelOrderWithInfoSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": {
				"success": true,
				"remainingAmount": 0.0
			}
		}`),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	response, err := order.CancelOrderWithInfo(12345)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if !response.Data.Success {
		t.Error("Expected cancel order with info to return success true")
	}

	if response.Data.RemainingAmount != 0.0 {
		t.Errorf("Expected remaining amount to be 0.0, got %f", response.Data.RemainingAmount)
	}
}

func TestCancelOrderWithInfoHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte("Bad Request"),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	_, err := order.CancelOrderWithInfo(12345)

	if err == nil {
		t.Errorf("Expected error for non-200 response")
	}
}

func TestBuyInstantSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": 12349
		}`),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	response, err := order.BuyInstant(1000.0, "BTC_EUR", 0)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if response.OrderId != 12349 {
		t.Errorf("Expected order ID to be 12349, got %d", response.OrderId)
	}
}

func TestBuyInstantHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte("Bad Request"),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	_, err := order.BuyInstant(1000.0, "BTC_EUR", 0)

	if err == nil {
		t.Errorf("Expected error for non-200 response")
	}
}

func TestSellInstantSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": 12350
		}`),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	response, err := order.SellInstant(0.5, "BTC_EUR", 0)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if response.OrderId != 12350 {
		t.Errorf("Expected order ID to be 12350, got %d", response.OrderId)
	}
}

func TestSellInstantHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte("Bad Request"),
	}

	mockClient := &MockSecureClient{response: mockResponse}
	order := &Order{Client: mockClient}

	_, err := order.SellInstant(0.5, "BTC_EUR", 0)

	if err == nil {
		t.Errorf("Expected error for non-200 response")
	}
}

func TestOrderDataStructures(t *testing.T) {
	// Test OrderHistoryData marshaling/unmarshaling
	historyData := OrderHistoryData{
		Id:              12345,
		Timestamp:       1640995200,
		Type:            "BUY",
		Price:           50000.0,
		RemainingAmount: 0.0,
		OriginalAmount:  1.0,
		Status:          "FILLED",
		StopPrice:       0.0,
		OrderTradeType:  "LIMIT",
		Hidden:          false,
	}

	jsonData, err := json.Marshal(historyData)
	if err != nil {
		t.Errorf("Failed to marshal OrderHistoryData: %v", err)
	}

	var unmarshaledData OrderHistoryData
	err = json.Unmarshal(jsonData, &unmarshaledData)
	if err != nil {
		t.Errorf("Failed to unmarshal OrderHistoryData: %v", err)
	}

	if unmarshaledData.Id != historyData.Id {
		t.Errorf("Expected ID to be %d, got %d", historyData.Id, unmarshaledData.Id)
	}

	// Test OpenOrdersData marshaling/unmarshaling
	openOrderData := OpenOrdersData{
		Id:             12346,
		Timestamp:      1640995200,
		Type:           "SELL",
		CurrencyPair:   "BTC_EUR",
		Price:          51000.0,
		Amount:         0.5,
		OrderTradeType: "LIMIT",
		StopPrice:      0.0,
		Hidden:         false,
	}

	jsonData, err = json.Marshal(openOrderData)
	if err != nil {
		t.Errorf("Failed to marshal OpenOrdersData: %v", err)
	}

	var unmarshaledOpenData OpenOrdersData
	err = json.Unmarshal(jsonData, &unmarshaledOpenData)
	if err != nil {
		t.Errorf("Failed to unmarshal OpenOrdersData: %v", err)
	}

	if unmarshaledOpenData.Id != openOrderData.Id {
		t.Errorf("Expected ID to be %d, got %d", openOrderData.Id, unmarshaledOpenData.Id)
	}
}
