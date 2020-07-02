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

package testtrace_test

import (
	"fmt"

	"go.opentelemetry.io/otel/api/trace"
)

type extractTest struct {
	name    string
	headers map[string]string
	wantSc  trace.SpanContext
}

var (
	traceID64bitPadded = mustTraceIDFromHex("0000000000000000a3ce929d0e0e4736")
)

var extractHeaders = []extractTest{
	{
		name:    "empty",
		headers: map[string]string{},
		wantSc:  trace.EmptySpanContext(),
	},
	{
		name: "multiple: sampling state defer",
		headers: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  spanIDStr,
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
	},
	{
		name: "multiple: sampling state deny",
		headers: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  spanIDStr,
			trace.B3SampledHeader: "0",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsNotSampled,
		},
	},
	{
		name: "multiple: sampling state accept",
		headers: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  spanIDStr,
			trace.B3SampledHeader: "1",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "multiple: sampling state as a boolean",
		headers: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  spanIDStr,
			trace.B3SampledHeader: "true",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "multiple: debug flag set",
		headers: map[string]string{
			trace.B3TraceIDHeader:   traceIDStr,
			trace.B3SpanIDHeader:    spanIDStr,
			trace.B3DebugFlagHeader: "1",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
	},
	{
		name: "multiple: debug flag set to not 1 (ignored)",
		headers: map[string]string{
			trace.B3TraceIDHeader:   traceIDStr,
			trace.B3SpanIDHeader:    spanIDStr,
			trace.B3SampledHeader:   "1",
			trace.B3DebugFlagHeader: "2",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		// spec explicitly states "Debug implies an accept decision, so don't
		// also send the X-B3-Sampled header", make sure sampling is
		// deferred.
		name: "multiple: debug flag set and sampling state is deny",
		headers: map[string]string{
			trace.B3TraceIDHeader:   traceIDStr,
			trace.B3SpanIDHeader:    spanIDStr,
			trace.B3SampledHeader:   "0",
			trace.B3DebugFlagHeader: "1",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
	},
	{
		name: "multiple: with parent span id",
		headers: map[string]string{
			trace.B3TraceIDHeader:      traceIDStr,
			trace.B3SpanIDHeader:       spanIDStr,
			trace.B3SampledHeader:      "1",
			trace.B3ParentSpanIDHeader: "00f067aa0ba90200",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "multiple: with only sampled state header",
		headers: map[string]string{
			trace.B3SampledHeader: "0",
		},
		wantSc: trace.EmptySpanContext(),
	},
	{
		name: "multiple: left-padding 64-bit traceID",
		headers: map[string]string{
			trace.B3TraceIDHeader: "a3ce929d0e0e4736",
			trace.B3SpanIDHeader:  spanIDStr,
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID64bitPadded,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
	},
	{
		name: "single: sampling state defer",
		headers: map[string]string{
			trace.B3SingleHeader: fmt.Sprintf("%s-%s", traceIDStr, spanIDStr),
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
	},
	{
		name: "single: sampling state deny",
		headers: map[string]string{
			trace.B3SingleHeader: fmt.Sprintf("%s-%s-0", traceIDStr, spanIDStr),
		},
		wantSc: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
	},
	{
		name: "single: sampling state accept",
		headers: map[string]string{
			trace.B3SingleHeader: fmt.Sprintf("%s-%s-1", traceIDStr, spanIDStr),
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "single: sampling state debug",
		headers: map[string]string{
			trace.B3SingleHeader: fmt.Sprintf("%s-%s-d", traceIDStr, spanIDStr),
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
	},
	{
		name: "single: with parent span id",
		headers: map[string]string{
			trace.B3SingleHeader: fmt.Sprintf("%s-%s-1-00000000000000cd", traceIDStr, spanIDStr),
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "single: with only sampling state deny",
		headers: map[string]string{
			trace.B3SingleHeader: "0",
		},
		wantSc: trace.EmptySpanContext(),
	},
	{
		name: "single: left-padding 64-bit traceID",
		headers: map[string]string{
			trace.B3SingleHeader: fmt.Sprintf("a3ce929d0e0e4736-%s", spanIDStr),
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID64bitPadded,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
	},
	{
		name: "both single and multiple: single priority",
		headers: map[string]string{
			trace.B3SingleHeader:  fmt.Sprintf("%s-%s-1", traceIDStr, spanIDStr),
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  spanIDStr,
			trace.B3SampledHeader: "0",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	// An invalid Single Headers should fallback to multiple.
	{
		name: "both single and multiple: invalid single",
		headers: map[string]string{
			trace.B3SingleHeader:  fmt.Sprintf("%s-%s-", traceIDStr, spanIDStr),
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  spanIDStr,
			trace.B3SampledHeader: "0",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsNotSampled,
		},
	},
	// Invalid Mult Header should not be noticed as Single takes precedence.
	{
		name: "both single and multiple: invalid multiple",
		headers: map[string]string{
			trace.B3SingleHeader:  fmt.Sprintf("%s-%s-1", traceIDStr, spanIDStr),
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  spanIDStr,
			trace.B3SampledHeader: "invalid",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
}

var extractInvalidHeaders = []extractTest{
	{
		name: "multiple: trace ID length > 32",
		headers: map[string]string{
			trace.B3TraceIDHeader: "ab00000000000000000000000000000000",
			trace.B3SpanIDHeader:  "cd00000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: trace ID length >16 and <32",
		headers: map[string]string{
			trace.B3TraceIDHeader: "ab0000000000000000000000000000",
			trace.B3SpanIDHeader:  "cd00000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: trace ID length <16",
		headers: map[string]string{
			trace.B3TraceIDHeader: "ab0000000000",
			trace.B3SpanIDHeader:  "cd00000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: wrong span ID length",
		headers: map[string]string{
			trace.B3TraceIDHeader: "ab000000000000000000000000000000",
			trace.B3SpanIDHeader:  "cd0000000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: wrong sampled flag length",
		headers: map[string]string{
			trace.B3TraceIDHeader: "ab000000000000000000000000000000",
			trace.B3SpanIDHeader:  "cd00000000000000",
			trace.B3SampledHeader: "10",
		},
	},
	{
		name: "multiple: bogus trace ID",
		headers: map[string]string{
			trace.B3TraceIDHeader: "qw000000000000000000000000000000",
			trace.B3SpanIDHeader:  "cd00000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: bogus span ID",
		headers: map[string]string{
			trace.B3TraceIDHeader: "ab000000000000000000000000000000",
			trace.B3SpanIDHeader:  "qw00000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: bogus sampled flag",
		headers: map[string]string{
			trace.B3TraceIDHeader: "ab000000000000000000000000000000",
			trace.B3SpanIDHeader:  "cd00000000000000",
			trace.B3SampledHeader: "d",
		},
	},
	{
		name: "multiple: upper case trace ID",
		headers: map[string]string{
			trace.B3TraceIDHeader: "AB000000000000000000000000000000",
			trace.B3SpanIDHeader:  "cd00000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: upper case span ID",
		headers: map[string]string{
			trace.B3TraceIDHeader: "ab000000000000000000000000000000",
			trace.B3SpanIDHeader:  "CD00000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: zero trace ID",
		headers: map[string]string{
			trace.B3TraceIDHeader: "00000000000000000000000000000000",
			trace.B3SpanIDHeader:  "cd00000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: zero span ID",
		headers: map[string]string{
			trace.B3TraceIDHeader: "ab000000000000000000000000000000",
			trace.B3SpanIDHeader:  "0000000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: missing span ID",
		headers: map[string]string{
			trace.B3TraceIDHeader: "ab000000000000000000000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: missing trace ID",
		headers: map[string]string{
			trace.B3SpanIDHeader:  "cd00000000000000",
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "multiple: sampled header set to 1 but trace ID and span ID are missing",
		headers: map[string]string{
			trace.B3SampledHeader: "1",
		},
	},
	{
		name: "single: wrong trace ID length",
		headers: map[string]string{
			trace.B3SingleHeader: "ab00000000000000000000000000000000-cd00000000000000-1",
		},
	},
	{
		name: "single: wrong span ID length",
		headers: map[string]string{
			trace.B3SingleHeader: "ab000000000000000000000000000000-cd0000000000000000-1",
		},
	},
	{
		name: "single: wrong sampled state length",
		headers: map[string]string{
			trace.B3SingleHeader: "00-ab000000000000000000000000000000-cd00000000000000-01",
		},
	},
	{
		name: "single: wrong parent span ID length",
		headers: map[string]string{
			trace.B3SingleHeader: "ab000000000000000000000000000000-cd00000000000000-1-cd0000000000000000",
		},
	},
	{
		name: "single: bogus trace ID",
		headers: map[string]string{
			trace.B3SingleHeader: "qw000000000000000000000000000000-cd00000000000000-1",
		},
	},
	{
		name: "single: bogus span ID",
		headers: map[string]string{
			trace.B3SingleHeader: "ab000000000000000000000000000000-qw00000000000000-1",
		},
	},
	{
		name: "single: bogus sampled flag",
		headers: map[string]string{
			trace.B3SingleHeader: "ab000000000000000000000000000000-cd00000000000000-q",
		},
	},
	{
		name: "single: bogus parent span ID",
		headers: map[string]string{
			trace.B3SingleHeader: "ab000000000000000000000000000000-cd00000000000000-1-qw00000000000000",
		},
	},
	{
		name: "single: upper case trace ID",
		headers: map[string]string{
			trace.B3SingleHeader: "AB000000000000000000000000000000-cd00000000000000-1",
		},
	},
	{
		name: "single: upper case span ID",
		headers: map[string]string{
			trace.B3SingleHeader: "ab000000000000000000000000000000-CD00000000000000-1",
		},
	},
	{
		name: "single: upper case parent span ID",
		headers: map[string]string{
			trace.B3SingleHeader: "ab000000000000000000000000000000-cd00000000000000-1-EF00000000000000",
		},
	},
	{
		name: "single: zero trace ID and span ID",
		headers: map[string]string{
			trace.B3SingleHeader: "00000000000000000000000000000000-0000000000000000-1",
		},
	},
	{
		name: "single: with sampling set to true",
		headers: map[string]string{
			trace.B3SingleHeader: "ab000000000000000000000000000000-cd00000000000000-true",
		},
	},
}

type injectTest struct {
	name             string
	encoding         trace.B3Encoding
	parentSc         trace.SpanContext
	wantHeaders      map[string]string
	doNotWantHeaders []string
}

var injectHeader = []injectTest{
	{
		name: "none: sampled",
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  "0000000000000001",
			trace.B3SampledHeader: "1",
		},
		doNotWantHeaders: []string{
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
			trace.B3SingleHeader,
		},
	},
	{
		name: "none: not sampled",
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsNotSampled,
		},
		wantHeaders: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  "0000000000000002",
			trace.B3SampledHeader: "0",
		},
		doNotWantHeaders: []string{
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
			trace.B3SingleHeader,
		},
	},
	{
		name: "none: unset sampled",
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
		wantHeaders: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  "0000000000000003",
		},
		doNotWantHeaders: []string{
			trace.B3SampledHeader,
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
			trace.B3SingleHeader,
		},
	},
	{
		name:     "multiple: sampled",
		encoding: trace.MultipleHeader,
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  "0000000000000004",
			trace.B3SampledHeader: "1",
		},
		doNotWantHeaders: []string{
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
			trace.B3SingleHeader,
		},
	},
	{
		name:     "multiple: not sampled",
		encoding: trace.MultipleHeader,
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsNotSampled,
		},
		wantHeaders: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  "0000000000000005",
			trace.B3SampledHeader: "0",
		},
		doNotWantHeaders: []string{
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
			trace.B3SingleHeader,
		},
	},
	{
		name:     "multiple: unset sampled",
		encoding: trace.MultipleHeader,
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
		wantHeaders: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  "0000000000000006",
		},
		doNotWantHeaders: []string{
			trace.B3SampledHeader,
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
			trace.B3SingleHeader,
		},
	},
	{
		name:     "single: sampled",
		encoding: trace.SingleHeader,
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			trace.B3SingleHeader: fmt.Sprintf("%s-0000000000000007-1", traceIDStr),
		},
		doNotWantHeaders: []string{
			trace.B3TraceIDHeader,
			trace.B3SpanIDHeader,
			trace.B3SampledHeader,
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
		},
	},
	{
		name:     "single: not sampled",
		encoding: trace.SingleHeader,
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsNotSampled,
		},
		wantHeaders: map[string]string{
			trace.B3SingleHeader: fmt.Sprintf("%s-0000000000000008-0", traceIDStr),
		},
		doNotWantHeaders: []string{
			trace.B3TraceIDHeader,
			trace.B3SpanIDHeader,
			trace.B3SampledHeader,
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
		},
	},
	{
		name:     "single: unset sampled",
		encoding: trace.SingleHeader,
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
		wantHeaders: map[string]string{
			trace.B3SingleHeader: fmt.Sprintf("%s-0000000000000009", traceIDStr),
		},
		doNotWantHeaders: []string{
			trace.B3TraceIDHeader,
			trace.B3SpanIDHeader,
			trace.B3SampledHeader,
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
		},
	},
	{
		name:     "single+multiple: sampled",
		encoding: trace.SingleHeader | trace.MultipleHeader,
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  "000000000000000a",
			trace.B3SampledHeader: "1",
			trace.B3SingleHeader:  fmt.Sprintf("%s-000000000000000a-1", traceIDStr),
		},
		doNotWantHeaders: []string{
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
		},
	},
	{
		name:     "single+multiple: not sampled",
		encoding: trace.SingleHeader | trace.MultipleHeader,
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsNotSampled,
		},
		wantHeaders: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  "000000000000000b",
			trace.B3SampledHeader: "0",
			trace.B3SingleHeader:  fmt.Sprintf("%s-000000000000000b-0", traceIDStr),
		},
		doNotWantHeaders: []string{
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
		},
	},
	{
		name:     "single+multiple: unset sampled",
		encoding: trace.SingleHeader | trace.MultipleHeader,
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsUnset,
		},
		wantHeaders: map[string]string{
			trace.B3TraceIDHeader: traceIDStr,
			trace.B3SpanIDHeader:  "000000000000000c",
			trace.B3SingleHeader:  fmt.Sprintf("%s-000000000000000c", traceIDStr),
		},
		doNotWantHeaders: []string{
			trace.B3SampledHeader,
			trace.B3ParentSpanIDHeader,
			trace.B3DebugFlagHeader,
		},
	},
}

var injectInvalidHeaderGenerator = []injectTest{
	{
		name:     "empty",
		parentSc: trace.SpanContext{},
	},
	{
		name: "missing traceID",
		parentSc: trace.SpanContext{
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "missing spanID",
		parentSc: trace.SpanContext{
			TraceID:    traceID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "missing traceID and spanID",
		parentSc: trace.SpanContext{
			TraceFlags: trace.FlagsSampled,
		},
	},
}

var injectInvalidHeader []injectTest

func init() {
	injectInvalidHeader = make([]injectTest, 0, len(injectInvalidHeaderGenerator)*4)
	allHeaders := []string{
		trace.B3TraceIDHeader,
		trace.B3SpanIDHeader,
		trace.B3SampledHeader,
		trace.B3ParentSpanIDHeader,
		trace.B3DebugFlagHeader,
		trace.B3SingleHeader,
	}
	// Nothing should be set for any header regardless of encoding.
	for _, t := range injectInvalidHeaderGenerator {
		injectInvalidHeader = append(injectInvalidHeader, injectTest{
			name:             "none: " + t.name,
			parentSc:         t.parentSc,
			doNotWantHeaders: allHeaders,
		})
		injectInvalidHeader = append(injectInvalidHeader, injectTest{
			name:             "multiple: " + t.name,
			encoding:         trace.MultipleHeader,
			parentSc:         t.parentSc,
			doNotWantHeaders: allHeaders,
		})
		injectInvalidHeader = append(injectInvalidHeader, injectTest{
			name:             "single: " + t.name,
			encoding:         trace.SingleHeader,
			parentSc:         t.parentSc,
			doNotWantHeaders: allHeaders,
		})
		injectInvalidHeader = append(injectInvalidHeader, injectTest{
			name:             "single+multiple: " + t.name,
			encoding:         trace.SingleHeader | trace.MultipleHeader,
			parentSc:         t.parentSc,
			doNotWantHeaders: allHeaders,
		})
	}
}
