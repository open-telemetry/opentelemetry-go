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

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
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

// Valid sampling decisions
const (
	// Drop will not record the span and all attributes/events will be dropped
	Drop SamplingDecision = iota

	// Record indicates the span's `IsRecording() == true`, but `Sampled` flag
	// *must not* be set
	RecordOnly

	// RecordAndSample has span's `IsRecording() == true` and `Sampled` flag
	// *must* be set
	RecordAndSample
)

// OTEP 168 sampling constants
const (
	otelTraceStateKey                    = "otel"
	otelTraceStateProbabilityValueSubkey = "p"
	otelTraceStateRandomValueSubkey      = "r"
	otelTraceStateNumberBase             = 16
	otelSamplingAttributePrefix          = "sampler."
	otelSamplingAdjustedCountKey         = otelSamplingAttributePrefix + "adjusted_count"
	otelSamplingNameKey                  = otelSamplingAttributePrefix + "name"
	otelSamplingParentSampler            = "parent"

	// TODO: the zero value is the maximum permissible value of
	// the OTel TraceState `random` value, which equates with zero
	otelSamplingZeroValue  = 0x3f
	otelTraceStateBitWidth = 6
)

// OTEP 168 sampling variables
var (
	errTraceStateSyntax = fmt.Errorf("otel tracestate: invalid syntax")
)

// SamplingResult conveys a SamplingDecision, set of Attributes and a Tracestate.
type SamplingResult struct {
	Decision   SamplingDecision
	Attributes []attribute.KeyValue
	Tracestate trace.TraceState
}

type otelTraceState struct {
	random      int
	probability int
	unknown     []string
}

type probabilitySampler struct {
	traceIDUpperBound uint64
	description       string
}

func (ts probabilitySampler) ShouldSample(p SamplingParameters) SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	x := binary.BigEndian.Uint64(p.TraceID[0:8]) >> 1
	if x < ts.traceIDUpperBound {
		return SamplingResult{
			Decision:   RecordAndSample,
			Tracestate: psc.TraceState(),
		}
	}
	return SamplingResult{
		Decision:   Drop,
		Tracestate: psc.TraceState(),
	}
}

func (ts probabilitySampler) Description() string {
	return ts.description
}

// ProbabilityBased samples a given fraction of traces, supports
// arbitrary fractions.
// - Fractions >= 1 will always sample.
// - Fractions < 0 are treated as zero.
//
// Note: This Sampler implements the legacy behavior of
// TraceIDRatioSampler prior to standardizing how to propagate
// sampling probability.  This sampler does not guarantee consistent
// sampling when used with other ProbabilityBased implementations.
//
// To respect the parent trace's `SampledFlag`, the `ProbabilityBased`
// sampler should be used as a delegate of a `Parent` sampler.
func ProbabilityBased(fraction float64) Sampler {
	if fraction >= 1 {
		fraction = 1
	}

	if fraction <= 0 {
		fraction = 0
	}

	return &probabilitySampler{
		traceIDUpperBound: uint64(fraction * (1 << 63)),
		description:       fmt.Sprintf("ProbabilityBased{%g}", fraction),
	}
}

// AlwaysSample returns a Sampler that samples every trace.
// Be careful about using this sampler in a production application with
// significant traffic: a new trace will be started and exported for every
// request.
func AlwaysSample() Sampler {
	return TraceIDRatioBased(0)
}

// NeverSample returns a Sampler that samples no traces.
func NeverSample() Sampler {
	return TraceIDRatioBased(otelSamplingZeroValue)
}

// ParentBased returns a composite sampler which behaves differently,
// based on the parent of the span. If the span has no parent,
// the root(Sampler) is used to make sampling decision. If the span has
// a parent, depending on whether the parent is remote and whether it
// is sampled, one of the following samplers will apply:
// - remoteParentSampled(Sampler) (default: AlwaysOn)
// - remoteParentNotSampled(Sampler) (default: AlwaysOff)
// - localParentSampled(Sampler) (default: AlwaysOn)
// - localParentNotSampled(Sampler) (default: AlwaysOff)
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
		remoteParentSampled:    PropagateSampler(),
		remoteParentNotSampled: PropagateSampler(),
		localParentSampled:     PropagateSampler(),
		localParentNotSampled:  PropagateSampler(),
	}

	for _, so := range samplers {
		so.apply(&c)
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
	apply(*samplerConfig)
}

// WithRemoteParentSampled sets the sampler for the case of sampled remote parent.
func WithRemoteParentSampled(s Sampler) ParentBasedSamplerOption {
	return remoteParentSampledOption{s}
}

type remoteParentSampledOption struct {
	s Sampler
}

func (o remoteParentSampledOption) apply(config *samplerConfig) {
	config.remoteParentSampled = o.s
}

// WithRemoteParentNotSampled sets the sampler for the case of remote parent
// which is not sampled.
func WithRemoteParentNotSampled(s Sampler) ParentBasedSamplerOption {
	return remoteParentNotSampledOption{s}
}

type remoteParentNotSampledOption struct {
	s Sampler
}

func (o remoteParentNotSampledOption) apply(config *samplerConfig) {
	config.remoteParentNotSampled = o.s
}

// WithLocalParentSampled sets the sampler for the case of sampled local parent.
func WithLocalParentSampled(s Sampler) ParentBasedSamplerOption {
	return localParentSampledOption{s}
}

type localParentSampledOption struct {
	s Sampler
}

func (o localParentSampledOption) apply(config *samplerConfig) {
	config.localParentSampled = o.s
}

// WithLocalParentNotSampled sets the sampler for the case of local parent
// which is not sampled.
func WithLocalParentNotSampled(s Sampler) ParentBasedSamplerOption {
	return localParentNotSampledOption{s}
}

type localParentNotSampledOption struct {
	s Sampler
}

func (o localParentNotSampledOption) apply(config *samplerConfig) {
	config.localParentNotSampled = o.s
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

type propagateSampler struct{}

func PropagateSampler() Sampler {
	return propagateSampler{}
}

func (ps propagateSampler) ShouldSample(p SamplingParameters) SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)

	if !psc.IsSampled() {
		// For unsampled contexts we skip validating the OTel
		// tracestate key; `attrs` are unset because they will
		// not be recorded.
		return SamplingResult{
			Decision:   Drop,
			Tracestate: psc.TraceState(),
		}
	}

	otts, err := parseOTelTraceState(
		psc.TraceState().Get(otelTraceStateKey),
	)

	var attrs []attribute.KeyValue
	if err == nil && otts.hasProbability() {
		// We have known inclusion probability.  Set an
		// attribute when the count is not 1.
		if otts.probability != 0 {
			attrs = append(attrs,
				attribute.Int64(otelSamplingAdjustedCountKey, 1<<otts.probability),
			)
		}
	} else {
		// Set the sampler name to indicate an unknown
		// adjusted count.
		attrs = append(attrs,
			attribute.String(otelSamplingNameKey, otelSamplingParentSampler),
		)

		// Spec error handling behavior (TODO).
		if err != nil {
			otel.Handle(err)
		}
	}

	return SamplingResult{
		Decision:   RecordAndSample,
		Attributes: attrs,
		Tracestate: psc.TraceState(),
	}
}

func newOTelTraceState() otelTraceState {
	return otelTraceState{
		random:      -1,
		probability: -1,
	}
}

func parseOTelTraceState(ts string) (otelTraceState, error) {
	// TODO: Spec limit on trace state value length?
	otts := newOTelTraceState()
	for len(ts) > 0 {
		eqPos := 0
		for ; eqPos < len(ts); eqPos++ {
			if ts[eqPos] >= 'a' && ts[eqPos] <= 'z' {
				continue
			}
			break
		}
		if eqPos == 0 || eqPos == len(ts) || ts[eqPos] != ':' {
			return otts, errTraceStateSyntax
		}

		key := ts[0:eqPos]
		tail := ts[eqPos+1:]

		sepPos := 0

		if key == otelTraceStateProbabilityValueSubkey ||
			key == otelTraceStateRandomValueSubkey {

			for ; sepPos < len(tail); sepPos++ {
				if (tail[sepPos] >= '0' && tail[sepPos] <= '9') ||
					(tail[sepPos] >= 'a' && tail[sepPos] <= 'f') ||
					(tail[sepPos] >= 'A' && tail[sepPos] <= 'F') {
					continue
				}
				break
			}
			value, err := strconv.ParseUint(
				tail[0:sepPos],
				otelTraceStateNumberBase,
				otelTraceStateBitWidth,
			)
			if err != nil {
				return otts, fmt.Errorf("%w: %s", errTraceStateSyntax, err)
			}
			if key == otelTraceStateProbabilityValueSubkey {
				otts.probability = int(value)
			} else if key == otelTraceStateRandomValueSubkey {
				otts.random = int(value)
			}

		} else {
			// TODO: Spec valid character set for forward compatibility.
			// Should this be the base64 characters?
			for ; sepPos < len(tail); sepPos++ {
				if (tail[sepPos] >= '0' && tail[sepPos] <= '9') ||
					(tail[sepPos] >= 'a' && tail[sepPos] <= 'z') ||
					(tail[sepPos] >= 'A' && tail[sepPos] <= 'Z') {
					continue
				}
			}
			otts.unknown = append(otts.unknown, ts[0:sepPos])
		}

		if sepPos == 0 || (sepPos < len(tail) && tail[sepPos] != ';') {
			return otts, errTraceStateSyntax
		}

		if sepPos == len(tail) {
			break
		}

		ts = tail[sepPos+1:]
	}

	return otts, nil
}

func (ps propagateSampler) Description() string {
	return "Propagate"
}

func (otts otelTraceState) hasProbability() bool {
	return otts.probability >= 0
}

func (otts otelTraceState) hasRandom() bool {
	return otts.random >= 0
}

type TraceIDRatioBasedRandomSource func() int

type traceIDRatioBasedConfig struct {
	source TraceIDRatioBasedRandomSource
}

type TraceIDRatioBasedOption interface {
	apply(*traceIDRatioBasedConfig)
}

type traceIDRatioBasedRandomSource TraceIDRatioBasedRandomSource

func WithRandomSource(source TraceIDRatioBasedRandomSource) TraceIDRatioBasedOption {
	return traceIDRatioBasedRandomSource(source)
}

func (s traceIDRatioBasedRandomSource) apply(cfg *traceIDRatioBasedConfig) {
	cfg.source = TraceIDRatioBasedRandomSource(s)
}

func TraceIDRatioBased(logAdjCnt int, opts ...TraceIDRatioBasedOption) Sampler {
	cfg := traceIDRatioBasedConfig{
		source: func() int {
			// TODO: Optimize me; This wastes 61 bits of
			// randomness on average.
			var x int64

			for x = rand.Int63(); x == 0; {
			}

			cnt := 0
			for (x & 1) == 0 {
				cnt++
				x >>= 1
			}
			return cnt
		},
	}
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	if logAdjCnt < 0 || logAdjCnt > otelSamplingZeroValue {
		// Zero probability
		logAdjCnt = otelSamplingZeroValue
	}
	return traceIDRatioSampler{
		logAdjCnt: logAdjCnt,
		source:    cfg.source,
	}
}

type traceIDRatioSampler struct {
	logAdjCnt int
	source    TraceIDRatioBasedRandomSource
}

func (t traceIDRatioSampler) ShouldSample(p SamplingParameters) SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	var otts otelTraceState
	var state trace.TraceState

	if !psc.IsValid() {
		// A new root is happening.  Compute `random`.
		otts = newOTelTraceState()
		otts.random = t.source()
	} else {
		// A valid parent context.
		// It does not matter if psc.IsSampled().
		state = psc.TraceState()

		var err error
		otts, err = parseOTelTraceState(state.Get(otelTraceStateKey))

		if err != nil {
			// TODO: Spec error treatment.
		}
		// TODO: Spec behavior when `random` is set badly or missing.
		// e.g., what if random == otelSamplingZeroValue?
	}

	var decision SamplingDecision
	var cnt int64

	// Calculate the adjusted count.  Treat the max-value bucket
	// as zero probability, thus adjusted count zero.
	if t.logAdjCnt < otelSamplingZeroValue {
		cnt = 1 << t.logAdjCnt
	}

	if cnt != 0 && t.logAdjCnt <= otts.random {
		decision = RecordAndSample
	} else {
		decision = Drop
	}

	otts.probability = t.logAdjCnt

	var attrs []attribute.KeyValue

	if cnt != 1 {
		attrs = append(attrs, attribute.Int64(otelSamplingAdjustedCountKey, cnt))
	}

	state, err := state.Insert(otelTraceStateKey, otts.serialize())
	if err != nil {
		// TODO: Spec error treatment
	}

	return SamplingResult{
		Decision:   decision,
		Attributes: attrs,
		Tracestate: state,
	}
}

func (ts traceIDRatioSampler) Description() string {
	return fmt.Sprintf("TraceIDRatioBased{%g}", math.Pow(2, float64(-ts.logAdjCnt)))
}

func (otts otelTraceState) serialize() string {
	var sb strings.Builder
	if otts.hasProbability() {
		_, _ = sb.WriteString(fmt.Sprintf("p:%02x;", otts.probability))
	}
	if otts.hasRandom() {
		_, _ = sb.WriteString(fmt.Sprintf("r:%02x;", otts.random))
	}
	for _, unk := range otts.unknown {
		_, _ = sb.WriteString(unk)
		_, _ = sb.WriteString(";")
	}
	return sb.String()
}
