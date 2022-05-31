package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel/example/buybooks/common/trace"
)

const (
	service     = "client"
	environment = "dev"
)

var bookstore = os.Getenv("bookstore")

func main() {
	if bookstore == "" {
		bookstore = "localhost:8080"
	}
	// tracing
	jaegerURL := os.Getenv("jaeger")
	if jaegerURL == "" {
		jaegerURL = "http://localhost:14268/api/traces"
	}
	provider, shutdown := trace.NewTracerProvider(
		service,
		environment,
		trace.WithJaegerExporter(fmt.Sprintf("%s/api/traces", jaegerURL)),
		trace.WithSamplingRatio(trace.AlwaysSample),
	)
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	tracer := provider.GetNamedTracer("client")

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C

		err := SendRequest(context.Background(), tracer, client)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
func SendRequest(ctx context.Context, tr trace.Tracer, client http.Client) error {
	ctx, span := tr.Start(ctx, "client.SendRequest")
	defer span.End()
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/order/%d", bookstore, time.Now().Second()), nil)

	span.AddEvent("sample event")
	span.RecordError(fmt.Errorf("sample error"))
	span.SetStatus(codes.Error, "sample desc")
	fmt.Printf("Sending request...\n")
	_, err := client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
