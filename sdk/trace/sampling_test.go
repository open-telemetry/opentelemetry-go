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
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	api "go.opentelemetry.io/otel/api/trace"
)

func TestAlwaysParentSampleWithParentSampled(t *testing.T) {
	sampler := ParentSample(AlwaysSample())
	traceID, _ := api.IDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := api.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := api.SpanContext{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: api.FlagsSampled,
	}
	if sampler.ShouldSample(SamplingParameters{ParentContext: parentCtx}).Decision != RecordAndSampled {
		t.Error("Sampling decision should be RecordAndSampled")
	}
}

func TestNeverParentSampleWithParentSampled(t *testing.T) {
	sampler := ParentSample(NeverSample())
	traceID, _ := api.IDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := api.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := api.SpanContext{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: api.FlagsSampled,
	}
	if sampler.ShouldSample(SamplingParameters{ParentContext: parentCtx}).Decision != RecordAndSampled {
		t.Error("Sampling decision should be RecordAndSampled")
	}
}

func TestAlwaysParentSampleWithParentNotSampled(t *testing.T) {
	sampler := ParentSample(AlwaysSample())
	traceID, _ := api.IDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := api.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := api.SpanContext{
		TraceID: traceID,
		SpanID:  spanID,
	}
	if sampler.ShouldSample(SamplingParameters{ParentContext: parentCtx}).Decision != NotRecord {
		t.Error("Sampling decision should be NotRecord")
	}
}

func TestParentSampleWithNoParent(t *testing.T) {
	params := SamplingParameters{}

	sampler := ParentSample(AlwaysSample())
	if sampler.ShouldSample(params).Decision != RecordAndSampled {
		t.Error("Sampling decision should be RecordAndSampled")
	}

	sampler = ParentSample(NeverSample())
	if sampler.ShouldSample(params).Decision != NotRecord {
		t.Error("Sampling decision should be NotRecord")
	}
}

// TraceIDRatioBased sampler requirements state
//  "A TraceIDRatioBased sampler with a given sampling rate MUST also sample
//   all traces that any TraceIDRatioBased sampler with a lower sampling rate
//   would sample."
func TestTraceIdRatioSamplesInclusively(t *testing.T) {
	const (
		numSamplers = 1000
		numTraces   = 100
	)
	idg := defIDGenerator()

	for i := 0; i < numSamplers; i++ {
		ratioLo, ratioHi := rand.Float64(), rand.Float64()
		if ratioHi < ratioLo {
			ratioLo, ratioHi = ratioHi, ratioLo
		}
		samplerHi := TraceIDRatioBased(ratioHi)
		samplerLo := TraceIDRatioBased(ratioLo)
		for j := 0; j < numTraces; j++ {
			traceID := idg.NewTraceID()

			params := SamplingParameters{TraceID: traceID}
			if samplerLo.ShouldSample(params).Decision == RecordAndSampled {
				require.Equal(t, RecordAndSampled, samplerHi.ShouldSample(params).Decision,
					"%s sampled but %s did not", samplerLo.Description(), samplerHi.Description())
			}
		}
	}
}
