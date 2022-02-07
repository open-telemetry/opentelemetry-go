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

package otlpmetric // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric"

import "go.opentelemetry.io/otel/sdk/metric/export/aggregation"

// Option are setting options passed to an Exporter on creation.
type Option interface {
	apply(config) config
}

type exporterOptionFunc func(config) config

func (fn exporterOptionFunc) apply(cfg config) config {
	return fn(cfg)
}

type config struct {
	temporalitySelector aggregation.TemporalitySelector
}

// WithMetricAggregationTemporalitySelector defines the aggregation.TemporalitySelector used
// for selecting aggregation.Temporality (i.e., Cumulative vs. Delta
// aggregation). If not specified otherwise, exporter will use a
// cumulative temporality selector.
func WithMetricAggregationTemporalitySelector(selector aggregation.TemporalitySelector) Option {
	return exporterOptionFunc(func(cfg config) config {
		cfg.temporalitySelector = selector
		return cfg
	})
}
