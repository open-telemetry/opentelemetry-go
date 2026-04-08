// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x contains experimental trace features.
package x // import "go.opentelemetry.io/otel/sdk/trace/x"

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	// defaultSamplingPrecision is the default precision for threshold encoding.
	defaultSamplingPrecision = 4
	maxAdjustedCount         = 1 << 56
	// randomnessMask masks the least significant 56 bits of the trace ID per
	// W3C Trace Context Level 2 Random Trace ID Flag.
	// https://www.w3.org/TR/trace-context-2/#random-trace-id-flag
	randomnessMask = maxAdjustedCount - 1

	probabilityZeroThreshold = 1 / float64(maxAdjustedCount)
	probabilityOneThreshold  = 1 - 0x1p-52
)

// traceIDRatioSampler is the sdktrace.Sampler implementation used by
// TraceIDRatioBased.
type traceIDRatioSampler struct {
	threshold   uint64
	thkv        string
	description string
}

// ShouldSample implements sdktrace.Sampler.
func (ts *traceIDRatioSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	state := psc.TraceState()

	existingOtts := state.Get("ot")

	var randomness uint64
	var hasRandomness bool
	if existingOtts != "" {
		randomness, hasRandomness = tracestateRandomness(existingOtts)
	}

	if !hasRandomness {
		randomness = binary.BigEndian.Uint64(p.TraceID[8:16]) & randomnessMask
	}

	if ts.threshold > randomness {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.Drop,
			Tracestate: state,
		}
	}

	var newOtts string
	// Only insert/update th when randomness is available (either from
	// explicit rv value or trace ID with the random flag). Otherwise,
	// erase any existing th to signal the span is not guaranteed to be
	// statistically representative.
	// See https://opentelemetry.io/docs/specs/otel/trace/tracestate-probability-sampling/#general-requirements
	if hasRandomness || psc.TraceFlags().IsRandom() {
		newOtts = InsertOrUpdateTraceStateThKeyValue(existingOtts, ts.thkv)
	} else {
		newOtts = eraseTraceStateThKeyValue(existingOtts)
	}

	if newOtts == "" {
		state = state.Delete("ot")
		return sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample, Tracestate: state}
	}
	combined, err := state.Insert("ot", newOtts)
	if err != nil {
		otel.Handle(fmt.Errorf("could not combine tracestate: %w", err))
		return sdktrace.SamplingResult{Decision: sdktrace.Drop, Tracestate: state}
	}
	return sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample, Tracestate: combined}
}

// Description implements sdktrace.Sampler.
func (ts *traceIDRatioSampler) Description() string {
	return ts.description
}

// TraceIDRatioBased samples a given fraction of traces. Fractions >= 1 will
// always sample. Fractions < 0 are treated as zero. To respect the parent
// trace's SampledFlag, the TraceIDRatioBased sampler should be used as a
// delegate of a ParentBased sampler.
//
//nolint:revive // revive complains about stutter of `x.TraceIDRatioBased`
func TraceIDRatioBased(fraction float64) sdktrace.Sampler {
	const (
		maxp  = 14
		defp  = defaultSamplingPrecision
		hbits = 4
	)
	if fraction > probabilityOneThreshold {
		return sdktrace.AlwaysSample()
	}
	if fraction < probabilityZeroThreshold {
		return sdktrace.NeverSample()
	}

	_, expF := math.Frexp(fraction)
	_, expR := math.Frexp(1 - fraction)
	precision := min(maxp, max(defp+expF/-hbits, defp+expR/-hbits))

	scaled := uint64(math.Round(fraction * float64(maxAdjustedCount)))
	threshold := maxAdjustedCount - scaled

	if shift := hbits * (maxp - precision); shift != 0 {
		half := uint64(1) << (shift - 1)
		threshold += half
		threshold >>= shift
		threshold <<= shift
	}

	tvalue := strings.TrimRight(strconv.FormatUint(maxAdjustedCount+threshold, 16)[1:], "0")
	return &traceIDRatioSampler{
		threshold:   threshold,
		thkv:        "th:" + tvalue,
		description: fmt.Sprintf("TraceIDRatioBased{%g}", fraction),
	}
}
