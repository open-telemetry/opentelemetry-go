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

package asyncfloat64 // import "go.opentelemetry.io/otel/metric/instrument/asyncfloat64"

import "go.opentelemetry.io/otel/metric/unit"

// Config contains options for Asynchronous instruments that observe float64
// values.
type Config struct {
	description string
	unit        unit.Unit
	callbacks   []Callback
}

// NewConfig returns a new Config with all opts applied.
func NewConfig(opts ...Option) Config {
	var config Config
	for _, o := range opts {
		config = o.apply(config)
	}
	return config
}

// Description returns the Config description.
func (c Config) Description() string {
	return c.description
}

// Unit returns the Config unit.
func (c Config) Unit() unit.Unit {
	return c.unit
}

// Callbacks returns the Config callbacks.
func (c Config) Callbacks() []Callback {
	return c.callbacks
}

// Option applies options to the instruments in this package.
type Option interface {
	apply(Config) Config
}

type optionFunc func(Config) Config

func (fn optionFunc) apply(cfg Config) Config {
	return fn(cfg)
}

// WithDescription sets the instrument description.
func WithDescription(desc string) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.description = desc
		return cfg
	})
}

// WithUnit sets the instrument unit.
func WithUnit(u unit.Unit) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.unit = u
		return cfg
	})
}

// WithCallback adds callback to be called for an instrument.
func WithCallback(callback Callback) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.callbacks = append(cfg.callbacks, callback)
		return cfg
	})
}
