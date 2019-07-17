// Copyright 2019, OpenTelemetry Authors
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
	"go.opentelemetry.io/api/core"
)

const defaultSamplingProbability = 1e-4

// Sampler decides whether a trace should be sampled and exported.
type Sampler func(SamplingParameters) SamplingDecision

// SamplingParameters contains the values passed to a Sampler.
type SamplingParameters struct {
	ParentContext   core.SpanContext
	TraceID         core.TraceID
	SpanID          uint64
	Name            string
	HasRemoteParent bool
}

// SamplingDecision is the value returned by a Sampler.
type SamplingDecision struct {
	Sample bool
}

// ProbabilitySampler returns a Sampler that samples a given fraction of traces.
//
// It also samples spans whose parents are sampled.
func ProbabilitySampler(fraction float64) Sampler {
	if !(fraction >= 0) {
		fraction = 0
	} else if fraction >= 1 {
		return AlwaysSample()
	}

	traceIDUpperBound := uint64(fraction * (1 << 63))
	return Sampler(func(p SamplingParameters) SamplingDecision {
		if p.ParentContext.IsSampled() {
			return SamplingDecision{Sample: true}
		}
		x := p.TraceID.High >> 1
		return SamplingDecision{Sample: x < traceIDUpperBound}
	})
}

// AlwaysSample returns a Sampler that samples every trace.
// Be careful about using this sampler in a production application with
// significant traffic: a new trace will be started and exported for every
// request.
func AlwaysSample() Sampler {
	return func(p SamplingParameters) SamplingDecision {
		return SamplingDecision{Sample: true}
	}
}

// NeverSample returns a Sampler that samples no traces.
func NeverSample() Sampler {
	return func(p SamplingParameters) SamplingDecision {
		return SamplingDecision{Sample: false}
	}
}
