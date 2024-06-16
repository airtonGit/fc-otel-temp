package domain

import (
	"context"
	"github.com/airtongit/fc-otel-temp/infra/http"
	"go.opentelemetry.io/otel/trace"
)

type CEPByTempClient interface {
	DoRequest(ctx context.Context, cep string) (http.TempByCEPResponse, error)
}

type TempByCEPService interface {
	GetTempByCEP(ctx context.Context, cep string) (http.TempByCEPResponse, error)
}

type tempByCEPService struct {
	cepByTempClient CEPByTempClient
	OTELTracer      trace.Tracer
}

func NewTempByCEPService(cepByTempClient CEPByTempClient, tracer trace.Tracer) *tempByCEPService {
	return &tempByCEPService{
		cepByTempClient: cepByTempClient,
		OTELTracer:      tracer,
	}
}

func (s *tempByCEPService) GetTempByCEP(ctx context.Context, cep string) (http.TempByCEPResponse, error) {

	ctx, span := s.OTELTracer.Start(ctx, "Chama service-b")
	defer span.End()

	return s.cepByTempClient.DoRequest(ctx, cep)
}
