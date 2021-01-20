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
	"github.com/uber/jaeger-client-go/utils"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

// SamplerTypeRateLimiting is a RateLimitingSampler default parameter.
const SamplerTypeRateLimiting string = "ratelimiting"

const defaultItemCost = float64(1)

// RateLimitingSampler samples at most maxTracesPerSecond. The distribution of sampled traces follows
// burstiness of the service, i.e. a service with uniformly distributed requests will have those
// requests sampled uniformly as well, but if requests are bursty, especially sub-second, then a
// number of sequential requests can be sampled each second.
type RateLimitingSampler struct {
	maxTracesPerSecond float64
	rateLimiter        *utils.ReconfigurableRateLimiter
	tags               []label.KeyValue
}

// NewRateLimitingSampler creates new RateLimitingSampler.
func NewRateLimitingSampler(maxTracesPerSecond float64) *RateLimitingSampler {
	rateLimiter := utils.NewRateLimiter(maxTracesPerSecond, maxTracesPerSecond)

	return &RateLimitingSampler{
		maxTracesPerSecond: maxTracesPerSecond,
		rateLimiter:        rateLimiter,
		tags: []label.KeyValue{
			label.String(SamplerTypeTagKey, SamplerTypeRateLimiting),
			label.Float64(SamplerParamTagKey, maxTracesPerSecond),
		},
	}
}

// IsSampled makes sampling decision.
// ReconfigurableRateLimiter CheckCredit function is based on leaky bucket algorithm.
func (sampler *RateLimitingSampler) IsSampled(id trace.TraceID) (bool, []label.KeyValue) {
	return sampler.rateLimiter.CheckCredit(defaultItemCost), sampler.tags
}
