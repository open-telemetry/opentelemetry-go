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

package sampling_test // import "go.opentelemetry.io/otel/sdk/trace/sampling_test"

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/trace/sampling"
	"go.opentelemetry.io/otel/trace"
)

func TestProbabilisticSampler_IsSampled(t *testing.T) {
	tests := []struct {
		name         string
		samplingRate float64
		idString     string
		want         bool
	}{
		{
			name:         "test new probabilistic sampler, rate 0, isSampled == false",
			samplingRate: 0,
			idString:     "e44efd7c19b5445988f6799233bb4008",
			want:         false,
		},
		{
			name:         "test new probabilistic sampler, rate 0.5, isSampled == true",
			samplingRate: 0.5,
			idString:     "3be9447364324b86a20cdf6de12c61be", // == 4317056980204276614
			want:         true,
		},
		{
			name:         "test new probabilistic sampler, rate 0.5, isSampled == false",
			samplingRate: 0.5,
			idString:     "4ed01488b9804ebeb9562f067d16c347", // == 5679061707574496958
			want:         false,
		},
		{
			name:         "test new probabilistic sampler, rate 1, isSampled == true",
			samplingRate: 1,
			idString:     "f86372bf27ba47a0903b267ee5535a5a", // == 8674903472576546720
			want:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler := sampling.NewProbabilisticSampler(tt.samplingRate)
			traceID, err := trace.TraceIDFromHex(tt.idString)
			assert.NoError(t, err)
			got, _ := sampler.IsSampled(traceID)
			if got != tt.want {
				t.Errorf("IsSampled() got = %v, want %v", got, tt.want)
			}
		})
	}
}
