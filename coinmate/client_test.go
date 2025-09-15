package coinmate

import (
	"testing"
	"time"
)

func TestGetCoinmateClient(t *testing.T) {
	clientId := "test-client-id"
	apiKey := "test-api-key"
	privateKey := "test-private-key"

	client := GetCoinmateClient(clientId, apiKey, privateKey)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.ClientID != clientId {
		t.Errorf("Expected ClientID to be %s, got %s", clientId, client.ClientID)
	}

	if client.ApiKey != apiKey {
		t.Errorf("Expected ApiKey to be %s, got %s", apiKey, client.ApiKey)
	}

	if client.PrivateKey != privateKey {
		t.Errorf("Expected PrivateKey to be %s, got %s", privateKey, client.PrivateKey)
	}
}

func TestGetNonce(t *testing.T) {
	client := GetCoinmateClient("test", "test", "test")

	nonce1 := client.GetNonce()
	nonce2 := client.GetNonce()

	if nonce1 == "" {
		t.Error("Expected nonce to be non-empty")
	}

	if nonce2 == "" {
		t.Error("Expected nonce to be non-empty")
	}

	// Nonces should be different (time-based)
	if nonce1 == nonce2 {
		t.Error("Expected nonces to be different")
	}
}

func TestGetSignature(t *testing.T) {
	clientId := "test-client-id"
	apiKey := "test-api-key"
	privateKey := "test-private-key"
	nonce := "1234567890"

	client := GetCoinmateClient(clientId, apiKey, privateKey)
	signature := client.GetSignature(clientId, apiKey, nonce, privateKey)

	if signature == "" {
		t.Error("Expected signature to be non-empty")
	}

	// Signature should be consistent for same inputs
	signature2 := client.GetSignature(clientId, apiKey, nonce, privateKey)
	if signature != signature2 {
		t.Error("Expected signatures to be consistent for same inputs")
	}
}

func TestGetBaseUrl(t *testing.T) {
	client := GetCoinmateClient("test", "test", "test")

	baseUrl := client.GetBaseUrl()
	expectedUrl := "https://coinmate.io/api"

	if baseUrl != expectedUrl {
		t.Errorf("Expected base URL to be %s, got %s", expectedUrl, baseUrl)
	}
}

func TestGetRequestBody(t *testing.T) {
	client := GetCoinmateClient("test", "test", "test")

	additionalParams := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	body := client.GetRequestBody(additionalParams)

	if len(body) == 0 {
		t.Error("Expected request body to be non-empty")
	}

	bodyStr := string(body)

	// Should contain required fields
	if !contains(bodyStr, "clientId=test") {
		t.Error("Expected body to contain clientId")
	}

	if !contains(bodyStr, "publicKey=test") {
		t.Error("Expected body to contain publicKey")
	}

	if !contains(bodyStr, "nonce=") {
		t.Error("Expected body to contain nonce")
	}

	if !contains(bodyStr, "signature=") {
		t.Error("Expected body to contain signature")
	}

	// Should contain additional parameters
	if !contains(bodyStr, "param1=value1") {
		t.Error("Expected body to contain additional param1")
	}

	if !contains(bodyStr, "param2=value2") {
		t.Error("Expected body to contain additional param2")
	}
}

func TestGetRequestBodyEmpty(t *testing.T) {
	client := GetCoinmateClient("test", "test", "test")

	body := client.GetRequestBody(nil)

	if len(body) == 0 {
		t.Error("Expected request body to be non-empty even with nil params")
	}

	bodyStr := string(body)

	// Should contain required fields even with empty params
	if !contains(bodyStr, "clientId=test") {
		t.Error("Expected body to contain clientId")
	}

	if !contains(bodyStr, "publicKey=test") {
		t.Error("Expected body to contain publicKey")
	}
}

func TestRequestTimeout(t *testing.T) {
	client := GetCoinmateClient("test", "test", "test")

	// The timeout should be 2 seconds as defined in the constants
	expectedTimeout := 2 * time.Second

	// We can't directly test the timeout, but we can verify the client was created
	if client.httpClient.Timeout != expectedTimeout {
		t.Errorf("Expected timeout to be %v, got %v", expectedTimeout, client.httpClient.Timeout)
	}
}

func TestSetTimeout(t *testing.T) {
	client := GetCoinmateClient("test", "test", "test")

	// default is 2s
	if client.httpClient.Timeout != 2*time.Second {
		t.Fatalf("expected default timeout 2s, got %v", client.httpClient.Timeout)
	}

	client.SetTimeout(5 * time.Second)
	if client.httpClient.Timeout != 5*time.Second {
		t.Fatalf("expected timeout 5s after SetTimeout, got %v", client.httpClient.Timeout)
	}

	// ignore non-positive values
	client.SetTimeout(0)
	if client.httpClient.Timeout != 5*time.Second {
		t.Fatalf("expected timeout unchanged at 5s when setting 0, got %v", client.httpClient.Timeout)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
