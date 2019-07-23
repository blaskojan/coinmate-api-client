package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"tourGo/coinmate"
)

const (
	orderBookEndpoint          = "/orderBook"
	currencyPairParamName      = "currencyPair"
	groupByPriceLimitParamName = "groupByPriceLimit"
)

type OrderBook struct {
	Client coinmate.ClientInterface
}

// Order book response
type OrderBookResponse struct {
	Error        bool
	ErrorMessage string
	Data         OrderBookData
}

type OrderBookData struct {
	Asks []OrderBookAsksBids `json:"asks"`
	Bids []OrderBookAsksBids `json:"bids"`
}

type OrderBookAsksBids struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

// Order book endpoint
func (o *OrderBook) GetOrderBook(currencyPair string, groupByPriceLimit bool) (OrderBookResponse, error) {
	orderBookResponse := OrderBookResponse{}

	// URL compose
	u, err := url.Parse(o.Client.GetBaseUrl() + orderBookEndpoint)
	if err != nil {
		return orderBookResponse, err
	}
	q := u.Query()
	q.Set(currencyPairParamName, currencyPair)
	q.Set(groupByPriceLimitParamName, strings.Title(strconv.FormatBool(groupByPriceLimit)))
	u.RawQuery = q.Encode()

	r := coinmate.Request{
		HTTPMethod: http.MethodGet,
		URL:        u.String(),
		Body:       nil,
	}
	response, err := o.Client.MakePublicRequest(r)
	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Println("Coinmate error: " + string(response.Body))
		return orderBookResponse, err
	}

	err = json.Unmarshal(response.Body, &orderBookResponse)
	if err != nil {
		fmt.Println(err)
		return orderBookResponse, err
	}

	return orderBookResponse, err

}
