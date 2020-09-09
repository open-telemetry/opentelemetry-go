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
	"encoding/binary"
	"fmt"

	api "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

// Sampler decides whether a trace should be sampled and exported.
type Sampler interface {
	ShouldSample(SamplingParameters) SamplingResult
	Description() string
}

// SamplingParameters contains the values passed to a Sampler.
type SamplingParameters struct {
	ParentContext   api.SpanContext
	TraceID         api.ID
	Name            string
	HasRemoteParent bool
	Kind            api.SpanKind
	Attributes      []label.KeyValue
	Links           []api.Link
}

// SamplingDecision indicates whether a span is recorded and sampled.
type SamplingDecision uint8

// Valid sampling decisions
const (
	NotRecord SamplingDecision = iota
	Record
	RecordAndSampled
)

// SamplingResult conveys a SamplingDecision and a set of Attributes.
type SamplingResult struct {
	Decision   SamplingDecision
	Attributes []label.KeyValue
}

type traceIDRatioSampler struct {
	traceIDUpperBound uint64
	description       string
}

func (ts traceIDRatioSampler) ShouldSample(p SamplingParameters) SamplingResult {
	x := binary.BigEndian.Uint64(p.TraceID[0:8]) >> 1
	if x < ts.traceIDUpperBound {
		return SamplingResult{Decision: RecordAndSampled}
	}
	return SamplingResult{Decision: NotRecord}
}

func (ts traceIDRatioSampler) Description() string {
	return ts.description
}

// TraceIDRatioBased samples a given fraction of traces. Fractions >= 1 will
// always sample. Fractions < 0 are treated as zero. To respect the
// parent trace's `SampledFlag`, the `TraceIDRatioBased` sampler should be used
// as a delegate of a `Parent` sampler.
//nolint:golint // golint complains about stutter of `trace.TraceIDRatioBased`
func TraceIDRatioBased(fraction float64) Sampler {
	if fraction >= 1 {
		return AlwaysSample()
	}

	if fraction <= 0 {
		fraction = 0
	}

	return &traceIDRatioSampler{
		traceIDUpperBound: uint64(fraction * (1 << 63)),
		description:       fmt.Sprintf("TraceIDRatioBased{%g}", fraction),
	}
}

type alwaysOnSampler struct{}

func (as alwaysOnSampler) ShouldSample(p SamplingParameters) SamplingResult {
	return SamplingResult{Decision: RecordAndSampled}
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
	return SamplingResult{Decision: NotRecord}
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
func ParentBased(root Sampler, samplers ...SamplerOption) Sampler {
	return parentBased{
		root:     root,
		samplers: configureSamplersForParentBased(samplers),
	}
}

type parentBased struct {
	root     Sampler
	samplers parentBasedSamplers
}

type SamplerOption func(*parentBasedSamplers)

type parentBasedSamplers struct {
	remoteParentSampled, remoteParentNotSampled Sampler
	localParentSampled, localParentNotSampled   Sampler
}

func configureSamplersForParentBased(samplers []SamplerOption) parentBasedSamplers {
	o := parentBasedSamplers{}
	for _, s := range samplers {
		s(&o)
	}

	if o.remoteParentSampled == nil {
		o.remoteParentSampled = AlwaysSample()
	}
	if o.remoteParentNotSampled == nil {
		o.remoteParentNotSampled = NeverSample()
	}
	if o.localParentSampled == nil {
		o.localParentSampled = AlwaysSample()
	}
	if o.localParentNotSampled == nil {
		o.localParentNotSampled = NeverSample()
	}

	return o
}

func WithRemoteParentSampled(s Sampler) SamplerOption {
	return func(o *parentBasedSamplers) {
		o.remoteParentSampled = s
	}
}
func WithRemoteParentNotSampled(s Sampler) SamplerOption {
	return func(o *parentBasedSamplers) {
		o.remoteParentNotSampled = s
	}
}
func WithLocalParentSampled(s Sampler) SamplerOption {
	return func(o *parentBasedSamplers) {
		o.localParentSampled = s
	}
}
func WithLocalParentNotSampled(s Sampler) SamplerOption {
	return func(o *parentBasedSamplers) {
		o.localParentNotSampled = s
	}
}

func (pb parentBased) ShouldSample(p SamplingParameters) SamplingResult {
	if p.ParentContext.IsValid() {
		if p.HasRemoteParent {
			if p.ParentContext.IsSampled() {
				return pb.samplers.remoteParentSampled.ShouldSample(p)
			}
			return pb.samplers.remoteParentNotSampled.ShouldSample(p)
		}

		if p.ParentContext.IsSampled() {
			return pb.samplers.localParentSampled.ShouldSample(p)
		}
		return pb.samplers.localParentNotSampled.ShouldSample(p)
	}
	return pb.root.ShouldSample(p)
}

func (pb parentBased) Description() string {
	return fmt.Sprintf("ParentBased{root:%s,remoteParentSampled:%s,"+
		"remoteParentNotSampled:%s,localParentSampled:%s,localParentNotSampled:%s}",
		pb.root.Description(),
		pb.samplers.remoteParentSampled.Description(),
		pb.samplers.remoteParentNotSampled.Description(),
		pb.samplers.localParentSampled.Description(),
		pb.samplers.localParentNotSampled.Description(),
	)
}
