package http

import (
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TempByCEPClient interface {
	GetTempByCEP(cep string) (float64, error)
}

type tempByCEPClient struct {
	client HTTPClient
}

func NewTempByCEPClient(client HTTPClient) *tempByCEPClient {

	return &tempByCEPClient{
		client: client,
	}
}

func (c *tempByCEPClient) GetTempByCEP(cep string) (float64, error) {

	return 0, nil
}
