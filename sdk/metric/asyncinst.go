package metric

import (
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/sdk/metric/internal/asyncstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

type (
	asyncint64Instruments   struct{ *meter }
	asyncfloat64Instruments struct{ *meter }
)

func (i asyncint64Instruments) Counter(name string, opts ...instrument.Option) (asyncint64.Counter, error) {
	inst, err := i.asynchronousInstrument(name, opts, number.Int64Kind, sdkinstrument.CounterObserverKind)
	return asyncstate.NewObserver[int64, number.Int64Traits](inst), err
}

func (i asyncint64Instruments) UpDownCounter(name string, opts ...instrument.Option) (asyncint64.UpDownCounter, error) {
	inst, err := i.asynchronousInstrument(name, opts, number.Int64Kind, sdkinstrument.UpDownCounterObserverKind)
	return asyncstate.NewObserver[int64, number.Int64Traits](inst), err
}

func (i asyncint64Instruments) Gauge(name string, opts ...instrument.Option) (asyncint64.Gauge, error) {
	inst, err := i.asynchronousInstrument(name, opts, number.Int64Kind, sdkinstrument.GaugeObserverKind)
	return asyncstate.NewObserver[int64, number.Int64Traits](inst), err
}

func (f asyncfloat64Instruments) Counter(name string, opts ...instrument.Option) (asyncfloat64.Counter, error) {
	inst, err := f.asynchronousInstrument(name, opts, number.Float64Kind, sdkinstrument.CounterObserverKind)
	return asyncstate.NewObserver[float64, number.Float64Traits](inst), err
}

func (f asyncfloat64Instruments) UpDownCounter(name string, opts ...instrument.Option) (asyncfloat64.UpDownCounter, error) {
	inst, err := f.asynchronousInstrument(name, opts, number.Float64Kind, sdkinstrument.UpDownCounterObserverKind)
	return asyncstate.NewObserver[float64, number.Float64Traits](inst), err
}

func (f asyncfloat64Instruments) Gauge(name string, opts ...instrument.Option) (asyncfloat64.Gauge, error) {
	inst, err := f.asynchronousInstrument(name, opts, number.Float64Kind, sdkinstrument.GaugeObserverKind)
	return asyncstate.NewObserver[float64, number.Float64Traits](inst), err
}
