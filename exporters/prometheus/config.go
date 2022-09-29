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
) // config is added here to allow for options expansion in the future.
type config struct {
	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer
}

func newConfig(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.gatherer == nil {
		cfg.gatherer = prometheus.DefaultGatherer
	}
	if cfg.registerer == nil {
		cfg.registerer = prometheus.DefaultRegisterer
	}

	return cfg
}

// Option may be used in the future to apply options to a Prometheus Exporter config.
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

// WithRegisterer configures which prometheus Gatherer the Exporter will
// Gather from.  If no gatherer is used the prometheus DefaultGatherer is
// used.
func WithGatherer(gatherer prometheus.Gatherer) Option {
	return optionFunc(func(cfg config) config {
		cfg.gatherer = gatherer
		return cfg
	})
}
