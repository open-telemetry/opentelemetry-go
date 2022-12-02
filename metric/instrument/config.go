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

package instrument // import "go.opentelemetry.io/otel/metric/instrument"

import "go.opentelemetry.io/otel/metric/unit"

// config contains options for all instruments.
type config struct {
	description string
	unit        unit.Unit
}

// Description returns the Config description.
func (c config) Description() string {
	return c.description
}

// Unit returns the Config unit.
func (c config) Unit() unit.Unit {
	return c.unit
}

// Option applies options to all instrument configuration.
type Option interface {
	AsynchronousOption
	SynchronousOption
}

type descriptionOption string

func (o descriptionOption) applySynchronous(cfg SynchronousConfig) SynchronousConfig {
	cfg.description = string(o)
	return cfg
}

func (o descriptionOption) applyAsynchronous(cfg AsynchronousConfig) AsynchronousConfig {
	cfg.description = string(o)
	return cfg
}

// WithDescription sets the instrument description.
func WithDescription(desc string) Option {
	return descriptionOption(desc)
}

type unitOption unit.Unit

func (o unitOption) applySynchronous(cfg SynchronousConfig) SynchronousConfig {
	cfg.unit = unit.Unit(o)
	return cfg
}

func (o unitOption) applyAsynchronous(cfg AsynchronousConfig) AsynchronousConfig {
	cfg.unit = unit.Unit(o)
	return cfg
}

// WithUnit sets the instrument unit.
func WithUnit(u unit.Unit) Option {
	return unitOption(u)
}

// SynchronousConfig contains options for Synchronous instruments.
type SynchronousConfig struct {
	config
}

// NewSynchronousConfig returns a new SynchronousConfig with all opts applied.
func NewSynchronousConfig(opts ...SynchronousOption) SynchronousConfig {
	var config SynchronousConfig
	for _, o := range opts {
		config = o.applySynchronous(config)
	}
	return config
}

// SynchronousOption applies options to Synchronous instruments.
type SynchronousOption interface {
	applySynchronous(SynchronousConfig) SynchronousConfig
}

// AsynchronousConfig contains options for Asynchronous instruments.
type AsynchronousConfig struct {
	config

	callbacks []Callback
}

// NewAsynchronousConfig returns a new AsynchronousConfig with all opts applied.
func NewAsynchronousConfig(opts ...AsynchronousOption) AsynchronousConfig {
	var config AsynchronousConfig
	for _, o := range opts {
		config = o.applyAsynchronous(config)
	}
	return config
}

// Callbacks returns the AsynchronousConfig callbacks.
func (c AsynchronousConfig) Callbacks() []Callback {
	return c.callbacks
}

// AsynchronousOption applies options to Asynchronous instruments.
type AsynchronousOption interface {
	applyAsynchronous(AsynchronousConfig) AsynchronousConfig
}

type callbackOption Callback

func (o callbackOption) applyAsynchronous(cfg AsynchronousConfig) AsynchronousConfig {
	cfg.callbacks = append(cfg.callbacks, (Callback)(o))
	return cfg
}

// WithCallback adds callback to be called for an Asynchronous instrument.
func WithCallback(callback Callback) AsynchronousOption {
	return callbackOption(callback)
}
