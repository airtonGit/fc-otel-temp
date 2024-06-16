package domain

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TempByCEPRequest struct {
	CEP string `json:"cep"`
}

func MakeRequestTempByCEPHandler(tempByCEPService TempByCEPService) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		request := TempByCEPRequest{}
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		if request.CEP == "" {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		carrier := propagation.HeaderCarrier(req.Header)
		ctx := req.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

		temp, err := tempByCEPService.GetTempByCEP(ctx, request.CEP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(temp)
	}
}
