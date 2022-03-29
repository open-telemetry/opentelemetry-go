package metric

import (
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/sdk/metric/internal/asyncstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type (
	asyncint64Instruments   struct{ *meter }
	asyncfloat64Instruments struct{ *meter }
)

func (m *meter) AsyncInt64() asyncint64.InstrumentProvider {
	return asyncint64Instruments{m}
}

func (m *meter) AsyncFloat64() asyncfloat64.InstrumentProvider {
	return asyncfloat64Instruments{m}
}

func (m *meter) newAsyncInst(name string, opts []instrument.Option, nk number.Kind, ik sdkapi.InstrumentKind) (*asyncstate.Instrument, error) {
	return nameLookup(
		m, name, opts, nk, ik,
		func(desc sdkapi.Descriptor) *asyncstate.Instrument {
			compiled := m.views.Compile(desc)
			inst := asyncstate.NewInstrument(desc, compiled)
			m.instruments = append(m.instruments, inst)
			return inst
		})
}

func (i asyncint64Instruments) Counter(name string, opts ...instrument.Option) (asyncint64.Counter, error) {
	inst, err := i.newAsyncInst(name, opts, number.Int64Kind, sdkapi.CounterObserverInstrumentKind)
	return asyncstate.NewObserver[int64, traits.Int64](inst), err
}

func (i asyncint64Instruments) UpDownCounter(name string, opts ...instrument.Option) (asyncint64.UpDownCounter, error) {
	inst, err := i.newAsyncInst(name, opts, number.Int64Kind, sdkapi.UpDownCounterObserverInstrumentKind)
	return asyncstate.NewObserver[int64, traits.Int64](inst), err
}

func (i asyncint64Instruments) Gauge(name string, opts ...instrument.Option) (asyncint64.Gauge, error) {
	inst, err := i.newAsyncInst(name, opts, number.Int64Kind, sdkapi.GaugeObserverInstrumentKind)
	return asyncstate.NewObserver[int64, traits.Int64](inst), err
}

func (f asyncfloat64Instruments) Counter(name string, opts ...instrument.Option) (asyncfloat64.Counter, error) {
	inst, err := f.newAsyncInst(name, opts, number.Float64Kind, sdkapi.CounterObserverInstrumentKind)
	return asyncstate.NewObserver[float64, traits.Float64](inst), err
}

func (f asyncfloat64Instruments) UpDownCounter(name string, opts ...instrument.Option) (asyncfloat64.UpDownCounter, error) {
	inst, err := f.newAsyncInst(name, opts, number.Float64Kind, sdkapi.UpDownCounterObserverInstrumentKind)
	return asyncstate.NewObserver[float64, traits.Float64](inst), err
}

func (f asyncfloat64Instruments) Gauge(name string, opts ...instrument.Option) (asyncfloat64.Gauge, error) {
	inst, err := f.newAsyncInst(name, opts, number.Float64Kind, sdkapi.GaugeObserverInstrumentKind)
	return asyncstate.NewObserver[float64, traits.Float64](inst), err
}
