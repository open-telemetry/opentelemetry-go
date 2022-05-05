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

import (
	"go.opentelemetry.io/otel/sdk/metric/view"
	"go.opentelemetry.io/otel/sdk/resource"
)

// config contains configuration options for a MeterProvider.
type config struct {
	res     *resource.Resource
	readers []Reader
	views   []*view.Views
}

// Option applies a configuration option value to a MeterProvider.
type Option interface {
	apply(config) config
}

// optionFunction makes a functional Option out of a function object.
type optionFunction func(cfg config) config

// apply implements Option.
func (of optionFunction) apply(in config) config {
	return of(in)
}

// WithResource associates a Resource with a new MeterProvider.
func WithResource(res *resource.Resource) Option {
	return optionFunction(func(cfg config) config {
		cfg.res = res
		return cfg
	})
}

// WithReader associates a new Reader and associated View options with
// a new MeterProvider
func WithReader(r Reader, opts ...view.Option) Option {
	return optionFunction(func(cfg config) config {
		cfg.readers = append(cfg.readers, r)
		cfg.views = append(cfg.views, view.New(r.String(), opts...))
		return cfg
	})
}
