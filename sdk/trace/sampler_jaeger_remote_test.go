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
	"net/http"
	"net/http/httptest"
	"testing"

	jaeger_api_v2 "github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"github.com/stretchr/testify/assert"
)

func TestJaegerRemoteSampler_updateSamplingStrategies(t *testing.T) {
	sampler := JaegerRemoteSampler("foo", "http://localhost:5778")
	jaegerRemoteSampler := sampler.(*jaegerRemoteSampler)

	defaultSampler := TraceIDRatioBased(jaegerRemoteDefaultSamplerFraction)
	assert.Equal(t, defaultSampler, jaegerRemoteSampler.sampler)

	tests := []struct {
		name        string
		strategy    jaeger_api_v2.SamplingStrategyResponse
		expectedErr string
		sampler     Sampler
	}{
		{
			name:     "update strategy without changes",
			strategy: jaeger_api_v2.SamplingStrategyResponse{},
			sampler:  defaultSampler,
		},
		{
			name: "update strategy with PROBABILISTIC and sampling rate 0.8",
			strategy: jaeger_api_v2.SamplingStrategyResponse{
				StrategyType: jaeger_api_v2.SamplingStrategyType_PROBABILISTIC,
				ProbabilisticSampling: &jaeger_api_v2.ProbabilisticSamplingStrategy{
					SamplingRate: 0.8,
				},
			},
			sampler: TraceIDRatioBased(0.8),
		},
		{
			name: "update strategy with RATE_LIMITING",
			strategy: jaeger_api_v2.SamplingStrategyResponse{
				StrategyType: jaeger_api_v2.SamplingStrategyType_RATE_LIMITING,
				RateLimitingSampling: &jaeger_api_v2.RateLimitingSamplingStrategy{
					MaxTracesPerSecond: 100,
				},
			},
			expectedErr: "loading failed: only strategy type PROBABILISTC is supported, got RATE_LIMITING",
			sampler:     TraceIDRatioBased(0.8),
		},
		{
			name: "update strategy with per operation sampling",
			strategy: jaeger_api_v2.SamplingStrategyResponse{
				StrategyType: jaeger_api_v2.SamplingStrategyType_PROBABILISTIC,
				ProbabilisticSampling: &jaeger_api_v2.ProbabilisticSamplingStrategy{
					SamplingRate: 1,
				},
				OperationSampling: &jaeger_api_v2.PerOperationSamplingStrategies{
					DefaultSamplingProbability: 1,
				},
			},
			expectedErr: "loading failed: per operation sampling is not supported",
			sampler:     TraceIDRatioBased(0.8),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jaegerRemoteSampler.fetcher = mockStrategyFetcher{
				response: tt.strategy,
			}

			err := jaegerRemoteSampler.updateSamplingStrategies()
			// TODO this feels awkward
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.sampler, jaegerRemoteSampler.sampler)
		})
	}
}

type mockStrategyFetcher struct {
	response jaeger_api_v2.SamplingStrategyResponse
	err      error
}

func (m mockStrategyFetcher) Fetch() (jaeger_api_v2.SamplingStrategyResponse, error) {
	return m.response, m.err
}

func Test_jaegerStrategiesFetcher_Fetch(t *testing.T) {
	tests := []struct {
		name               string
		responseStatusCode int
		responseBody       string
		expectedErr        string
		expectedStrategy   jaeger_api_v2.SamplingStrategyResponse
	}{
		{
			name:               "RequestOK",
			responseStatusCode: http.StatusOK,
			responseBody: `{
  "strategyType": 0,
  "probabilisticSampling": {
    "samplingRate": 0.5
  }
}`,
			expectedStrategy: jaeger_api_v2.SamplingStrategyResponse{
				StrategyType: jaeger_api_v2.SamplingStrategyType_PROBABILISTIC,
				ProbabilisticSampling: &jaeger_api_v2.ProbabilisticSamplingStrategy{
					SamplingRate: 0.5,
				},
			},
		},
		{
			name:               "RequestError",
			responseStatusCode: http.StatusTooManyRequests,
			responseBody:       "you are sending too many requests",
			expectedErr:        "request failed (429): you are sending too many requests",
		},
		{
			name:               "InvalidResponseData",
			responseStatusCode: http.StatusOK,
			responseBody:       `{"strategy`,
			expectedErr:        "unexpected end of JSON input",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/?service=foo", r.URL.RequestURI())

				w.WriteHeader(tt.responseStatusCode)
				_, err := w.Write([]byte(tt.responseBody))
				assert.NoError(t, err)
			}))
			defer server.Close()

			fetcher := jaegerSamplingStrategyFetcherImpl{
				serviceName:       "foo",
				samplingServerUrl: server.URL,
			}

			strategyResponse, err := fetcher.Fetch()
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStrategy, strategyResponse)
			}
		})
	}
}
