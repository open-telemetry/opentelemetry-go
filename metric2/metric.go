package metric2

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric2/asyncmetric"
	"go.opentelemetry.io/otel/metric2/batch"
	"go.opentelemetry.io/otel/metric2/float64metric"
	"go.opentelemetry.io/otel/metric2/int64metric"
)

type MeterProvider interface {
	Meter(instrumentationName string /*, opts ...MeterOption*/) Meter
}

type Meter struct {
}

func (m Meter) Integer() int64metric.Meter {
	return int64metric.Meter{}
}

func (m Meter) FloatingPoint() float64metric.Meter {
	return float64metric.Meter{}
}

func (m Meter) Asynchronous() asyncmetric.Meter {
	return asyncmetric.Meter{}
}

func (m Meter) ProcessBatch(
	ctx context.Context,
	attrs []attribute.KeyValue,
	batch ...batch.Measurement) {
}

func (m Meter) Process(
	ctx context.Context,
	ms batch.Measurement,
	attrs ...attribute.KeyValue) {
	// Make this a singleton batch.
	m.ProcessBatch(ctx, attrs, ms)
}
