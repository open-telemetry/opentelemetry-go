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

package metric

// Config contains configuration for an SDK.
type Config struct {
	// If provided, MetricsLabelsEnricher is executed each time a metric is recorded
	// by the Accumulator's sync instrument implementation
	MetricsLabelsEnricher MetricsLabelsEnricher
}

// Option is the interface that applies the value to a configuration option.
type Option interface {
	// Apply sets the Option value of a Config.
	Apply(*Config)
}

func WithMetricsLabelsEnricher(e MetricsLabelsEnricher) Option {
	return metricsLabelsEnricherOption(e)
}

type metricsLabelsEnricherOption MetricsLabelsEnricher

func (e metricsLabelsEnricherOption) Apply(config *Config) {
	config.MetricsLabelsEnricher = MetricsLabelsEnricher(e)
}
