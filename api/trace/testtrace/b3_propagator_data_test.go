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

const (
	b3Context      = "b3"
	b3Flags        = "x-b3-flags"
	b3TraceID      = "x-b3-traceid"
	b3SpanID       = "x-b3-spanid"
	b3Sampled      = "x-b3-sampled"
	b3ParentSpanID = "x-b3-parentspanid"
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
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDeferred,
		},
	},
	{
		name: "multiple: sampling state deny",
		headers: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "0",
		},
		wantSc: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
	},
	{
		name: "multiple: sampling state accept",
		headers: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "1",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "multiple: sampling state as a boolean: true",
		headers: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "true",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "multiple: sampling state as a boolean: false",
		headers: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "false",
		},
		wantSc: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
	},
	{
		name: "multiple: debug flag set",
		headers: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Flags:   "1",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDeferred | trace.FlagsDebug,
		},
	},
	{
		name: "multiple: debug flag set to not 1 (ignored)",
		headers: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "1",
			b3Flags:   "2",
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
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "0",
			b3Flags:   "1",
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDebug,
		},
	},
	{
		name: "multiple: with parent span id",
		headers: map[string]string{
			b3TraceID:      traceIDStr,
			b3SpanID:       spanIDStr,
			b3Sampled:      "1",
			b3ParentSpanID: "00f067aa0ba90200",
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
			b3Sampled: "0",
		},
		wantSc: trace.EmptySpanContext(),
	},
	{
		name: "multiple: left-padding 64-bit traceID",
		headers: map[string]string{
			b3TraceID: "a3ce929d0e0e4736",
			b3SpanID:  spanIDStr,
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID64bitPadded,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDeferred,
		},
	},
	{
		name: "single: sampling state defer",
		headers: map[string]string{
			b3Context: fmt.Sprintf("%s-%s", traceIDStr, spanIDStr),
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDeferred,
		},
	},
	{
		name: "single: sampling state deny",
		headers: map[string]string{
			b3Context: fmt.Sprintf("%s-%s-0", traceIDStr, spanIDStr),
		},
		wantSc: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
	},
	{
		name: "single: sampling state accept",
		headers: map[string]string{
			b3Context: fmt.Sprintf("%s-%s-1", traceIDStr, spanIDStr),
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
			b3Context: fmt.Sprintf("%s-%s-d", traceIDStr, spanIDStr),
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDebug,
		},
	},
	{
		name: "single: with parent span id",
		headers: map[string]string{
			b3Context: fmt.Sprintf("%s-%s-1-00000000000000cd", traceIDStr, spanIDStr),
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
			b3Context: "0",
		},
		wantSc: trace.EmptySpanContext(),
	},
	{
		name: "single: left-padding 64-bit traceID",
		headers: map[string]string{
			b3Context: fmt.Sprintf("a3ce929d0e0e4736-%s", spanIDStr),
		},
		wantSc: trace.SpanContext{
			TraceID:    traceID64bitPadded,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDeferred,
		},
	},
	{
		name: "both single and multiple: single priority",
		headers: map[string]string{
			b3Context: fmt.Sprintf("%s-%s-1", traceIDStr, spanIDStr),
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "0",
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
			b3Context: fmt.Sprintf("%s-%s-", traceIDStr, spanIDStr),
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "0",
		},
		wantSc: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
	},
	// Invalid Mult Header should not be noticed as Single takes precedence.
	{
		name: "both single and multiple: invalid multiple",
		headers: map[string]string{
			b3Context: fmt.Sprintf("%s-%s-1", traceIDStr, spanIDStr),
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "invalid",
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
			b3TraceID: "ab00000000000000000000000000000000",
			b3SpanID:  "cd00000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: trace ID length >16 and <32",
		headers: map[string]string{
			b3TraceID: "ab0000000000000000000000000000",
			b3SpanID:  "cd00000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: trace ID length <16",
		headers: map[string]string{
			b3TraceID: "ab0000000000",
			b3SpanID:  "cd00000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: wrong span ID length",
		headers: map[string]string{
			b3TraceID: "ab000000000000000000000000000000",
			b3SpanID:  "cd0000000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: wrong sampled flag length",
		headers: map[string]string{
			b3TraceID: "ab000000000000000000000000000000",
			b3SpanID:  "cd00000000000000",
			b3Sampled: "10",
		},
	},
	{
		name: "multiple: bogus trace ID",
		headers: map[string]string{
			b3TraceID: "qw000000000000000000000000000000",
			b3SpanID:  "cd00000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: bogus span ID",
		headers: map[string]string{
			b3TraceID: "ab000000000000000000000000000000",
			b3SpanID:  "qw00000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: bogus sampled flag",
		headers: map[string]string{
			b3TraceID: "ab000000000000000000000000000000",
			b3SpanID:  "cd00000000000000",
			b3Sampled: "d",
		},
	},
	{
		name: "multiple: upper case trace ID",
		headers: map[string]string{
			b3TraceID: "AB000000000000000000000000000000",
			b3SpanID:  "cd00000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: upper case span ID",
		headers: map[string]string{
			b3TraceID: "ab000000000000000000000000000000",
			b3SpanID:  "CD00000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: zero trace ID",
		headers: map[string]string{
			b3TraceID: "00000000000000000000000000000000",
			b3SpanID:  "cd00000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: zero span ID",
		headers: map[string]string{
			b3TraceID: "ab000000000000000000000000000000",
			b3SpanID:  "0000000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: missing span ID",
		headers: map[string]string{
			b3TraceID: "ab000000000000000000000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: missing trace ID",
		headers: map[string]string{
			b3SpanID:  "cd00000000000000",
			b3Sampled: "1",
		},
	},
	{
		name: "multiple: sampled header set to 1 but trace ID and span ID are missing",
		headers: map[string]string{
			b3Sampled: "1",
		},
	},
	{
		name: "single: wrong trace ID length",
		headers: map[string]string{
			b3Context: "ab00000000000000000000000000000000-cd00000000000000-1",
		},
	},
	{
		name: "single: wrong span ID length",
		headers: map[string]string{
			b3Context: "ab000000000000000000000000000000-cd0000000000000000-1",
		},
	},
	{
		name: "single: wrong sampled state length",
		headers: map[string]string{
			b3Context: "00-ab000000000000000000000000000000-cd00000000000000-01",
		},
	},
	{
		name: "single: wrong parent span ID length",
		headers: map[string]string{
			b3Context: "ab000000000000000000000000000000-cd00000000000000-1-cd0000000000000000",
		},
	},
	{
		name: "single: bogus trace ID",
		headers: map[string]string{
			b3Context: "qw000000000000000000000000000000-cd00000000000000-1",
		},
	},
	{
		name: "single: bogus span ID",
		headers: map[string]string{
			b3Context: "ab000000000000000000000000000000-qw00000000000000-1",
		},
	},
	{
		name: "single: bogus sampled flag",
		headers: map[string]string{
			b3Context: "ab000000000000000000000000000000-cd00000000000000-q",
		},
	},
	{
		name: "single: bogus parent span ID",
		headers: map[string]string{
			b3Context: "ab000000000000000000000000000000-cd00000000000000-1-qw00000000000000",
		},
	},
	{
		name: "single: upper case trace ID",
		headers: map[string]string{
			b3Context: "AB000000000000000000000000000000-cd00000000000000-1",
		},
	},
	{
		name: "single: upper case span ID",
		headers: map[string]string{
			b3Context: "ab000000000000000000000000000000-CD00000000000000-1",
		},
	},
	{
		name: "single: upper case parent span ID",
		headers: map[string]string{
			b3Context: "ab000000000000000000000000000000-cd00000000000000-1-EF00000000000000",
		},
	},
	{
		name: "single: zero trace ID and span ID",
		headers: map[string]string{
			b3Context: "00000000000000000000000000000000-0000000000000000-1",
		},
	},
	{
		name: "single: with sampling set to true",
		headers: map[string]string{
			b3Context: "ab000000000000000000000000000000-cd00000000000000-true",
		},
	},
}

type injectTest struct {
	name             string
	encoding         trace.B3Encoding
	sc               trace.SpanContext
	wantHeaders      map[string]string
	doNotWantHeaders []string
}

var injectHeader = []injectTest{
	{
		name: "none: sampled",
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "1",
		},
		doNotWantHeaders: []string{
			b3ParentSpanID,
			b3Flags,
			b3Context,
		},
	},
	{
		name: "none: not sampled",
		sc: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "0",
		},
		doNotWantHeaders: []string{
			b3ParentSpanID,
			b3Flags,
			b3Context,
		},
	},
	{
		name: "none: unset sampled",
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDeferred,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
		},
		doNotWantHeaders: []string{
			b3Sampled,
			b3ParentSpanID,
			b3Flags,
			b3Context,
		},
	},
	{
		name: "none: sampled only",
		sc: trace.SpanContext{
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			b3Sampled: "1",
		},
		doNotWantHeaders: []string{
			b3TraceID,
			b3SpanID,
			b3ParentSpanID,
			b3Flags,
			b3Context,
		},
	},
	{
		name: "none: debug",
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDebug,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Flags:   "1",
		},
		doNotWantHeaders: []string{
			b3Sampled,
			b3ParentSpanID,
			b3Context,
		},
	},
	{
		name: "none: debug omitting sample",
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled | trace.FlagsDebug,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Flags:   "1",
		},
		doNotWantHeaders: []string{
			b3Sampled,
			b3ParentSpanID,
			b3Context,
		},
	},
	{
		name:     "multiple: sampled",
		encoding: trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "1",
		},
		doNotWantHeaders: []string{
			b3ParentSpanID,
			b3Flags,
			b3Context,
		},
	},
	{
		name:     "multiple: not sampled",
		encoding: trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "0",
		},
		doNotWantHeaders: []string{
			b3ParentSpanID,
			b3Flags,
			b3Context,
		},
	},
	{
		name:     "multiple: unset sampled",
		encoding: trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDeferred,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
		},
		doNotWantHeaders: []string{
			b3Sampled,
			b3ParentSpanID,
			b3Flags,
			b3Context,
		},
	},
	{
		name:     "multiple: sampled only",
		encoding: trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			b3Sampled: "1",
		},
		doNotWantHeaders: []string{
			b3TraceID,
			b3SpanID,
			b3ParentSpanID,
			b3Flags,
			b3Context,
		},
	},
	{
		name:     "multiple: debug",
		encoding: trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDebug,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Flags:   "1",
		},
		doNotWantHeaders: []string{
			b3Sampled,
			b3ParentSpanID,
			b3Context,
		},
	},
	{
		name:     "multiple: debug omitting sample",
		encoding: trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled | trace.FlagsDebug,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Flags:   "1",
		},
		doNotWantHeaders: []string{
			b3Sampled,
			b3ParentSpanID,
			b3Context,
		},
	},
	{
		name:     "single: sampled",
		encoding: trace.B3SingleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			b3Context: fmt.Sprintf("%s-%s-1", traceIDStr, spanIDStr),
		},
		doNotWantHeaders: []string{
			b3TraceID,
			b3SpanID,
			b3Sampled,
			b3ParentSpanID,
			b3Flags,
		},
	},
	{
		name:     "single: not sampled",
		encoding: trace.B3SingleHeader,
		sc: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		wantHeaders: map[string]string{
			b3Context: fmt.Sprintf("%s-%s-0", traceIDStr, spanIDStr),
		},
		doNotWantHeaders: []string{
			b3TraceID,
			b3SpanID,
			b3Sampled,
			b3ParentSpanID,
			b3Flags,
		},
	},
	{
		name:     "single: unset sampled",
		encoding: trace.B3SingleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDeferred,
		},
		wantHeaders: map[string]string{
			b3Context: fmt.Sprintf("%s-%s", traceIDStr, spanIDStr),
		},
		doNotWantHeaders: []string{
			b3TraceID,
			b3SpanID,
			b3Sampled,
			b3ParentSpanID,
			b3Flags,
		},
	},
	{
		name:     "single: sampled only",
		encoding: trace.B3SingleHeader,
		sc: trace.SpanContext{
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			b3Context: "1",
		},
		doNotWantHeaders: []string{
			b3Sampled,
			b3TraceID,
			b3SpanID,
			b3ParentSpanID,
			b3Flags,
			b3Context,
		},
	},
	{
		name:     "single: debug",
		encoding: trace.B3SingleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDebug,
		},
		wantHeaders: map[string]string{
			b3Context: fmt.Sprintf("%s-%s-d", traceIDStr, spanIDStr),
		},
		doNotWantHeaders: []string{
			b3TraceID,
			b3SpanID,
			b3Flags,
			b3Sampled,
			b3ParentSpanID,
			b3Context,
		},
	},
	{
		name:     "single: debug omitting sample",
		encoding: trace.B3SingleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled | trace.FlagsDebug,
		},
		wantHeaders: map[string]string{
			b3Context: fmt.Sprintf("%s-%s-d", traceIDStr, spanIDStr),
		},
		doNotWantHeaders: []string{
			b3TraceID,
			b3SpanID,
			b3Flags,
			b3Sampled,
			b3ParentSpanID,
			b3Context,
		},
	},
	{
		name:     "single+multiple: sampled",
		encoding: trace.B3SingleHeader | trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "1",
			b3Context: fmt.Sprintf("%s-%s-1", traceIDStr, spanIDStr),
		},
		doNotWantHeaders: []string{
			b3ParentSpanID,
			b3Flags,
		},
	},
	{
		name:     "single+multiple: not sampled",
		encoding: trace.B3SingleHeader | trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Sampled: "0",
			b3Context: fmt.Sprintf("%s-%s-0", traceIDStr, spanIDStr),
		},
		doNotWantHeaders: []string{
			b3ParentSpanID,
			b3Flags,
		},
	},
	{
		name:     "single+multiple: unset sampled",
		encoding: trace.B3SingleHeader | trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDeferred,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Context: fmt.Sprintf("%s-%s", traceIDStr, spanIDStr),
		},
		doNotWantHeaders: []string{
			b3Sampled,
			b3ParentSpanID,
			b3Flags,
		},
	},
	{
		name:     "single+multiple: sampled only",
		encoding: trace.B3SingleHeader | trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceFlags: trace.FlagsSampled,
		},
		wantHeaders: map[string]string{
			b3Context: "1",
			b3Sampled: "1",
		},
		doNotWantHeaders: []string{
			b3TraceID,
			b3SpanID,
			b3ParentSpanID,
			b3Flags,
		},
	},
	{
		name:     "single+multiple: debug",
		encoding: trace.B3SingleHeader | trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsDebug,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Flags:   "1",
			b3Context: fmt.Sprintf("%s-%s-d", traceIDStr, spanIDStr),
		},
		doNotWantHeaders: []string{
			b3Sampled,
			b3ParentSpanID,
		},
	},
	{
		name:     "single+multiple: debug omitting sample",
		encoding: trace.B3SingleHeader | trace.B3MultipleHeader,
		sc: trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled | trace.FlagsDebug,
		},
		wantHeaders: map[string]string{
			b3TraceID: traceIDStr,
			b3SpanID:  spanIDStr,
			b3Flags:   "1",
			b3Context: fmt.Sprintf("%s-%s-d", traceIDStr, spanIDStr),
		},
		doNotWantHeaders: []string{
			b3Sampled,
			b3ParentSpanID,
		},
	},
}

var injectInvalidHeaderGenerator = []injectTest{
	{
		name: "empty",
		sc:   trace.SpanContext{},
	},
	{
		name: "missing traceID",
		sc: trace.SpanContext{
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "missing spanID",
		sc: trace.SpanContext{
			TraceID:    traceID,
			TraceFlags: trace.FlagsSampled,
		},
	},
	{
		name: "missing traceID and spanID",
		sc: trace.SpanContext{
			TraceFlags: trace.FlagsSampled,
		},
	},
}

var injectInvalidHeader []injectTest

func init() {
	// Preform a test for each invalid injectTest with all combinations of
	// encoding values.
	injectInvalidHeader = make([]injectTest, 0, len(injectInvalidHeaderGenerator)*4)
	allHeaders := []string{
		b3TraceID,
		b3SpanID,
		b3Sampled,
		b3ParentSpanID,
		b3Flags,
		b3Context,
	}
	// Nothing should be set for any header regardless of encoding.
	for _, t := range injectInvalidHeaderGenerator {
		injectInvalidHeader = append(injectInvalidHeader, injectTest{
			name:             "none: " + t.name,
			sc:               t.sc,
			doNotWantHeaders: allHeaders,
		})
		injectInvalidHeader = append(injectInvalidHeader, injectTest{
			name:             "multiple: " + t.name,
			encoding:         trace.B3MultipleHeader,
			sc:               t.sc,
			doNotWantHeaders: allHeaders,
		})
		injectInvalidHeader = append(injectInvalidHeader, injectTest{
			name:             "single: " + t.name,
			encoding:         trace.B3SingleHeader,
			sc:               t.sc,
			doNotWantHeaders: allHeaders,
		})
		injectInvalidHeader = append(injectInvalidHeader, injectTest{
			name:             "single+multiple: " + t.name,
			encoding:         trace.B3SingleHeader | trace.B3MultipleHeader,
			sc:               t.sc,
			doNotWantHeaders: allHeaders,
		})
	}
}
