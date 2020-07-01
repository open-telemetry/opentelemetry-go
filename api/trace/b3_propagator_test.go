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

package trace

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	traceID    = ID{0, 0, 0, 0, 0, 0, 0, 0x7b, 0, 0, 0, 0, 0, 0, 0x1, 0xc8}
	traceIDStr = "000000000000007b00000000000001c8"
	spanID     = SpanID{0, 0, 0, 0, 0, 0, 0, 0x7b}
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
	case B3SingleHeader:
		return s.SingleHeader
	case B3DebugFlagHeader:
		return s.DebugFlagHeader
	case B3TraceIDHeader:
		return s.TraceIDHeader
	case B3SpanIDHeader:
		return s.SpanIDHeader
	case B3SampledHeader:
		return s.SampledHeader
	case B3ParentSpanIDHeader:
		return s.ParentSpanIDHeader
	}
	return ""
}

func (s *Supplier) Set(key, value string) {
	fmt.Println("called", key, value)
	switch key {
	case B3SingleHeader:
		s.SingleHeader = value
	case B3DebugFlagHeader:
		s.DebugFlagHeader = value
	case B3TraceIDHeader:
		s.TraceIDHeader = value
	case B3SpanIDHeader:
		s.SpanIDHeader = value
	case B3SampledHeader:
		s.SampledHeader = value
	case B3ParentSpanIDHeader:
		s.ParentSpanIDHeader = value
	}
}

type TestSpan struct {
	NoopSpan
	sc SpanContext
}

func (s TestSpan) SpanContext() SpanContext {
	return s.sc
}

func TestInject(t *testing.T) {
	tests := []struct {
		sc       SpanContext
		b3       B3
		expected *Supplier
	}{
		{
			sc:       SpanContext{},
			b3:       B3{},
			expected: &Supplier{},
		},
		{
			sc: SpanContext{
				SpanID:     spanID,
				TraceFlags: FlagsSampled,
			},
			b3:       B3{},
			expected: &Supplier{},
		},
		{
			sc: SpanContext{
				TraceID:    traceID,
				TraceFlags: FlagsSampled,
			},
			b3:       B3{},
			expected: &Supplier{},
		},
		{
			sc: SpanContext{
				TraceFlags: FlagsSampled,
			},
			b3:       B3{},
			expected: &Supplier{},
		},
		{
			sc: SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: FlagsSampled,
			},
			b3: B3{SingleHeader: true},
			expected: &Supplier{
				SingleHeader: "000000000000007b00000000000001c8-000000000000007b-1",
			},
		},
		{
			sc: SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
			b3: B3{SingleHeader: true},
			expected: &Supplier{
				SingleHeader: "000000000000007b00000000000001c8-000000000000007b-0",
			},
		},
		{
			sc: SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: FlagsUnused,
			},
			b3: B3{SingleHeader: true},
			expected: &Supplier{
				SingleHeader: "000000000000007b00000000000001c8-000000000000007b",
			},
		},
		{
			sc: SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: FlagsSampled,
			},
			b3: B3{SingleHeader: true, SingleAndMultiHeader: true},
			expected: &Supplier{
				TraceIDHeader: traceIDStr,
				SpanIDHeader:  spanIDStr,
				SampledHeader: "1",
				SingleHeader:  "000000000000007b00000000000001c8-000000000000007b-1",
			},
		},
		{
			sc: SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: FlagsUnused,
			},
			b3: B3{},
			expected: &Supplier{
				TraceIDHeader: traceIDStr,
				SpanIDHeader:  spanIDStr,
			},
		},
		{
			sc: SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
			b3: B3{},
			expected: &Supplier{
				TraceIDHeader: traceIDStr,
				SpanIDHeader:  spanIDStr,
				SampledHeader: "0",
			},
		},
		{
			sc: SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: FlagsSampled,
			},
			b3: B3{},
			expected: &Supplier{
				TraceIDHeader: traceIDStr,
				SpanIDHeader:  spanIDStr,
				SampledHeader: "1",
			},
		},
	}

	for _, test := range tests {
		ctx := ContextWithSpan(context.Background(), TestSpan{sc: test.sc})
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
		expected     SpanContext
	}{
		{
			traceID:      "",
			spanID:       "",
			parentSpanID: "",
			sampled:      "",
			flags:        "",
			single:       "",
			expected:     empty,
		},
		{
			traceID:      "",
			spanID:       "",
			parentSpanID: "",
			sampled:      "",
			flags:        "",
			single:       "000000000000007b00000000000001c8-000000000000007b-1-00000000000001c8",
			expected:     SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsSampled},
		},
		{
			traceID:      traceIDStr,
			spanID:       spanIDStr,
			parentSpanID: "00000000000001c8",
			sampled:      "1",
			flags:        "",
			single:       "",
			expected:     SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsSampled},
		},
		{
			traceID:      traceIDStr,
			spanID:       spanIDStr,
			parentSpanID: "00000000000001c8",
			sampled:      "1",
			flags:        "",
			single:       "000000000000007b00000000000001c8-000000000000007b-1-00000000000001c8",
			expected:     SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsSampled},
		},
	}

	b3 := B3{}
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
		actual := RemoteSpanContextFromContext(b3.Extract(ctx, supplier))
		assert.Equal(t, test.expected, actual, info...)
	}
}

func TestExtractMultiple(t *testing.T) {
	tests := []struct {
		traceID      string
		spanID       string
		parentSpanID string
		sampled      string
		flags        string
		expected     SpanContext
		err          error
	}{
		{
			"", "", "", "0", "",
			SpanContext{},
			nil,
		},
		{
			"", "", "", "", "",
			SpanContext{TraceFlags: FlagsUnused},
			nil,
		},
		{
			"", "", "", "1", "",
			SpanContext{TraceFlags: FlagsSampled},
			nil,
		},
		{
			traceIDStr, spanIDStr, "", "", "",
			SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsUnused},
			nil,
		},
		{
			traceIDStr, spanIDStr, "", "0", "",
			SpanContext{TraceID: traceID, SpanID: spanID},
			nil,
		},
		// Ensure backwards compatibility.
		{
			traceIDStr, spanIDStr, "", "false", "",
			SpanContext{TraceID: traceID, SpanID: spanID},
			nil,
		},
		{
			traceIDStr, spanIDStr, "", "1", "",
			SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsSampled},
			nil,
		},
		// Ensure backwards compatibility.
		{
			traceIDStr, spanIDStr, "", "true", "",
			SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsSampled},
			nil,
		},
		{
			traceIDStr, spanIDStr, "", "a", "",
			empty,
			errInvalidSampledHeader,
		},
		{
			traceIDStr, spanIDStr, "", "1", "1",
			SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsUnused},
			nil,
		},
		// Invalid flags are discarded.
		{
			traceIDStr, spanIDStr, "", "1", "invalid",
			SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsSampled},
			nil,
		},
		// Support short trace IDs.
		{
			"00000000000001c8", spanIDStr, "", "0", "",
			SpanContext{
				TraceID: ID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x1, 0xc8},
				SpanID:  spanID,
			},
			nil,
		},
		{
			"00000000000001c", spanIDStr, "", "0", "",
			empty,
			errInvalidTraceIDHeader,
		},
		{
			"00000000000001c80", spanIDStr, "", "0", "",
			empty,
			errInvalidTraceIDHeader,
		},
		{
			traceIDStr[:len(traceIDStr)-2], spanIDStr, "", "0", "",
			empty,
			errInvalidTraceIDHeader,
		},
		{
			traceIDStr + "0", spanIDStr, "", "0", "",
			empty,
			errInvalidTraceIDHeader,
		},
		{
			traceIDStr, "00000000000001c", "", "0", "",
			empty,
			errInvalidSpanIDHeader,
		},
		{
			traceIDStr, "00000000000001c80", "", "0", "",
			empty,
			errInvalidSpanIDHeader,
		},
		{
			traceIDStr, "", "", "0", "",
			empty,
			errInvalidScope,
		},
		{
			"", spanIDStr, "", "0", "",
			empty,
			errInvalidScope,
		},
		{
			"", "", spanIDStr, "0", "",
			empty,
			errInvalidScopeParent,
		},
		{
			traceIDStr, spanIDStr, "00000000000001c8", "0", "",
			SpanContext{TraceID: traceID, SpanID: spanID},
			nil,
		},
		{
			traceIDStr, spanIDStr, "00000000000001c", "0", "",
			empty,
			errInvalidParentSpanIDHeader,
		},
		{
			traceIDStr, spanIDStr, "00000000000001c80", "0", "",
			empty,
			errInvalidParentSpanIDHeader,
		},
	}

	for _, test := range tests {
		actual, err := extractMultiple(
			test.traceID,
			test.spanID,
			test.parentSpanID,
			test.sampled,
			test.flags,
		)
		info := []interface{}{
			"trace ID: %q, span ID: %q, parent span ID: %q, sampled: %q, flags: %q",
			test.traceID,
			test.spanID,
			test.parentSpanID,
			test.sampled,
			test.flags,
		}
		if !assert.Equal(t, test.err, err, info...) {
			continue
		}
		assert.Equal(t, test.expected, actual, info...)
	}
}

func TestExtractSingle(t *testing.T) {
	tests := []struct {
		header   string
		expected SpanContext
		err      error
	}{
		{"0", SpanContext{}, nil},
		{"1", SpanContext{TraceFlags: FlagsSampled}, nil},
		// debug flag is valid, but ignored.
		{"d", SpanContext{}, nil},
		{"a", empty, errInvalidSampledByte},
		{"3", empty, errInvalidSampledByte},
		{"000000000000007b", empty, errInvalidScope},
		{"000000000000007b00000000000001c8", empty, errInvalidScope},
		// Support short trace IDs.
		{
			"00000000000001c8-000000000000007b",
			SpanContext{
				TraceID: ID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x1, 0xc8},
				SpanID:  spanID,
			},
			nil,
		},
		{
			"000000000000007b00000000000001c8-000000000000007b",
			SpanContext{TraceID: traceID, SpanID: spanID},
			nil,
		},
		{
			"000000000000007b00000000000001c8-000000000000007b-",
			empty,
			errInvalidSampledByte,
		},
		{
			"000000000000007b00000000000001c8-000000000000007b-3",
			empty,
			errInvalidSampledByte,
		},
		{
			"000000000000007b00000000000001c8-000000000000007b-00000000000001c8",
			empty,
			errInvalidScopeParentSingle,
		},
		{
			"000000000000007b00000000000001c8-000000000000007b-1",
			SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsSampled},
			nil,
		},
		// ParentSpanID is discarded, but should still restult in a parsable
		// header.
		{
			"000000000000007b00000000000001c8-000000000000007b-1-00000000000001c8",
			SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsSampled},
			nil,
		},
		{
			"000000000000007b00000000000001c8-000000000000007b-1-00000000000001c",
			empty,
			errInvalidParentSpanIDValue,
		},
		{"", empty, errEmptyContext},
	}

	for _, test := range tests {
		actual, err := extractSingle(test.header)
		if !assert.Equal(t, test.err, err, "header: %s", test.header) {
			continue
		}
		assert.Equal(t, test.expected, actual, "header: %s", test.header)
	}
}
