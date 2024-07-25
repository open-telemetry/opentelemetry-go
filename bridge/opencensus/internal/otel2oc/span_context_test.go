// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otel2oc

import (
	"go.opencensus.io/trace/tracestate"
	"strings"
	"testing"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/trace"
)

func TestSpanContextConversion(t *testing.T) {
	tsOc, _ := tracestate.New(nil,
		tracestate.Entry{"key1", "value1"},
		tracestate.Entry{"key2", "value2"},
	)
	tsOtel := trace.TraceState{}
	tsOtel, _ = tsOtel.Insert("key1", "value1")
	tsOtel, _ = tsOtel.Insert("key2", "value2")

	for _, tc := range []struct {
		description string
		input       trace.SpanContext
		expected    octrace.SpanContext
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
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			output := SpanContext(tc.input)
			if !equal(output, tc.expected) {
				t.Fatalf("Got %+v spancontext, expected %+v.", toString(output.Tracestate), toString(tc.expected.Tracestate))
			}
		})
	}
}

func equal(t1, t2 octrace.SpanContext) bool {
	return t1.IsSampled() == t2.IsSampled() &&
		t1.SpanID == t2.SpanID &&
		t1.TraceID == t2.TraceID &&
		t1.TraceOptions == t2.TraceOptions &&
		toString(t1.Tracestate) == toString(t2.Tracestate)
}

func toString(t *tracestate.Tracestate) string {
	result := new(strings.Builder)
	for _, e := range t.Entries() {
		result.WriteString(e.Key)
		result.WriteString("=")
		result.WriteString(e.Value)
		result.WriteString(",")
	}
	return result.String()
}
