package public

import (
	"net/http"
	"testing"
	"tourGo/coinmate"
)

func TestGetTickerAllSuccess(t *testing.T) {
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       []byte(`{"error":false,"errorMessage":"","data":{"BTC_EUR":{"last":50000,"high":51000,"low":49000,"bid":49950,"ask":50050,"change":500}}}`),
	}
	mockClient := &MockClient{response: mockResponse}
	ta := &TickerAll{Client: mockClient}

	resp, err := ta.GetTickerAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Error || resp.Data["BTC_EUR"].Last != 50000 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetTickerAllHTTPError(t *testing.T) {
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Body:       []byte("bad"),
	}
	mockClient := &MockClient{response: mockResponse}
	ta := &TickerAll{Client: mockClient}

	_, err := ta.GetTickerAll()
	if err == nil {
		t.Fatalf("expected error")
	}
}
