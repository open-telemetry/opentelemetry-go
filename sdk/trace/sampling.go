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

	// Record indicates the span's `IsRecording() == true`, but `Sampled` flag
	// *must not* be set.
	RecordOnly

	// RecordAndSample has span's `IsRecording() == true` and `Sampled` flag
	// *must* be set.
	RecordAndSample
)

// SamplingResult conveys a SamplingDecision, set of Attributes and a Tracestate.
type SamplingResult struct {
	Decision   SamplingDecision
	Attributes []attribute.KeyValue
	Tracestate trace.TraceState
}

type traceIDRatioSampler struct {
	// threshold is a rejection threshold.
	// Select when (T <= R)
	// Drop when (T > R)
	// Range is [0, 1<<56).
	threshold uint64

	// otts is the encoded OTel trace state field, containing "th:<tvalue>"
	otts string

	description string
}

// tracestateHasRandomness determines whether there is a "rv" sub-key
// in `otts` which is the OTel tracestate value (i.e., the top-level "ot" value).
func tracestateHasRandomness(otts string) (randomness uint64, hasRandom bool) {
	var low int
	if has := strings.HasPrefix(otts, "rv:"); has {
		low = 3
	} else if pos := strings.Index(otts, ";rv:"); pos > 0 {
		low = pos + 4
	} else {
		return 0, false
	}
	if len(otts) < low+14 {
		otel.Handle(fmt.Errorf("could not parse tracestate randomness: %q: %w", otts, strconv.ErrSyntax))
	} else if len(otts) > low+14 && otts[low+14] != ';' {
		otel.Handle(fmt.Errorf("could not parse tracestate randomness: %q: %w", otts, strconv.ErrSyntax))
	} else {
		randomIn := otts[low : low+14]
		if rv, err := strconv.ParseUint(randomIn, 16, 64); err == nil {
			randomness = rv
			hasRandom = true
		} else {
			otel.Handle(fmt.Errorf("could not parse tracestate randomness: %q: %w", randomIn, err))
		}
	}
	return
}

func (ts traceIDRatioSampler) ShouldSample(p SamplingParameters) SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	state := psc.TraceState()

	existOtts := state.Get("ot")

	var randomness uint64
	var hasRandom bool
	if existOtts != "" {
		// When the OTel trace state field exists, we will
		// inspect for a "rv", otherwise assume that the
		// TraceID is random.
		randomness, hasRandom = tracestateHasRandomness(existOtts)
	}
	if !hasRandom {
		// Interpret the least-significant 8-bytes as an
		// unsigned number, then zero the top 8 bits using
		// randomnessMask, yielding the least-significant 56
		// bits of randomness, as specified in W3C Trace
		// Context Level 2.
		randomness = binary.BigEndian.Uint64(p.TraceID[8:16]) & randomnessMask
	}
	if ts.threshold > randomness {
		return SamplingResult{
			Decision:   Drop,
			Tracestate: state,
		}
	}

	if mod, err := state.Insert("ot", combineTracestate(existOtts, ts.otts)); err == nil {
		state = mod
	} else {
		otel.Handle(fmt.Errorf("could not update tracestate: %q", err))
	}
	return SamplingResult{
		Decision:   RecordAndSample,
		Tracestate: state,
	}
}

// combineTracestate combines an existing OTel tracestate fragment,
// which is the value of a top-level "ot" tracestate vendor tag.
func combineTracestate(incoming, updated string) string {
	// `incoming` is formatted according to the OTel tracestate
	// spec, with colon separating two-byte key and value, with
	// semi-colon separating key-value pairs.
	//
	// `updated` should be a single two-byte key:value to modify
	// or insert therefore colonOffset is 2 bytes, valueOffset is
	// 3 bytes into `incoming`.
	const colonOffset = 2
	const valueOffset = colonOffset + 1

	if incoming == "" {
		return updated
	}
	var out strings.Builder

	// The update is expected to be a single key-value of the form
	// `XX:value` for with two-character key.
	upkey := updated[:colonOffset]

	// In this case, there is an existing field under "ot" and we
	// need to combine.  We will pass the parts of "incoming"
	// through except the field we are updating, which we will
	// modify if it is found.
	foundUp := false

	for count := 0; len(incoming) != 0; count++ {
		key, rest, hasCol := strings.Cut(incoming, ":")
		if !hasCol {
			// return the updated value, ignore invalid inputs
			return updated
		}
		value, next, _ := strings.Cut(rest, ";")

		if key == upkey {
			value = updated[valueOffset:]
			foundUp = true
		}
		if count != 0 {
			out.WriteString(";")
		}
		out.WriteString(key)
		out.WriteString(":")
		out.WriteString(value)

		incoming = next
	}
	if !foundUp {
		out.WriteString(";")
		out.WriteString(updated)
	}
	return out.String()
}

func (ts traceIDRatioSampler) Description() string {
	return ts.description
}

const (
	// DefaultSamplingPrecision is the number of hexadecimal
	// digits of precision used to expressed the samplling probability.
	DefaultSamplingPrecision = 4

	// MinSupportedProbability is the smallest probability that
	// can be encoded by this implementation, and it defines the
	// smallest interval between probabilities across the range.
	// The largest supported probability is (1-MinSupportedProbability).
	//
	// This value corresponds with the size of a float64
	// significand, because it simplifies this implementation to
	// restrict the probability to use 52 bits (vs 56 bits).
	minSupportedProbability float64 = 1 / float64(maxAdjustedCount)

	// maxSupportedProbability is the number closest to 1.0 (i.e.,
	// near 99.999999%) that is not equal to 1.0 in terms of the
	// float64 representation, having 52 bits of significand.
	// Other ways to express this number:
	//
	//   0x1.ffffffffffffe0p-01
	//   0x0.fffffffffffff0p+00
	//   math.Nextafter(1.0, 0.0)
	maxSupportedProbability float64 = 1 - 0x1p-52

	// maxAdjustedCount is the inverse of the smallest
	// representable sampling probability, it is the number of
	// distinct 56 bit values.
	maxAdjustedCount uint64 = 1 << 56

	// randomnessMask is a mask that selects the least-significant
	// 56 bits of a uint64.
	randomnessMask uint64 = maxAdjustedCount - 1
)

// TraceIDRatioBased samples a given fraction of traces. Fractions >= 1 will
// always sample. Fractions < 0 are treated as zero. To respect the
// parent trace's `SampledFlag`, the `TraceIDRatioBased` sampler should be used
// as a delegate of a `Parent` sampler.
//
//nolint:revive // revive complains about stutter of `trace.TraceIDRatioBased`
func TraceIDRatioBased(fraction float64) Sampler {
	const (
		maxp  = 14                       // maximum precision is 56 bits
		defp  = DefaultSamplingPrecision // default precision
		hbits = 4                        // bits per hex digit
	)

	if fraction > 1-0x1p-52 {
		return AlwaysSample()
	}

	if fraction < minSupportedProbability {
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
		otts:        fmt.Sprint("th:", tvalue),
		description: fmt.Sprintf("TraceIDRatioBased{%g}", fraction),
	}
}

type alwaysOnSampler struct{}

func (as alwaysOnSampler) ShouldSample(p SamplingParameters) SamplingResult {
	ts := trace.SpanContextFromContext(p.ParentContext).TraceState()
	// 100% sampling equals zero rejection threshold.
	if mod, err := ts.Insert("ot", combineTracestate(ts.Get("ot"), "th:0")); err == nil {
		ts = mod
	} else {
		otel.Handle(fmt.Errorf("could not update tracestate: %w", err))
	}
	return SamplingResult{
		Decision:   RecordAndSample,
		Tracestate: ts,
	}
}

func (as alwaysOnSampler) Description() string {
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

func (as alwaysOffSampler) ShouldSample(p SamplingParameters) SamplingResult {
	return SamplingResult{
		Decision:   Drop,
		Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
	}
}

func (as alwaysOffSampler) Description() string {
	return "AlwaysOffSampler"
}

// NeverSample returns a Sampler that samples no traces.
func NeverSample() Sampler {
	return alwaysOffSampler{}
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
