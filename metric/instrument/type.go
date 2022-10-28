package instrument

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// SyncInstrumentProvider provides access to synchronous instruments.
type SyncInstrumentProvider[T int64 | float64] interface {
	// Counter creates an instrument for recording increasing values.
	Counter(name string, opts ...Option) (SyncCounter[T], error)
	// UpDownCounter creates an instrument for recording changes of a value.
	UpDownCounter(name string, opts ...Option) (SyncUpDownCounter[T], error)
	// Histogram creates an instrument for recording a distribution of values.
	Histogram(name string, opts ...Option) (SyncHistogram[T], error)
}

// Counter is an instrument that records increasing values.
type SyncCounter[T int64 | float64] interface {
	// Add records a change to the counter.
	Add(ctx context.Context, incr T, attrs ...attribute.KeyValue)

	Synchronous
}

// UpDownCounter is an instrument that records increasing or decreasing values.
type SyncUpDownCounter[T int64 | float64] interface {
	// Add records a change to the counter.
	Add(ctx context.Context, incr T, attrs ...attribute.KeyValue)

	Synchronous
}

// Histogram is an instrument that records a distribution of values.
type SyncHistogram[T int64 | float64] interface {
	// Record adds an additional value to the distribution.
	Record(ctx context.Context, incr T, attrs ...attribute.KeyValue)

	Synchronous
}

// AsyncInstrumentProvider provides access to asynchronous instruments.
type AsyncInstrumentProvider[T int64 | float64] interface {
	// Counter creates an instrument for recording increasing values.
	Counter(name string, opts ...Option) (AsyncCounter[T], error)
	// UpDownCounter creates an instrument for recording changes of a value.
	UpDownCounter(name string, opts ...Option) (AsyncUpDownCounter[T], error)
	// Histogram creates an instrument for recording a distribution of values.
	Gauge(name string, opts ...Option) (AsyncGauge[T], error)
}

// Counter is an instrument that records increasing values.
type AsyncCounter[T int64 | float64] interface {
	// Add records a change to the counter.
	Observe(ctx context.Context, x T, attrs ...attribute.KeyValue)

	Asynchronous
}

// UpDownCounter is an instrument that records increasing or decreasing values.
type AsyncUpDownCounter[T int64 | float64] interface {
	// Add records a change to the counter.
	Observe(ctx context.Context, x T, attrs ...attribute.KeyValue)

	Asynchronous
}

// Gauge is an instrument that records independent readings.
type AsyncGauge[T int64 | float64] interface {
	// Observe records the state of the instrument to be x.
	//
	// It is only valid to call this within a callback. If called outside of the
	// registered callback it should have no effect on the instrument, and an
	// error will be reported via the error handler.
	Observe(ctx context.Context, x T, attrs ...attribute.KeyValue)

	Asynchronous
}
