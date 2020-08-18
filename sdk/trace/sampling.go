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

type probabilitySampler struct {
	traceIDUpperBound uint64
	description       string
}

func (ps probabilitySampler) ShouldSample(p SamplingParameters) SamplingResult {
	if p.ParentContext.IsSampled() {
		return SamplingResult{Decision: RecordAndSampled}
	}

	x := binary.BigEndian.Uint64(p.TraceID[0:8]) >> 1
	if x < ps.traceIDUpperBound {
		return SamplingResult{Decision: RecordAndSampled}
	}
	return SamplingResult{Decision: NotRecord}
}

func (ps probabilitySampler) Description() string {
	return ps.description
}

// ProbabilitySampler samples a given fraction of traces. Fractions >= 1 will
// always sample. If the parent span is sampled, then it's child spans will
// automatically be sampled. Fractions < 0 are treated as zero, but spans may
// still be sampled if their parent is.
func ProbabilitySampler(fraction float64) Sampler {
	if fraction >= 1 {
		return AlwaysSample()
	}

	if fraction <= 0 {
		fraction = 0
	}

	return &probabilitySampler{
		traceIDUpperBound: uint64(fraction * (1 << 63)),
		description:       fmt.Sprintf("ProbabilitySampler{%g}", fraction),
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

// ParentSample returns a Sampler that samples a trace only
// if the the span has a parent span and it is sampled. If the span has
// parent span but it is not sampled, neither will this span. If the span
// does not have a parent the fallback Sampler is used to determine if the
// span should be sampled.
func ParentSample(fallback Sampler) Sampler {
	return parentSampler{fallback}
}

type parentSampler struct {
	fallback Sampler
}

func (ps parentSampler) ShouldSample(p SamplingParameters) SamplingResult {
	if p.ParentContext.IsValid() {
		if p.ParentContext.IsSampled() {
			return SamplingResult{Decision: RecordAndSampled}
		}
		return SamplingResult{Decision: NotRecord}
	}
	return ps.fallback.ShouldSample(p)
}

func (ps parentSampler) Description() string {
	return fmt.Sprintf("ParentOrElse{%s}", ps.fallback.Description())
}
