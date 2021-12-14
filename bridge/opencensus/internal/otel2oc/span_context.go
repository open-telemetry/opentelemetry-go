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

package otel2oc // import "go.opentelemetry.io/otel/bridge/opencensus/internal/otel2oc"

import (
	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/trace"
)

func SpanContext(sc trace.SpanContext) octrace.SpanContext {
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
