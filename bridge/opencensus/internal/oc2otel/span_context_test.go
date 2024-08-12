// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oc2otel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"

	"go.opentelemetry.io/otel/bridge/opencensus/internal/otel2oc"

	octrace "go.opencensus.io/trace"
	"go.opencensus.io/trace/tracestate"

	"go.opentelemetry.io/otel/trace"
)

func TestSpanContextConversion(t *testing.T) {
	tsOc, _ := tracestate.New(nil,
		tracestate.Entry{Key: "key1", Value: "value1"},
		tracestate.Entry{Key: "key2", Value: "value2"},
	)
	tsOtel := trace.TraceState{}
	tsOtel, _ = tsOtel.Insert("key2", "value2")
	tsOtel, _ = tsOtel.Insert("key1", "value1")

	httpFormatOc := &tracecontext.HTTPFormat{}

	for _, tc := range []struct {
		description        string
		input              octrace.SpanContext
		expected           trace.SpanContext
		expectedTracestate string
	}{
		{
			description: "empty",
		},
		{
			description: "sampled",
			input: octrace.SpanContext{
				TraceID:      octrace.TraceID([16]byte{1}),
				SpanID:       octrace.SpanID([8]byte{2}),
				TraceOptions: octrace.TraceOptions(0x1),
			},
			expected: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    trace.TraceID([16]byte{1}),
				SpanID:     trace.SpanID([8]byte{2}),
				TraceFlags: trace.FlagsSampled,
			}),
		},
		{
			description: "not sampled",
			input: octrace.SpanContext{
				TraceID:      octrace.TraceID([16]byte{1}),
				SpanID:       octrace.SpanID([8]byte{2}),
				TraceOptions: octrace.TraceOptions(0),
			},
			expected: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID([16]byte{1}),
				SpanID:  trace.SpanID([8]byte{2}),
			}),
		},
		{
			description: "trace state should be propagated",
			input: octrace.SpanContext{
				TraceID:    octrace.TraceID([16]byte{1}),
				SpanID:     octrace.SpanID([8]byte{2}),
				Tracestate: tsOc,
			},
			expected: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    trace.TraceID([16]byte{1}),
				SpanID:     trace.SpanID([8]byte{2}),
				TraceState: tsOtel,
			}),
			expectedTracestate: "key1=value1,key2=value2",
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			output := SpanContext(tc.input)
			assert.Equal(t, tc.expected, output)

			// Ensure the otel tracestate and oc tracestate has the same header output
			_, ts := httpFormatOc.SpanContextToHeaders(tc.input)
			assert.Equal(t, tc.expectedTracestate, ts)
			assert.Equal(t, tc.expectedTracestate, tc.expected.TraceState().String())

			// The reverse conversion should yield the original input
			input := otel2oc.SpanContext(output)
			assert.Equal(t, tc.input, input)
		})
	}
}
