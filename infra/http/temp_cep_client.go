package http

import (
	"context"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TempByCEPClient interface {
	DoRequest(ctx context.Context, cep string) (float64, error)
}

type tempByCEPClient struct {
	client  HTTPClient
	enpoint string
}

func NewTempByCEPClient(client HTTPClient) *tempByCEPClient {

	return &tempByCEPClient{
		client: client,
	}
}

func (c *tempByCEPClient) DoRequest(ctx context.Context, cep string) (float64, error) {

	return 0, nil
}
