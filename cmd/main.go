package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	infrahttp "github.com/airtongit/fc-otel-temp/infra/http"
	"github.com/airtongit/fc-otel-temp/internal/domain"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func initProvider(serviceName, collectorURL string) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

type config struct {
	OtelServiceName      string
	OtelExporterEndpoing string
}

func loadConfig() (config, error) {

	viper.AutomaticEnv()

	cfg := config{
		OtelServiceName:      viper.GetString("OTEL_SERVICE_NAME"),
		OtelExporterEndpoing: viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"),
	}

	if cfg.OtelServiceName == "" {
		return config{}, fmt.Errorf("otel service name empty")
	}

	if cfg.OtelExporterEndpoing == "" {
		return config{}, fmt.Errorf("otel endpoint empty")
	}
	return cfg, nil
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	shutdown, err := initProvider(cfg.OtelServiceName, cfg.OtelExporterEndpoing)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("temp-by-cep-otel-tracer")
	tempbyCEPClient := infrahttp.NewTempByCEPClient(http.DefaultClient, "localhost:7070")
	tempByCEPService := domain.NewTempByCEPService(tempbyCEPClient, tracer)

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Listening, READY"))
	})
	r.Post("/cep", domain.MakeRequestTempByCEPHandler(tempByCEPService))

	// Determine port for HTTP service.
	port := viper.GetString("PORT") //os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}

}
