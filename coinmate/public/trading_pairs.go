package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tourGo/coinmate"
)

const (
	tradingPairsEndpoint = "/tradingPairs"
)

type TradingPairs struct {
	Client coinmate.ClientInterface
}

type TradingPairsResponse struct {
	Error        bool
	ErrorMessage string
	Data         []TradingPairsData
}

type TradingPairsData struct {
	Name                              string  `json:"name"`
	FirstCurrency                     string  `json:"firstCurrency"`
	SecondCurrency                    string  `json:"secondCurrency"`
	PriceDecimals                     uint64  `json:"priceDecimals"`
	LotDecimals                       uint64  `json:"lotDecimals"`
	MinAmount                         float64 `json:"minAmount"`
	TradesWebSocketChannelId          string  `json:"tradesWebSocketChannelId"`
	OrderBookWebSocketChannelId       string  `json:"orderBookWebSocketChannelId"`
	TradeStatisticsWebSocketChannelId string  `json:"tradeStatisticsWebSocketChannelId"`
}

// Trading pairs endpoint
func (t *TradingPairs) GetTradingPairs() (TradingPairsResponse, error) {
	tpr := TradingPairsResponse{}

	r := coinmate.Request{
		HTTPMethod: http.MethodGet,
		URL:        t.Client.GetBaseUrl() + tradingPairsEndpoint,
		Body:       nil,
	}
	response, err := t.Client.MakePublicRequest(r)
	if err != nil {
		return tpr, fmt.Errorf("trading pairs request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return tpr, fmt.Errorf("trading pairs request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	err = json.Unmarshal(response.Body, &tpr)
	if err != nil {
		return tpr, fmt.Errorf("failed to decode trading pairs response: %w", err)
	}

	return tpr, err
}
