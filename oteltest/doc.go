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

/*
Package oteltest provides testing utilities for the otel package.

This package is currently in a Release Candidate phase. Backwards incompatible changes
may be introduced prior to v1.0.0, but we believe the current API is ready to stabilize.

API Validation

The Harness can be used to validate an implementation of the OpenTelemetry API
defined by the `otel` package.

	func TestCustomSDKTracingImplementation(t *testing.T) {
		yourTraceProvider := NewTracerProvider()
		subjectFactory := func() otel.Tracer {
			return yourTraceProvider.Tracer("testing")
		}

		oteltest.NewHarness(t).TestTracer(subjectFactory)
	}

Currently the Harness only provides testing of the trace portion of the
OpenTelemetry API.

Trace Testing

To test tracing functionality a full testing implementation of the
OpenTelemetry tracing API are provided. The provided TracerProvider, Tracer,
and Span all implement their related interface and are designed to allow
introspection of their state and history. Additionally, a SpanRecorder can be
provided to the TracerProvider to record all Spans started and ended by the
testing structures.

	sr := new(oteltest.SpanRecorder)
	tp := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))

Deprecated: This package contains an alternate implementation of the
OpenTelemetry SDK. This means it will diverge from the one commonly used in
real world operations and therefore will not provide adequate testing
guarantees for users. Because of this, this package should not be used to test
performance or behavior of code running OpenTelemetry. Instead, the
SpanRecorder from the go.opentelemetry.io/otel/sdk/trace/tracetest package can
be registered with the default SDK (go.opentelemetry.io/otel/sdk/trace) as a
SpanProcessor and used to test. This will ensure code will work with the
default SDK. If users do not want to include a dependency on the default SDK
it is recommended to run integration tests in their own module to isolate the
dependency (see go.opentelemetry.io/otel/bridge/opencensus/test as an
example).
*/
package oteltest // import "go.opentelemetry.io/otel/oteltest"
