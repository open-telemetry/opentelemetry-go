package meter

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric2/meter/asyncfloat64"
	"go.opentelemetry.io/otel/metric2/meter/asyncint64"
	"go.opentelemetry.io/otel/metric2/meter/syncfloat64"
	"go.opentelemetry.io/otel/metric2/meter/syncint64"
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

func (m Meter) AsyncInt64() asyncint64.Builder {
	return asyncint64.Builder{}
}

func (m Meter) AsyncFloat64() asyncfloat64.Builder {
	return asyncfloat64.Builder{}
}

func (m Meter) SyncInt64() syncint64.Builder {
	return syncint64.Builder{}
}

func (m Meter) SyncFloat64() syncfloat64.Builder {
	return syncfloat64.Builder{}
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
