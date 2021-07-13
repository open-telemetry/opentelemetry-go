package meter

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric2/asyncfloat64"
	"go.opentelemetry.io/otel/metric2/asyncint64"
	"go.opentelemetry.io/otel/metric2/sdkapi"
	"go.opentelemetry.io/otel/metric2/syncfloat64"
	"go.opentelemetry.io/otel/metric2/syncint64"
)

type MeterOption = metric.MeterOption
type InstrumentOption = metric.InstrumentOption

// Provider supports creating named Meter instances, for instrumenting
// an application containing multiple libraries of code.
type Provider interface {
	Meter(instrumentationName string, opts ...MeterOption) Meter
}

// Meter is an instance of an OpenTelemetry metrics interface for an
// individual named library of code.  This is the top-level entry
// point for creating instruments.
type Meter struct {
}

type AsyncFloat64Instruments struct{}
type AsyncInt64Instruments struct{}
type SyncFloat64Instruments struct{}
type SyncInt64Instruments struct{}

func (m Meter) AsyncInt64() AsyncInt64Instruments {
	return AsyncInt64Instruments{}
}

func (m Meter) AsyncFloat64() AsyncFloat64Instruments {
	return AsyncFloat64Instruments{}
}

func (m Meter) SyncInt64() SyncInt64Instruments {
	return SyncInt64Instruments{}
}

func (m Meter) SyncFloat64() SyncFloat64Instruments {
	return SyncFloat64Instruments{}
}

// ProcessBatch processes a batch of measurements as a single logical
// event.
func (m Meter) ProcessBatch(
	ctx context.Context,
	attrs []attribute.KeyValue,
	batch ...sdkapi.Measurement) {
}

// Process processes a single measurement.  This offers the
// convenience of passing a variable length list of attributes for a
// processing a single measurement.
func (m Meter) Process(
	ctx context.Context,
	ms sdkapi.Measurement,
	attrs ...attribute.KeyValue) {
	// Process a singleton batch
	m.ProcessBatch(ctx, attrs, ms)
}

func (m AsyncFloat64Instruments) Counter(name string, opts ...InstrumentOption) (asyncfloat64.Counter, error) {
	return asyncfloat64.Counter{}, nil
}

func (m AsyncFloat64Instruments) UpDownCounter(name string, opts ...InstrumentOption) (asyncfloat64.UpDownCounter, error) {
	return asyncfloat64.UpDownCounter{}, nil
}

func (m AsyncFloat64Instruments) Gauge(name string, opts ...InstrumentOption) (asyncfloat64.Gauge, error) {
	return asyncfloat64.Gauge{}, nil
}

func (m AsyncInt64Instruments) Counter(name string, opts ...InstrumentOption) (asyncint64.Counter, error) {
	return asyncint64.Counter{}, nil
}

func (m AsyncInt64Instruments) UpDownCounter(name string, opts ...InstrumentOption) (asyncint64.UpDownCounter, error) {
	return asyncint64.UpDownCounter{}, nil
}

func (m AsyncInt64Instruments) Gauge(name string, opts ...InstrumentOption) (asyncint64.Gauge, error) {
	return asyncint64.Gauge{}, nil
}

func (m SyncFloat64Instruments) Counter(name string, opts ...InstrumentOption) (syncfloat64.Counter, error) {
	return syncfloat64.Counter{}, nil
}

func (m SyncFloat64Instruments) UpDownCounter(name string, opts ...InstrumentOption) (syncfloat64.UpDownCounter, error) {
	return syncfloat64.UpDownCounter{}, nil
}

func (m SyncFloat64Instruments) Histogram(name string, opts ...InstrumentOption) (syncfloat64.Histogram, error) {
	return syncfloat64.Histogram{}, nil
}

func (m SyncInt64Instruments) Counter(name string, opts ...InstrumentOption) (syncint64.Counter, error) {
	return syncint64.Counter{}, nil
}

func (m SyncInt64Instruments) UpDownCounter(name string, opts ...InstrumentOption) (syncint64.UpDownCounter, error) {
	return syncint64.UpDownCounter{}, nil
}

func (m SyncInt64Instruments) Histogram(name string, opts ...InstrumentOption) (syncint64.Histogram, error) {
	return syncint64.Histogram{}, nil
}
