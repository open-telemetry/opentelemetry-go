// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"fmt"
	"math/rand/v2"
	"strings"
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

func TestTracestateRandomness(t *testing.T) {
	// Valid 14 hex digit value for testing (0x0123456789abcd = 320159098735149)
	const validRV = "0123456789abcd"
	const validRVValue uint64 = 0x0123456789abcd

	testCases := []struct {
		name       string
		otts       string
		wantRandom uint64
		wantHasRV  bool
	}{
		// Correct cases - rv at beginning
		{
			name:       "rv at beginning",
			otts:       "rv:" + validRV,
			wantRandom: validRVValue,
			wantHasRV:  true,
		},
		{
			name:       "rv at beginning with more keys after",
			otts:       "rv:" + validRV + ";th:0;other:value",
			wantRandom: validRVValue,
			wantHasRV:  true,
		},
		// Correct cases - rv in middle
		{
			name:       "rv in middle",
			otts:       "th:0;rv:" + validRV + ";other:value",
			wantRandom: validRVValue,
			wantHasRV:  true,
		},
		{
			name:       "rv in middle with multiple keys",
			otts:       "foo:bar;rv:" + validRV + ";baz:qux",
			wantRandom: validRVValue,
			wantHasRV:  true,
		},
		// Correct cases - rv at end
		{
			name:       "rv at end",
			otts:       "th:0;other:value;rv:" + validRV,
			wantRandom: validRVValue,
			wantHasRV:  true,
		},
		{
			name:       "rv at end as only key after first",
			otts:       "th:0;rv:" + validRV,
			wantRandom: validRVValue,
			wantHasRV:  true,
		},
		// Correct case - max 56-bit value
		{
			name:       "rv with max 56-bit hex value",
			otts:       "rv:0fffffffffffff",
			wantRandom: 0x0fffffffffffff,
			wantHasRV:  true,
		},
		// Incorrect cases
		{
			name:       "no rv key",
			otts:       "th:0;other:value",
			wantRandom: 0,
			wantHasRV:  false,
		},
		{
			name:       "empty string",
			otts:       "",
			wantRandom: 0,
			wantHasRV:  false,
		},
		{
			name:       "rv value too short at beginning",
			otts:       "rv:0123456789abc",
			wantRandom: 0,
			wantHasRV:  false,
		},
		{
			name:       "rv value too short in middle",
			otts:       "th:0;rv:0123456789abc;other:value",
			wantRandom: 0,
			wantHasRV:  false,
		},
		{
			name:       "rv value too short at end",
			otts:       "th:0;rv:0123456789abc",
			wantRandom: 0,
			wantHasRV:  false,
		},
		{
			name:       "rv value too long at beginning without semicolon",
			otts:       "rv:0123456789abcdef",
			wantRandom: 0,
			wantHasRV:  false,
		},
		{
			name:       "rv value too long in middle without semicolon",
			otts:       "th:0;rv:0123456789abcdef;other:value",
			wantRandom: 0,
			wantHasRV:  false,
		},
		{
			name:       "rv value too long at end",
			otts:       "th:0;rv:0123456789abcdef",
			wantRandom: 0,
			wantHasRV:  false,
		},
		{
			name:       "rv with invalid hex character",
			otts:       "rv:0123456789abcg",
			wantRandom: 0,
			wantHasRV:  false,
		},
		{
			name:       "rv as substring of another key",
			otts:       "rvv:0123456789abcd",
			wantRandom: 0,
			wantHasRV:  false,
		},
		{
			name:       "rv as substring in middle",
			otts:       "th:0;rvv:0123456789abcd;other:value",
			wantRandom: 0,
			wantHasRV:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotRandom, gotHasRV := tracestateRandomness(tc.otts)
			assert.Equal(t, tc.wantHasRV, gotHasRV, "hasRandomness mismatch")
			if tc.wantHasRV {
				assert.Equal(t, tc.wantRandom, gotRandom, "randomness value mismatch")
			}
		})
	}
}

func TestEraseTraceStateThKeyValue(t *testing.T) {
	testCases := []struct {
		name string
		otts string
		want string
	}{
		{
			name: "empty string",
			otts: "",
			want: "",
		},
		{
			name: "no th in existing returns unchanged",
			otts: "rv:0123456789abcd;other:value",
			want: "rv:0123456789abcd;other:value",
		},
		{
			name: "only th returns empty",
			otts: "th:0ad",
			want: "",
		},
		{
			name: "th at front",
			otts: "th:0ad;rv:0123456789abcd",
			want: "rv:0123456789abcd",
		},
		{
			name: "th in middle removes th",
			otts: "rv:0123456789abcd;th:0ad;other:value",
			want: "rv:0123456789abcd;other:value",
		},
		{
			name: "th at end removes th",
			otts: "rv:0123456789abcd;th:0ad",
			want: "rv:0123456789abcd",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := eraseTraceStateThKeyValue(tc.otts)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestInsertOrUpdateTraceStateThKeyValue(t *testing.T) {
	testCases := []struct {
		name         string
		existingOtts string
		thkv         string
		want         string
	}{
		{
			name:         "empty existing adds th at front",
			existingOtts: "",
			thkv:         "th:123456789abcd",
			want:         "th:123456789abcd",
		},
		{
			name:         "no th in existing adds th at front",
			existingOtts: "rv:0123456789abcd;other:value",
			thkv:         "th:fedcba987654321",
			want:         "th:fedcba987654321;rv:0123456789abcd;other:value",
		},
		{
			name:         "existing th is replaced and moved to front",
			existingOtts: "rv:0123456789abcd;th:0ad;other:value",
			thkv:         "th:0e1",
			want:         "th:0e1;rv:0123456789abcd;other:value",
		},
		{
			name:         "th at front is replaced",
			existingOtts: "th:0ad;rv:0123456789abcd",
			thkv:         "th:0e1",
			want:         "th:0e1;rv:0123456789abcd",
		},
		{
			name:         "only th in existing is replaced",
			existingOtts: "th:0ad",
			thkv:         "th:0e1",
			want:         "th:0e1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := insertOrUpdateTraceStateThKeyValue(tc.existingOtts, tc.thkv)
			assert.Equal(t, tc.want, got)
		})
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

			// FlagsRandom is required for the sampler to use trace ID for randomness
			params := SamplingParameters{
				ParentContext: trace.ContextWithSpanContext(
					t.Context(),
					trace.NewSpanContext(trace.SpanContextConfig{
						TraceID:    traceID,
						TraceFlags: trace.FlagsRandom,
					}),
				),
				TraceID: traceID,
			}
			if samplerLo.ShouldSample(params).Decision == RecordAndSample {
				require.Equal(t, RecordAndSample, samplerHi.ShouldSample(params).Decision,
					"%s sampled but %s did not", samplerLo.Description(), samplerHi.Description())
			}
		}
	}
}

func TestTracestateIsPassed(t *testing.T) {
	testCases := []struct {
		name                string
		sampler             Sampler
		resultingTracestate string
	}{
		{
			"notSampled",
			NeverSample(),
			"k=v",
		},
		{
			"sampled",
			AlwaysSample(),
			"ot=th:0,k=v",
		},
		{
			"parentSampled",
			ParentBased(AlwaysSample()),
			"ot=th:0,k=v",
		},
		{
			"parentNotSampled",
			ParentBased(NeverSample()),
			"k=v",
		},
		{
			"traceIDRatioSampler",
			TraceIDRatioBased(.5),
			"k=v",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			traceState, err := trace.ParseTraceState("k=v")
			if err != nil {
				t.Error(err)
			}

			expectedTracestate, _ := trace.ParseTraceState(tc.resultingTracestate)
			require.NotNil(t, expectedTracestate)
			require.NotEmpty(t, expectedTracestate)
			traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")

			params := SamplingParameters{
				ParentContext: trace.ContextWithSpanContext(
					t.Context(),
					trace.NewSpanContext(trace.SpanContextConfig{
						TraceState: traceState,
					}),
				),
				TraceID: traceID,
			}

			require.Equal(t, expectedTracestate, tc.sampler.ShouldSample(params).Tracestate, "TraceState is not equal")
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

func TestSamplerDescriptions(t *testing.T) {
	// Direct sampler descriptions
	assert.Equal(t, "AlwaysOnSampler", AlwaysSample().Description())
	assert.Equal(t, "AlwaysOffSampler", NeverSample().Description())

	// TraceIDRatioBased descriptions
	for _, tc := range []struct {
		prob float64
		desc string
	}{
		// Well-known values
		{0.5, "TraceIDRatioBased{0.5}"},
		{1. / 3, "TraceIDRatioBased{0.3333333333333333}"},
		{1. / 10000, "TraceIDRatioBased{0.0001}"},

		// Values very close to 1.0 round up to AlwaysOnSampler
		{1, "AlwaysOnSampler"},
		{1.5, "AlwaysOnSampler"},
		{1 - 0x1p-55, "AlwaysOnSampler"},
		{1 - 0x1p-54, "AlwaysOnSampler"},
		{1 - 0x1p-53, "AlwaysOnSampler"},

		// Values very close to 0 round down to AlwaysOffSampler
		{0, "AlwaysOffSampler"},
		{-0.5, "AlwaysOffSampler"},
		{probabilityZeroThreshold / 2, "AlwaysOffSampler"},
	} {
		require.Equal(t, tc.desc, TraceIDRatioBased(tc.prob).Description())
	}
}

// TestTraceIDRatioBasedThreshold tests the unsigned threshold value to ensure
// it is calculated correctly. The test inputs are the same as used in
// TestSamplerDescriptions.
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

func TestAlwaysSampleTracestate(t *testing.T) {
	sampler := AlwaysSample()
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := trace.ContextWithSpanContext(
		t.Context(),
		trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: 0,
			TraceState: trace.TraceState{},
		}),
	)
	samplingResult := sampler.ShouldSample(SamplingParameters{ParentContext: parentCtx})
	assert.Equal(t, RecordAndSample, samplingResult.Decision)
	assert.Equal(t, "th:0", samplingResult.Tracestate.Get("ot"))
}

func TestTraceIDRatioSamplerShouldSample(t *testing.T) {
	// TraceIDRatioBased(0.5) has threshold = 1<<55 (0x80000000000000). Randomness is derived from
	// traceID[8:16] as big-endian uint64 & 0x00ffffffffffffff.
	// - Trace ID ending in 0080000000000000 gives randomness 0x80000000000000 (>= threshold) -> RecordAndSample
	// - Trace ID ending in 007fffffffffffff gives randomness 0x7fffffffffffff (< threshold) -> Drop
	const (
		traceIDWillSample = "00000000000000000080000000000000"
		traceIDWillDrop   = "0000000000000000007fffffffffffff"
	)
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7") // non-zero span ID - it does not matter for sampling

	// FlagsRandom (0x02) signals trace ID has randomness - required for sampler to use traceID
	flagsWithRandom := trace.FlagsRandom

	t.Run("RecordAndSample adds ot.th to tracestate", func(t *testing.T) {
		sampler := TraceIDRatioBased(0.5)
		traceID, _ := trace.TraceIDFromHex(traceIDWillSample)
		initialState, err := trace.ParseTraceState("vendor=value")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: flagsWithRandom,
				TraceState: initialState,
			}),
		)
		params := SamplingParameters{
			ParentContext: parentCtx,
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, RecordAndSample, result.Decision, "expected RecordAndSample when randomness >= threshold")
		ot := result.Tracestate.Get("ot")
		require.NotEmpty(t, ot, "ot key should be set when sampler is applied")
		assert.True(t, strings.HasPrefix(ot, "th:"), "ot value should contain th (threshold) key, got %q", ot)
		// Verify vendor key is preserved
		assert.Equal(t, "value", result.Tracestate.Get("vendor"))
	})

	t.Run("RecordAndSample with empty tracestate adds ot.th", func(t *testing.T) {
		sampler := TraceIDRatioBased(0.5)
		traceID, _ := trace.TraceIDFromHex(traceIDWillSample)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: flagsWithRandom,
				TraceState: trace.TraceState{},
			}),
		)
		params := SamplingParameters{ParentContext: parentCtx, TraceID: traceID}

		result := sampler.ShouldSample(params)

		assert.Equal(t, RecordAndSample, result.Decision)
		ot := result.Tracestate.Get("ot")
		require.NotEmpty(t, ot)
		assert.True(t, strings.HasPrefix(ot, "th:"), "ot value should contain threshold, got %q", ot)
	})

	t.Run("Drop leaves tracestate unchanged when randomness < threshold", func(t *testing.T) {
		sampler := TraceIDRatioBased(0.5)
		traceID, _ := trace.TraceIDFromHex(traceIDWillDrop)
		initialState, err := trace.ParseTraceState("ot=th:0;rv:0123456789abcd,vendor=value")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: flagsWithRandom,
				TraceState: initialState,
			}),
		)
		params := SamplingParameters{ParentContext: parentCtx, TraceID: traceID}

		result := sampler.ShouldSample(params)

		assert.Equal(t, Drop, result.Decision, "expected Drop when randomness < threshold")
		assert.Equal(t, initialState, result.Tracestate, "tracestate should be unchanged when decision is Drop")
	})

	t.Run("Drop leaves tracestate unchanged when no randomness", func(t *testing.T) {
		sampler := TraceIDRatioBased(0.5)
		traceID, _ := trace.TraceIDFromHex(traceIDWillSample)
		initialState, err := trace.ParseTraceState("vendor=value")
		require.NoError(t, err)

		// No FlagsRandom - sampler cannot use trace ID, returns Drop without applying
		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: 0,
				TraceState: initialState,
			}),
		)
		params := SamplingParameters{ParentContext: parentCtx, TraceID: traceID}

		result := sampler.ShouldSample(params)

		assert.Equal(t, Drop, result.Decision, "expected Drop when trace has no randomness")
		assert.Equal(t, initialState, result.Tracestate, "tracestate should be unchanged when sampler is not applied")
	})

	t.Run("Drop with rv in tracestate leaves tracestate unchanged", func(t *testing.T) {
		sampler := TraceIDRatioBased(0.5)
		traceID, _ := trace.TraceIDFromHex("00000000000000010000000000000000")
		// rv:00000000000000 = 0, always < threshold
		initialState, err := trace.ParseTraceState("ot=rv:00000000000000,vendor=value")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: 0,
				TraceState: initialState,
			}),
		)
		params := SamplingParameters{ParentContext: parentCtx, TraceID: traceID}

		result := sampler.ShouldSample(params)

		assert.Equal(t, Drop, result.Decision)
		assert.Equal(t, initialState, result.Tracestate, "tracestate should be unchanged when decision is Drop")
	})

}
