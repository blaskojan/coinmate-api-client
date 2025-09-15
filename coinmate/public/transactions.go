package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"tourGo/coinmate"
)

type Transactions struct {
	Client coinmate.ClientInterface
}

const (
	transactionsEndpoint        = "/transactions"
	minutesIntoHistoryParamName = "minutesIntoHistory"
)

// Ticker endpoint response
type TransactionsResponse struct {
	Error        bool               `json:"error"`
	ErrorMessage string             `json:"errorMessage"`
	Data         []TransactionsData `json:"data"`
}

type TransactionsData struct {
	Timestamp     int64   `json:"timestamp"`
	TransactionId string  `json:"transactionId"`
	Price         float64 `json:"price"`
	Amount        float64 `json:"amount"`
	CurrencyPair  string  `json:"currencyPair"`
	TradeType     string  `json:"tradeType"`
}

// Transactions
func (t *Transactions) GetTransactions(currencyPair string, minutesIntoHistory uint64) (TransactionsResponse, error) {
	transactionsResponse := TransactionsResponse{}

	if currencyPair == "" {
		return transactionsResponse, fmt.Errorf("currencyPair must not be empty")
	}

	// URL compose
	u, err := url.Parse(t.Client.GetBaseUrl() + transactionsEndpoint)
	if err != nil {
		return transactionsResponse, fmt.Errorf("failed to parse transactions URL: %w", err)
	}
	q := u.Query()
	q.Set(currencyPairParamName, currencyPair)
	q.Set(minutesIntoHistoryParamName, strconv.FormatUint(minutesIntoHistory, 10))
	u.RawQuery = q.Encode()

	r := coinmate.Request{
		HTTPMethod: http.MethodGet,
		URL:        u.String(),
		Body:       nil,
	}
	response, err := t.Client.MakePublicRequest(r)
	if err != nil {
		return transactionsResponse, fmt.Errorf("transactions request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return transactionsResponse, fmt.Errorf("transactions request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	err = json.Unmarshal(response.Body, &transactionsResponse)
	if err != nil {
		return transactionsResponse, fmt.Errorf("failed to decode transactions response: %w", err)
	}

	return transactionsResponse, err
}
