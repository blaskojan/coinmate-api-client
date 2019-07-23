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

	r := coinmate.Request{
		HTTPMethod: http.MethodGet,
		URL:        t.Client.GetBaseUrl() + "/ticker?currencyPair=" + currencyPair,
		Body:       nil,
	}
	response, err := t.Client.MakePublicRequest(r)
	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println("Coinmate error: " + string(response.Body))
		return tickerResponse, err
	}

	err = json.Unmarshal(response.Body, &tickerResponse)
	if err != nil {
		fmt.Println(err)
		return tickerResponse, err
	}

	return tickerResponse, err
}
