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
) // config is added here to allow for options expansion in the future.
type config struct {
	reader metric.Reader

	registry   *prometheus.Registry
	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer
}

func newConfig(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.reader == nil {
		cfg.reader = metric.NewManualReader()
	}

	if cfg.registry != nil {
		cfg.registerer = cfg.registry
		cfg.gatherer = cfg.registry
	} else {
		cfg.registerer = prometheus.DefaultRegisterer
		cfg.gatherer = prometheus.DefaultGatherer
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

func WithReader(rdr metric.Reader) Option {
	return optionFunc(func(cfg config) config {
		cfg.reader = rdr
		return cfg
	})
}

func WithRegistry(reg *prometheus.Registry) Option {
	return optionFunc(func(cfg config) config {
		cfg.registry = reg
		return cfg
	})
}
