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

package utils

import (
	"testing"

	"go.opencensus.io/trace/tracestate"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/trace"
)

func TestOTelSpanContextToOC(t *testing.T) {
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
			input: trace.SpanContext{
				TraceID:    trace.TraceID([16]byte{1}),
				SpanID:     trace.SpanID([8]byte{2}),
				TraceFlags: trace.FlagsSampled,
			},
			expected: octrace.SpanContext{
				TraceID:      octrace.TraceID([16]byte{1}),
				SpanID:       octrace.SpanID([8]byte{2}),
				TraceOptions: octrace.TraceOptions(0x1),
			},
		},
		{
			description: "not sampled",
			input: trace.SpanContext{
				TraceID: trace.TraceID([16]byte{1}),
				SpanID:  trace.SpanID([8]byte{2}),
			},
			expected: octrace.SpanContext{
				TraceID:      octrace.TraceID([16]byte{1}),
				SpanID:       octrace.SpanID([8]byte{2}),
				TraceOptions: octrace.TraceOptions(0),
			},
		},
		{
			description: "debug flag",
			input: trace.SpanContext{
				TraceID:    trace.TraceID([16]byte{1}),
				SpanID:     trace.SpanID([8]byte{2}),
				TraceFlags: trace.FlagsDebug,
			},
			expected: octrace.SpanContext{
				TraceID:      octrace.TraceID([16]byte{1}),
				SpanID:       octrace.SpanID([8]byte{2}),
				TraceOptions: octrace.TraceOptions(0),
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			output := OTelSpanContextToOC(tc.input)
			if output != tc.expected {
				t.Fatalf("Got %+v spancontext, exepected %+v.", output, tc.expected)
			}
		})
	}
}

func TestOCSpanContextToOTel(t *testing.T) {
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
			expected: trace.SpanContext{
				TraceID:    trace.TraceID([16]byte{1}),
				SpanID:     trace.SpanID([8]byte{2}),
				TraceFlags: trace.FlagsSampled,
			},
		},
		{
			description: "not sampled",
			input: octrace.SpanContext{
				TraceID:      octrace.TraceID([16]byte{1}),
				SpanID:       octrace.SpanID([8]byte{2}),
				TraceOptions: octrace.TraceOptions(0),
			},
			expected: trace.SpanContext{
				TraceID: trace.TraceID([16]byte{1}),
				SpanID:  trace.SpanID([8]byte{2}),
			},
		},
		{
			description: "trace state is ignored",
			input: octrace.SpanContext{
				TraceID:    octrace.TraceID([16]byte{1}),
				SpanID:     octrace.SpanID([8]byte{2}),
				Tracestate: &tracestate.Tracestate{},
			},
			expected: trace.SpanContext{
				TraceID: trace.TraceID([16]byte{1}),
				SpanID:  trace.SpanID([8]byte{2}),
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			output := OCSpanContextToOTel(tc.input)
			if output != tc.expected {
				t.Fatalf("Got %+v spancontext, exepected %+v.", output, tc.expected)
			}
		})
	}
}
