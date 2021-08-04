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
		{"r:3e", -1, 0x3e, nil},
		{"r:3f", -1, 0x3f, nil},
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

func TestProbabilitySamplingBasic(t *testing.T) {
	ctx0 := context.Background()
	te := NewTestExporter()

	oneFunction := func() int { return 1 }

	// 50% sampling at the root span
	sampler1 := TraceIDRatioBased(1, WithRandomSource(oneFunction))
	provider1 := NewTracerProvider(WithSyncer(te), WithSampler(sampler1))

	ctx1, span1 := provider1.Tracer("test").Start(ctx0, "hello")

	require.True(t, span1.IsRecording())
	require.True(t, span1.SpanContext().IsSampled())
	require.Equal(t, "p:01;r:01;", span1.SpanContext().TraceState().Get("otel"))

	span1.End()

	require.Equal(t, 1, te.Len())

	got := te.Spans()[0]

	require.True(t, got.SpanContext().SpanID().IsValid())
	require.EqualValues(t, []attribute.KeyValue{
		attribute.Int64("sampler.adjusted_count", 2),
	}, got.Attributes())

	te.Reset()

	// Parent sampling using propgated probability
	sampler2 := ParentBased(NeverSample())
	provider2 := NewTracerProvider(WithSyncer(te), WithSampler(sampler2))

	ctx2, span2 := provider2.Tracer("test").Start(ctx1, "hello")

	require.True(t, span2.IsRecording())
	require.True(t, span2.SpanContext().IsSampled())
	require.Equal(t, "p:01;r:01;", span2.SpanContext().TraceState().Get("otel"))

	span2.End()

	require.Equal(t, 1, te.Len())

	got = te.Spans()[0]

	require.True(t, got.SpanContext().SpanID().IsValid())
	require.EqualValues(t, []attribute.KeyValue{
		attribute.Int64("sampler.adjusted_count", 2),
	}, got.Attributes())

	// Compare contexts
	require.Equal(t, span1.SpanContext().TraceID(), span2.SpanContext().TraceID())
	require.NotEqual(t, span1.SpanContext().SpanID(), span2.SpanContext().SpanID())

	_ = ctx2
}
