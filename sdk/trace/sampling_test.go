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
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
)

func TestParentBasedDefaultLocalParentSampled(t *testing.T) {
	sampler := ParentBased(AlwaysSample())
	traceID, _ := otel.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := otel.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := otel.SpanContext{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: otel.FlagsSampled,
	}
	if sampler.ShouldSample(SamplingParameters{ParentContext: parentCtx}).Decision != RecordAndSample {
		t.Error("Sampling decision should be RecordAndSample")
	}
}

func TestParentBasedDefaultLocalParentNotSampled(t *testing.T) {
	sampler := ParentBased(AlwaysSample())
	traceID, _ := otel.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := otel.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := otel.SpanContext{
		TraceID: traceID,
		SpanID:  spanID,
	}
	if sampler.ShouldSample(SamplingParameters{ParentContext: parentCtx}).Decision != Drop {
		t.Error("Sampling decision should be Drop")
	}
}

func TestParentBasedWithNoParent(t *testing.T) {
	params := SamplingParameters{}

	sampler := ParentBased(AlwaysSample())
	if sampler.ShouldSample(params).Decision != RecordAndSample {
		t.Error("Sampling decision should be RecordAndSample")
	}

	sampler = ParentBased(NeverSample())
	if sampler.ShouldSample(params).Decision != Drop {
		t.Error("Sampling decision should be Drop")
	}
}

func TestParentBasedWithSamplerOptions(t *testing.T) {
	testCases := []struct {
		name                            string
		samplerOption                   ParentBasedSamplerOption
		isParentRemote, isParentSampled bool
		expectedDecision                SamplingDecision
	}{
		{
			"localParentSampled",
			WithLocalParentSampled(NeverSample()),
			false,
			true,
			Drop,
		},
		{
			"localParentNotSampled",
			WithLocalParentNotSampled(AlwaysSample()),
			false,
			false,
			RecordAndSample,
		},
		{
			"remoteParentSampled",
			WithRemoteParentSampled(NeverSample()),
			true,
			true,
			Drop,
		},
		{
			"remoteParentNotSampled",
			WithRemoteParentNotSampled(AlwaysSample()),
			true,
			false,
			RecordAndSample,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			traceID, _ := otel.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
			spanID, _ := otel.SpanIDFromHex("00f067aa0ba902b7")
			parentCtx := otel.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			}

			if tc.isParentSampled {
				parentCtx.TraceFlags = otel.FlagsSampled
			}

			params := SamplingParameters{ParentContext: parentCtx}
			if tc.isParentRemote {
				params.HasRemoteParent = true
			}

			sampler := ParentBased(
				nil,
				tc.samplerOption,
			)

			switch tc.expectedDecision {
			case RecordAndSample:
				if sampler.ShouldSample(params).Decision != tc.expectedDecision {
					t.Error("Sampling decision should be RecordAndSample")
				}
			case Drop:
				if sampler.ShouldSample(params).Decision != tc.expectedDecision {
					t.Error("Sampling decision should be Drop")
				}
			}
		})
	}
}

func TestParentBasedDefaultDescription(t *testing.T) {
	sampler := ParentBased(AlwaysSample())

	expectedDescription := fmt.Sprintf("ParentBased{root:%s,remoteParentSampled:%s,"+
		"remoteParentNotSampled:%s,localParentSampled:%s,localParentNotSampled:%s}",
		AlwaysSample().Description(),
		AlwaysSample().Description(),
		NeverSample().Description(),
		AlwaysSample().Description(),
		NeverSample().Description())

	if sampler.Description() != expectedDescription {
		t.Error(fmt.Sprintf("Sampler description should be %s, got '%s' instead",
			expectedDescription,
			sampler.Description(),
		))
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
			if samplerLo.ShouldSample(params).Decision == RecordAndSample {
				require.Equal(t, RecordAndSample, samplerHi.ShouldSample(params).Decision,
					"%s sampled but %s did not", samplerLo.Description(), samplerHi.Description())
			}
		}
	}
}
