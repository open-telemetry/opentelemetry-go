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

// ProbabilisticSamplingStrategy defines probabilistic strategy option.
type ProbabilisticSamplingStrategy struct {
	SamplingRate float64 `json:"samplingRate"`
}

// RateLimitingSamplingStrategy defines rate limiting strategy option.
type RateLimitingSamplingStrategy struct {
	MaxTracesPerSecond float64 `json:"maxTracesPerSecond"`
}

// OperationSamplingStrategy defines operation strategy options.
type OperationSamplingStrategy struct {
	Operation             string                         `json:"operation"`
	ProbabilisticSampling *ProbabilisticSamplingStrategy `json:"probabilisticSampling"`
}

// PerOperationSamplingStrategies defines per operation strategies options.
type PerOperationSamplingStrategies struct {
	DefaultSamplingProbability       float64                      `json:"defaultSamplingProbability"`
	DefaultLowerBoundTracesPerSecond float64                      `json:"defaultLowerBoundTracesPerSecond"`
	PerOperationStrategies           []*OperationSamplingStrategy `json:"perOperationStrategies"`
	DefaultUpperBoundTracesPerSecond *float64                     `json:"defaultUpperBoundTracesPerSecond,omitempty"`
}

// ServiceStrategy defines a service specific sampling Strategy.
type ServiceStrategy struct {
	StrategyType          int                             `json:"strategyType"`
	ProbabilisticSampling *ProbabilisticSamplingStrategy  `json:"probabilisticSampling,omitempty"`
	RateLimitingSampling  *RateLimitingSamplingStrategy   `json:"rateLimitingSampling,omitempty"`
	OperationSampling     *PerOperationSamplingStrategies `json:"operationSampling,omitempty"`
}
