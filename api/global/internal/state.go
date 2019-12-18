package internal

import (
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
)

type (
	traceProviderHolder struct {
		tp trace.Provider
	}

	meterProviderHolder struct {
		mp metric.Provider
	}
)

var (
	globalTracer = defaultTracerValue()
	globalMeter  = defaultMeterValue()

	delegateMeterOnce sync.Once
)

func TraceProvider() trace.Provider {
	return globalTracer.Load().(traceProviderHolder).tp
}

func SetTraceProvider(tp trace.Provider) {
	globalTracer.Store(traceProviderHolder{tp: tp})
}

func MeterProvider() metric.Provider {
	return globalMeter.Load().(meterProviderHolder).mp
}

func SetMeterProvider(mp metric.Provider) {
	delegateMeterOnce.Do(func() {
		current := MeterProvider()
		if def, ok := current.(*meterProvider); ok {
			def.delegateTo(mp)
		}
	})
	globalMeter.Store(meterProviderHolder{mp: mp})
}

func defaultTracerValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(traceProviderHolder{tp: trace.NoopProvider{}})
	return v
}

func defaultMeterValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(meterProviderHolder{mp: &meterProvider{}})
	return v
}
