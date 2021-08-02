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
	"strconv"

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

// Probability sampling constants
const (
	otelTraceStateKey                    = "otel"
	otelTraceStateProbabilityValueSubkey = "p"
	otelSamplingAttributePrefix          = "sampler."
	otelSamplingAdjustedCountKey         = otelSamplingAttributePrefix + "adjusted_count"
	otelSamplingNameKey                  = otelSamplingAttributePrefix + "name"
	otelSamplingParentSampler            = "parent"
)

// Probability sampling errors
var (
	errTraceStateSyntax   = fmt.Errorf("otel tracestate: invalid syntax")
	errTraceStateNotFound = fmt.Errorf("otel tracestate: subkey not found")
)

// SamplingResult conveys a SamplingDecision, set of Attributes and a Tracestate.
type SamplingResult struct {
	Decision   SamplingDecision
	Attributes []attribute.KeyValue
	Tracestate trace.TraceState
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
		return AlwaysSample()
	}

	if fraction <= 0 {
		return NeverSample()
	}

	return &probabilitySampler{
		traceIDUpperBound: uint64(fraction * (1 << 63)),
		description:       fmt.Sprintf("ProbabilityBased{%g}", fraction),
	}
}

type alwaysOnSampler struct{}

func (as alwaysOnSampler) ShouldSample(p SamplingParameters) SamplingResult {
	return SamplingResult{
		Decision:   RecordAndSample,
		Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
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
		remoteParentSampled:    propagateSampler{},
		remoteParentNotSampled: propagateSampler{},
		localParentSampled:     propagateSampler{},
		localParentNotSampled:  propagateSampler{},
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

var _ Sampler = propagateSampler{}

func PropagateSampler() Sampler {
	return propagateSampler{}
}

func (ps propagateSampler) ShouldSample(p SamplingParameters) SamplingResult {
	var attrs []attribute.KeyValue
	var decision SamplingDecision

	psc := trace.SpanContextFromContext(p.ParentContext)

	if !psc.IsSampled() {
		// For unsampled contexts we skip validating the OTel
		// tracestate key; `attrs` are unset because they will
		// not be recorded.
		decision = Drop
	} else {
		decision = RecordAndSample

		probVal, err := parseTraceStateInt(
			psc.TraceState().Get(otelTraceStateKey),
			otelTraceStateProbabilityValueSubkey,
		)

		if err == nil {
			// We have known inclusion probability.  Set an
			// attribute when the count is not 1.
			if probVal != 0 {
				attrs = append(attrs,
					attribute.Int64(otelSamplingAdjustedCountKey, 1<<probVal),
				)
			}
		} else {
			// Set the sampler name to indicate an unknown
			// adjusted count.
			attrs = append(attrs,
				attribute.String(otelSamplingNameKey, otelSamplingParentSampler),
			)

			// Spec error handling behavior (TODO).
			if err != errTraceStateNotFound {
				otel.Handle(err)
			}
		}
	}

	return SamplingResult{
		Decision:   decision,
		Attributes: attrs,
		Tracestate: psc.TraceState(),
	}
}

func parseTraceStateInt(ts, key string) (int, error) {
	for {
		if len(ts) == 0 {
			return 0, errTraceStateNotFound
		}
		eqPos := 0
		for ; eqPos < len(ts); eqPos++ {
			if ts[eqPos] >= 'a' && ts[eqPos] <= 'z' {
				continue
			}
			break
		}
		if eqPos == 0 || eqPos == len(ts) || ts[eqPos] != ':' {
			return 0, errTraceStateSyntax
		}

		isMatch := key == ts[0:eqPos]
		ts = ts[eqPos+1:]

		sepPos := 0
		for ; sepPos < len(ts); sepPos++ {
			if ts[sepPos] >= '0' && ts[sepPos] <= '9' {
				continue
			}
			if ts[sepPos] >= 'a' && ts[sepPos] <= 'f' {
				continue
			}
			break
		}
		value, err := strconv.ParseUint(ts[0:sepPos], 16, 32)
		if err != nil {
			return 0, err
		}

		if sepPos == 0 || (sepPos < len(ts) && ts[sepPos] != ';') {
			return 0, errTraceStateSyntax
		}
		if !isMatch {
			if sepPos == len(ts) {
				return 0, errTraceStateNotFound
			}
			ts = ts[sepPos+1:]
			continue
		}

		return int(value), nil
	}
}

func (ps propagateSampler) Description() string {
	return "Propagate"
}
