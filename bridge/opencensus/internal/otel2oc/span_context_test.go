// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otel2oc

import (
	"testing"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/trace"
)

func TestSpanContextConversion(t *testing.T) {
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
	} {
		t.Run(tc.description, func(t *testing.T) {
			output := SpanContext(tc.input)
			if output != tc.expected {
				t.Fatalf("Got %+v spancontext, expected %+v.", output, tc.expected)
			}
		})
	}
}
