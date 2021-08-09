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
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TestParentBasedDefaultLocalParentSampled(t *testing.T) {
	sampler := ParentBased(AlwaysSample())
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := trace.ContextWithSpanContext(
		context.Background(),
		trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		}),
	)
	if sampler.ShouldSample(SamplingParameters{ParentContext: parentCtx}).Decision != RecordAndSample {
		t.Error("Sampling decision should be RecordAndSample")
	}
}

func TestParentBasedDefaultLocalParentNotSampled(t *testing.T) {
	sampler := ParentBased(AlwaysSample())
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := trace.ContextWithSpanContext(
		context.Background(),
		trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  spanID,
		}),
	)
	if sampler.ShouldSample(SamplingParameters{ParentContext: parentCtx}).Decision != Drop {
		t.Error("Sampling decision should be Drop")
	}
}

func TestParentBasedWithNoParent(t *testing.T) {
	params := SamplingParameters{}

	sampler := ParentBased(AlwaysSample())
	if sampler.ShouldSample(params).Decision != RecordAndSample {
		t.Error("Sampling decision should be RecordAndSample")
	}

	sampler = ParentBased(NeverSample())
	if sampler.ShouldSample(params).Decision != Drop {
		t.Error("Sampling decision should be Drop")
	}
}

func TestParentBasedWithSamplerOptions(t *testing.T) {
	testCases := []struct {
		name                            string
		samplerOption                   ParentBasedSamplerOption
		isParentRemote, isParentSampled bool
		expectedDecision                SamplingDecision
	}{
		{
			"localParentSampled",
			WithLocalParentSampled(NeverSample()),
			false,
			true,
			Drop,
		},
		{
			"localParentNotSampled",
			WithLocalParentNotSampled(AlwaysSample()),
			false,
			false,
			RecordAndSample,
		},
		{
			"remoteParentSampled",
			WithRemoteParentSampled(NeverSample()),
			true,
			true,
			Drop,
		},
		{
			"remoteParentNotSampled",
			WithRemoteParentNotSampled(AlwaysSample()),
			true,
			false,
			RecordAndSample,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
			spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
			pscc := trace.SpanContextConfig{
				TraceID: traceID,
				SpanID:  spanID,
				Remote:  tc.isParentRemote,
			}
			if tc.isParentSampled {
				pscc.TraceFlags = trace.FlagsSampled
			}

			params := SamplingParameters{
				ParentContext: trace.ContextWithSpanContext(
					context.Background(),
					trace.NewSpanContext(pscc),
				),
			}

			sampler := ParentBased(
				nil,
				tc.samplerOption,
			)

			var wantStr, gotStr string
			switch tc.expectedDecision {
			case RecordAndSample:
				wantStr = "RecordAndSample"
			case Drop:
				wantStr = "Drop"
			default:
				wantStr = "unknown"
			}

			actualDecision := sampler.ShouldSample(params).Decision
			switch actualDecision {
			case RecordAndSample:
				gotStr = "RecordAndSample"
			case Drop:
				gotStr = "Drop"
			default:
				gotStr = "unknown"
			}

			assert.Equalf(t, tc.expectedDecision, actualDecision, "want %s, got %s", wantStr, gotStr)
		})
	}
}

func TestParentBasedDefaultDescription(t *testing.T) {
	sampler := ParentBased(AlwaysSample())

	expectedDescription := fmt.Sprintf("ParentBased{root:%s,remoteParentSampled:%s,"+
		"remoteParentNotSampled:%s,localParentSampled:%s,localParentNotSampled:%s}",
		AlwaysSample().Description(),
		AlwaysSample().Description(),
		NeverSample().Description(),
		AlwaysSample().Description(),
		NeverSample().Description())

	if sampler.Description() != expectedDescription {
		t.Errorf("Sampler description should be %s, got '%s' instead",
			expectedDescription,
			sampler.Description(),
		)
	}

}

// TraceIDRatioBased sampler requirements state
//  "A TraceIDRatioBased sampler with a given sampling rate MUST also sample
//   all traces that any TraceIDRatioBased sampler with a lower sampling rate
//   would sample."
func TestTraceIdRatioSamplesInclusively(t *testing.T) {
	const (
		numSamplers = 1000
		numTraces   = 100
	)
	idg := defaultIDGenerator()

	for i := 0; i < numSamplers; i++ {
		ratioLo, ratioHi := rand.Float64(), rand.Float64()
		if ratioHi < ratioLo {
			ratioLo, ratioHi = ratioHi, ratioLo
		}
		samplerHi := ProbabilityBased(ratioHi)
		samplerLo := ProbabilityBased(ratioLo)
		for j := 0; j < numTraces; j++ {
			traceID, _ := idg.NewIDs(context.Background())

			params := SamplingParameters{TraceID: traceID}
			if samplerLo.ShouldSample(params).Decision == RecordAndSample {
				require.Equal(t, RecordAndSample, samplerHi.ShouldSample(params).Decision,
					"%s sampled but %s did not", samplerLo.Description(), samplerHi.Description())
			}
		}
	}
}

func TestTracestateIsPassed(t *testing.T) {
	testCases := []struct {
		name    string
		sampler Sampler
	}{
		{
			"notSampled",
			NeverSample(),
		},
		{
			"sampled",
			AlwaysSample(),
		},
		{
			"parentSampled",
			ParentBased(AlwaysSample()),
		},
		{
			"parentNotSampled",
			ParentBased(NeverSample()),
		},
		{
			"probabilitySampler",
			ProbabilityBased(.5),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			traceState, err := trace.ParseTraceState("k=v")
			if err != nil {
				t.Error(err)
			}

			params := SamplingParameters{
				ParentContext: trace.ContextWithSpanContext(
					context.Background(),
					trace.NewSpanContext(trace.SpanContextConfig{
						TraceState: traceState,
					}),
				),
			}

			require.Equal(t, traceState, tc.sampler.ShouldSample(params).Tracestate, "TraceState is not equal")
		})
	}
}

func TestParseTraceState(t *testing.T) {
	type testCase struct {
		in         string
		pval, rval int
		expectErr  error
	}
	for _, test := range []testCase{
		{"r:1;p:2", 2, 1, nil},
		{"r:1;p:2;", 2, 1, nil},
		{"p:2;r:1;", 2, 1, nil},
		{"p:2;r:1", 2, 1, nil},
		{"r:1;", -1, 1, nil},
		{"r:1", -1, 1, nil},
		{"r:1=p:2", -1, -1, errTraceStateSyntax},
		{"r:1;p:2=s:3", -1, -1, errTraceStateSyntax},
		{":1;p:2=s:3", -1, -1, errTraceStateSyntax},
		{":;p:2=s:3", -1, -1, errTraceStateSyntax},
		{":;:", -1, -1, errTraceStateSyntax},
		{":", -1, -1, errTraceStateSyntax},
		{"", -1, -1, nil},
		{"r:", -1, -1, errTraceStateSyntax},
		{"r:;p=1", -1, -1, errTraceStateSyntax},
		{"r:1", -1, 1, nil},
		{"r:10", -1, 0x10, nil},
		{"r:3e", -1, 62, nil}, // Smallest non-zero (62)
		{"r:3f", -1, 63, nil}, // The zero value (63)
		{"r:40", -1, -1, errTraceStateSyntax},
	} {
		t.Run(test.in, func(t *testing.T) {
			otts, err := parseOTelTraceState(test.in)

			if test.expectErr != nil {
				require.True(t, errors.Is(err, test.expectErr))
			} else {
				require.Equal(t, test.rval, otts.random)
				require.Equal(t, test.pval, otts.probability)
			}
		})
	}
}

func TestProbabilitySamplingParentBased(t *testing.T) {
	type testCase struct {
		rand        uint
		prob        int
		expectCount int64
	}

	for _, tc := range []testCase{
		// 1-in-1
		{0, 0, 1},
		{1, 0, 1},
		{2, 0, 1},

		// 1-in-2
		{0, 1, 0},
		{1, 1, 2},
		{2, 1, 2},

		// 1-in-4
		{0, 2, 0},
		{1, 2, 0},
		{2, 2, 4},
		{3, 2, 4},

		// 1-in-(2^62)
		{0, 62, 0},
		{61, 62, 0},
		{62, 62, 1 << 62},

		// Zero is a special case
		{0, 63, 0},
		{62, 63, 0},
	} {
		t.Run(fmt.Sprint(tc.rand, "/", tc.prob), func(t *testing.T) {
			ctx0 := context.Background()
			te := NewTestExporter()
			expectTS := fmt.Sprintf("p:%02x;r:%02x;", tc.prob, tc.rand)
			expectSampled := tc.expectCount != 0
			var numSampled int
			if tc.expectCount != 0 {
				numSampled = 1
			}

			var expectAttrs []attribute.KeyValue

			if tc.expectCount != 1 {
				expectAttrs = []attribute.KeyValue{
					attribute.Int64("sampler.adjusted_count", tc.expectCount),
				}
			}

			// Root span uses a TraceIDRatioBased Sampler
			sampler1 := TraceIDRatioBased(tc.prob, WithRandomSource(func() uint { return tc.rand }))
			provider1 := NewTracerProvider(WithSyncer(te), WithSampler(sampler1))

			ctx1, span1 := provider1.Tracer("test").Start(ctx0, "hello")

			require.Equal(t, expectSampled, span1.IsRecording())
			require.Equal(t, expectSampled, span1.SpanContext().IsSampled())

			require.Equal(t, expectTS, span1.SpanContext().TraceState().Get("otel"))

			span1.End()

			require.Equal(t, numSampled, te.Len())

			if expectSampled {
				got := te.Spans()[0]

				require.True(t, got.SpanContext().SpanID().IsValid())
				require.EqualValues(t, expectAttrs, got.Attributes())

			}
			te.Reset()

			// ParentBased sampler
			sampler2 := ParentBased(NeverSample())
			provider2 := NewTracerProvider(WithSyncer(te), WithSampler(sampler2))

			ctx2, span2 := provider2.Tracer("test").Start(ctx1, "hello")

			require.Equal(t, expectSampled, span2.IsRecording())
			require.Equal(t, expectSampled, span2.SpanContext().IsSampled())

			require.Equal(t, expectTS, span2.SpanContext().TraceState().Get("otel"))

			span2.End()

			require.Equal(t, numSampled, te.Len())

			if expectSampled {
				got := te.Spans()[0]

				require.True(t, got.SpanContext().SpanID().IsValid())
				require.Equal(t, span1.SpanContext().SpanID(), got.Parent().SpanID())
				require.EqualValues(t, expectAttrs, got.Attributes())
			}
			te.Reset()

			require.Equal(t, span1.SpanContext().TraceID(), span2.SpanContext().TraceID())
			require.NotEqual(t, span1.SpanContext().SpanID(), span2.SpanContext().SpanID())

			// Repeat with a grandchild
			sampler3 := ParentBased(NeverSample())
			provider3 := NewTracerProvider(WithSyncer(te), WithSampler(sampler3))

			ctx3, span3 := provider3.Tracer("test").Start(ctx2, "hello")

			require.Equal(t, expectSampled, span3.IsRecording())
			require.Equal(t, expectSampled, span3.SpanContext().IsSampled())

			require.Equal(t, expectTS, span3.SpanContext().TraceState().Get("otel"))

			span3.End()

			require.Equal(t, numSampled, te.Len())

			if expectSampled {
				got := te.Spans()[0]

				require.True(t, got.SpanContext().SpanID().IsValid())
				require.Equal(t, span2.SpanContext().SpanID(), got.Parent().SpanID())
				require.EqualValues(t, expectAttrs, got.Attributes())
			}
			te.Reset()

			require.Equal(t, span1.SpanContext().TraceID(), span3.SpanContext().TraceID())
			require.NotEqual(t, span1.SpanContext().SpanID(), span3.SpanContext().SpanID())
			require.NotEqual(t, span2.SpanContext().SpanID(), span3.SpanContext().SpanID())

			_ = ctx3
		})
	}
}

func TestProbabilitySamplingTraceIDRatioBased(t *testing.T) {
	testWith := func(t *testing.T, rval uint) {
		te := NewTestExporter()

		var pvalFunc func(context.Context, int)

		pvalFunc = func(ctx context.Context, pval int) {
			if pval == otelSamplingZeroValue {
				return
			}

			sampler := TraceIDRatioBased(pval, WithRandomSource(func() uint { return rval }))
			provider := NewTracerProvider(WithSyncer(te), WithSampler(sampler))

			childCtx, span := provider.Tracer("test").Start(ctx, "hello")

			expectTS := fmt.Sprintf("p:%02x;r:%02x;", pval, rval)

			require.Equal(t, expectTS, span.SpanContext().TraceState().Get("otel"))

			pvalFunc(childCtx, pval+1)

			span.End()
		}
		pvalFunc(context.Background(), 0)
		expect := int(rval + 1) // count of values such that p <= rval
		require.Equal(t, expect, te.Len())

		var total int64
		for _, got := range te.Spans() {
			var this int64 = 1
			for _, attr := range got.Attributes() {
				if attr.Key == "sampler.adjusted_count" {
					this = attr.Value.AsInt64()
				}
			}
			total += this
		}
		// Note this is a geometric series, e.g., if we
		// expected 4 spans the value is 1+2+4+8 == 15 == (2^4)-1.
		require.Equal(t, (int64(1)<<expect)-1, total)
	}

	for i := uint(0); i < otelSamplingZeroValue; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			testWith(t, i)
		})
	}
}
