package internal

import (
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

type (
	traceProviderHolder struct {
		tp trace.Provider
	}

	meterProviderHolder struct {
		mp metric.Provider
	}

	propagatorsHolder struct {
		pr propagation.Propagators
	}
)

var (
	globalTracer      = defaultTracerValue()
	globalMeter       = defaultMeterValue()
	globalPropagators = defaultPropagatorsValue()

	delegateMeterOnce sync.Once
	delegateTraceOnce sync.Once
)

// TraceProvider is the internal implementation for global.TraceProvider.
func TraceProvider() trace.Provider {
	return globalTracer.Load().(traceProviderHolder).tp
}

// SetTraceProvider is the internal implementation for global.SetTraceProvider.
func SetTraceProvider(tp trace.Provider) {
	delegateTraceOnce.Do(func() {
		current := TraceProvider()
		if current == tp {
			// Setting the provider to the prior default is nonsense, panic.
			// Panic is acceptable because we are likely still early in the
			// process lifetime.
			panic("invalid Provider, the global instance cannot be reinstalled")
		} else if def, ok := current.(*traceProvider); ok {
			def.setDelegate(tp)
		}

	})
	globalTracer.Store(traceProviderHolder{tp: tp})
}

// MeterProvider is the internal implementation for global.MeterProvider.
func MeterProvider() metric.Provider {
	return globalMeter.Load().(meterProviderHolder).mp
}

// SetMeterProvider is the internal implementation for global.SetMeterProvider.
func SetMeterProvider(mp metric.Provider) {
	delegateMeterOnce.Do(func() {
		current := MeterProvider()

		if current == mp {
			// Setting the provider to the prior default is nonsense, panic.
			// Panic is acceptable because we are likely still early in the
			// process lifetime.
			panic("invalid Provider, the global instance cannot be reinstalled")
		} else if def, ok := current.(*meterProvider); ok {
			def.setDelegate(mp)
		}
	})
	globalMeter.Store(meterProviderHolder{mp: mp})
}

// Propagators is the internal implementation for global.Propagators.
func Propagators() propagation.Propagators {
	return globalPropagators.Load().(propagatorsHolder).pr
}

// SetPropagators is the internal implementation for global.SetPropagators.
func SetPropagators(pr propagation.Propagators) {
	globalPropagators.Store(propagatorsHolder{pr: pr})
}

func defaultTracerValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(traceProviderHolder{tp: &traceProvider{}})
	return v
}

func defaultMeterValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(meterProviderHolder{mp: &meterProvider{}})
	return v
}

func defaultPropagatorsValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(propagatorsHolder{pr: getDefaultPropagators()})
	return v
}

// getDefaultPropagators returns a default Propagators, configured
// with W3C trace and correlation context propagation.
func getDefaultPropagators() propagation.Propagators {
	tcPropagator := trace.TraceContext{}
	ccPropagator := correlation.CorrelationContext{}
	return propagation.New(
		propagation.WithExtractors(tcPropagator, ccPropagator),
		propagation.WithInjectors(tcPropagator, ccPropagator),
	)
}

// ResetForTest restores the initial global state, for testing purposes.
func ResetForTest() {
	globalTracer = defaultTracerValue()
	globalMeter = defaultMeterValue()
	globalPropagators = defaultPropagatorsValue()
	delegateMeterOnce = sync.Once{}
	delegateTraceOnce = sync.Once{}
}
