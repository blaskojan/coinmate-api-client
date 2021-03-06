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
	Error        bool
	ErrorMessage string
	Data         map[string]BalanceCurrency
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
	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println("Coinmate error: " + string(response.Body))
		return balancesResponse, err
	}

	err = json.Unmarshal(response.Body, &balancesResponse)
	if err != nil {
		fmt.Println(err)
		return balancesResponse, err
	}

	return balancesResponse, err
}
