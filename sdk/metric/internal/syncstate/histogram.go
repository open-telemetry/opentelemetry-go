package syncstate

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/metric/number"
)

// histogram is a synchronous instrument having a Record() method.
type histogram[N number.Any, Traits number.Traits[N]] struct {
	instrument.SynchronousStruct

	inst *Instrument
}

// histogram satisfies 2 instrument APIs.
var (
	_ syncint64.Histogram   = histogram[int64, number.Int64Traits]{}
	_ syncfloat64.Histogram = histogram[float64, number.Float64Traits]{}
)

// NewCounter returns a value that implements the Histogram API.
func NewHistogram[N number.Any, Traits number.Traits[N]](inst *Instrument) histogram[N, Traits] {
	return histogram[N, Traits]{inst: inst}
}

// Record records a Histogram observation.
func (h histogram[N, Traits]) Record(ctx context.Context, incr N, attrs ...attribute.KeyValue) {
	capture[N, Traits](ctx, h.inst, incr, attrs)
}
