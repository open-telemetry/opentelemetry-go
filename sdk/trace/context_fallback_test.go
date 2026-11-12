// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

type contextFallbackTestOption struct {
	TracerProviderOption
	parent trace.SpanContext
}

func (contextFallbackTestOption) Experimental() {}

func (o contextFallbackTestOption) ContextFallback() (trace.SpanContext, baggage.Baggage, int) {
	return o.parent, baggage.Baggage{}, 0
}

func TestExperimentalContextFallback(t *testing.T) {
	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{1}, SpanID: trace.SpanID{2}, Remote: true,
		TraceFlags: trace.FlagsSampled,
	})
	tp := NewTracerProvider(contextFallbackTestOption{parent: parent})
	t.Cleanup(func() { assert.NoError(t, tp.Shutdown(t.Context())) })
	_, span := tp.Tracer("test").Start(t.Context(), "job")
	defer span.End()
	assert.Equal(t, parent.TraceID(), span.SpanContext().TraceID())
	assert.True(t, span.SpanContext().IsSampled())
}
