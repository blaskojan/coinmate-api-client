package public

import (
	"encoding/json"
	"net/http"
	"testing"
	"tourGo/coinmate"
)

func TestGetTradingPairsSuccess(t *testing.T) {
	// Create mock response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": false,
			"errorMessage": "",
			"data": [
				{
					"name": "BTC_EUR",
					"firstCurrency": "BTC",
					"secondCurrency": "EUR",
					"priceDecimals": 2,
					"lotDecimals": 8,
					"minAmount": 0.001,
					"tradesWebSocketChannelId": "trades-BTC_EUR",
					"orderBookWebSocketChannelId": "orderBook-BTC_EUR",
					"tradeStatisticsWebSocketChannelId": "tradeStatistics-BTC_EUR"
				},
				{
					"name": "ETH_EUR",
					"firstCurrency": "ETH",
					"secondCurrency": "EUR",
					"priceDecimals": 2,
					"lotDecimals": 8,
					"minAmount": 0.01,
					"tradesWebSocketChannelId": "trades-ETH_EUR",
					"orderBookWebSocketChannelId": "orderBook-ETH_EUR",
					"tradeStatisticsWebSocketChannelId": "tradeStatistics-ETH_EUR"
				}
			]
		}`),
	}

	mockClient := &MockClient{response: mockResponse}
	tradingPairs := &TradingPairs{Client: mockClient}

	response, err := tradingPairs.GetTradingPairs()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if len(response.Data) != 2 {
		t.Errorf("Expected 2 trading pairs, got %d", len(response.Data))
	}

	// Check first trading pair
	btcEur := response.Data[0]
	if btcEur.Name != "BTC_EUR" {
		t.Errorf("Expected name to be 'BTC_EUR', got '%s'", btcEur.Name)
	}

	if btcEur.FirstCurrency != "BTC" {
		t.Errorf("Expected first currency to be 'BTC', got '%s'", btcEur.FirstCurrency)
	}

	if btcEur.SecondCurrency != "EUR" {
		t.Errorf("Expected second currency to be 'EUR', got '%s'", btcEur.SecondCurrency)
	}

	if btcEur.PriceDecimals != 2 {
		t.Errorf("Expected price decimals to be 2, got %d", btcEur.PriceDecimals)
	}

	if btcEur.LotDecimals != 8 {
		t.Errorf("Expected lot decimals to be 8, got %d", btcEur.LotDecimals)
	}

	if btcEur.MinAmount != 0.001 {
		t.Errorf("Expected min amount to be 0.001, got %f", btcEur.MinAmount)
	}

	if btcEur.TradesWebSocketChannelId != "trades-BTC_EUR" {
		t.Errorf("Expected trades channel ID to be 'trades-BTC_EUR', got '%s'", btcEur.TradesWebSocketChannelId)
	}

	if btcEur.OrderBookWebSocketChannelId != "orderBook-BTC_EUR" {
		t.Errorf("Expected order book channel ID to be 'orderBook-BTC_EUR', got '%s'", btcEur.OrderBookWebSocketChannelId)
	}

	if btcEur.TradeStatisticsWebSocketChannelId != "tradeStatistics-BTC_EUR" {
		t.Errorf("Expected trade statistics channel ID to be 'tradeStatistics-BTC_EUR', got '%s'", btcEur.TradeStatisticsWebSocketChannelId)
	}

	// Check second trading pair
	ethEur := response.Data[1]
	if ethEur.Name != "ETH_EUR" {
		t.Errorf("Expected name to be 'ETH_EUR', got '%s'", ethEur.Name)
	}

	if ethEur.FirstCurrency != "ETH" {
		t.Errorf("Expected first currency to be 'ETH', got '%s'", ethEur.FirstCurrency)
	}

	if ethEur.SecondCurrency != "EUR" {
		t.Errorf("Expected second currency to be 'EUR', got '%s'", ethEur.SecondCurrency)
	}

	if ethEur.PriceDecimals != 2 {
		t.Errorf("Expected price decimals to be 2, got %d", ethEur.PriceDecimals)
	}

	if ethEur.LotDecimals != 8 {
		t.Errorf("Expected lot decimals to be 8, got %d", ethEur.LotDecimals)
	}

	if ethEur.MinAmount != 0.01 {
		t.Errorf("Expected min amount to be 0.01, got %f", ethEur.MinAmount)
	}
}

func TestGetTradingPairsErrorResponse(t *testing.T) {
	// Create mock error response
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body: []byte(`{
			"error": true,
			"errorMessage": "Service temporarily unavailable",
			"data": []
		}`),
	}

	mockClient := &MockClient{response: mockResponse}
	tradingPairs := &TradingPairs{Client: mockClient}

	response, err := tradingPairs.GetTradingPairs()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !response.Error {
		t.Error("Expected error in response")
	}

	if response.ErrorMessage != "Service temporarily unavailable" {
		t.Errorf("Expected error message 'Service temporarily unavailable', got '%s'", response.ErrorMessage)
	}
}

func TestGetTradingPairsHTTPError(t *testing.T) {
	// Create mock HTTP error
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusInternalServerError,
		Status:     "500 Internal Server Error",
		Body:       []byte("Internal Server Error"),
	}

	mockClient := &MockClient{response: mockResponse}
	tradingPairs := &TradingPairs{Client: mockClient}

	response, err := tradingPairs.GetTradingPairs()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should return empty response when HTTP status is not OK
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestGetTradingPairsNetworkError(t *testing.T) {
	// Create mock network error
	mockClient := &MockClient{err: &http.ProtocolError{}}
	tradingPairs := &TradingPairs{Client: mockClient}

	response, err := tradingPairs.GetTradingPairs()

	if err == nil {
		t.Error("Expected error for network failure")
	}

	// Should return empty response on network error
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestGetTradingPairsInvalidJSON(t *testing.T) {
	// Create mock response with invalid JSON
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       []byte("invalid json"),
	}

	mockClient := &MockClient{response: mockResponse}
	tradingPairs := &TradingPairs{Client: mockClient}

	response, err := tradingPairs.GetTradingPairs()

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	// Should return empty response on JSON error
	if response.Error != false {
		t.Error("Expected default error state")
	}
}

func TestGetTradingPairsEmptyResponse(t *testing.T) {
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
	tradingPairs := &TradingPairs{Client: mockClient}

	response, err := tradingPairs.GetTradingPairs()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Error {
		t.Error("Expected no error in response")
	}

	if len(response.Data) != 0 {
		t.Errorf("Expected 0 trading pairs, got %d", len(response.Data))
	}
}

func TestTradingPairsDataStructure(t *testing.T) {
	// Test that TradingPairsData struct can be marshaled/unmarshaled correctly
	data := TradingPairsData{
		Name:                              "BTC_EUR",
		FirstCurrency:                     "BTC",
		SecondCurrency:                    "EUR",
		PriceDecimals:                     2,
		LotDecimals:                       8,
		MinAmount:                         0.001,
		TradesWebSocketChannelId:          "trades-BTC_EUR",
		OrderBookWebSocketChannelId:       "orderBook-BTC_EUR",
		TradeStatisticsWebSocketChannelId: "tradeStatistics-BTC_EUR",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal TradingPairsData: %v", err)
	}

	var unmarshaledData TradingPairsData
	err = json.Unmarshal(jsonData, &unmarshaledData)
	if err != nil {
		t.Errorf("Failed to unmarshal TradingPairsData: %v", err)
	}

	if unmarshaledData.Name != data.Name {
		t.Errorf("Expected name to be %s, got %s", data.Name, unmarshaledData.Name)
	}

	if unmarshaledData.FirstCurrency != data.FirstCurrency {
		t.Errorf("Expected first currency to be %s, got %s", data.FirstCurrency, unmarshaledData.FirstCurrency)
	}

	if unmarshaledData.SecondCurrency != data.SecondCurrency {
		t.Errorf("Expected second currency to be %s, got %s", data.SecondCurrency, unmarshaledData.SecondCurrency)
	}

	if unmarshaledData.PriceDecimals != data.PriceDecimals {
		t.Errorf("Expected price decimals to be %d, got %d", data.PriceDecimals, unmarshaledData.PriceDecimals)
	}

	if unmarshaledData.LotDecimals != data.LotDecimals {
		t.Errorf("Expected lot decimals to be %d, got %d", data.LotDecimals, unmarshaledData.LotDecimals)
	}

	if unmarshaledData.MinAmount != data.MinAmount {
		t.Errorf("Expected min amount to be %f, got %f", data.MinAmount, unmarshaledData.MinAmount)
	}

	if unmarshaledData.TradesWebSocketChannelId != data.TradesWebSocketChannelId {
		t.Errorf("Expected trades channel ID to be %s, got %s", data.TradesWebSocketChannelId, unmarshaledData.TradesWebSocketChannelId)
	}

	if unmarshaledData.OrderBookWebSocketChannelId != data.OrderBookWebSocketChannelId {
		t.Errorf("Expected order book channel ID to be %s, got %s", data.OrderBookWebSocketChannelId, unmarshaledData.OrderBookWebSocketChannelId)
	}

	if unmarshaledData.TradeStatisticsWebSocketChannelId != data.TradeStatisticsWebSocketChannelId {
		t.Errorf("Expected trade statistics channel ID to be %s, got %s", data.TradeStatisticsWebSocketChannelId, unmarshaledData.TradeStatisticsWebSocketChannelId)
	}
}

