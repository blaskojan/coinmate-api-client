package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tourGo/coinmate"
)

type Ticker struct {
	Client coinmate.ClientInterface
}

// Ticker endpoint response
type TickerResponse struct {
	Error        bool       `json:"error"`
	ErrorMessage string     `json:"errorMessage"`
	Data         TickerData `json:"data"`
}

type TickerData struct {
	Last      float64 `json:"last"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Amount    float64 `json:"amount"`
	Bid       float64 `json:"bid"`
	Ask       float64 `json:"ask"`
	Change    float64 `json:"change"`
	Open      float64 `json:"open"`
	Timestamp uint64  `json:"timestamp"`
}

// Ticker endpoint
func (t *Ticker) GetTicker(currencyPair string) (TickerResponse, error) {
	tickerResponse := TickerResponse{}

	if currencyPair == "" {
		return tickerResponse, fmt.Errorf("currencyPair must not be empty")
	}

	r := coinmate.Request{
		HTTPMethod: http.MethodGet,
		URL:        t.Client.GetBaseUrl() + "/ticker?currencyPair=" + currencyPair,
		Body:       nil,
	}
	response, err := t.Client.MakePublicRequest(r)
	if err != nil {
		return tickerResponse, fmt.Errorf("ticker request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return tickerResponse, fmt.Errorf("ticker request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	err = json.Unmarshal(response.Body, &tickerResponse)
	if err != nil {
		return tickerResponse, fmt.Errorf("failed to decode ticker response: %w", err)
	}

	return tickerResponse, err
}
