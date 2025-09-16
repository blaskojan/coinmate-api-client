package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tourGo/coinmate"
)

type ServerTime struct {
	Client coinmate.ClientInterface
}

type ServerTimeResponse struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	Data         int64  `json:"data"`
}

func (s *ServerTime) GetServerTime() (ServerTimeResponse, error) {
	st := ServerTimeResponse{}

	r := coinmate.Request{
		HTTPMethod: http.MethodGet,
		URL:        s.Client.GetBaseUrl() + "/system/get-server-time",
		Body:       nil,
	}
	response, err := s.Client.MakePublicRequest(r)
	if err != nil {
		return st, fmt.Errorf("server time request failed: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return st, fmt.Errorf("server time request failed: status=%d body=%s", response.StatusCode, string(response.Body))
	}

	if err := json.Unmarshal(response.Body, &st); err != nil {
		return st, fmt.Errorf("failed to decode server time response: %w", err)
	}
	return st, nil
}
