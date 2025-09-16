package public

import (
	"net/http"
	"testing"
	"tourGo/coinmate"
)

func TestGetCurrencyPairsSuccess(t *testing.T) {
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       []byte(`{"error":false,"errorMessage":"","data":[{"name":"BTC_EUR","firstCurrency":"BTC","secondCurrency":"EUR"}]}`),
	}
	mockClient := &MockClient{response: mockResponse}
	cp := &CurrencyPairs{Client: mockClient}

	resp, err := cp.GetCurrencyPairs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Error || len(resp.Data) != 1 || resp.Data[0].Name != "BTC_EUR" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetCurrencyPairsHTTPError(t *testing.T) {
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadGateway,
		Status:     "502 Bad Gateway",
		Body:       []byte("oops"),
	}
	mockClient := &MockClient{response: mockResponse}
	cp := &CurrencyPairs{Client: mockClient}

	_, err := cp.GetCurrencyPairs()
	if err == nil {
		t.Fatalf("expected error")
	}
}
