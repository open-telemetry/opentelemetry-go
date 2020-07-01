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

package trace_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/api/trace"
)

var (
	traceID    = trace.ID{0, 0, 0, 0, 0, 0, 0, 0x7b, 0, 0, 0, 0, 0, 0, 0x1, 0xc8}
	traceIDStr = "000000000000007b00000000000001c8"
	spanID     = trace.SpanID{0, 0, 0, 0, 0, 0, 0, 0x7b}
	spanIDStr  = "000000000000007b"
)

type Supplier struct {
	SingleHeader       string
	DebugFlagHeader    string
	TraceIDHeader      string
	SpanIDHeader       string
	SampledHeader      string
	ParentSpanIDHeader string
}

func (s *Supplier) Get(key string) string {
	switch key {
	case trace.B3SingleHeader:
		return s.SingleHeader
	case trace.B3DebugFlagHeader:
		return s.DebugFlagHeader
	case trace.B3TraceIDHeader:
		return s.TraceIDHeader
	case trace.B3SpanIDHeader:
		return s.SpanIDHeader
	case trace.B3SampledHeader:
		return s.SampledHeader
	case trace.B3ParentSpanIDHeader:
		return s.ParentSpanIDHeader
	}
	return ""
}

func (s *Supplier) Set(key, value string) {
	fmt.Println("called", key, value)
	switch key {
	case trace.B3SingleHeader:
		s.SingleHeader = value
	case trace.B3DebugFlagHeader:
		s.DebugFlagHeader = value
	case trace.B3TraceIDHeader:
		s.TraceIDHeader = value
	case trace.B3SpanIDHeader:
		s.SpanIDHeader = value
	case trace.B3SampledHeader:
		s.SampledHeader = value
	case trace.B3ParentSpanIDHeader:
		s.ParentSpanIDHeader = value
	}
}

type TestSpan struct {
	trace.NoopSpan
	sc trace.SpanContext
}

func (s TestSpan) SpanContext() trace.SpanContext {
	return s.sc
}

func TestInject(t *testing.T) {
	tests := []struct {
		sc       trace.SpanContext
		b3       trace.B3
		expected *Supplier
	}{
		{
			sc:       trace.SpanContext{},
			b3:       trace.B3{},
			expected: &Supplier{},
		},
		{
			sc: trace.SpanContext{
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
			b3:       trace.B3{},
			expected: &Supplier{},
		},
		{
			sc: trace.SpanContext{
				TraceID:    traceID,
				TraceFlags: trace.FlagsSampled,
			},
			b3:       trace.B3{},
			expected: &Supplier{},
		},
		{
			sc: trace.SpanContext{
				TraceFlags: trace.FlagsSampled,
			},
			b3:       trace.B3{},
			expected: &Supplier{},
		},
		{
			sc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
			b3: trace.B3{SingleHeader: true},
			expected: &Supplier{
				SingleHeader: "000000000000007b00000000000001c8-000000000000007b-1",
			},
		},
		{
			sc: trace.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
			b3: trace.B3{SingleHeader: true},
			expected: &Supplier{
				SingleHeader: "000000000000007b00000000000001c8-000000000000007b-0",
			},
		},
		{
			sc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsUnused,
			},
			b3: trace.B3{SingleHeader: true},
			expected: &Supplier{
				SingleHeader: "000000000000007b00000000000001c8-000000000000007b",
			},
		},
		{
			sc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
			b3: trace.B3{SingleHeader: true, SingleAndMultiHeader: true},
			expected: &Supplier{
				TraceIDHeader: traceIDStr,
				SpanIDHeader:  spanIDStr,
				SampledHeader: "1",
				SingleHeader:  "000000000000007b00000000000001c8-000000000000007b-1",
			},
		},
		{
			sc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsUnused,
			},
			b3: trace.B3{},
			expected: &Supplier{
				TraceIDHeader: traceIDStr,
				SpanIDHeader:  spanIDStr,
			},
		},
		{
			sc: trace.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
			b3: trace.B3{},
			expected: &Supplier{
				TraceIDHeader: traceIDStr,
				SpanIDHeader:  spanIDStr,
				SampledHeader: "0",
			},
		},
		{
			sc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
			b3: trace.B3{},
			expected: &Supplier{
				TraceIDHeader: traceIDStr,
				SpanIDHeader:  spanIDStr,
				SampledHeader: "1",
			},
		},
	}

	for _, test := range tests {
		ctx := trace.ContextWithSpan(context.Background(), TestSpan{sc: test.sc})
		actual := new(Supplier)
		test.b3.Inject(ctx, actual)
		assert.Equal(
			t,
			test.expected,
			actual,
			"B3: %#v, SpanContext: %#v", test.b3, test.sc,
		)
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		traceID      string
		spanID       string
		parentSpanID string
		sampled      string
		flags        string
		single       string
		expected     trace.SpanContext
	}{
		{
			traceID:      "",
			spanID:       "",
			parentSpanID: "",
			sampled:      "",
			flags:        "",
			single:       "",
			expected:     trace.EmptySpanContext(),
		},
		{
			traceID:      "",
			spanID:       "",
			parentSpanID: "",
			sampled:      "",
			flags:        "",
			single:       "000000000000007b00000000000001c8-000000000000007b-1-00000000000001c8",
			expected:     trace.SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled},
		},
		{
			traceID:      traceIDStr,
			spanID:       spanIDStr,
			parentSpanID: "00000000000001c8",
			sampled:      "1",
			flags:        "",
			single:       "",
			expected:     trace.SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled},
		},
		{
			traceID:      traceIDStr,
			spanID:       spanIDStr,
			parentSpanID: "00000000000001c8",
			sampled:      "1",
			flags:        "",
			single:       "000000000000007b00000000000001c8-000000000000007b-1-00000000000001c8",
			expected:     trace.SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled},
		},
	}

	b3 := trace.B3{}
	for _, test := range tests {
		ctx := context.Background()
		supplier := &Supplier{
			SingleHeader:       test.single,
			DebugFlagHeader:    test.flags,
			TraceIDHeader:      test.traceID,
			SpanIDHeader:       test.spanID,
			SampledHeader:      test.sampled,
			ParentSpanIDHeader: test.parentSpanID,
		}
		info := []interface{}{
			"trace ID: %q, span ID: %q, parent span ID: %q, sampled: %q, flags: %q, single: %q",
			test.traceID,
			test.spanID,
			test.parentSpanID,
			test.sampled,
			test.flags,
			test.single,
		}
		actual := trace.RemoteSpanContextFromContext(b3.Extract(ctx, supplier))
		assert.Equal(t, test.expected, actual, info...)
	}
}
