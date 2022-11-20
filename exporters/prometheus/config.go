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

package prometheus // import "go.opentelemetry.io/otel/exporters/prometheus"

import (
	"github.com/prometheus/client_golang/prometheus"

	"go.opentelemetry.io/otel/sdk/metric"
)

// config contains options for the exporter.
type config struct {
	registerer        prometheus.Registerer
	disableTargetInfo bool
	withoutUnits      bool
	aggregation       metric.AggregationSelector
	disableScopeInfo  bool
}

// newConfig creates a validated config configured with options.
func newConfig(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.registerer == nil {
		cfg.registerer = prometheus.DefaultRegisterer
	}

	return cfg
}

func (cfg config) manualReaderOptions() []metric.ManualReaderOption {
	opts := []metric.ManualReaderOption{}
	if cfg.aggregation != nil {
		opts = append(opts, metric.WithAggregationSelector(cfg.aggregation))
	}
	return opts
}

// Option sets exporter option values.
type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (fn optionFunc) apply(cfg config) config {
	return fn(cfg)
}

// WithRegisterer configures which prometheus Registerer the Exporter will
// register with.  If no registerer is used the prometheus DefaultRegisterer is
// used.
func WithRegisterer(reg prometheus.Registerer) Option {
	return optionFunc(func(cfg config) config {
		cfg.registerer = reg
		return cfg
	})
}

// WithAggregationSelector configure the Aggregation Selector the exporter will
// use. If no AggregationSelector is provided the DefaultAggregationSelector is
// used.
func WithAggregationSelector(agg metric.AggregationSelector) Option {
	return optionFunc(func(cfg config) config {
		cfg.aggregation = agg
		return cfg
	})
}

// WithoutTargetInfo configures the Exporter to not export the resource target_info metric.
// If not specified, the Exporter will create a target_info metric containing
// the metrics' resource.Resource attributes.
func WithoutTargetInfo() Option {
	return optionFunc(func(cfg config) config {
		cfg.disableTargetInfo = true
		return cfg
	})
}

// WithoutUnits disables exporter's addition of unit suffixes to metric names,
// and will also prevent unit comments from being added in OpenMetrics once
// unit comments are supported.
//
// By default, metric names include a unit suffix to follow Prometheus naming
// conventions. For example, the counter metric request.duration, with unit
// milliseconds would become request_duration_milliseconds_total.
// With this option set, the name would instead be request_duration_total.
func WithoutUnits() Option {
	return optionFunc(func(cfg config) config {
		cfg.withoutUnits = true
		return cfg
	})
}

// WithoutScopeInfo configures the Exporter to not export the otel_scope_info metric.
// If not specified, the Exporter will create a otel_scope_info metric containing
// the metrics' Instrumentation Scope, and also add labels about Instrumentation Scope to all metric points.
func WithoutScopeInfo() Option {
	return optionFunc(func(cfg config) config {
		cfg.disableScopeInfo = true
		return cfg
	})
}
