package metric

import (
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/metric/internal/syncstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

type (
	syncint64Instruments   struct{ *meter }
	syncfloat64Instruments struct{ *meter }
)

func (m *meter) SyncInt64() syncint64.InstrumentProvider {
	return syncint64Instruments{m}
}

func (m *meter) SyncFloat64() syncfloat64.InstrumentProvider {
	return syncfloat64Instruments{m}
}

func (m *meter) newSyncInst(name string, opts []instrument.Option, nk number.Kind, ik sdkinstrument.Kind) (*syncstate.Instrument, error) {
	return nameLookup(
		m, name, opts, nk, ik,
		func(desc sdkinstrument.Descriptor) (*syncstate.Instrument, error) {
			compiled, err := m.views.Compile(desc)
			inst := syncstate.NewInstrument(desc, compiled)

			// Hey!! Problem!!! @@@
			// How do we avoid over-collecting duplicate-registered instruments?
			m.instruments = append(m.instruments, inst)
			return inst, err
		})
}

func (i syncint64Instruments) Counter(name string, opts ...instrument.Option) (syncint64.Counter, error) {
	inst, err := i.newSyncInst(name, opts, number.Int64Kind, sdkinstrument.CounterKind)
	return syncstate.NewCounter[int64, traits.Int64](inst), err
}

func (i syncint64Instruments) UpDownCounter(name string, opts ...instrument.Option) (syncint64.UpDownCounter, error) {
	inst, err := i.newSyncInst(name, opts, number.Int64Kind, sdkinstrument.UpDownCounterKind)
	return syncstate.NewCounter[int64, traits.Int64](inst), err
}

func (i syncint64Instruments) Histogram(name string, opts ...instrument.Option) (syncint64.Histogram, error) {
	inst, err := i.newSyncInst(name, opts, number.Int64Kind, sdkinstrument.HistogramKind)
	return syncstate.NewHistogram[int64, traits.Int64](inst), err
}

func (f syncfloat64Instruments) Counter(name string, opts ...instrument.Option) (syncfloat64.Counter, error) {
	inst, err := f.newSyncInst(name, opts, number.Float64Kind, sdkinstrument.CounterKind)
	return syncstate.NewCounter[float64, traits.Float64](inst), err
}

func (f syncfloat64Instruments) UpDownCounter(name string, opts ...instrument.Option) (syncfloat64.UpDownCounter, error) {
	inst, err := f.newSyncInst(name, opts, number.Float64Kind, sdkinstrument.UpDownCounterKind)
	return syncstate.NewCounter[float64, traits.Float64](inst), err
}

func (f syncfloat64Instruments) Histogram(name string, opts ...instrument.Option) (syncfloat64.Histogram, error) {
	inst, err := f.newSyncInst(name, opts, number.Float64Kind, sdkinstrument.HistogramKind)
	return syncstate.NewHistogram[float64, traits.Float64](inst), err
}
