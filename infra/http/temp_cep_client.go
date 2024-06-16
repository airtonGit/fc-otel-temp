package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TempByCEPResponse struct {
	TempC      float64 `json:"temp_C,omitempty"`
	TempF      float64 `json:"temp_F,omitempty"`
	TempK      float64 `json:"temp_K,omitempty"`
	Localidade string  `json:"localidade"`
}

type TempByCEPClient interface {
	DoRequest(ctx context.Context, cep string) (TempByCEPResponse, error)
}

type tempByCEPClient struct {
	client   HTTPClient
	endpoint string
}

func NewTempByCEPClient(client HTTPClient, endpoint string) *tempByCEPClient {

	return &tempByCEPClient{
		client:   client,
		endpoint: endpoint,
	}
}

func (c *tempByCEPClient) DoRequest(ctx context.Context, cep string) (TempByCEPResponse, error) {

	endpoint := fmt.Sprintf("http://%s/cep/%s", c.endpoint, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return TempByCEPResponse{}, fmt.Errorf("fail make request err=%w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return TempByCEPResponse{}, fmt.Errorf("fail request err=%w", err)
	}

	var localTempResp TempByCEPResponse
	err = json.NewDecoder(resp.Body).Decode(&localTempResp)
	if err != nil {
		return TempByCEPResponse{}, fmt.Errorf("fail response decode err=%w", err)
	}

	return localTempResp, nil
}
