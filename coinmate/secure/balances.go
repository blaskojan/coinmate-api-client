package secure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tourGo/coinmate"
)

const endpoint = "/balances"

type Balances struct {
	Client coinmate.ClientInterface
}

// Order book response
type BalancesResponse struct {
	Error        bool                       `json:"error"`
	ErrorMessage string                     `json:"errorMessage"`
	Data         map[string]BalanceCurrency `json:"data"`
}

// Balance currency data
type BalanceCurrency struct {
	Currency  string  `json:"currency"`
	Balance   float32 `json:"balance"`
	Reserved  float32 `json:"reserved"`
	Available float32 `json:"available"`
}

// Balances endpoint
func (b *Balances) GetBalances() (BalancesResponse, error) {
	balancesResponse := BalancesResponse{}

	ap := map[string]string{}

	r := coinmate.Request{
		HTTPMethod: http.MethodPost,
		URL:        b.Client.GetBaseUrl() + endpoint,
		Body:       b.Client.GetRequestBody(ap),
	}
	response, err := b.Client.MakeSecureRequest(r)
	if err != nil {
		return balancesResponse, fmt.Errorf("balances request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return balancesResponse, fmt.Errorf("balances request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	err = json.Unmarshal(response.Body, &balancesResponse)
	if err != nil {
		return balancesResponse, fmt.Errorf("failed to decode balances response: %w", err)
	}

	return balancesResponse, err
}
