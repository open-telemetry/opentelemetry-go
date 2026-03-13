// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Sampler decides whether a trace should be sampled and exported.
type Sampler interface {
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// ShouldSample returns a SamplingResult based on a decision made from the
	// passed parameters.
	ShouldSample(parameters SamplingParameters) SamplingResult
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// Description returns information describing the Sampler.
	Description() string
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.
}

// SamplingParameters contains the values passed to a Sampler.
type SamplingParameters struct {
	ParentContext context.Context
	TraceID       trace.TraceID
	Name          string
	Kind          trace.SpanKind
	Attributes    []attribute.KeyValue
	Links         []trace.Link
}

// SamplingDecision indicates whether a span is dropped, recorded and/or sampled.
type SamplingDecision uint8

// Valid sampling decisions.
const (
	// Drop will not record the span and all attributes/events will be dropped.
	Drop SamplingDecision = iota

	// RecordOnly indicates the span's IsRecording method returns true, but trace.FlagsSampled flag
	// must not be set.
	RecordOnly

	// RecordAndSample indicates the span's IsRecording method returns true and trace.FlagsSampled flag
	// must be set.
	RecordAndSample
)

// SamplingResult conveys a SamplingDecision, set of Attributes and a Tracestate.
type SamplingResult struct {
	Decision   SamplingDecision
	Attributes []attribute.KeyValue
	Tracestate trace.TraceState
}

type traceIDRatioSampler struct {
	// threshold T is a rejection threshold with range [0, 1<<56).
	// It is used to compare with a random value R, such as a random trace ID.
	// Select when (R >= T).
	// Drop when (R < T).
	threshold uint64

	// thkv is the encoded OTel tracestate threshold key-value pair in the "ot" key.
	// Applying traceIdRatioBased sampler to a trace context will result in tracestate "ot" key value containing `th:<tvalue>`
	// where <tvalue> is a hex digit string representing the threshold T.
	thkv string

	description string
}

// Determines whether there is a randomness "rv" sub-key in `otts` (the topc level OTel tracestate field).
// If present, "rv" is a 56-bit unsigned integer, encoded in 14 hexdigits.
func tracestateRandomness(otts string) (randomness uint64, hasRandomness bool) {
	var start int // start index of the "rv" sub-key
	if strings.HasPrefix(otts, "rv:") {
		start = 3
	} else if idx := strings.Index(otts, ";rv:"); idx != -1 {
		start = idx + 4
	} else {
		return 0, false
	}

	if len(otts) < start+14 || (len(otts) > start+14 && otts[start+14] != ';') {
		otel.Handle(fmt.Errorf("could not parse tracestate randomness: %s", otts))
		return 0, false
	}

	if rv, err := strconv.ParseUint(otts[start:start+14], 16, 56); err != nil {
		otel.Handle(fmt.Errorf("could not parse tracestate randomness: %s", otts))
		return 0, false
	} else {
		randomness = rv
		hasRandomness = true
	}
	return
}

func (ts traceIDRatioSampler) ShouldSample(p SamplingParameters) SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	state := psc.TraceState()

	existingOtts := state.Get("ot")

	var randomness uint64
	var hasRandomness bool
	// When there is an explicit rv value in tracestate, we always make use of it.
	// Otherwise, we use the traceID if trace flags indicates it has randomness.
	if existingOtts != "" {
		randomness, hasRandomness = tracestateRandomness(existingOtts)
	}

	// TODO: rethink this and compare with the alternative of just assuming trace id is random.
	if !hasRandomness {
		randomness = binary.BigEndian.Uint64(p.TraceID[8:16]) & randomnessMask
	}

	if ts.threshold > randomness {
		return SamplingResult{
			Decision:   Drop,
			Tracestate: state,
		}
	}
	// If the trace flags do not indicate that the trace ID is random, we proceed with
	// trace ID as the source of randomness regardless. However, this will not ensure that
	// the span is sampled with a known probability.
	// As per the general requirement of the spec, https://opentelemetry.io/docs/specs/otel/trace/tracestate-probability-sampling/#general-requirements
	// "Sampling stages that yield spans with unknown sampling probability, ..., must erase
	// the OpenTelemetry threshold value in their output."
	var newOtts string
	if !psc.TraceFlags().IsRandom() {
		newOtts = eraseTraceStateThKeyValue(existingOtts)
	} else {
		newOtts = insertOrUpdateTraceStateThKeyValue(existingOtts, ts.thkv)
	}

	// If we are not able to update the tracestate, we drop the span.
	// Think of the following scenario:
	// If a span was previously sampled by a traceIdRatioBased sampler with a lower threshold,
	// it should be updated to reflect to the higher threshold of the current sampler.
	// All spans at this staged is already filtered by the higher threshold. If we are not able to
	// update the tracestate of the SamplingResult and we do not drop but instead record and sample the span,
	// when the span is exported and used to compute metrics, it'll be inappropriately adjusted.
	// This skews the metrics and the failure is not loud.
	// The alternative we chose (to DROP) is not perfect either. However, we do not expect the insert
	// error to happen except perhapswith some internal error. In this case, fail-closed will at least be loud.
	if combined, err := state.Insert("ot", newOtts); err != nil {
		otel.Handle(fmt.Errorf("could not combine tracestate: %w", err))
		return SamplingResult{Decision: Drop, Tracestate: state}
	} else {
		state = combined
	}
	return SamplingResult{Decision: RecordAndSample, Tracestate: state}
}

func eraseTraceStateThKeyValue(otts string) string {
	start := strings.Index(otts, "th:")
	if start == -1 {
		return otts
	}
	if start > 0 && otts[start-1] == ';' {
		start--
	}
	end := -1
	for end = start + 1; end < len(otts); end++ {
		if otts[end] == ';' {
			if start == 0 {
				end++
			}
			break
		}
	}
	if end == len(otts) {
		return otts[0:start]
	}
	return otts[0:start] + otts[end:]
}

func insertOrUpdateTraceStateThKeyValue(existingOtts, thkv string) string {
	if existingOtts == "" {
		return thkv
	}

	start := -1
	end := -1
	if strings.HasPrefix(existingOtts, "th:") {
		start = 0
	} else if idx := strings.Index(existingOtts, ";th:"); idx != -1 {
		start = idx + 1
	}
	if start == -1 {
		return thkv + ";" + existingOtts
	}

	for end = start; end < len(existingOtts); end++ {
		if existingOtts[end] == ';' {
			end++
			break
		}
	}

	if end == len(existingOtts) {
		return strings.TrimSuffix(thkv+";"+existingOtts[0:start], ";")
	}
	return thkv + ";" + existingOtts[0:start] + existingOtts[end:]
}

const (
	DefaultSamplingPrecision = 4
	maxAdjustedCount         = 1 << 56
	// Mask the least significant 56 bits of the trace ID as per W3C Trace Context Level 2 Random Trace ID Flag.
	// https://www.w3.org/TR/trace-context-2/#random-trace-id-flag
	randomnessMask = maxAdjustedCount - 1

	// probabilityZeroThreshold is the smallest probability that
	// can be encoded by this implementation, and it defines the
	// smallest interval between probabilities across the range.
	// Probabilities below this threshold are treated as zero.
	//
	// This value corresponds with the size of a float64
	// significand, because it simplifies this implementation to
	// restrict the probability to use 52 bits (vs 56 bits).
	probabilityZeroThreshold float64 = 1 / float64(maxAdjustedCount)

	// probabilityOneThreshold is the number closest to 1.0 (i.e.,
	// near 99.999999%) that is not equal to 1.0 in terms of the
	// float64 representation, with 52 bits of significand.
	// This is the largest float64 representable with only the 52 bits of significand.
	// Probabilities above this threshold are treated as one.
	// Other ways to express this number:
	//
	//   0x1.ffffffffffffe0p-01
	//   0x0.fffffffffffff0p+00
	//   math.Nextafter(1.0, 0.0)
	probabilityOneThreshold float64 = 1 - 0x1p-52
)

func (ts traceIDRatioSampler) Description() string {
	return ts.description
}

// TraceIDRatioBased samples a given fraction of traces. Fractions >= 1 will
// always sample. Fractions < 0 are treated as zero. To respect the
// parent trace's `SampledFlag`, the `TraceIDRatioBased` sampler should be used
// as a delegate of a `Parent` sampler.
//
//nolint:revive // revive complains about stutter of `trace.TraceIDRatioBased`
func TraceIDRatioBased(fraction float64) Sampler {
	const (
		maxp  = 14 // maximum precision
		defp  = DefaultSamplingPrecision
		hbits = 4 // bits per hex digit
	)
	if fraction > probabilityOneThreshold {
		return AlwaysSample()
	}
	if fraction < probabilityZeroThreshold {
		return NeverSample()
	}

	// Calculate the amount of precision needed to encode the
	// threshold with reasonable precision.
	//
	// 13 hex digits is the maximum reasonable precision, since
	// that equals 52 bits, the number of bits in the float64
	// significand.
	//
	// Frexp() normalizes both the fraction and one-minus the
	// fraction, because more digits of precision are needed in
	// both cases -- in these cases the threshold has all leading
	// '0' or 'f' characters.
	//
	// We know that `exp <= 0`.  If `exp <= -4`, there will be a
	// leading hex `0` or `f`.  For every multiple of -4, another
	// leading `0` or `f` appears, so this raises precision
	// accordingly.
	_, expF := math.Frexp(fraction)
	_, expR := math.Frexp(1 - fraction)
	precision := min(maxp, max(defp+expF/-hbits, defp+expR/-hbits))

	// Compute the threshold
	scaled := uint64(math.Round(fraction * float64(maxAdjustedCount)))
	threshold := maxAdjustedCount - scaled

	// Round to the specified precision, if less than the maximum.
	if shift := hbits * (maxp - precision); shift != 0 {
		half := uint64(1) << (shift - 1)
		threshold += half
		threshold >>= shift
		threshold <<= shift
	}

	// Add maxAdjustedCount so that leading-zeros are formatted by
	// the strconv library after an artificial leading "1".  Then,
	// strip the leadingt "1", then remove trailing zeros.
	tvalue := strings.TrimRight(strconv.FormatUint(maxAdjustedCount+threshold, 16)[1:], "0")
	return &traceIDRatioSampler{
		threshold:   threshold,
		thkv:        "th:" + tvalue,
		description: fmt.Sprintf("TraceIDRatioBased{%g}", fraction),
	}
}

type alwaysOnSampler struct{}

func (alwaysOnSampler) ShouldSample(p SamplingParameters) SamplingResult {
	ts := trace.SpanContextFromContext(p.ParentContext).TraceState()
	if mod, err := ts.Insert("ot", insertOrUpdateTraceStateThKeyValue(ts.Get("ot"), "th:0")); err != nil {
		otel.Handle(fmt.Errorf("could not update threshold (`ot.th`) in tracestate: %w", err))
		// I feel this is a contentional decision, but I'm putting this here for discussion.
		// The contention is:
		// If we do this, we kind of violate what "alwaysOn" sampling means.
		// On the other hand, I don't know whether the semantics should apply in internal error scenarios.
		// In addition, if we are unable to update the threshold and the ts has an existing threshold,
		// downstream span metrics will be wrong.
		// Finally, I want to point out that the contention is probably moot as this path is not likely to happen.
		// return SamplingResult{Decision: Drop, Tracestate: ts}
	} else {
		ts = mod
	}
	return SamplingResult{Decision: RecordAndSample, Tracestate: ts}
}

func (alwaysOnSampler) Description() string {
	// https://opentelemetry.io/docs/specs/otel/trace/sdk/#alwayson
	return "AlwaysOnSampler"
}

// AlwaysSample returns a Sampler that samples every trace.
// Be careful about using this sampler in a production application with
// significant traffic: a new trace will be started and exported for every
// request.
func AlwaysSample() Sampler {
	return alwaysOnSampler{}
}

type alwaysOffSampler struct{}

func (alwaysOffSampler) ShouldSample(p SamplingParameters) SamplingResult {
	return SamplingResult{
		Decision:   Drop,
		Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
	}
}

func (alwaysOffSampler) Description() string {
	// https://opentelemetry.io/docs/specs/otel/trace/sdk/#alwaysoff
	return "AlwaysOffSampler"
}

// NeverSample returns a Sampler that samples no traces.
func NeverSample() Sampler {
	return alwaysOffSampler{}
}

type predeterminedSampler struct {
	description string
	decision    SamplingDecision
}

func (s predeterminedSampler) ShouldSample(p SamplingParameters) SamplingResult {
	return SamplingResult{
		Decision:   s.decision,
		Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
	}
}

func (s predeterminedSampler) Description() string {
	return s.description
}

// ParentBased returns a sampler decorator which behaves differently,
// based on the parent of the span. If the span has no parent,
// the decorated sampler is used to make sampling decision. If the span has
// a parent, depending on whether the parent is remote and whether it
// is sampled, one of the following samplers will apply:
//   - remoteParentSampled(Sampler) (default: AlwaysOn)
//   - remoteParentNotSampled(Sampler) (default: AlwaysOff)
//   - localParentSampled(Sampler) (default: AlwaysOn)
//   - localParentNotSampled(Sampler) (default: AlwaysOff)
func ParentBased(root Sampler, samplers ...ParentBasedSamplerOption) Sampler {
	return parentBased{
		root:   root,
		config: configureSamplersForParentBased(samplers),
	}
}

type parentBased struct {
	root   Sampler
	config samplerConfig
}

func configureSamplersForParentBased(samplers []ParentBasedSamplerOption) samplerConfig {
	c := samplerConfig{
		remoteParentSampled:    AlwaysSample(),
		remoteParentNotSampled: NeverSample(),
		localParentSampled:     AlwaysSample(),
		localParentNotSampled:  NeverSample(),
	}

	for _, so := range samplers {
		c = so.apply(c)
	}

	return c
}

// samplerConfig is a group of options for parentBased sampler.
type samplerConfig struct {
	remoteParentSampled, remoteParentNotSampled Sampler
	localParentSampled, localParentNotSampled   Sampler
}

// ParentBasedSamplerOption configures the sampler for a particular sampling case.
type ParentBasedSamplerOption interface {
	apply(samplerConfig) samplerConfig
}

// WithRemoteParentSampled sets the sampler for the case of sampled remote parent.
func WithRemoteParentSampled(s Sampler) ParentBasedSamplerOption {
	return remoteParentSampledOption{s}
}

type remoteParentSampledOption struct {
	s Sampler
}

func (o remoteParentSampledOption) apply(config samplerConfig) samplerConfig {
	config.remoteParentSampled = o.s
	return config
}

// WithRemoteParentNotSampled sets the sampler for the case of remote parent
// which is not sampled.
func WithRemoteParentNotSampled(s Sampler) ParentBasedSamplerOption {
	return remoteParentNotSampledOption{s}
}

type remoteParentNotSampledOption struct {
	s Sampler
}

func (o remoteParentNotSampledOption) apply(config samplerConfig) samplerConfig {
	config.remoteParentNotSampled = o.s
	return config
}

// WithLocalParentSampled sets the sampler for the case of sampled local parent.
func WithLocalParentSampled(s Sampler) ParentBasedSamplerOption {
	return localParentSampledOption{s}
}

type localParentSampledOption struct {
	s Sampler
}

func (o localParentSampledOption) apply(config samplerConfig) samplerConfig {
	config.localParentSampled = o.s
	return config
}

// WithLocalParentNotSampled sets the sampler for the case of local parent
// which is not sampled.
func WithLocalParentNotSampled(s Sampler) ParentBasedSamplerOption {
	return localParentNotSampledOption{s}
}

type localParentNotSampledOption struct {
	s Sampler
}

func (o localParentNotSampledOption) apply(config samplerConfig) samplerConfig {
	config.localParentNotSampled = o.s
	return config
}

func (pb parentBased) ShouldSample(p SamplingParameters) SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	if psc.IsValid() {
		if psc.IsRemote() {
			if psc.IsSampled() {
				return pb.config.remoteParentSampled.ShouldSample(p)
			}
			return pb.config.remoteParentNotSampled.ShouldSample(p)
		}

		if psc.IsSampled() {
			return pb.config.localParentSampled.ShouldSample(p)
		}
		return pb.config.localParentNotSampled.ShouldSample(p)
	}
	return pb.root.ShouldSample(p)
}

func (pb parentBased) Description() string {
	return fmt.Sprintf("ParentBased{root:%s,remoteParentSampled:%s,"+
		"remoteParentNotSampled:%s,localParentSampled:%s,localParentNotSampled:%s}",
		pb.root.Description(),
		pb.config.remoteParentSampled.Description(),
		pb.config.remoteParentNotSampled.Description(),
		pb.config.localParentSampled.Description(),
		pb.config.localParentNotSampled.Description(),
	)
}

// AlwaysRecord returns a sampler decorator which ensures that every span
// is passed to the SpanProcessor, even those that would be normally dropped.
// It converts `Drop` decisions from the root sampler into `RecordOnly` decisions,
// allowing processors to see all spans without sending them to exporters. This is
// typically used to enable accurate span-to-metrics processing.
func AlwaysRecord(root Sampler) Sampler {
	return alwaysRecord{root}
}

type alwaysRecord struct {
	root Sampler
}

func (ar alwaysRecord) ShouldSample(p SamplingParameters) SamplingResult {
	rootSamplerSamplingResult := ar.root.ShouldSample(p)
	if rootSamplerSamplingResult.Decision == Drop {
		return SamplingResult{
			Decision:   RecordOnly,
			Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
		}
	}
	return rootSamplerSamplingResult
}

func (ar alwaysRecord) Description() string {
	return "AlwaysRecord{root:" + ar.root.Description() + "}"
}
