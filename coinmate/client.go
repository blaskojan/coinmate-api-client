package coinmate

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	baseUrl                = "https://coinmate.io/api"
	contentType            = "Content-Type"
	secureContentTypeValue = "application/x-www-form-urlencoded"

	// Http request timeout to 2s
	requestTimeout = 2 * time.Second
)

// Request data
type Request struct {
	HTTPMethod string
	URL        string
	Body       []byte
}

// Response data
type Response struct {
	Status     string
	StatusCode int
	Body       []byte
}

type ClientInterface interface {
	GetNonce() string
	GetSignature(clientId, apiKey, nonce, privateKey string) string
	GetBaseUrl() string
	MakePublicRequest(r Request) (Response, error)
	MakeSecureRequest(r Request) (Response, error)
	GetRequestBody(map[string]string) []byte
}

type CoinmateClient struct {
	ClientID   string
	ApiKey     string
	PrivateKey string
	Nonce      string
	Signature  string
	httpClient http.Client
	lastNonce  int64
}

type CoinmateResponse struct {
	Error        bool        `json:"error"`
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data"`
}

// Balance currency data
type BalanceCurrency struct {
	Currency  string  `json:"currency"`
	Balance   float32 `json:"balance"`
	Reserved  float32 `json:"reserved"`
	Available float32 `json:"available"`
}

// Return Coinmate client
func GetCoinmateClient(clientId, publicKey, privateKey string) *CoinmateClient {
	client := new(CoinmateClient)

	client.ClientID = clientId
	client.ApiKey = publicKey
	client.PrivateKey = privateKey
	client.httpClient = http.Client{
		Timeout: time.Duration(requestTimeout),
	}
	return client
}

// Return nonce (security)
func (c *CoinmateClient) GetNonce() string {
	now := time.Now().UnixNano()
	if now <= c.lastNonce {
		now = c.lastNonce + 1
	}
	c.lastNonce = now
	return strconv.FormatInt(now, 10)
}

// Return signature (security)
func (c *CoinmateClient) GetSignature(clientId, apiKey, nonce, privateKey string) string {
	secret := privateKey
	data := nonce + clientId + apiKey
	//fmt.Printf("Secret: %s Data: %s\n", secret, data)

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(data))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	//fmt.Println("Result: " + strings.ToUpper(sha))
	return strings.ToUpper(sha)
}

// Return url prefix
func (c *CoinmateClient) GetBaseUrl() string {
	return baseUrl
}

// Return request body due to security
func (c *CoinmateClient) GetRequestBody(additionalParams map[string]string) []byte {
	nonce := c.GetNonce()

	u := url.URL{}
	q := u.Query()

	// Additional parameters into body
	if len(additionalParams) > 0 {
		for paramName, paramValue := range additionalParams {
			q.Set(paramName, paramValue)
		}
	}
	q.Set("clientId", c.ClientID)
	q.Set("publicKey", c.ApiKey)
	q.Set("nonce", nonce)
	q.Set("signature", c.GetSignature(c.ClientID, c.ApiKey, nonce, c.PrivateKey))
	v := q.Encode()
	fmt.Println(v)

	return []byte(v)
}

// Make public request
func (c *CoinmateClient) MakePublicRequest(r Request) (Response, error) {
	var rb io.Reader
	if r.Body != nil {
		rb = bytes.NewBuffer(r.Body)
	}
	request, err := http.NewRequest(r.HTTPMethod, r.URL, rb)
	if err != nil {
		return Response{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Response{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: response.StatusCode,
		Status:     response.Status,
		Body:       body,
	}, nil
}

// Make secure request
func (c *CoinmateClient) MakeSecureRequest(r Request) (Response, error) {
	var rb io.Reader

	if r.Body != nil {
		rb = bytes.NewBuffer(r.Body)
	}
	request, err := http.NewRequest(r.HTTPMethod, r.URL, rb)
	if err != nil {
		return Response{}, err
	}

	request.Header.Add(contentType, secureContentTypeValue)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Response{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: response.StatusCode,
		Status:     response.Status,
		Body:       body,
	}, nil
}
