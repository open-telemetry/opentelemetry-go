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

This package is currently in a pre-GA phase. Backwards incompatible changes
may be introduced in subsequent minor version releases as we work to track the
evolving OpenTelemetry specification and user feedback.

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

	sr := new(oteltest.StandardSpanRecorder)
	tp := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
*/
package oteltest // import "go.opentelemetry.io/otel/oteltest"
