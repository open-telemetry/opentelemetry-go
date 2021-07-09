package meter

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric2/async"
	"go.opentelemetry.io/otel/metric2/sync"
)

// MeterProvider supports creating named Meter instances, for
// instrumenting an application containing multiple libraries of code.
type MeterProvider interface {
	Meter(instrumentationName string /*, opts ...MeterOption*/) Meter
}

// Meter is an instance of an OpenTelemetry metrics interface for an
// individual named library of code.  This is the top-level entry
// point for creating instruments.
type Meter struct {
}

// Asynchronous provides access to an async.Meter for constructing
// asynchronous metric instruments.
func (m Meter) Asynchronous() async.Meter {
	return async.Meter{}
}

// Synchronous provides access to a sync.Meter for constructing
// synchronous metric instruments.
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
