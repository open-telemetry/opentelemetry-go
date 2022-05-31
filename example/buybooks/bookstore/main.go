package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/example/buybooks/common/api"
	"go.opentelemetry.io/otel/example/buybooks/common/event"
	"go.opentelemetry.io/otel/example/buybooks/common/trace"
)

const (
	service     = "bookstore"
	environment = "dev"
)

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
	tracer := provider.GetNamedTracer("bookstore/main")

	// NSQ
	config := nsq.NewConfig()
	nsqd := os.Getenv("nsqd")
	if nsqd == "" {
		nsqd = "127.0.0.1:4150"
	}
	producer, err := nsq.NewProducer(nsqd, config)
	if err != nil {
		log.Fatal(err)
	}
	topic := "order"

	// grpc
	paymentservice := os.Getenv("paymentservice")
	if paymentservice == "" {
		paymentservice = "127.0.0.1:8082"
	}
	conn, err := grpc.Dial(paymentservice, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor(
			otelgrpc.WithTracerProvider(provider.TracerProvider),
			otelgrpc.WithPropagators(provider.Propagators),
		)),
	)
	if err != nil {
		log.Fatalf("not connected : %v", err)
	}
	defer conn.Close()
	pc := api.NewPaymentClient(conn)

	// gin
	r := gin.New()
	r.Use(otelgin.Middleware("bookstore", otelgin.WithTracerProvider(provider.TracerProvider)))
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	r.GET("/order/:bookID", func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "order")
		defer span.End()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		bookID := c.Param("bookID")

		res, err := pc.Charge(ctx, &api.ChargeRequest{BookID: bookID})
		if err != nil {
			// todo: better status code
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		fmt.Println("resp: ", res)
		fmt.Printf("$%v charged\n", res.Amount)

		err = publish(ctx, tracer, bookID, producer, topic)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, nil)
	})
	if err := r.Run(); err != nil {
		log.Fatalf(err.Error())
	}
	producer.Stop()
}

func publish(ctx context.Context, tracer trace.Tracer, bookID string, producer *nsq.Producer, topic string) error {
	ctx, span := tracer.Start(ctx, "publish", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()
	body, _ := json.Marshal(struct {
		BookID string `json:"bookID"`
	}{
		BookID: bookID,
	})

	e := &event.Event{
		EventType:  event.OrderEvent,
		Attributes: map[string]string{},
		Body:       body,
	}

	otel.GetTextMapPropagator().Inject(ctx, event.NewEventCarrier(e))
	bs, _ := json.Marshal(e)
	spanContext := span.SpanContext()
	fmt.Println("trace id:", spanContext.TraceID())
	err := producer.Publish(topic, bs)
	return err
}
