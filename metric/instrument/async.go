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

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/embedded"
)

// Observable represents all instruments that make measurements
// asynchronously via a Callback. It is used as a grouping mechanism.
type Observable interface {
	observable()
}

// ObservableT represents Observable instruments that make observations for a
// value of type N via a Callback. It is used as a grouping mechanism.
type ObservableT[N int64 | float64] interface {
	Observable

	observableT()
}

// ObservableCounter is an instrument used to asynchronously record increasing
// measurements once per collection cycle. Observations are only made within a
// callback for this instrument. The value observed is assumed the to be the
// cumulative sum of the count.
//
// Warning: Methods may be added to this interface in minor releases. See
// [go.opentelemetry.io/otel/metric] package documentation on API
// implementation for information on how to set default behavior for
// unimplemented methods.
type ObservableCounter[N int64 | float64] interface {
	embedded.ObservableCounter[N]

	ObservableT[N]
}

// ObservableCounterConfig contains options for asynchronous counter
// instruments.
type ObservableCounterConfig[N int64 | float64] struct {
	description string
	unit        string
	callbacks   []Callback[N]
}

// NewObservableCounterConfig returns a new [ObservableCounterConfig]
// with all opts applied.
func NewObservableCounterConfig[N int64 | float64](opts ...ObservableCounterOption[N]) ObservableCounterConfig[N] {
	var config ObservableCounterConfig[N]
	for _, o := range opts {
		config = o.applyObservableCounter(config)
	}
	return config
}

// Description returns the configured description.
func (c ObservableCounterConfig[N]) Description() string {
	return c.description
}

// Unit returns the configured unit.
func (c ObservableCounterConfig[N]) Unit() string {
	return c.unit
}

// Callbacks returns the configured callbacks.
func (c ObservableCounterConfig[N]) Callbacks() []Callback[N] {
	return c.callbacks
}

// ObservableCounterOption applies options to a
// [ObservableCounterConfig]. See [ObservableOption] and [Option] for
// other options that can be used as an ObservableCounterOption.
type ObservableCounterOption[N int64 | float64] interface {
	applyObservableCounter(ObservableCounterConfig[N]) ObservableCounterConfig[N]
}

// ObservableUpDownCounter is an instrument used to asynchronously record
// measurements once per collection cycle. Observations are only made within a
// callback for this instrument. The value observed is assumed the to be the
// cumulative sum of the count.
//
// Warning: Methods may be added to this interface in minor releases. See
// [go.opentelemetry.io/otel/metric] package documentation on API
// implementation for information on how to set default behavior for
// unimplemented methods.
type ObservableUpDownCounter[N int64 | float64] interface {
	embedded.ObservableUpDownCounter[N]

	ObservableT[N]
}

// ObservableUpDownCounterConfig contains options for asynchronous
// up-down-counter instruments.
type ObservableUpDownCounterConfig[N int64 | float64] struct {
	description string
	unit        string
	callbacks   []Callback[N]
}

// NewObservableUpDownCounterConfig returns a new
// [ObservableUpDownCounterConfig] with all opts applied.
func NewObservableUpDownCounterConfig[N int64 | float64](opts ...ObservableUpDownCounterOption[N]) ObservableUpDownCounterConfig[N] {
	var config ObservableUpDownCounterConfig[N]
	for _, o := range opts {
		config = o.applyObservableUpDownCounter(config)
	}
	return config
}

// Description returns the configured description.
func (c ObservableUpDownCounterConfig[N]) Description() string {
	return c.description
}

// Unit returns the configured unit.
func (c ObservableUpDownCounterConfig[N]) Unit() string {
	return c.unit
}

// Callbacks returns the configured callbacks.
func (c ObservableUpDownCounterConfig[N]) Callbacks() []Callback[N] {
	return c.callbacks
}

// ObservableUpDownCounterOption applies options to a
// [ObservableUpDownCounterConfig]. See [ObservableOption] and
// [Option] for other options that can be used as an
// ObservableUpDownCounterOption.
type ObservableUpDownCounterOption[N int64 | float64] interface {
	applyObservableUpDownCounter(ObservableUpDownCounterConfig[N]) ObservableUpDownCounterConfig[N]
}

// ObservableGauge is an instrument used to asynchronously record instantaneous
// measurements once per collection cycle. Observations are only made within a
// callback for this instrument.
//
// Warning: Methods may be added to this interface in minor releases. See
// [go.opentelemetry.io/otel/metric] package documentation on API
// implementation for information on how to set default behavior for
// unimplemented methods.
type ObservableGauge[N int64 | float64] interface {
	embedded.ObservableGauge[N]

	ObservableT[N]
}

// ObservableGaugeConfig contains options for asynchronous gauge instruments.
type ObservableGaugeConfig[N int64 | float64] struct {
	description string
	unit        string
	callbacks   []Callback[N]
}

// NewObservableGaugeConfig returns a new [ObservableGaugeConfig]
// with all opts applied.
func NewObservableGaugeConfig[N int64 | float64](opts ...ObservableGaugeOption[N]) ObservableGaugeConfig[N] {
	var config ObservableGaugeConfig[N]
	for _, o := range opts {
		config = o.applyObservableGauge(config)
	}
	return config
}

// Description returns the configured description.
func (c ObservableGaugeConfig[N]) Description() string {
	return c.description
}

// Unit returns the configured unit.
func (c ObservableGaugeConfig[N]) Unit() string {
	return c.unit
}

// Callbacks returns the configured callbacks.
func (c ObservableGaugeConfig[N]) Callbacks() []Callback[N] {
	return c.callbacks
}

// ObservableGaugeOption applies options to a
// [ObservableGaugeConfig]. See [ObservableOption] and [Option] for
// other options that can be used as an ObservableGaugeOption.
type ObservableGaugeOption[N int64 | float64] interface {
	applyObservableGauge(ObservableGaugeConfig[N]) ObservableGaugeConfig[N]
}

// ObserverT is a recorder of measurements for an ObservableT that measures
// values of type N.
//
// Warning: Methods may be added to this interface in minor releases. See
// [go.opentelemetry.io/otel/metric] package documentation on API
// implementation for information on how to set default behavior for
// unimplemented methods.
type ObserverT[N int64 | float64] interface {
	embedded.ObserverT[N]

	// Observe records the value with attributes.
	Observe(value N, attributes ...attribute.KeyValue)
}

// Callback is a function registered with a Meter that makes observations
// for an Observerable instrument it is registered with. Calls to the
// Observer record measurement values for the Observable.
//
// The function needs to complete in a finite amount of time and the deadline
// of the passed context is expected to be honored.
//
// The function needs to make unique observations across all registered
// Callbacks. Meaning, it should not report measurements with the same
// attributes as another Callbacks also registered for the same
// instrument.
//
// The function needs to be concurrent safe.
type Callback[N int64 | float64] func(context.Context, ObserverT[N]) error

// ObservableOption applies options to Observable instruments.
type ObservableOption[N int64 | float64] interface {
	ObservableCounterOption[N]
	ObservableUpDownCounterOption[N]
	ObservableGaugeOption[N]
}

type callbackOpt[N int64 | float64] struct {
	cback Callback[N]
}

func (o callbackOpt[N]) applyObservableCounter(cfg ObservableCounterConfig[N]) ObservableCounterConfig[N] {
	cfg.callbacks = append(cfg.callbacks, o.cback)
	return cfg
}

func (o callbackOpt[N]) applyObservableUpDownCounter(cfg ObservableUpDownCounterConfig[N]) ObservableUpDownCounterConfig[N] {
	cfg.callbacks = append(cfg.callbacks, o.cback)
	return cfg
}

func (o callbackOpt[N]) applyObservableGauge(cfg ObservableGaugeConfig[N]) ObservableGaugeConfig[N] {
	cfg.callbacks = append(cfg.callbacks, o.cback)
	return cfg
}

// WithCallback adds callback to be called for an instrument.
func WithCallback[N int64 | float64](callback Callback[N]) ObservableOption[N] {
	return callbackOpt[N]{callback}
}
