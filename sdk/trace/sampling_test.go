// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"fmt"
	"math/rand"
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

// TestTraceIDRatioBasedDescription tests the formatted description and
// the corresponding threshold.
func TestTraceIDRatioBasedDescription(t *testing.T) {
	for _, tc := range []struct {
		prob float64
		desc string
	}{
		// Some well-known values
		{0.5, "TraceIDRatioBased{0.5;th:8}"},
		{1 / 3.0, "TraceIDRatioBased{0.3333333333333333;th:aaab}"},
		{2 / 3.0, "TraceIDRatioBased{0.6666666666666666;th:5555}"},
		{1 / 10.0, "TraceIDRatioBased{0.1;th:e666}"},

		// Small powers of two
		{1 / 256.0, "TraceIDRatioBased{0.00390625;th:ff}"},
		{1 / 65536.0, "TraceIDRatioBased{1.52587890625e-05;th:ffff}"},
		{1 / 1048576.0, "TraceIDRatioBased{9.5367431640625e-07;th:fffff}"},

		// Threshold precision automatically rises for small values
		{1 / 100.0, "TraceIDRatioBased{0.01;th:fd70a}"},                       // precision 5
		{1 / 1000.0, "TraceIDRatioBased{0.001;th:ffbe77}"},                    // precision 6
		{1 / 10000.0, "TraceIDRatioBased{0.0001;th:fff9724}"},                 // precision 7
		{1 / 100000.0, "TraceIDRatioBased{1e-05;th:ffff583a}"},                // precision 8
		{1 / 1000000.0, "TraceIDRatioBased{1e-06;th:ffffef39}"},               // precision 8
		{1 / 10000000.0, "TraceIDRatioBased{1e-07;th:fffffe528}"},             // precision 9
		{1 / 100000000.0, "TraceIDRatioBased{1e-08;th:ffffffd50d}"},           // precision 10
		{1 / 1000000000.0, "TraceIDRatioBased{1e-09;th:fffffffbb48}"},         // precision 11
		{1 / 10000000000.0, "TraceIDRatioBased{1e-10;th:ffffffff920d}"},       // precision 12
		{1 / 100000000000.0, "TraceIDRatioBased{1e-11;th:fffffffff5014}"},     // precision 13
		{1 / 1000000000000.0, "TraceIDRatioBased{1e-12;th:fffffffffee68}"},    // precision 13
		{1 / 10000000000000.0, "TraceIDRatioBased{1e-13;th:ffffffffffe3da}"},  // precision 14
		{1 / 100000000000000.0, "TraceIDRatioBased{1e-14;th:fffffffffffd2f}"}, // precision 14

		// Note this has 13 'f' digits.
		{0x1p-52, "TraceIDRatioBased{2.220446049250313e-16;th:fffffffffffff}"},

		// This has 12 '0' digits and a 1.
		{1 - 0x1p-52, "TraceIDRatioBased{0.9999999999999998;th:0000000000001}"},

		// Values very close to 0.0
		{0x1p-53, "TraceIDRatioBased{1.1102230246251565e-16;th:fffffffffffff8}"},
		{0x1p-54, "TraceIDRatioBased{5.551115123125783e-17;th:fffffffffffffc}"},
		{0x1p-55, "TraceIDRatioBased{2.7755575615628914e-17;th:fffffffffffffe}"},
		{0x1p-56, "TraceIDRatioBased{1.3877787807814457e-17;th:ffffffffffffff}"},

		// Values very close to 1.0 round up to 1.0
		{1, "AlwaysOnSampler"},
		{1 - 0x1p-55, "AlwaysOnSampler"},
		{1 - 0x1p-54, "AlwaysOnSampler"},
		{1 - 0x1p-53, "AlwaysOnSampler"},
	} {
		sampler := TraceIDRatioBased(tc.prob)

		require.Equal(t, tc.desc, sampler.Description())
	}
}

// TestTraceIDRatioBasedThreshold tests the unsigned threshold value to ensure
// it is calculated correctly, separately from the printed threshold tested as
// part of the description.  The test inputs are some of same as are used in
// TestTraceIDRatioBasedDescription.
func TestTraceIDRatioBasedThreshold(t *testing.T) {
	for _, tc := range []struct {
		prob      float64
		threshold uint64
	}{
		// Some well-known values
		{0.5, 0x80000000000000},
		{1 / 3.0, 0xaaab0000000000},
		{2 / 3.0, 0x55550000000000},
		{1 / 10.0, 0xe6660000000000},

		// Small powers of two
		{1 / 256.0, 0xff000000000000},
		{1 / 65536.0, 0xffff0000000000},
		{1 / 1048576.0, 0xfffff000000000},
	} {
		sampler := TraceIDRatioBased(tc.prob).(*traceIDRatioSampler)

		require.Equal(t, tc.threshold, sampler.threshold)
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

	for i := 0; i < numSamplers; i++ {
		ratioLo, ratioHi := rand.Float64(), rand.Float64()
		if ratioHi < ratioLo {
			ratioLo, ratioHi = ratioHi, ratioLo
		}
		samplerHi := TraceIDRatioBased(ratioHi)
		samplerLo := TraceIDRatioBased(ratioLo)
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

type unusedSampler struct{}

var _ Sampler = unusedSampler{}

func (unusedSampler) ShouldSample(parameters SamplingParameters) SamplingResult {
	panic("unused sampler should not be called")
}

func (unusedSampler) Description() string {
	return ""
}

// TestTracestateIsPassed exercises a variety of sampler
// configurations and ensures their tracestate output is correct, with
// and without selecting the item for sampling.  For non-100%, non-0%
// configurations, this is tested using the explicit R-value logic
// which makes the test deterministic.
func TestTracestateIsPassed(t *testing.T) {
	type outcome struct {
		sampled bool
		ts      string
	}
	// Note: Inputs always have valid span context (TraceID and SpanID)
	// so ParentBased always takes the always/never sampled of
	// the incoming trace flags.
	testCases := []struct {
		name    string
		sampler Sampler

		// invalidCtx, if true, indicates not to set TraceID
		// and SpanID, which will cause a ParentBased sampler
		// to call the root sampler.
		invalidCtx bool

		// inputTs is the arriving encoded TraceState
		inputTs string

		// ifSampled is the outcome when the incoming context is sampled
		ifSampled outcome

		// ifUnsampled is the outcome when the incoming context is unsampled.
		ifUnsampled outcome
	}{
		{
			// NeverSample() passes trace state.
			name:        "neverSample",
			sampler:     NeverSample(),
			inputTs:     "k=v",
			ifSampled:   outcome{false, "k=v"},
			ifUnsampled: outcome{false, "k=v"},
		},
		{
			// AlwaysSample() passes trace state.
			name:        "alwaysSample",
			sampler:     AlwaysSample(),
			inputTs:     "k=v",
			ifSampled:   outcome{true, "ot=th:0,k=v"},
			ifUnsampled: outcome{true, "ot=th:0,k=v"},
		},
		{
			// ParentBased() passes trace state to the
			// Always- or NeverSample().
			name:        "parentBasedDefaults",
			sampler:     ParentBased(unusedSampler{}),
			inputTs:     "k=v",
			ifSampled:   outcome{true, "ot=th:0,k=v"},
			ifUnsampled: outcome{false, "k=v"},
		},
		{
			// ParentBased passes trace state to the
			// root-based sampler
			name:        "parentBasedRootAlways",
			sampler:     ParentBased(AlwaysSample()),
			invalidCtx:  true,
			inputTs:     "k=v",
			ifSampled:   outcome{true, "ot=th:0,k=v"},
			ifUnsampled: outcome{true, "ot=th:0,k=v"},
		},
		{
			// TraceIDRatioBased ignores parent decision,
			// 50% sampler w/ sampled R-value.
			name:        "fiftyPctSampled",
			sampler:     TraceIDRatioBased(0.5),
			inputTs:     "k=v,ot=rv:ababababababab",
			ifSampled:   outcome{true, "ot=rv:ababababababab;th:8,k=v"},
			ifUnsampled: outcome{true, "ot=rv:ababababababab;th:8,k=v"},
		},
		{
			// TraceIDRatioBased ignores parent decision,
			// 50% sampler w/ unsampled R-value.
			name:        "fiftyPctUnsampled",
			sampler:     TraceIDRatioBased(0.5),
			inputTs:     "k=v,ot=rv:12121212121212",
			ifSampled:   outcome{false, "k=v,ot=rv:12121212121212"},
			ifUnsampled: outcome{false, "k=v,ot=rv:12121212121212"},
		},
	}

	generator := defaultIDGenerator()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, inputSampled := range []bool{true, false} {
				t.Run(fmt.Sprint("sampled=", inputSampled), func(t *testing.T) {
					// Since the TraceID generation step is
					// randomized, repeating the test ensures it
					// is deterministic.  The outcome should not
					// be probabilistic, the repetition is to
					// verify that.
					const repeats = 10
					for i := 0; i < repeats; i++ {
						traceState, err := trace.ParseTraceState(tc.inputTs)
						if err != nil {
							t.Error(err)
						}

						var scc trace.SpanContextConfig
						scc.TraceState = traceState

						if !tc.invalidCtx {
							randTid, randSid := generator.NewIDs(context.Background())
							scc.TraceID = randTid
							scc.SpanID = randSid
						}

						var expect outcome
						if inputSampled {
							scc.TraceFlags = trace.FlagsSampled
							expect = tc.ifSampled
						} else {
							expect = tc.ifUnsampled
						}

						expectState, err := trace.ParseTraceState(expect.ts)
						if err != nil {
							t.Error(err)
						}

						params := SamplingParameters{
							ParentContext: trace.ContextWithSpanContext(
								context.Background(),
								trace.NewSpanContext(scc),
							),
						}

						decision := tc.sampler.ShouldSample(params)
						require.Equal(t, expect.sampled, decision.Decision == RecordAndSample, "Sampler decision is unexpected")
						require.Equal(t, expectState, decision.Tracestate, "TraceState is unexpected")
					}
				})
			}
		})
	}
}

// TestCombineTracestate exercises combineTraceState in a variety of ways
func TestCombineTracestate(t *testing.T) {
	for _, tc := range []struct {
		orig, add, out string
	}{
		// R-value exists : T-value added
		{"rv:ababababababab", "th:123", "rv:ababababababab;th:123"},
		// Ex + R-value : T-value added
		{"ex:xyz;rv:ababababababab", "th:123", "ex:xyz;rv:ababababababab;th:123"},
		// R-value + Ex : T-value added
		{"rv:ababababababab;ex:xyz", "th:123", "rv:ababababababab;ex:xyz;th:123"},
		// Ex : T-value added
		{"ex:xyz", "th:123", "ex:xyz;th:123"},
		// T-value, Ex : T-value overwritten
		{"th:456;ex:xyz", "th:12345", "th:12345;ex:xyz"},
		// Ex, T-value : T-value overwritten
		{"ex:xyz;th:456", "th:12345", "ex:xyz;th:12345"},
		// Ex1, T-value, Ex2 : T-value overwritten
		{"ex1:xyz;th:456;ex2:zyx", "th:12345", "ex1:xyz;th:12345;ex2:zyx"},
		// Ex1, Ex2 : T-value added
		{"ex1:xyz;ex2:zyx", "th:12345", "ex1:xyz;ex2:zyx;th:12345"},
	} {
		require.Equal(t, tc.out, combineTracestate(tc.orig, tc.add))
	}
}
