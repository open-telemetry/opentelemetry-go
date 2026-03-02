// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/trace"
)

func TestParentBasedDefaultLocalParentSampled(t *testing.T) {
	sampler := ParentBased(AlwaysSample())
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := trace.ContextWithSpanContext(
		t.Context(),
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
		t.Context(),
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
					t.Context(),
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
//
//	"A TraceIDRatioBased sampler with a given sampling rate MUST also sample
//	 all traces that any TraceIDRatioBased sampler with a lower sampling rate
//	 would sample."
func TestTraceIdRatioSamplesInclusively(t *testing.T) {
	const (
		numSamplers = 1000
		numTraces   = 100
	)
	idg := defaultIDGenerator()

	for range numSamplers {
		ratioLo, ratioHi := rand.Float64(), rand.Float64()
		if ratioHi < ratioLo {
			ratioLo, ratioHi = ratioHi, ratioLo
		}
		samplerHi := TraceIDRatioBased(ratioHi)
		samplerLo := TraceIDRatioBased(ratioLo)
		for range numTraces {
			traceID, _ := idg.NewIDs(t.Context())

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
			"traceIDRatioSampler",
			TraceIDRatioBased(.5),
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
					t.Context(),
					trace.NewSpanContext(trace.SpanContextConfig{
						TraceState: traceState,
					}),
				),
			}

			require.Equal(t, traceState, tc.sampler.ShouldSample(params).Tracestate, "TraceState is not equal")
		})
	}
}

func TestAlwaysRecordSamplingDecision(t *testing.T) {
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")

	testCases := []struct {
		name             string
		rootSampler      Sampler
		expectedDecision SamplingDecision
	}{
		{
			name:             "when root sampler decision is RecordAndSample, AlwaysRecord returns RecordAndSample",
			rootSampler:      AlwaysSample(),
			expectedDecision: RecordAndSample,
		},
		{
			name:             "when root sampler decision is Drop, AlwaysRecord returns RecordOnly",
			rootSampler:      NeverSample(),
			expectedDecision: RecordOnly,
		},
		{
			name:             "when root sampler decision is RecordOnly, AlwaysRecord returns RecordOnly",
			rootSampler:      RecordingOnly(),
			expectedDecision: RecordOnly,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sampler := AlwaysRecord(tc.rootSampler)
			parentCtx := trace.ContextWithSpanContext(
				t.Context(),
				trace.NewSpanContext(trace.SpanContextConfig{
					TraceID:    traceID,
					SpanID:     spanID,
					TraceFlags: trace.FlagsSampled,
				}),
			)
			samplingResult := sampler.ShouldSample(SamplingParameters{ParentContext: parentCtx})
			if samplingResult.Decision != tc.expectedDecision {
				t.Errorf("Sampling decision should be %v, got %v instead",
					tc.expectedDecision,
					samplingResult.Decision,
				)
			}
		})
	}
}

func TestAlwaysRecordDefaultDescription(t *testing.T) {
	sampler := AlwaysRecord(NeverSample())

	expectedDescription := fmt.Sprintf("AlwaysRecord{root:%s}", NeverSample().Description())

	if sampler.Description() != expectedDescription {
		t.Errorf("Sampler description should be %s, got '%s' instead",
			expectedDescription,
			sampler.Description(),
		)
	}
}
