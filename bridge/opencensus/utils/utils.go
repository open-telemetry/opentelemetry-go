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

package utils // import "go.opentelemetry.io/otel/bridge/opencensus/utils"

import (
	"fmt"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// OTelSpanContextToOC converts from an OpenTelemetry SpanContext to an
// OpenCensus SpanContext, and handles any incompatibilities with the global
// error handler.
func OTelSpanContextToOC(sc trace.SpanContext) octrace.SpanContext {
	if sc.IsDebug() || sc.IsDeferred() {
		otel.Handle(fmt.Errorf("ignoring OpenTelemetry Debug or Deferred trace flags for span %q because they are not supported by OpenCensus", sc.SpanID()))
	}
	var to octrace.TraceOptions
	if sc.IsSampled() {
		// OpenCensus doesn't expose functions to directly set sampled
		to = 0x1
	}
	return octrace.SpanContext{
		TraceID:      octrace.TraceID(sc.TraceID()),
		SpanID:       octrace.SpanID(sc.SpanID()),
		TraceOptions: to,
	}
}

// OCSpanContextToOTel converts from an OpenCensus SpanContext to an
// OpenTelemetry SpanContext.
func OCSpanContextToOTel(sc octrace.SpanContext) trace.SpanContext {
	var traceFlags byte
	if sc.IsSampled() {
		traceFlags = trace.FlagsSampled
	}
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID(sc.TraceID),
		SpanID:     trace.SpanID(sc.SpanID),
		TraceFlags: traceFlags,
	})
}
