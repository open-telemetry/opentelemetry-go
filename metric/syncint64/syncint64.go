package syncint64

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
)

type Instruments interface {
	Counter(name string, opts ...instrument.Option) (Counter, error)
	UpDownCounter(name string, opts ...instrument.Option) (UpDownCounter, error)
	Histogram(name string, opts ...instrument.Option) (Histogram, error)
}

type Counter interface {
	Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue)

	instrument.Synchronous
}

type UpDownCounter interface {
	Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue)

	instrument.Synchronous
}

type Histogram interface {
	Record(ctx context.Context, incr int64, attrs ...attribute.KeyValue)

	instrument.Synchronous
}
