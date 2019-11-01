// Copyright 2019, OpenTelemetry Authors
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

package propagation_test

import (
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/propagation"
)

type extractTest struct {
	name    string
	headers map[string]string
	wantSc  core.SpanContext
}

var extractMultipleHeaders = []extractTest{
	{
		name: "sampling state defer",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:  "00f067aa0ba902b7",
		},
		wantSc: core.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
	},
	{
		name: "sampling state deny",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:  "00f067aa0ba902b7",
			propagation.B3SampledHeader: "0",
		},
		wantSc: core.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
	},
	{
		name: "sampling state accept",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:  "00f067aa0ba902b7",
			propagation.B3SampledHeader: "1",
		},
		wantSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		},
	},
	{
		name: "sampling state as a boolean",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:  "00f067aa0ba902b7",
			propagation.B3SampledHeader: "true",
		},
		wantSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		},
	},
	{
		name: "debug flag set",
		headers: map[string]string{
			propagation.B3TraceIDHeader:   "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:    "00f067aa0ba902b7",
			propagation.B3DebugFlagHeader: "1",
		},
		wantSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		},
	},
	{
		// spec is not clear on the behavior for this case. If debug flag is set
		// then sampled state should not be set. From that perspective debug
		// takes precedence. Hence, it is sampled.
		name: "debug flag set and sampling state is deny",
		headers: map[string]string{
			propagation.B3TraceIDHeader:   "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:    "00f067aa0ba902b7",
			propagation.B3SampledHeader:   "0",
			propagation.B3DebugFlagHeader: "1",
		},
		wantSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		},
	},
	{
		name: "with parent span id",
		headers: map[string]string{
			propagation.B3TraceIDHeader:      "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:       "00f067aa0ba902b7",
			propagation.B3SampledHeader:      "1",
			propagation.B3ParentSpanIDHeader: "00f067aa0ba90200",
		},
		wantSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		},
	},
	{
		name: "with only sampled state header",
		headers: map[string]string{
			propagation.B3SampledHeader: "0",
		},
		wantSc: core.EmptySpanContext(),
	},
}

var extractSingleHeader = []extractTest{
	{
		name: "sampling state defer",
		headers: map[string]string{
			propagation.B3SingleHeader: "4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7",
		},
		wantSc: core.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
	},
	{
		name: "sampling state deny",
		headers: map[string]string{
			propagation.B3SingleHeader: "4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-0",
		},
		wantSc: core.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
	},
	{
		name: "sampling state accept",
		headers: map[string]string{
			propagation.B3SingleHeader: "4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-1",
		},
		wantSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		},
	},
	{
		name: "sampling state debug",
		headers: map[string]string{
			propagation.B3SingleHeader: "4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-d",
		},
		wantSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		},
	},
	{
		name: "with parent span id",
		headers: map[string]string{
			propagation.B3SingleHeader: "4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-1-00000000000000cd",
		},
		wantSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		},
	},
	{
		name: "with only sampling state deny",
		headers: map[string]string{
			propagation.B3SingleHeader: "0",
		},
		wantSc: core.EmptySpanContext(),
	},
}

var extractInvalidB3MultipleHeaders = []extractTest{
	{
		name: "trace ID length > 32",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "ab00000000000000000000000000000000",
			propagation.B3SpanIDHeader:  "cd00000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "trace ID length >16 and <32",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "ab0000000000000000000000000000",
			propagation.B3SpanIDHeader:  "cd00000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "trace ID length <16",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "ab0000000000",
			propagation.B3SpanIDHeader:  "cd00000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "wrong span ID length",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "ab000000000000000000000000000000",
			propagation.B3SpanIDHeader:  "cd0000000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "wrong sampled flag length",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "ab000000000000000000000000000000",
			propagation.B3SpanIDHeader:  "cd00000000000000",
			propagation.B3SampledHeader: "10",
		},
	},
	{
		name: "bogus trace ID",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "qw000000000000000000000000000000",
			propagation.B3SpanIDHeader:  "cd00000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "bogus span ID",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "ab000000000000000000000000000000",
			propagation.B3SpanIDHeader:  "qw00000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "bogus sampled flag",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "ab000000000000000000000000000000",
			propagation.B3SpanIDHeader:  "cd00000000000000",
			propagation.B3SampledHeader: "d",
		},
	},
	{
		name: "bogus debug flag (string)",
		headers: map[string]string{
			propagation.B3TraceIDHeader:   "ab000000000000000000000000000000",
			propagation.B3SpanIDHeader:    "cd00000000000000",
			propagation.B3SampledHeader:   "1",
			propagation.B3DebugFlagHeader: "d",
		},
	},
	{
		name: "bogus debug flag (number)",
		headers: map[string]string{
			propagation.B3TraceIDHeader:   "ab000000000000000000000000000000",
			propagation.B3SpanIDHeader:    "cd00000000000000",
			propagation.B3SampledHeader:   "1",
			propagation.B3DebugFlagHeader: "10",
		},
	},
	{
		name: "upper case trace ID",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "AB000000000000000000000000000000",
			propagation.B3SpanIDHeader:  "cd00000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "upper case span ID",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "ab000000000000000000000000000000",
			propagation.B3SpanIDHeader:  "CD00000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "zero trace ID",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "00000000000000000000000000000000",
			propagation.B3SpanIDHeader:  "cd00000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "zero span ID",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "ab000000000000000000000000000000",
			propagation.B3SpanIDHeader:  "0000000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "missing span ID",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "ab000000000000000000000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "missing trace ID",
		headers: map[string]string{
			propagation.B3SpanIDHeader:  "cd00000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "missing trace ID with valid single header",
		headers: map[string]string{
			propagation.B3SingleHeader:  "4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-1",
			propagation.B3SpanIDHeader:  "cd00000000000000",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "sampled header set to 1 but trace ID and span ID are missing",
		headers: map[string]string{
			propagation.B3SampledHeader: "1",
		},
	},
}

var extractInvalidB3SingleHeader = []extractTest{
	{
		name: "wrong trace ID length",
		headers: map[string]string{
			propagation.B3SingleHeader: "ab00000000000000000000000000000000-cd00000000000000-1",
		},
	},
	{
		name: "wrong span ID length",
		headers: map[string]string{
			propagation.B3SingleHeader: "ab000000000000000000000000000000-cd0000000000000000-1",
		},
	},
	{
		name: "wrong sampled state length",
		headers: map[string]string{
			propagation.B3SingleHeader: "00-ab000000000000000000000000000000-cd00000000000000-01",
		},
	},
	{
		name: "wrong parent span ID length",
		headers: map[string]string{
			propagation.B3SingleHeader: "ab000000000000000000000000000000-cd00000000000000-1-cd0000000000000000",
		},
	},
	{
		name: "bogus trace ID",
		headers: map[string]string{
			propagation.B3SingleHeader: "qw000000000000000000000000000000-cd00000000000000-1",
		},
	},
	{
		name: "bogus span ID",
		headers: map[string]string{
			propagation.B3SingleHeader: "ab000000000000000000000000000000-qw00000000000000-1",
		},
	},
	{
		name: "bogus sampled flag",
		headers: map[string]string{
			propagation.B3SingleHeader: "ab000000000000000000000000000000-cd00000000000000-q",
		},
	},
	{
		name: "bogus parent span ID",
		headers: map[string]string{
			propagation.B3SingleHeader: "ab000000000000000000000000000000-cd00000000000000-1-qw00000000000000",
		},
	},
	{
		name: "upper case trace ID",
		headers: map[string]string{
			propagation.B3SingleHeader: "AB000000000000000000000000000000-cd00000000000000-1",
		},
	},
	{
		name: "upper case span ID",
		headers: map[string]string{
			propagation.B3SingleHeader: "ab000000000000000000000000000000-CD00000000000000-1",
		},
	},
	{
		name: "upper case parent span ID",
		headers: map[string]string{
			propagation.B3SingleHeader: "ab000000000000000000000000000000-cd00000000000000-1-EF00000000000000",
		},
	},
	{
		name: "zero trace ID and span ID",
		headers: map[string]string{
			propagation.B3SingleHeader: "00000000000000000000000000000000-0000000000000000-1",
		},
	},
	{
		name: "missing single header with valid separate headers",
		headers: map[string]string{
			propagation.B3TraceIDHeader: "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:  "00f067aa0ba902b7",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "upper case span ID with valid separate headers",
		headers: map[string]string{
			propagation.B3SingleHeader:  "ab000000000000000000000000000000-CD00000000000000-1",
			propagation.B3TraceIDHeader: "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:  "00f067aa0ba902b7",
			propagation.B3SampledHeader: "1",
		},
	},
	{
		name: "with sampling set to true",
		headers: map[string]string{
			propagation.B3SingleHeader: "ab000000000000000000000000000000-cd00000000000000-true",
		},
	},
}

type injectTest struct {
	name             string
	parentSc         core.SpanContext
	wantHeaders      map[string]string
	doNotWantHeaders []string
}

var injectB3MultipleHeader = []injectTest{
	{
		name: "valid spancontext, sampled",
		parentSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		},
		wantHeaders: map[string]string{
			propagation.B3TraceIDHeader: "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:  "0000000000000001",
			propagation.B3SampledHeader: "1",
		},
		doNotWantHeaders: []string{
			propagation.B3ParentSpanIDHeader,
		},
	},
	{
		name: "valid spancontext, not sampled",
		parentSc: core.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		wantHeaders: map[string]string{
			propagation.B3TraceIDHeader: "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:  "0000000000000002",
			propagation.B3SampledHeader: "0",
		},
		doNotWantHeaders: []string{
			propagation.B3ParentSpanIDHeader,
		},
	},
	{
		name: "valid spancontext, with unsupported bit set in traceflags",
		parentSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: 0xff,
		},
		wantHeaders: map[string]string{
			propagation.B3TraceIDHeader: "4bf92f3577b34da6a3ce929d0e0e4736",
			propagation.B3SpanIDHeader:  "0000000000000003",
			propagation.B3SampledHeader: "1",
		},
		doNotWantHeaders: []string{
			propagation.B3ParentSpanIDHeader,
		},
	},
}

var injectB3SingleleHeader = []injectTest{
	{
		name: "valid spancontext, sampled",
		parentSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		},
		wantHeaders: map[string]string{
			propagation.B3SingleHeader: "4bf92f3577b34da6a3ce929d0e0e4736-0000000000000001-1",
		},
		doNotWantHeaders: []string{
			propagation.B3TraceIDHeader,
			propagation.B3SpanIDHeader,
			propagation.B3SampledHeader,
			propagation.B3ParentSpanIDHeader,
		},
	},
	{
		name: "valid spancontext, not sampled",
		parentSc: core.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		wantHeaders: map[string]string{
			propagation.B3SingleHeader: "4bf92f3577b34da6a3ce929d0e0e4736-0000000000000002-0",
		},
		doNotWantHeaders: []string{
			propagation.B3TraceIDHeader,
			propagation.B3SpanIDHeader,
			propagation.B3SampledHeader,
			propagation.B3ParentSpanIDHeader,
		},
	},
	{
		name: "valid spancontext, with unsupported bit set in traceflags",
		parentSc: core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: 0xff,
		},
		wantHeaders: map[string]string{
			propagation.B3SingleHeader: "4bf92f3577b34da6a3ce929d0e0e4736-0000000000000003-1",
		},
		doNotWantHeaders: []string{
			propagation.B3TraceIDHeader,
			propagation.B3SpanIDHeader,
			propagation.B3SampledHeader,
			propagation.B3ParentSpanIDHeader,
		},
	},
}
