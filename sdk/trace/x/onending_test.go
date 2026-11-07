// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// enrichingProcessor adds an attribute in OnEnding and records spans in OnEnd.
type enrichingProcessor struct {
	ended []sdktrace.ReadOnlySpan
}

func (*enrichingProcessor) OnStart(context.Context, sdktrace.ReadWriteSpan) {}
func (p *enrichingProcessor) OnEnd(s sdktrace.ReadOnlySpan)                 { p.ended = append(p.ended, s) }
func (*enrichingProcessor) Shutdown(context.Context) error                  { return nil }
func (*enrichingProcessor) ForceFlush(context.Context) error                { return nil }

// OnEnding implements OnEndingSpanProcessor. Registering a processor that
// satisfies this interface with sdktrace.TracerProvider causes it to be called
// during Span.End before the span becomes read-only.
func (*enrichingProcessor) OnEnding(s sdktrace.ReadWriteSpan, _ trace.Span) {
	s.SetAttributes(attribute.Bool("enriched", true))
}

// Compile-time check: enrichingProcessor satisfies OnEndingSpanProcessor.
var _ OnEndingSpanProcessor = (*enrichingProcessor)(nil)

func TestOnEndingSpanProcessorEndToEnd(t *testing.T) {
	exp := tracetest.NewInMemoryExporter()
	p := &enrichingProcessor{}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exp),
	)
	tp.RegisterSpanProcessor(p)

	_, span := tp.Tracer("test").Start(t.Context(), "op")
	span.End()

	spans := exp.GetSpans()
	require.Len(t, spans, 1)

	var found bool
	for _, kv := range spans[0].Attributes {
		if kv.Key == "enriched" && kv.Value.AsBool() {
			found = true
		}
	}
	assert.True(t, found, "attribute set in OnEnding should be present in exported span")
}
