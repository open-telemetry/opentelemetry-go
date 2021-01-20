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

package sampling // import "go.opentelemetry.io/otel/sdk/trace/sampling"

import (
	"encoding/binary"
	"math"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

const (
	maxRandomNumber = ^(uint64(1) << 63)
	maxRate         = float64(1)
	minRate         = float64(0)
)

// ProbabilisticSampler default parameters.
const (
	SamplerTypeTagKey        string = "sampler.type"
	SamplerTypeProbabilistic string = "probabilistic"
	SamplerParamTagKey       string = "sampler.param"
)

// ProbabilisticSampler is a sampler that randomly samples a certain percentage of traces.
type ProbabilisticSampler struct {
	// samplingRate     float64
	samplingBoundary uint64
	tags             []label.KeyValue
}

// NewProbabilisticSampler constructs new ProbabilisticSampler with sampling rate given.
// Equals samplingRate > 1 to 1, samplingRate < 0 to 0.
func NewProbabilisticSampler(samplingRate float64) *ProbabilisticSampler {
	samplingRate = math.Max(minRate, math.Min(samplingRate, maxRate))

	return &ProbabilisticSampler{
		samplingBoundary: uint64(float64(maxRandomNumber) * samplingRate),
		tags: []label.KeyValue{
			label.String(SamplerTypeTagKey, SamplerTypeProbabilistic),
			label.Float64(SamplerParamTagKey, samplingRate),
		},
	}
}

// IsSampled makes sampling decision considering the fact that trace IDs are 63bit random numbers themselves,
// thus simply calculating if traceID < (samplingRate * 2^63).
func (sampler *ProbabilisticSampler) IsSampled(id trace.TraceID) (bool, []label.KeyValue) {
	traceIDInt := binary.BigEndian.Uint64(id[0:8])

	return sampler.samplingBoundary >= traceIDInt&maxRandomNumber, sampler.tags
}
