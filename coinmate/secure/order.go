package secure

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
	orderHistoryEndpoint        = "/orderHistory"
	openOrdersEndpoint          = "/openOrders"
	cancelOrderEndpoint         = "/cancelOrder"
	cancelOrderWithInfoEndpoint = "/cancelOrderWithInfo"
	buyLimitOrderEndpoint       = "/buyLimit"
	sellLimitOrderEndpoint      = "/sellLimit"
	buyInstantOrderEndpoint     = "/buyInstant"
	sellInstantOrderEndpoint    = "/sellInstant"
	currencyPairParamName       = "currencyPair"
	limitReturnedOrders         = "limit"
	orderIdParamName            = "orderId"
	amountParamName             = "amount"
	priceParamName              = "price"
	stopPriceParamName          = "stopPrice"
	hiddenParamName             = "hidden"
	immediateOrCancelParamName  = "immediateOrCancel"
	clientOrderIdParamName      = "clientOrderId"
	totalParamName              = "total"
)

type Order struct {
	Client coinmate.ClientInterface
}

// Order history response
type OrderHistoryResponse struct {
	Error        bool               `json:"error"`
	ErrorMessage string             `json:"errorMessage"`
	Data         []OrderHistoryData `json:"data"`
}

// Order history data
type OrderHistoryData struct {
	Id              uint64  `json:"id"`
	Timestamp       int64   `json:"timestamp"`
	Type            string  `json:"type"`
	Price           float64 `json:"price"`
	RemainingAmount float64 `json:"remainingAmount"`
	OriginalAmount  float64 `json:"originalAmount"`
	Status          string  `json:"status"`
	StopPrice       float64 `json:"stopPrice"`
	OrderTradeType  string  `json:"orderTradeType"`
	Hidden          bool    `json:"hidden"`
}

// Open orders history response
type OpenOrdersResponse struct {
	Error        bool             `json:"error"`
	ErrorMessage string           `json:"errorMessage"`
	Data         []OpenOrdersData `json:"data"`
}

// Open orders data
type OpenOrdersData struct {
	Id             uint64  `json:"id"`
	Timestamp      int64   `json:"timestamp"`
	Type           string  `json:"type"`
	CurrencyPair   string  `json:"currencyPair"`
	Price          float64 `json:"price"`
	Amount         float64 `json:"amount"`
	OrderTradeType string  `json:"orderTradeType"`
	StopPrice      float64 `json:"stopPrice"`
	Hidden         bool    `json:"hidden"`
}

// Cancel order
type CancelOrderResponse struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	Data         bool   `json:"data"`
}

// Cancel order with info response
type CancelOrderWithInfoResponse struct {
	Error        bool                    `json:"error"`
	ErrorMessage string                  `json:"errorMessage"`
	Data         CancelOrderWithInfoData `json:"data"`
}

// Cancel order info data
type CancelOrderWithInfoData struct {
	Success         bool    `json:"success"`
	RemainingAmount float64 `json:"remainingAmount"`
}

// Buy limit response
type BuyLimitResponse struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	Data         uint64 `json:"data"`
}

type BuyLimit struct {
	Error        bool
	ErrorMessage string
	OrderId      uint64
}

// Sell limit response
type SellLimitResponse struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	Data         uint64 `json:"data"`
}

type SellLimit struct {
	Error        bool
	ErrorMessage string
	OrderId      uint64
}

// Buy instant response
type BuyInstantResponse struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	Data         uint64 `json:"data"`
}

type BuySell struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	Data         uint64 `json:"data"`
}

type BuyAndSellResponse struct {
	Error        bool
	ErrorMessage string
	OrderId      uint64
}

// Order history
func (o *Order) GetHistory(currencyPair string, limit int64) (OrderHistoryResponse, error) {
	orderHistoryResponse := OrderHistoryResponse{}

	// URL compose
	u, _ := url.Parse(o.Client.GetBaseUrl() + orderHistoryEndpoint)
	q := u.Query()

	// Adding optional GET params
	if currencyPair != "" {
		q.Set(currencyPairParamName, currencyPair)
	}
	if limit != 0 {
		q.Set(limitReturnedOrders, strconv.FormatInt(limit, 10))
	}

	u.RawQuery = q.Encode()

	ap := map[string]string{}

	r := coinmate.Request{
		HTTPMethod: http.MethodPost,
		URL:        u.String(),
		Body:       o.Client.GetRequestBody(ap),
	}
	response, err := o.Client.MakeSecureRequest(r)
	if err != nil {
		return orderHistoryResponse, fmt.Errorf("order history request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return orderHistoryResponse, fmt.Errorf("order history request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	err = json.Unmarshal(response.Body, &orderHistoryResponse)
	if err != nil {
		return orderHistoryResponse, fmt.Errorf("failed to decode order history response: %w", err)
	}

	return orderHistoryResponse, err
}

// Order history
func (o *Order) GetOpenOrders(currencyPair string) (OpenOrdersResponse, error) {
	openOrdersResponse := OpenOrdersResponse{}

	// URL compose
	u, _ := url.Parse(o.Client.GetBaseUrl() + openOrdersEndpoint)
	q := u.Query()

	// Adding optional GET param
	if currencyPair != "" {
		q.Set(currencyPairParamName, currencyPair)
	}

	u.RawQuery = q.Encode()

	ap := map[string]string{}

	r := coinmate.Request{
		HTTPMethod: http.MethodPost,
		URL:        u.String(),
		Body:       o.Client.GetRequestBody(ap),
	}
	response, err := o.Client.MakeSecureRequest(r)
	if err != nil {
		return openOrdersResponse, fmt.Errorf("open orders request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return openOrdersResponse, fmt.Errorf("open orders request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	err = json.Unmarshal(response.Body, &openOrdersResponse)
	if err != nil {
		return openOrdersResponse, fmt.Errorf("failed to decode open orders response: %w", err)
	}

	return openOrdersResponse, err
}

// Buy limit
func (o *Order) BuyLimit(amount, price, stopPrice float64, currencyPair string, hidden, immediateOrCancel bool, clientOrderId uint64) (SellLimit, error) {
	buyLimitResponse := BuyLimitResponse{}
	sellLimit := SellLimit{}

	response, err := limitOrders(o, amount, price, currencyPair, buyLimitOrderEndpoint, stopPrice, hidden, immediateOrCancel, clientOrderId)
	if err != nil {
		return sellLimit, fmt.Errorf("buy limit request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return sellLimit, fmt.Errorf("buy limit request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}
	err = json.Unmarshal(response.Body, &buyLimitResponse)
	if err != nil {
		return sellLimit, fmt.Errorf("failed to decode buy limit response: %w", err)
	}

	sellLimit.Error = buyLimitResponse.Error
	sellLimit.ErrorMessage = buyLimitResponse.ErrorMessage
	sellLimit.OrderId = buyLimitResponse.Data

	return sellLimit, err
}

// Sell limit
func (o *Order) SellLimit(amount, price, stopPrice float64, currencyPair string, hidden, immediateOrCancel bool, clientOrderId uint64) (SellLimit, error) {
	sellLimitResponse := SellLimitResponse{}
	sellLimit := SellLimit{}

	response, err := limitOrders(o, amount, price, currencyPair, sellLimitOrderEndpoint, stopPrice, hidden, immediateOrCancel, clientOrderId)
	if err != nil {
		return sellLimit, fmt.Errorf("sell limit request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return sellLimit, fmt.Errorf("sell limit request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}
	err = json.Unmarshal(response.Body, &sellLimitResponse)
	if err != nil {
		return sellLimit, fmt.Errorf("failed to decode sell limit response: %w", err)
	}

	sellLimit.Error = sellLimitResponse.Error
	sellLimit.ErrorMessage = sellLimitResponse.ErrorMessage
	sellLimit.OrderId = sellLimitResponse.Data

	return sellLimit, err
}

// Buy instantly
func (o *Order) BuyInstant(total float64, cp string, clientOrderId uint64) (BuyAndSellResponse, error) {
	return buySellInstantRequest(o, buyInstantOrderEndpoint, total, cp, clientOrderId)
}

// Sell instantly
func (o *Order) SellInstant(total float64, cp string, clientOrderId uint64) (BuyAndSellResponse, error) {
	return buySellInstantRequest(o, sellInstantOrderEndpoint, total, cp, clientOrderId)
}

// Cancel order
func (o *Order) CancelOrder(orderId uint64) (CancelOrderResponse, error) {
	cancelOrderResponse := CancelOrderResponse{}

	response, err := cancelOrderRequest(o, cancelOrderEndpoint, orderId)
	if err != nil {
		return cancelOrderResponse, fmt.Errorf("cancel order request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return cancelOrderResponse, fmt.Errorf("cancel order request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	err = json.Unmarshal(response.Body, &cancelOrderResponse)
	if err != nil {
		return cancelOrderResponse, fmt.Errorf("failed to decode cancel order response: %w", err)
	}

	return cancelOrderResponse, err
}

// Cancel order with info
func (o *Order) CancelOrderWithInfo(orderId uint64) (CancelOrderWithInfoResponse, error) {
	cancelOrderWithInfoResponse := CancelOrderWithInfoResponse{}

	response, err := cancelOrderRequest(o, cancelOrderWithInfoEndpoint, orderId)
	if err != nil {
		return cancelOrderWithInfoResponse, fmt.Errorf("cancel order with info request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return cancelOrderWithInfoResponse, fmt.Errorf("cancel order with info request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	err = json.Unmarshal(response.Body, &cancelOrderWithInfoResponse)
	if err != nil {
		return cancelOrderWithInfoResponse, fmt.Errorf("failed to decode cancel order with info response: %w", err)
	}

	return cancelOrderWithInfoResponse, err
}

// Helper functions

// Calling limit orders endpoints
func limitOrders(o *Order, amount float64, price float64, currencyPair, endpoint string, stopPrice float64, hidden bool, immediateOrCancel bool, clientOrderId uint64) (coinmate.Response, error) {
	// URL compose
	u, _ := url.Parse(o.Client.GetBaseUrl() + endpoint)
	ap := make(map[string]string)
	ap[amountParamName] = strconv.FormatFloat(amount, 'f', 2, 64)
	ap[priceParamName] = strconv.FormatFloat(price, 'f', 2, 64)
	ap[currencyPairParamName] = strings.ToLower(currencyPair)
	if stopPrice > 0 {
		ap[stopPriceParamName] = strconv.FormatFloat(stopPrice, 'f', 2, 64)
	}
	if hidden == true {
		ap[hiddenParamName] = "1"
	}
	if immediateOrCancel == true {
		ap[immediateOrCancelParamName] = "1"
	}
	if clientOrderId > 0 {
		ap[clientOrderIdParamName] = strconv.FormatUint(clientOrderId, 10)
	}
	r := coinmate.Request{
		HTTPMethod: http.MethodPost,
		URL:        u.String(),
		Body:       o.Client.GetRequestBody(ap),
	}
	response, err := o.Client.MakeSecureRequest(r)
	return response, err
}

// Buy or sell instant request
func buySellInstantRequest(o *Order, endpoint string, total float64, currencyPair string, clientOrderId uint64) (BuyAndSellResponse, error) {
	bir := BuySell{}
	basr := BuyAndSellResponse{}
	u, _ := url.Parse(o.Client.GetBaseUrl() + endpoint)
	ap := make(map[string]string)
	if endpoint == sellInstantOrderEndpoint {
		ap[amountParamName] = strconv.FormatFloat(total, 'f', 8, 64)
	} else {
		ap[totalParamName] = strconv.FormatFloat(total, 'f', 2, 64)
	}
	ap[currencyPairParamName] = strings.ToLower(currencyPair)
	if clientOrderId > 0 {
		ap[clientOrderIdParamName] = strconv.FormatUint(clientOrderId, 10)
	}
	r := coinmate.Request{
		HTTPMethod: http.MethodPost,
		URL:        u.String(),
		Body:       o.Client.GetRequestBody(ap),
	}
	response, err := o.Client.MakeSecureRequest(r)
	if err != nil {
		return basr, fmt.Errorf("%s request failed: %w", endpoint, err)
	}
	if response.StatusCode != http.StatusOK {
		return basr, fmt.Errorf("%s request failed: status=%d body=%s", endpoint, response.StatusCode, string(response.Body))
	}
	err = json.Unmarshal(response.Body, &bir)
	if err != nil {
		return basr, fmt.Errorf("failed to decode %s response: %w", endpoint, err)
	}

	basr.Error = bir.Error
	basr.ErrorMessage = bir.ErrorMessage
	basr.OrderId = bir.Data

	return basr, err
}

func cancelOrderRequest(o *Order, endpoint string, orderId uint64) (coinmate.Response, error) {
	// URL compose
	u, _ := url.Parse(o.Client.GetBaseUrl() + endpoint)
	ap := make(map[string]string)
	ap[orderIdParamName] = strconv.FormatUint(orderId, 10)
	r := coinmate.Request{
		HTTPMethod: http.MethodPost,
		URL:        u.String(),
		Body:       o.Client.GetRequestBody(ap),
	}
	response, err := o.Client.MakeSecureRequest(r)
	return response, err
}
