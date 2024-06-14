package domain

import (
	"encoding/json"
	infrahttp "github.com/airtongit/fc-otel-temp/infra/http"
	"net/http"
)

type TempByCEPRequest struct {
	CEP string `json:"cep"`
}

type handler struct {
	client infrahttp.TempByCEPClient
}

func RequestTempByCEP(w http.ResponseWriter, r *http.Request) {
	request := TempByCEPRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	if request.CEP == "" {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

}
