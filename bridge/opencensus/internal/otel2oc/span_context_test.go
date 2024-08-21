// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otel2oc

import (
	"testing"

	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"

	"go.opentelemetry.io/otel/bridge/opencensus/internal/oc2otel"

	"github.com/stretchr/testify/assert"

	"go.opencensus.io/trace/tracestate"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/trace"
)

func TestSpanContextConversion(t *testing.T) {
	tsOc, _ := tracestate.New(nil,
		// Oc has a reverse order of TraceState entries compared to OTel
		tracestate.Entry{Key: "key1", Value: "value1"},
		tracestate.Entry{Key: "key2", Value: "value2"},
	)
	tsOtel := trace.TraceState{}
	tsOtel, _ = tsOtel.Insert("key2", "value2")
	tsOtel, _ = tsOtel.Insert("key1", "value1")

	httpFormatOc := &tracecontext.HTTPFormat{}

	for _, tc := range []struct {
		description        string
		input              trace.SpanContext
		expected           octrace.SpanContext
		expectedTracestate string
	}{
		{
			description: "empty",
		},
		{
			description: "sampled",
			input: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    trace.TraceID([16]byte{1}),
				SpanID:     trace.SpanID([8]byte{2}),
				TraceFlags: trace.FlagsSampled,
			}),
			expected: octrace.SpanContext{
				TraceID:      octrace.TraceID([16]byte{1}),
				SpanID:       octrace.SpanID([8]byte{2}),
				TraceOptions: octrace.TraceOptions(0x1),
			},
		},
		{
			description: "not sampled",
			input: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID([16]byte{1}),
				SpanID:  trace.SpanID([8]byte{2}),
			}),
			expected: octrace.SpanContext{
				TraceID:      octrace.TraceID([16]byte{1}),
				SpanID:       octrace.SpanID([8]byte{2}),
				TraceOptions: octrace.TraceOptions(0),
			},
		},
		{
			description: "trace state should be propagated",
			input: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    trace.TraceID([16]byte{1}),
				SpanID:     trace.SpanID([8]byte{2}),
				TraceState: tsOtel,
			}),
			expected: octrace.SpanContext{
				TraceID:      octrace.TraceID([16]byte{1}),
				SpanID:       octrace.SpanID([8]byte{2}),
				TraceOptions: octrace.TraceOptions(0),
				Tracestate:   tsOc,
			},
			expectedTracestate: "key1=value1,key2=value2",
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			output := SpanContext(tc.input)
			assert.Equal(t, tc.expected, output)

			// Ensure the otel tracestate and oc tracestate has the same header output
			_, ts := httpFormatOc.SpanContextToHeaders(tc.expected)
			assert.Equal(t, tc.expectedTracestate, ts)
			assert.Equal(t, tc.expectedTracestate, tc.input.TraceState().String())

			// The reverse conversion should yield the original input
			input := oc2otel.SpanContext(output)
			assert.Equal(t, tc.input, input)
		})
	}
}
