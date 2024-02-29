// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oc2otel

import (
	"testing"

	octrace "go.opencensus.io/trace"
	"go.opencensus.io/trace/tracestate"

	"go.opentelemetry.io/otel/trace"
)

func TestSpanContextConversion(t *testing.T) {
	for _, tc := range []struct {
		description string
		input       octrace.SpanContext
		expected    trace.SpanContext
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
			description: "trace state is ignored",
			input: octrace.SpanContext{
				TraceID:    octrace.TraceID([16]byte{1}),
				SpanID:     octrace.SpanID([8]byte{2}),
				Tracestate: &tracestate.Tracestate{},
			},
			expected: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: trace.TraceID([16]byte{1}),
				SpanID:  trace.SpanID([8]byte{2}),
			}),
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			output := SpanContext(tc.input)
			if !output.Equal(tc.expected) {
				t.Fatalf("Got %+v spancontext, expected %+v.", output, tc.expected)
			}
		})
	}
}
