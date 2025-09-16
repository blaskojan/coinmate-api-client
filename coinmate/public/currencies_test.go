package public

import (
	"net/http"
	"testing"
	"tourGo/coinmate"
)

func TestGetCurrenciesSuccess(t *testing.T) {
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       []byte(`{"error":false,"errorMessage":"","data":["BTC","ETH","EUR"]}`),
	}
	mockClient := &MockClient{response: mockResponse}
	c := &Currencies{Client: mockClient}

	resp, err := c.GetCurrencies()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Error || len(resp.Data) != 3 || resp.Data[0] != "BTC" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetCurrenciesHTTPError(t *testing.T) {
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusInternalServerError,
		Status:     "500 Internal Server Error",
		Body:       []byte("boom"),
	}
	mockClient := &MockClient{response: mockResponse}
	c := &Currencies{Client: mockClient}

	_, err := c.GetCurrencies()
	if err == nil {
		t.Fatalf("expected error")
	}
}
