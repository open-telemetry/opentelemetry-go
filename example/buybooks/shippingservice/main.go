package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"

	"go.opentelemetry.io/otel/example/buybooks/common/event"
	"go.opentelemetry.io/otel/example/buybooks/common/trace"
)

const (
	service     = "shipping-service"
	environment = "dev"
)

type myMessageHandler struct {
	traceProvider *trace.Provider
}

func newMyMessageHandler(provider *trace.Provider) *myMessageHandler {
	return &myMessageHandler{traceProvider: provider}
}

// HandleMessage implements the Handler interface.
func (h myMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	// do whatever actual message processing is desired
	err := h.processMessage(m.Body)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return err
}

func (h myMessageHandler) processMessage(body []byte) error {
	var e event.Event
	var message struct {
		BookID string `json:"bookID"`
	}
	err := json.Unmarshal(body, &e)
	if err != nil {
		return err
	}
	ctx := h.traceProvider.GetPropagators().Extract(context.Background(), event.NewEventCarrier(&e))
	_, span := h.traceProvider.GetDefaultTracer().Start(ctx, "processMessage", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()
	err = json.Unmarshal(e.Body, &message)
	if err != nil {
		return err
	}
	// simulate time to handle event
	time.Sleep(1 * time.Millisecond)
	fmt.Println("trace id:", span.SpanContext().TraceID())
	fmt.Printf("receive BookID: %s\n", message.BookID)
	return nil
}

func main() {
	// tracing
	jaegerURL := os.Getenv("jaeger")
	if jaegerURL == "" {
		jaegerURL = "http://localhost:14268"
	}
	tracerProvider, shutdown := trace.NewTracerProvider(
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
	// Instantiate a consumer that will subscribe to the provided channel.
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("order", "channel", config)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Handler for messages received by this Consumer. Can be called multiple times.
	// See also AddConcurrentHandlers.
	consumer.AddHandler(newMyMessageHandler(tracerProvider))

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	nsqlookupdURL := os.Getenv("nsqlookupd")
	if nsqlookupdURL == "" {
		nsqlookupdURL = "localhost:4161"
	}
	err = consumer.ConnectToNSQLookupd(nsqlookupdURL)
	if err != nil {
		log.Fatal(err)
	}

	// wait for signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	consumer.Stop()
}
