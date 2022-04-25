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

package metrictest // import "go.opentelemetry.io/otel/sdk/metric/metrictest"

import "go.opentelemetry.io/otel/sdk/metric/export/aggregation"

type config struct {
	temporalitySelector aggregation.TemporalitySelector
}

func newConfig(opts ...Option) config {
	cfg := config{
		temporalitySelector: aggregation.CumulativeTemporalitySelector(),
	}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	return cfg
}

// Option allow for control of details of the TestMeterProvider created.
type Option interface {
	apply(config) config
}

type functionOption func(config) config

func (f functionOption) apply(cfg config) config {
	return f(cfg)
}

// WithTemporalitySelector allows for the use of either cumulative (default) or
// delta metrics.
//
// Warning: the current SDK does not convert async instruments into delta
// temporality.
func WithTemporalitySelector(ts aggregation.TemporalitySelector) Option {
	return functionOption(func(cfg config) config {
		if ts == nil {
			return cfg
		}
		cfg.temporalitySelector = ts
		return cfg
	})
}
