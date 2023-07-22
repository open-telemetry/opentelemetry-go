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

package prometheus // import "go.opentelemetry.io/otel/bridge/prometheus"

import (
	"github.com/prometheus/client_golang/prometheus"
)

// config contains options for the exporter.
type config struct {
	gatherers []prometheus.Gatherer
}

// newConfig creates a validated config configured with options.
func newConfig(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if len(cfg.gatherers) == 0 {
		cfg.gatherers = []prometheus.Gatherer{prometheus.DefaultGatherer}
	}

	return cfg
}

// Option sets exporter option values.
type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (fn optionFunc) apply(cfg config) config {
	return fn(cfg)
}

// WithGatherer configures which prometheus Gatherer the Bridge will gather
// from. If no registerer is used the prometheus DefaultGatherer is used.
func WithGatherer(gatherer prometheus.Gatherer) Option {
	return optionFunc(func(cfg config) config {
		cfg.gatherers = append(cfg.gatherers, gatherer)
		return cfg
	})
}
