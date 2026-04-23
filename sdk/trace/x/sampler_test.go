// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestProbabilitySampler(t *testing.T) {
	t.Run("description", func(t *testing.T) {
		for _, tc := range []struct {
			prob float64
			desc string
		}{
			{0.5, "ProbabilitySampler{0.5}"},
			{1. / 3, "ProbabilitySampler{0.3333333333333333}"},
			{1. / 10000, "ProbabilitySampler{0.0001}"},
			{1, "ProbabilitySampler{1}"},
			{1.5, "ProbabilitySampler{1}"},
			{0, "AlwaysOffSampler"},
			{-0.5, "AlwaysOffSampler"},
		} {
			require.Equal(t, tc.desc, ProbabilitySampler(tc.prob).Description())
		}
	})

	t.Run("threshold", func(t *testing.T) {
		for _, tc := range []struct {
			prob      float64
			threshold uint64
		}{
			{0.5, 0x80000000000000},
			{1 / 3.0, 0xaaab0000000000},
			{2 / 3.0, 0x55550000000000},
			{1 / 10.0, 0xe6660000000000},
			{1 / 256.0, 0xff000000000000},
			{1 / 65536.0, 0xffff0000000000},
			{1 / 1048576.0, 0xfffff000000000},
		} {
			sampler := ProbabilitySampler(tc.prob).(*probabilitySampler)
			require.Equal(t, tc.threshold, sampler.threshold)
		}
	})

	t.Run("probability one uses probabilitySampler with th:0", func(t *testing.T) {
		for _, prob := range []float64{1, 1.5} {
			t.Run(fmt.Sprintf("%g", prob), func(t *testing.T) {
				t.Helper()
				sampler := ProbabilitySampler(prob).(*probabilitySampler)
				require.Equal(t, uint64(0), sampler.threshold)
				require.Equal(t, "th:0", sampler.thkv)
				require.Equal(t, "ProbabilitySampler{1}", sampler.Description())
			})
		}
	})

	t.Run("probability one always samples including zero randomness", func(t *testing.T) {
		// Trace ID all zeros yields randomness 0; ProbabilitySampler(0.5) drops this case.
		sampler := ProbabilitySampler(1)
		traceID, _ := trace.TraceIDFromHex("00000000000000000000000000000000")
		spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
		initialState, err := trace.ParseTraceState("vendor=value")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsRandom,
				TraceState: initialState,
			}),
		)
		params := sdktrace.SamplingParameters{
			ParentContext: parentCtx,
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.RecordAndSample, result.Decision)
		ot := result.Tracestate.Get("ot")
		require.NotEmpty(t, ot)
		assert.Contains(t, ot, "th:0")
		assert.Equal(t, "value", result.Tracestate.Get("vendor"))
	})

	t.Run("inclusive sampling", func(t *testing.T) {
		const numSamplers = 100
		const numTraces = 50
		for range numSamplers {
			ratioLo, ratioHi := mrand.Float64(), mrand.Float64()
			if ratioHi < ratioLo {
				ratioLo, ratioHi = ratioHi, ratioLo
			}
			samplerHi := ProbabilitySampler(ratioHi)
			samplerLo := ProbabilitySampler(ratioLo)
			for range numTraces {
				traceID := trace.TraceID{}
				_, _ = rand.Read(traceID[:])
				params := sdktrace.SamplingParameters{
					ParentContext: trace.ContextWithSpanContext(
						t.Context(),
						trace.NewSpanContext(trace.SpanContextConfig{
							TraceID:    traceID,
							TraceFlags: trace.FlagsRandom,
						}),
					),
					TraceID: traceID,
				}
				if samplerLo.ShouldSample(params).Decision == sdktrace.RecordAndSample {
					assert.Equal(t, sdktrace.RecordAndSample, samplerHi.ShouldSample(params).Decision,
						"%s sampled but %s did not", samplerLo.Description(), samplerHi.Description())
				}
			}
		}
	})

	t.Run("RecordAndSample adds ot.th to tracestate", func(t *testing.T) {
		const traceIDWillSample = "00000000000000000080000000000000"
		sampler := ProbabilitySampler(0.5)
		traceID, _ := trace.TraceIDFromHex(traceIDWillSample)
		spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
		initialState, err := trace.ParseTraceState("vendor=value")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsRandom,
				TraceState: initialState,
			}),
		)
		params := sdktrace.SamplingParameters{
			ParentContext: parentCtx,
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.RecordAndSample, result.Decision)
		ot := result.Tracestate.Get("ot")
		require.NotEmpty(t, ot)
		assert.True(t, strings.HasPrefix(ot, "th:"), "ot value should contain th key, got %q", ot)
		assert.Equal(t, "value", result.Tracestate.Get("vendor"))
	})

	t.Run("RecordAndSample with explicit rv and no randomness flag inserts th in tracestate", func(t *testing.T) {
		sampler := ProbabilitySampler(0.5)
		traceID, _ := trace.TraceIDFromHex("00000000000000000000000000000001")
		spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
		initialState, err := trace.ParseTraceState("ot=rv:80000000000000,vendor=value")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.TraceFlags(0),
				TraceState: initialState,
			}),
		)
		params := sdktrace.SamplingParameters{
			ParentContext: parentCtx,
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.RecordAndSample, result.Decision, "rv value should be used for sampling decision")
		ot := result.Tracestate.Get("ot")
		require.NotEmpty(t, ot)
		assert.Contains(t, ot, "th:", "ot value should contain th when rv is present, got %q", ot)
		assert.Equal(t, "value", result.Tracestate.Get("vendor"))
	})

	t.Run("RecordAndSample without randomness flag erases ot.th from tracestate", func(t *testing.T) {
		sampler := ProbabilitySampler(0.5)
		traceID, _ := trace.TraceIDFromHex("00000000000000000080000000000000")
		spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
		initialState, err := trace.ParseTraceState("ot=th:0ad;other:value,vendor=v")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.TraceFlags(0),
				TraceState: initialState,
			}),
		)
		params := sdktrace.SamplingParameters{
			ParentContext: parentCtx,
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.RecordAndSample, result.Decision)
		ot := result.Tracestate.Get("ot")
		assert.NotContains(
			t,
			ot,
			"th:",
			"ot value should not contain th when TraceFlags has no randomness flag and no rv in tracestate, got %q",
			ot,
		)
		assert.Equal(t, "v", result.Tracestate.Get("vendor"))
	})

	t.Run("RecordAndSample when ot becomes empty deletes ot from tracestate", func(t *testing.T) {
		sampler := ProbabilitySampler(0.5)
		traceID, _ := trace.TraceIDFromHex("00000000000000000080000000000000")
		spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
		initialState, err := trace.ParseTraceState("ot=th:0ad,vendor=value")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.TraceFlags(0),
				TraceState: initialState,
			}),
		)
		params := sdktrace.SamplingParameters{
			ParentContext: parentCtx,
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.RecordAndSample, result.Decision)
		assert.Empty(t, result.Tracestate.Get("ot"))
		assert.Equal(t, "value", result.Tracestate.Get("vendor"))
	})

	t.Run("Drop when randomness < threshold", func(t *testing.T) {
		const traceIDWillDrop = "0000000000000000007fffffffffffff"
		sampler := ProbabilitySampler(0.5)
		traceID, _ := trace.TraceIDFromHex(traceIDWillDrop)
		spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
		initialState, err := trace.ParseTraceState("ot=th:0;rv:0123456789abcd,vendor=value")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsRandom,
				TraceState: initialState,
			}),
		)
		params := sdktrace.SamplingParameters{
			ParentContext: parentCtx,
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.Drop, result.Decision)
		assert.Equal(t, initialState, result.Tracestate)
	})

	t.Run("root span RecordAndSample", func(t *testing.T) {
		sampler := ProbabilitySampler(0.5)
		traceID, _ := trace.TraceIDFromHex("00000000000000000080000000000000")
		params := sdktrace.SamplingParameters{
			ParentContext: t.Context(),
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.RecordAndSample, result.Decision)
		ot := result.Tracestate.Get("ot")
		require.NotEmpty(t, ot)
		assert.Contains(t, ot, "th:")
	})

	t.Run("root span Drop", func(t *testing.T) {
		sampler := ProbabilitySampler(0.5)
		traceID, _ := trace.TraceIDFromHex("00000000000000000000000000000000")
		params := sdktrace.SamplingParameters{
			ParentContext: t.Context(),
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.Drop, result.Decision)
		assert.Empty(t, result.Tracestate.Get("ot"))
	})

	t.Run("RecordAndSample updates existing th in tracestate", func(t *testing.T) {
		sampler := ProbabilitySampler(0.5)
		traceID, _ := trace.TraceIDFromHex("00000000000000000080000000000000")
		spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
		initialState, err := trace.ParseTraceState("ot=th:0ad;other:value,vendor=v")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsRandom,
				TraceState: initialState,
			}),
		)
		params := sdktrace.SamplingParameters{
			ParentContext: parentCtx,
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.RecordAndSample, result.Decision)
		ot := result.Tracestate.Get("ot")
		require.NotEmpty(t, ot)
		assert.True(t, strings.HasPrefix(ot, "th:8"), "ot should have updated th for 0.5 sampler, got %q", ot)
		assert.Equal(t, "v", result.Tracestate.Get("vendor"))
	})

	t.Run("trace ID all zeros Drop", func(t *testing.T) {
		sampler := ProbabilitySampler(0.5)
		traceID, _ := trace.TraceIDFromHex("00000000000000000000000000000000")
		spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
		initialState, err := trace.ParseTraceState("vendor=value")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsRandom,
				TraceState: initialState,
			}),
		)
		params := sdktrace.SamplingParameters{
			ParentContext: parentCtx,
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.Drop, result.Decision)
	})

	t.Run("trace ID max randomness RecordAndSample", func(t *testing.T) {
		sampler := ProbabilitySampler(0.5)
		traceID, _ := trace.TraceIDFromHex("000000000000000000ffffffffffffff")
		spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
		initialState, err := trace.ParseTraceState("vendor=value")
		require.NoError(t, err)

		parentCtx := trace.ContextWithSpanContext(
			t.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsRandom,
				TraceState: initialState,
			}),
		)
		params := sdktrace.SamplingParameters{
			ParentContext: parentCtx,
			TraceID:       traceID,
		}

		result := sampler.ShouldSample(params)

		assert.Equal(t, sdktrace.RecordAndSample, result.Decision)
		ot := result.Tracestate.Get("ot")
		require.NotEmpty(t, ot)
		assert.True(t, strings.HasPrefix(ot, "th:"), "ot should contain th, got %q", ot)
	})
}
