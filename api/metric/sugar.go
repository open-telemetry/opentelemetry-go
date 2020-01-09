package metric

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/internal"
)

type provider interface {
	Meter() Meter
}

func getMeter(ctx context.Context) Meter {
	if ctx != nil {
		// ctx == nil is passed when the scope is only needed for a namespace
		// value.  these are intended for use in the global context.
		if p, ok := internal.ScopeImpl(ctx).(provider); ok {
			return p.Meter()
		}
	}

	if g, ok := (*atomic.Value)(atomic.LoadPointer(&internal.GlobalScope)).Load().(provider); ok {
		return g.Meter()
	}

	return NoopMeter{}
}

func NewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter {
	return getMeter(nil).NewInt64Counter(name, cos...)
}

func NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	return getMeter(nil).NewFloat64Counter(name, cos...)
}

func NewInt64Gauge(name string, gos ...GaugeOptionApplier) Int64Gauge {
	return getMeter(nil).NewInt64Gauge(name, gos...)
}

func NewFloat64Gauge(name string, gos ...GaugeOptionApplier) Float64Gauge {
	return getMeter(nil).NewFloat64Gauge(name, gos...)
}

func NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	return getMeter(nil).NewInt64Measure(name, mos...)
}

func NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	return getMeter(nil).NewFloat64Measure(name, mos...)
}

func RecordBatch(ctx context.Context, labels []core.KeyValue, ms ...Measurement) {
	getMeter(ctx).RecordBatch(ctx, labels, ms...)
}
