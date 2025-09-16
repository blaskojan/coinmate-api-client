package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tourGo/coinmate"
)

type TickerAll struct {
	Client coinmate.ClientInterface
}

type TickerAllItem struct {
	Last   float64 `json:"last"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Bid    float64 `json:"bid"`
	Ask    float64 `json:"ask"`
	Change float64 `json:"change"`
}

type TickerAllResponse struct {
	Error        bool                     `json:"error"`
	ErrorMessage string                   `json:"errorMessage"`
	Data         map[string]TickerAllItem `json:"data"`
}

func (t *TickerAll) GetTickerAll() (TickerAllResponse, error) {
	tr := TickerAllResponse{}

	r := coinmate.Request{
		HTTPMethod: http.MethodGet,
		URL:        t.Client.GetBaseUrl() + "/ticker-all",
		Body:       nil,
	}
	response, err := t.Client.MakePublicRequest(r)
	if err != nil {
		return tr, fmt.Errorf("ticker-all request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return tr, fmt.Errorf("ticker-all request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	if err := json.Unmarshal(response.Body, &tr); err != nil {
		return tr, fmt.Errorf("failed to decode ticker-all response: %w", err)
	}
	return tr, nil
}
