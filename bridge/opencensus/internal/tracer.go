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

package internal // import "go.opentelemetry.io/otel/bridge/opencensus/internal"

import (
	"context"
	"fmt"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/bridge/opencensus/internal/oc2otel"
	"go.opentelemetry.io/otel/trace"
)

// Tracer is an OpenCensus Tracer that wraps an OpenTelemetry Tracer.
type Tracer struct {
	otelTracer trace.Tracer
}

// NewTracer returns an OpenCensus Tracer that wraps the OpenTelemetry tracer.
func NewTracer(tracer trace.Tracer) octrace.Tracer {
	return &Tracer{otelTracer: tracer}
}

// StartSpan starts a new child span of the current span in the context. If
// there is no span in the context, it creates a new trace and span.
func (o *Tracer) StartSpan(ctx context.Context, name string, s ...octrace.StartOption) (context.Context, *octrace.Span) {
	otelOpts, err := oc2otel.StartOptions(s)
	if err != nil {
		Handle(fmt.Errorf("starting span %q: %w", name, err))
	}
	ctx, sp := o.otelTracer.Start(ctx, name, otelOpts...)
	return ctx, NewSpan(sp)
}

// StartSpanWithRemoteParent starts a new child span of the span from the
// given parent.
func (o *Tracer) StartSpanWithRemoteParent(ctx context.Context, name string, parent octrace.SpanContext, s ...octrace.StartOption) (context.Context, *octrace.Span) {
	// make sure span context is zero'd out so we use the remote parent
	ctx = trace.ContextWithSpan(ctx, nil)
	ctx = trace.ContextWithRemoteSpanContext(ctx, oc2otel.SpanContext(parent))
	return o.StartSpan(ctx, name, s...)
}

// FromContext returns the Span stored in a context.
func (o *Tracer) FromContext(ctx context.Context) *octrace.Span {
	return NewSpan(trace.SpanFromContext(ctx))
}

// NewContext returns a new context with the given Span attached.
func (o *Tracer) NewContext(parent context.Context, s *octrace.Span) context.Context {
	if otSpan, ok := s.Internal().(*Span); ok {
		return trace.ContextWithSpan(parent, otSpan.otelSpan)
	}
	Handle(fmt.Errorf("unable to create context with span %q, since it was created using a different tracer", s.String()))
	return parent
}
