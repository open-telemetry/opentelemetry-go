package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"go.opentelemetry.io/otel/example/buybooks/common/api"
	"go.opentelemetry.io/otel/example/buybooks/common/trace"
)

const (
	service     = "payment"
	environment = "dev"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

var _ api.PaymentServer = (*server)(nil)

type server struct {
	api.UnimplementedPaymentServer
}

func (s *server) Charge(ctx context.Context, request *api.ChargeRequest) (*api.ChargeReply, error) {
	log.Printf("Received: %s, all books have same price 10\n", request.BookID)
	return &api.ChargeReply{Amount: 10}, nil
}

func main() {
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

	flag.Parse()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(
				otelgrpc.WithTracerProvider(provider.TracerProvider),
				otelgrpc.WithPropagators(provider.Propagators))),
	)
	api.RegisterPaymentServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
