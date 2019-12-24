// Copyright 2019, OpenTelemetry Authors
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

package global

import (
	"go.opentelemetry.io/otel/api/global/internal"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
)

// TraceProvider returns the registered global trace provider.
// If none is registered then an instance of trace.NoopProvider is returned.
//
// Use the trace provider to create a named tracer. E.g.
//     tracer := global.TraceProvider().Tracer("example.com/foo")
func TraceProvider() trace.Provider {
	return internal.TraceProvider()
}

// SetTraceProvider registers `tp` as the global trace provider.
func SetTraceProvider(tp trace.Provider) {
	internal.SetTraceProvider(tp)
}

// MeterProvider returns the registered global meter provider.  If
// none is registered then a default meter provider is returned that
// forwards the Meter interface to the first registered Meter.
//
// Use the meter provider to create a named meter. E.g.
//     meter := global.MeterProvider().Meter("example.com/foo")
func MeterProvider() metric.Provider {
	return internal.MeterProvider()
}

// SetMeterProvider registers `mp` as the global meter provider.
func SetMeterProvider(mp metric.Provider) {
	internal.SetMeterProvider(mp)
}
