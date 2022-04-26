package syncstate

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/metric/number"
)

// counter is a synchronous instrument having an Add() method.
type counter[N number.Any, Traits number.Traits[N]] struct {
	instrument.SynchronousStruct

	inst *Instrument
}

// counter satisfies 4 instrument APIs.
var (
	_ syncint64.Counter         = counter[int64, number.Int64Traits]{}
	_ syncint64.UpDownCounter   = counter[int64, number.Int64Traits]{}
	_ syncfloat64.Counter       = counter[float64, number.Float64Traits]{}
	_ syncfloat64.UpDownCounter = counter[float64, number.Float64Traits]{}
)

// NewCounter returns a value that implements the Counter and UpDownCounter APIs.
func NewCounter[N number.Any, Traits number.Traits[N]](inst *Instrument) counter[N, Traits] {
	return counter[N, Traits]{inst: inst}
}

// Add increments a Counter or UpDownCounter.
func (c counter[N, Traits]) Add(ctx context.Context, incr N, attrs ...attribute.KeyValue) {
	capture[N, Traits](ctx, c.inst, incr, attrs)
}
