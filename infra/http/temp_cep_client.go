package http

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"log"
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
	StatusCode int     `json:"-"`
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

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := c.client.Do(req)
	if err != nil {
		return TempByCEPResponse{}, fmt.Errorf("fail request err=%w", err)
	}

	log.Println("service-b resp status", resp.Status)

	if resp.StatusCode > 201 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return TempByCEPResponse{}, err
		}
		log.Println("do_request response statuscode:", resp.StatusCode, "body", string(body))
		return TempByCEPResponse{StatusCode: resp.StatusCode}, nil
	}

	var localTempResp TempByCEPResponse
	err = json.NewDecoder(resp.Body).Decode(&localTempResp)
	if err != nil {
		return TempByCEPResponse{}, fmt.Errorf("fail response decode err=%w", err)
	}

	return localTempResp, nil
}
