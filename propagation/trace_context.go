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

package propagation // import "go.opentelemetry.io/otel/propagation"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

const (
	supportedVersion  = 0
	traceparentHeader = "traceparent"
	tracestateHeader  = "tracestate"
)

// TraceContext is a propagator that supports the W3C Trace Context format
// (https://www.w3.org/TR/trace-context/)
//
// This propagator will propagate the traceparent and tracestate headers to
// guarantee traces are not broken. It is up to the users of this propagator
// to choose if they want to participate in a trace by modifying the
// traceparent header and relevant parts of the tracestate header containing
// their proprietary information.
type TraceContext struct{}

var _ TextMapPropagator = TraceContext{}

// Inject set tracecontext from the Context into the carrier.
func (tc TraceContext) Inject(ctx context.Context, carrier TextMapCarrier) {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return
	}

	if ts := sc.TraceState().String(); ts != "" {
		carrier.Set(tracestateHeader, ts)
	}

	// Clear all flags other than the trace-context supported sampling bit.
	flags := sc.TraceFlags() & trace.FlagsSampled

	h := fmt.Sprintf("%.2x-%s-%s-%s",
		supportedVersion,
		sc.TraceID(),
		sc.SpanID(),
		flags)
	carrier.Set(traceparentHeader, h)
}

// Extract reads tracecontext from the carrier into a returned Context.
//
// The returned Context will be a copy of ctx and contain the extracted
// tracecontext as the remote SpanContext. If the extracted tracecontext is
// invalid, the passed ctx will be returned directly instead.
func (tc TraceContext) Extract(ctx context.Context, carrier TextMapCarrier) context.Context {
	sc := tc.extract(carrier)
	if !sc.IsValid() {
		return ctx
	}
	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func (tc TraceContext) extract(carrier TextMapCarrier) trace.SpanContext {
	traceparent := carrier.Get(traceparentHeader)
	tracestate := carrier.Get(tracestateHeader)
	return trace.NewSpanContextFromTrace(traceparent, tracestate)
}

// Fields returns the keys who's values are set with Inject.
func (tc TraceContext) Fields() []string {
	return []string{traceparentHeader, tracestateHeader}
}
