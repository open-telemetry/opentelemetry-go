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

package otel

import (
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

type Tracer = trace.Tracer
type TraceProvider = trace.Provider

type Meter = metric.Meter
type MeterProvider = metric.Provider

type Propagators = propagation.Propagators

func GlobalTracer(name string) Tracer {
	return global.Tracer(name)
}

// GlobalTraceProvider returns the registered global trace provider.
func GlobalTraceProvider() TraceProvider {
	return global.TraceProvider()
}

// SetGlobalTraceProvider registers `tp` as the global trace provider.
func SetGlobalTraceProvider(tp TraceProvider) {
	global.SetTraceProvider(tp)
}

func GlobalMeter(name string) Meter {
	return global.Meter(name)
}

// GlobalMeterProvider returns the registered global meter provider.
func GlobalMeterProvider() MeterProvider {
	return global.MeterProvider()
}

// SetGlobalMeterProvider registers `mp` as the global meter provider.
func SetGlobalMeterProvider(mp MeterProvider) {
	global.SetMeterProvider(mp)
}

// GlobalPropagators returns the registered global propagators instance.
func GlobalPropagators() Propagators {
	return global.Propagators()
}

// SetGlobalPropagators registers `p` as the global propagators instance.
func SetGlobalPropagators(p Propagators) {
	global.SetPropagators(p)
}
