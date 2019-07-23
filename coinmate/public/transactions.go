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

	// URL compose
	u, err := url.Parse(t.Client.GetBaseUrl() + transactionsEndpoint)
	if err != nil {
		return transactionsResponse, err
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
	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println("Coinmate error: " + string(response.Body))
		return transactionsResponse, err
	}

	err = json.Unmarshal(response.Body, &transactionsResponse)
	if err != nil {
		fmt.Println(err)
		return transactionsResponse, err
	}

	return transactionsResponse, err
}
