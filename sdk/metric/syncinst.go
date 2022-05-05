package metric

import (
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/metric/internal/syncstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

type (
	syncint64Instruments   struct{ *meter }
	syncfloat64Instruments struct{ *meter }
)

func (i syncint64Instruments) Counter(name string, opts ...instrument.Option) (syncint64.Counter, error) {
	inst, err := i.synchronousInstrument(name, opts, number.Int64Kind, sdkinstrument.CounterKind)
	return syncstate.NewCounter[int64, number.Int64Traits](inst), err
}

func (i syncint64Instruments) UpDownCounter(name string, opts ...instrument.Option) (syncint64.UpDownCounter, error) {
	inst, err := i.synchronousInstrument(name, opts, number.Int64Kind, sdkinstrument.UpDownCounterKind)
	return syncstate.NewCounter[int64, number.Int64Traits](inst), err
}

func (i syncint64Instruments) Histogram(name string, opts ...instrument.Option) (syncint64.Histogram, error) {
	inst, err := i.synchronousInstrument(name, opts, number.Int64Kind, sdkinstrument.HistogramKind)
	return syncstate.NewHistogram[int64, number.Int64Traits](inst), err
}

func (f syncfloat64Instruments) Counter(name string, opts ...instrument.Option) (syncfloat64.Counter, error) {
	inst, err := f.synchronousInstrument(name, opts, number.Float64Kind, sdkinstrument.CounterKind)
	return syncstate.NewCounter[float64, number.Float64Traits](inst), err
}

func (f syncfloat64Instruments) UpDownCounter(name string, opts ...instrument.Option) (syncfloat64.UpDownCounter, error) {
	inst, err := f.synchronousInstrument(name, opts, number.Float64Kind, sdkinstrument.UpDownCounterKind)
	return syncstate.NewCounter[float64, number.Float64Traits](inst), err
}

func (f syncfloat64Instruments) Histogram(name string, opts ...instrument.Option) (syncfloat64.Histogram, error) {
	inst, err := f.synchronousInstrument(name, opts, number.Float64Kind, sdkinstrument.HistogramKind)
	return syncstate.NewHistogram[float64, number.Float64Traits](inst), err
}
