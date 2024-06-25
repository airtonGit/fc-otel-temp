package domain

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type TempByCEPRequest struct {
	CEP string `json:"cep"`
}

func validate(cep string) error {
	log.Println("validate CEP", cep)
	matched, err := regexp.MatchString(`^\d{8}$`, cep)
	if err != nil {
		return err
	}
	if !matched {
		log.Println("matchstring not match with", cep)
		return errors.New("invalid zipcode")
	}
	return nil
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

		//if err := validate(request.CEP); err != nil {
		//	http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		//	return
		//}

		log.Println("request handling cep:", request.CEP)

		carrier := propagation.HeaderCarrier(req.Header)
		ctx := req.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		temp, err := tempByCEPService.GetTempByCEP(ctx, request.CEP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if temp.StatusCode > 201 {
			if temp.StatusCode == 404 {
				http.NotFound(w, req)
				return
			}
			w.WriteHeader(temp.StatusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(temp)
	}
}
