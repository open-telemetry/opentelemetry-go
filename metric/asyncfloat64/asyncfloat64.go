package asyncfloat64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
)

type Instruments interface {
	Counter(name string, opts ...instrument.Option) Counter
	UpDownCounter(name string, opts ...instrument.Option) UpDownCounter
	Gauge(name string, opts ...instrument.Option) Gauge
}

type Counter interface {
	Observe(ctx context.Context, incr float64, attrs ...attribute.KeyValue)

	instrument.Asynchronous
}

type UpDownCounter interface {
	Observe(ctx context.Context, incr float64, attrs ...attribute.KeyValue)

	instrument.Asynchronous
}

type Gauge interface {
	Observe(ctx context.Context, incr float64, attrs ...attribute.KeyValue)

	instrument.Asynchronous
}
