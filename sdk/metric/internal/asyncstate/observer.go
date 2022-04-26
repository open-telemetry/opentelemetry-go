package asyncstate

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/sdk/metric/number"
)

// observer is a generic (int64 or float64) instrument which
// satisfies any of the asynchronous instrument API interfaces.
type observer[N number.Any, Traits number.Traits[N]] struct {
	instrument.AsynchronousStruct

	inst *Instrument
}

// observer implements 6 instruments and memberInstrument.
var (
	_ asyncint64.Counter       = observer[int64, number.Int64Traits]{}
	_ asyncint64.UpDownCounter = observer[int64, number.Int64Traits]{}
	_ asyncint64.Gauge         = observer[int64, number.Int64Traits]{}
	_ memberInstrument         = observer[int64, number.Int64Traits]{}

	_ asyncfloat64.Counter       = observer[float64, number.Float64Traits]{}
	_ asyncfloat64.UpDownCounter = observer[float64, number.Float64Traits]{}
	_ asyncfloat64.Gauge         = observer[float64, number.Float64Traits]{}
	_ memberInstrument           = observer[float64, number.Float64Traits]{}
)

// memberInstrument indicates whether a user-provided
// instrument was returned by this SDK.
type memberInstrument interface {
	instrument() *Instrument
}

// NewObserver returns an generic value suitable for use as any of the
// asynchronous instrument APIs.
func NewObserver[N number.Any, Traits number.Traits[N]](inst *Instrument) observer[N, Traits] {
	return observer[N, Traits]{inst: inst}
}

func (o observer[N, Traits]) instrument() *Instrument {
	return o.inst
}

func (o observer[N, Traits]) Observe(ctx context.Context, value N, attrs ...attribute.KeyValue) {
	capture[N, Traits](ctx, o.inst, value, attrs)
}
