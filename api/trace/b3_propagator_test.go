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
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	traceID    = ID{0, 0, 0, 0, 0, 0, 0, 0x7b, 0, 0, 0, 0, 0, 0, 0x1, 0xc8}
	traceIDStr = "000000000000007b00000000000001c8"
	spanID     = SpanID{0, 0, 0, 0, 0, 0, 0, 0x7b}
	spanIDStr  = "000000000000007b"
)

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
			SpanContext{TraceFlags: FlagsUnset},
			nil,
		},
		{
			"", "", "", "1", "",
			SpanContext{TraceFlags: FlagsSampled},
			nil,
		},
		{
			traceIDStr, spanIDStr, "", "", "",
			SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsUnset},
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
			SpanContext{TraceID: traceID, SpanID: spanID, TraceFlags: FlagsUnset},
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
