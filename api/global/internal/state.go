// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/propagators"
)

type (
	tracerProviderHolder struct {
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

// TracerProvider is the internal implementation for global.TracerProvider.
func TracerProvider() trace.Provider {
	return globalTracer.Load().(tracerProviderHolder).tp
}

// SetTracerProvider is the internal implementation for global.SetTracerProvider.
func SetTracerProvider(tp trace.Provider) {
	delegateTraceOnce.Do(func() {
		current := TracerProvider()
		if current == tp {
			// Setting the provider to the prior default is nonsense, panic.
			// Panic is acceptable because we are likely still early in the
			// process lifetime.
			panic("invalid Provider, the global instance cannot be reinstalled")
		} else if def, ok := current.(*tracerProvider); ok {
			def.setDelegate(tp)
		}

	})
	globalTracer.Store(tracerProviderHolder{tp: tp})
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
	v.Store(tracerProviderHolder{tp: &tracerProvider{}})
	return v
}

func defaultMeterValue() *atomic.Value {
	v := &atomic.Value{}
	v.Store(meterProviderHolder{mp: newMeterProvider()})
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
	tcPropagator := propagators.TraceContext{}
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
