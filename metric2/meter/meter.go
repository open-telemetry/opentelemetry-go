package meter

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	metric "go.opentelemetry.io/otel/metric2"
	"go.opentelemetry.io/otel/metric2/async"
	"go.opentelemetry.io/otel/metric2/sync"
)

type MeterProvider interface {
	Meter(instrumentationName string /*, opts ...MeterOption*/) Meter
}

type Meter struct {
}

func (m Meter) Asynchronous() async.Meter {
	return async.Meter{}
}

func (m Meter) Synchronous() sync.Meter {
	return sync.Meter{}
}

// ProcessBatch processes a batch of measurements as a single logical
// event.
func (m Meter) ProcessBatch(
	ctx context.Context,
	attrs []attribute.KeyValue,
	batch ...metric.Measurement) {
}

// Process processes a single measurement.  This offers the
// convenience of passing a variable length list of attributes for a
// processing a single measurement.
func (m Meter) Process(
	ctx context.Context,
	ms metric.Measurement,
	attrs ...attribute.KeyValue) {
	// Process a singleton batch
	m.ProcessBatch(ctx, attrs, ms)
}
