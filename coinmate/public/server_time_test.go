package public

import (
	"net/http"
	"testing"
	"tourGo/coinmate"
)

func TestGetServerTimeSuccess(t *testing.T) {
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       []byte(`{"error":false,"errorMessage":"","data":1700000000}`),
	}
	mockClient := &MockClient{response: mockResponse}
	s := &ServerTime{Client: mockClient}

	resp, err := s.GetServerTime()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Error || resp.Data != 1700000000 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetServerTimeHTTPError(t *testing.T) {
	mockResponse := &coinmate.Response{
		StatusCode: http.StatusBadGateway,
		Status:     "502 Bad Gateway",
		Body:       []byte("down"),
	}
	mockClient := &MockClient{response: mockResponse}
	s := &ServerTime{Client: mockClient}

	_, err := s.GetServerTime()
	if err == nil {
		t.Fatalf("expected error")
	}
}
