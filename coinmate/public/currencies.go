package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tourGo/coinmate"
)

type Currencies struct {
	Client coinmate.ClientInterface
}

type CurrenciesResponse struct {
	Error        bool     `json:"error"`
	ErrorMessage string   `json:"errorMessage"`
	Data         []string `json:"data"`
}

func (c *Currencies) GetCurrencies() (CurrenciesResponse, error) {
	cr := CurrenciesResponse{}

	r := coinmate.Request{
		HTTPMethod: http.MethodGet,
		URL:        c.Client.GetBaseUrl() + "/currencies",
		Body:       nil,
	}
	response, err := c.Client.MakePublicRequest(r)
	if err != nil {
		return cr, fmt.Errorf("currencies request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return cr, fmt.Errorf("currencies request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	if err := json.Unmarshal(response.Body, &cr); err != nil {
		return cr, fmt.Errorf("failed to decode currencies response: %w", err)
	}
	return cr, nil
}
