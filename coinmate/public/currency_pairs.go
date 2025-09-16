package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tourGo/coinmate"
)

type CurrencyPairs struct {
	Client coinmate.ClientInterface
}

type CurrencyPair struct {
	Name           string `json:"name"`
	FirstCurrency  string `json:"firstCurrency"`
	SecondCurrency string `json:"secondCurrency"`
}

type CurrencyPairsResponse struct {
	Error        bool           `json:"error"`
	ErrorMessage string         `json:"errorMessage"`
	Data         []CurrencyPair `json:"data"`
}

func (c *CurrencyPairs) GetCurrencyPairs() (CurrencyPairsResponse, error) {
	resp := CurrencyPairsResponse{}

	r := coinmate.Request{
		HTTPMethod: http.MethodGet,
		URL:        c.Client.GetBaseUrl() + "/currency-pairs",
		Body:       nil,
	}
	response, err := c.Client.MakePublicRequest(r)
	if err != nil {
		return resp, fmt.Errorf("currency-pairs request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("currency-pairs request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	if err := json.Unmarshal(response.Body, &resp); err != nil {
		return resp, fmt.Errorf("failed to decode currency-pairs response: %w", err)
	}
	return resp, nil
}
