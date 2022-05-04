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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

// config contains configuration options for a MeterProvider.
type config struct {
	readers []Reader
}

// Option applies a configuration option value to a MeterProvider.
type Option interface {
	apply(config) config
}

type optionFunc func(cfg config) config

func (o optionFunc) apply(cfg config) config {
	return o(cfg)
}

func WithReader(rdr Reader) Option {
	return optionFunc(func(cfg config) config {
		cfg.readers = append(cfg.readers, rdr)
		return cfg
	})
}

// TODO (#2819): implement provider options.
