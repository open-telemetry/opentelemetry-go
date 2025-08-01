// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opencensus // import "go.opentelemetry.io/otel/bridge/opencensus"

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestNewTraceBridge(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))
	bridge := newTraceBridge([]TraceOption{WithTracerProvider(tp)})
	_, span := bridge.StartSpan(context.Background(), "foo")
	span.End()
	gotSpans := exporter.GetSpans()
	require.Len(t, gotSpans, 1)
	gotSpan := gotSpans[0]
	assert.Equal(t, scopeName, gotSpan.InstrumentationScope.Name)
	assert.Equal(t, gotSpan.InstrumentationScope.Version, Version())
}

func TestOCSpanContextToOTel(t *testing.T) {
	input := octrace.SpanContext{
		TraceID:      [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		SpanID:       [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
		TraceOptions: octrace.TraceOptions(1),
	}
	want := oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
		TraceID:    oteltrace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		SpanID:     oteltrace.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
		TraceFlags: oteltrace.TraceFlags(1),
	})
	got := OCSpanContextToOTel(input)
	assert.Equal(t, want, got)
}

func TestOTelSpanContextToOC(t *testing.T) {
	tests := []struct {
		name     string
		input    oteltrace.SpanContext
		expected octrace.SpanContext
	}{
		{
			name:  "empty span context",
			input: oteltrace.SpanContext{},
			expected: octrace.SpanContext{
				TraceID:      [16]byte{},
				SpanID:       [8]byte{},
				TraceOptions: octrace.TraceOptions(0),
			},
		},
		{
			name: "sampled span context",
			input: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
				TraceID:    oteltrace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				SpanID:     oteltrace.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
				TraceFlags: oteltrace.FlagsSampled,
			}),
			expected: octrace.SpanContext{
				TraceID:      [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				SpanID:       [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
				TraceOptions: octrace.TraceOptions(1),
			},
		},
		{
			name: "not sampled span context",
			input: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
				TraceID: oteltrace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				SpanID:  oteltrace.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
			}),
			expected: octrace.SpanContext{
				TraceID:      [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				SpanID:       [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
				TraceOptions: octrace.TraceOptions(0),
			},
		},
		{
			name: "span context with tracestate",
			input: func() oteltrace.SpanContext {
				ts := oteltrace.TraceState{}
				ts, _ = ts.Insert("key1", "value1")
				ts, _ = ts.Insert("key2", "value2")
				return oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
					TraceID:    oteltrace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
					SpanID:     oteltrace.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
					TraceState: ts,
				})
			}(),
			expected: octrace.SpanContext{
				TraceID:      [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				SpanID:       [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
				TraceOptions: octrace.TraceOptions(0),
				// Tracestate will be set by the conversion function
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OTelSpanContextToOC(tt.input)

			assert.Equal(t, tt.expected.TraceID, got.TraceID, "TraceID should be correctly converted")

			assert.Equal(t, tt.expected.SpanID, got.SpanID, "SpanID should be correctly converted")

			assert.Equal(t, tt.expected.TraceOptions, got.TraceOptions, "TraceOptions should be correctly converted")

			// Verify Tracestate is populated when input has tracestate
			if tt.input.TraceState().Len() > 0 {
				assert.NotNil(t, got.Tracestate, "Tracestate should be populated when input has tracestate")
				// Verify the tracestate entries are preserved
				expectedTraceState := tt.input.TraceState().String()
				gotEntries := got.Tracestate.Entries()

				// Convert entries back to a string representation for comparison
				var gotTraceStateEntries []string
				for _, entry := range gotEntries {
					gotTraceStateEntries = append(gotTraceStateEntries, entry.Key+"="+entry.Value)
				}
				gotTraceState := ""
				if len(gotTraceStateEntries) > 0 {
					gotTraceState = strings.Join(gotTraceStateEntries, ",")
				}
				assert.Equal(t, expectedTraceState, gotTraceState, "Tracestate should preserve entries")
			} else if got.Tracestate != nil {
				// For empty tracestate cases, ensure the field is properly handled
				entries := got.Tracestate.Entries()
				assert.Empty(t, entries, "Empty tracestate should result in empty entries")
			}
		})
	}
}

func TestInstallTraceBridge(t *testing.T) {
	originalTracer := octrace.DefaultTracer
	defer func() {
		octrace.DefaultTracer = originalTracer
	}()

	tests := []struct {
		name             string
		opts             []TraceOption
		expectValidSpans bool
	}{
		{
			name:             "install with default options",
			opts:             nil,
			expectValidSpans: false,
		},
		{
			name: "install with custom tracer provider",
			opts: []TraceOption{
				WithTracerProvider(trace.NewTracerProvider()),
			},
			expectValidSpans: true,
		},
		{
			name: "install with tracer provider with exporter",
			opts: []TraceOption{
				WithTracerProvider(
					trace.NewTracerProvider(
						trace.WithSyncer(tracetest.NewInMemoryExporter()),
					),
				),
			},
			expectValidSpans: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeTracer := octrace.DefaultTracer

			InstallTraceBridge(tt.opts...)

			assert.NotEqual(
				t,
				beforeTracer,
				octrace.DefaultTracer,
				"DefaultTracer should be updated",
			)
			assert.NotNil(
				t,
				octrace.DefaultTracer,
				"DefaultTracer should not be nil",
			)

			ctx, span := octrace.DefaultTracer.StartSpan(
				context.Background(),
				"test-span",
			)
			assert.NotNil(
				t,
				span,
				"Should be able to create spans",
			)
			assert.NotNil(t, ctx, "Should return a valid context")

			spanContext := span.SpanContext()
			if tt.expectValidSpans {
				assert.NotEqual(
					t,
					octrace.TraceID{},
					spanContext.TraceID,
					"Span should have a non-zero TraceID",
				)
				assert.NotEqual(
					t,
					octrace.SpanID{},
					spanContext.SpanID,
					"Span should have a non-zero SpanID",
				)
			}

			span.End()

			spanFromContext := octrace.DefaultTracer.FromContext(ctx)
			assert.NotNil(
				t,
				spanFromContext,
				"Should be able to get span from context",
			)
		})
	}
}
