package domain

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type CEPByTempClient interface {
	DoRequest(ctx context.Context, cep string) (float64, error)
}

type TempByCEPService interface {
	GetTempByCEP(ctx context.Context, cep string) (float64, error)
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

func (s *tempByCEPService) GetTempByCEP(ctx context.Context, cep string) (float64, error) {

	ctx, spanInicial := s.OTELTracer.Start(ctx, "SPAN_INICIAL_GetTempByCEP_Service")
	time.Sleep(time.Second)
	spanInicial.End()

	ctx, span := s.OTELTracer.Start(ctx, "Chama externa")
	defer span.End()

	return s.cepByTempClient.DoRequest(ctx, cep)
}
