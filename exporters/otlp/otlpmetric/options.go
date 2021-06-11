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

package otlpmetric

import metricsdk "go.opentelemetry.io/otel/sdk/export/metric"

// Option are setting options passed to an Exporter on creation.
type Option interface {
	apply(*config)
}

type exporterOptionFunc func(*config)

func (fn exporterOptionFunc) apply(cfg *config) {
	fn(cfg)
}

type config struct {
	exportKindSelector metricsdk.ExportKindSelector
}

// WithMetricExportKindSelector defines the ExportKindSelector used
// for selecting AggregationTemporality (i.e., Cumulative vs. Delta
// aggregation). If not specified otherwise, exporter will use a
// cumulative export kind selector.
func WithMetricExportKindSelector(selector metricsdk.ExportKindSelector) Option {
	return exporterOptionFunc(func(cfg *config) {
		cfg.exportKindSelector = selector
	})
}
