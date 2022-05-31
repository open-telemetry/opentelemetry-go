package trace

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracer = trace.Tracer

type TracerProviderShutdownFunc func(ctx context.Context) error

const (
	AlwaysSample float64 = 1
	NeverSample  float64 = 0
)

var (
	WithSpanKind     = trace.WithSpanKind
	SpanKindConsumer = trace.SpanKindConsumer
	SpanKindProducer = trace.SpanKindProducer
)

type providerConfig struct {
	exporter      tracesdk.SpanExporter
	samplingRatio float64
	attributes    map[string]string
}

func newTraceProviderConfig() *providerConfig {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	return &providerConfig{
		exporter:      exporter,
		samplingRatio: AlwaysSample,
		attributes:    map[string]string{},
	}
}

type ProviderOptionFunc func(config *providerConfig)

func (o ProviderOptionFunc) apply(t *providerConfig) {
	o(t)
}

// WithSamplingRatio set SamplingRatio, the value should between 0-1,
//      e.g.: 0:NeverSampling, 1:AlwaysSampling.
func WithSamplingRatio(ratio float64) ProviderOptionFunc {
	return func(config *providerConfig) {
		config.samplingRatio = ratio
	}
}

// WithJaegerExporter set jaeger URL.
func WithJaegerExporter(url string) ProviderOptionFunc {
	return func(config *providerConfig) {
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
		if err != nil {
			log.Fatalf("failed to create jaeger exporter: %+v", err)
		}
		config.exporter = exp
	}
}

// WithMetadata set jaeger metadata.
func WithMetadata(metadata map[string]string) ProviderOptionFunc {
	return func(config *providerConfig) {
		config.attributes = metadata
	}
}

type ProviderOption interface {
	apply(config *providerConfig)
}

// NewTracerProvider provider stdout trace exporter by default.
func NewTracerProvider(
	serviceName string,
	environment string,
	opt ...ProviderOption,
) (*Provider, TracerProviderShutdownFunc) {
	config := newTraceProviderConfig()
	for _, opt := range opt {
		opt.apply(config)
	}

	attributes := make([]attribute.KeyValue, 0, len(config.attributes)+2)
	for k, v := range config.attributes {
		attributes = append(attributes, attribute.String(k, v))
	}
	attributes = append(attributes, semconv.ServiceNameKey.String(serviceName))
	attributes = append(attributes, attribute.String("environment", environment))

	tp := tracesdk.NewTracerProvider(
		// In a production app, use trace.ParentBased(trace.TraceIDRatioBased) set at desired ratio
		tracesdk.WithSampler(tracesdk.ParentBased(
			tracesdk.TraceIDRatioBased(config.samplingRatio))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(config.exporter),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			attributes...,
		)),
	)
	// some middlewares use global instance as default, or in some cases needed, so set here
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return newTraceConfig(tp), tp.Shutdown
}

func NewNoopTracerProvider() (*Provider, TracerProviderShutdownFunc) {
	tp := trace.NewNoopTracerProvider()
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return newTraceConfig(tp), func(ctx context.Context) error { return nil }
}

// Provider holds the trace provider and propagator.
type Provider struct {
	TracerProvider trace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

func (t Provider) GetDefaultTracer() trace.Tracer {
	return t.TracerProvider.Tracer("")
}

func (t Provider) GetNamedTracer(name string) trace.Tracer {
	return t.TracerProvider.Tracer(name)
}

func (t Provider) GetPropagators() propagation.TextMapPropagator {
	return t.Propagators
}

func newTraceConfig(tp trace.TracerProvider) *Provider {
	return &Provider{
		TracerProvider: tp,
		Propagators:    otel.GetTextMapPropagator(),
	}
}
