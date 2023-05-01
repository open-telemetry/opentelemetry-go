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

// Package instrument provides the OpenTelemetry API instruments used to make
// measurements.
//
// Deprecated: Use go.opentelemetry.io/otel/metric instead.
package instrument // import "go.opentelemetry.io/otel/metric/instrument"

import (
	"go.opentelemetry.io/otel/metric"
)

// Float64Observable is an alias for [metric.Float64Observable].
//
// Deprecated: Use [metric.Float64Observable] instead.
type Float64Observable metric.Float64Observable

// Float64ObservableCounter is an alias for [metric.Float64ObservableCounter].
//
// Deprecated: Use [metric.Float64ObservableCounter] instead.
type Float64ObservableCounter metric.Float64ObservableCounter

// Float64ObservableCounterConfig is an alias for
// [metric.Float64ObservableCounterConfig].
//
// Deprecated: Use [metric.Float64ObservableCounterConfig] instead.
type Float64ObservableCounterConfig metric.Float64ObservableCounterConfig

// NewFloat64ObservableCounterConfig wraps
// [metric.NewFloat64ObservableCounterConfig].
//
// Deprecated: Use [metric.NewFloat64ObservableCounterConfig] instead.
func NewFloat64ObservableCounterConfig(opts ...Float64ObservableCounterOption) Float64ObservableCounterConfig {
	o := make([]metric.Float64ObservableCounterOption, len(opts))
	for i := range opts {
		o[i] = metric.Float64ObservableCounterOption(opts[i])
	}
	c := metric.NewFloat64ObservableCounterConfig(o...)
	return Float64ObservableCounterConfig(c)
}

// Float64ObservableCounterOption is an alias for
// [metric.Float64ObservableCounterOption].
//
// Deprecated: Use [metric.Float64ObservableCounterOption] instead.
type Float64ObservableCounterOption metric.Float64ObservableCounterOption

// Float64ObservableUpDownCounter is an alias for
// [metric.Float64ObservableUpDownCounter].
//
// Deprecated: Use [metric.Float64ObservableUpDownCounter] instead.
type Float64ObservableUpDownCounter metric.Float64ObservableUpDownCounter

// Float64ObservableUpDownCounterConfig is an alias for
// [metric.Float64ObservableUpDownCounterConfig].
//
// Deprecated: Use [metric.Float64ObservableUpDownCounterConfig] instead.
type Float64ObservableUpDownCounterConfig metric.Float64ObservableUpDownCounterConfig

// NewFloat64ObservableUpDownCounterConfig wraps
// [metric.NewFloat64ObservableUpDownCounterConfig].
//
// Deprecated: Use [metric.NewFloat64ObservableUpDownCounterConfig] instead.
func NewFloat64ObservableUpDownCounterConfig(opts ...Float64ObservableUpDownCounterOption) Float64ObservableUpDownCounterConfig {
	o := make([]metric.Float64ObservableUpDownCounterOption, len(opts))
	for i := range opts {
		o[i] = metric.Float64ObservableUpDownCounterOption(opts[i])
	}
	c := metric.NewFloat64ObservableUpDownCounterConfig(o...)
	return Float64ObservableUpDownCounterConfig(c)
}

// Float64ObservableUpDownCounterOption is an alias for
// [metric.Float64ObservableUpDownCounterOption].
//
// Deprecated: Use [metric.Float64ObservableUpDownCounterOption] instead.
type Float64ObservableUpDownCounterOption metric.Float64ObservableUpDownCounterOption

// Float64ObservableGauge is an alias for [metric.Float64ObservableGauge].
//
// Deprecated: Use [metric.Float64ObservableGauge] instead.
type Float64ObservableGauge metric.Float64ObservableGauge

// Float64ObservableGaugeConfig is an alias for
// [metric.Float64ObservableGaugeConfig].
//
// Deprecated: Use [metric.Float64ObservableGaugeConfig] instead.
type Float64ObservableGaugeConfig metric.Float64ObservableGaugeConfig

// NewFloat64ObservableGaugeConfig wraps
// [metric.NewFloat64ObservableGaugeConfig].
//
// Deprecated: Use [metric.NewFloat64ObservableGaugeConfig] instead.
func NewFloat64ObservableGaugeConfig(opts ...Float64ObservableGaugeOption) Float64ObservableGaugeConfig {
	o := make([]metric.Float64ObservableGaugeOption, len(opts))
	for i := range opts {
		o[i] = metric.Float64ObservableGaugeOption(opts[i])
	}
	c := metric.NewFloat64ObservableGaugeConfig(o...)
	return Float64ObservableGaugeConfig(c)
}

// Float64ObservableGaugeOption is an alias for
// [metric.Float64ObservableGaugeOption].
//
// Deprecated: Use [metric.Float64ObservableGaugeOption] instead.
type Float64ObservableGaugeOption metric.Float64ObservableGaugeOption

// Float64Observer is an alias for [metric.Float64Observer].
//
// Deprecated: Use [metric.Float64Observer] instead.
type Float64Observer metric.Float64Observer

// Float64Callback is an alias for [metric.Float64Callback].
//
// Deprecated: Use [metric.Float64Callback] instead.
type Float64Callback metric.Float64Callback

// Float64ObservableOption is an alias for [metric.Float64ObservableOption].
//
// Deprecated: Use [metric.Float64ObservableOption] instead.
type Float64ObservableOption metric.Float64ObservableOption

// WithFloat64Callback wraps [metric.WithFloat64Callback].
//
// Deprecated: Use [metric.WithFloat64Callback] instead.
func WithFloat64Callback(callback Float64Callback) Float64ObservableOption {
	cback := metric.Float64Callback(callback)
	opt := metric.WithFloat64Callback(cback)
	return Float64ObservableOption(opt)
}

// Int64Observable is an alias for [metric.Int64Observable].
//
// Deprecated: Use [metric.Int64Observable] instead.
type Int64Observable metric.Int64Observable

// Int64ObservableCounter is an alias for [metric.Int64ObservableCounter].
//
// Deprecated: Use [metric.Int64ObservableCounter] instead.
type Int64ObservableCounter metric.Int64ObservableCounter

// Int64ObservableCounterConfig is an alias for
// [metric.Int64ObservableCounterConfig].
//
// Deprecated: Use [metric.Int64ObservableCounterConfig] instead.
type Int64ObservableCounterConfig metric.Int64ObservableCounterConfig

// NewInt64ObservableCounterConfig wraps
// [metric.NewInt64ObservableCounterConfig].
//
// Deprecated: Use [metric.NewInt64ObservableCounterConfig] instead.
func NewInt64ObservableCounterConfig(opts ...Int64ObservableCounterOption) Int64ObservableCounterConfig {
	o := make([]metric.Int64ObservableCounterOption, len(opts))
	for i := range opts {
		o[i] = metric.Int64ObservableCounterOption(opts[i])
	}
	c := metric.NewInt64ObservableCounterConfig(o...)
	return Int64ObservableCounterConfig(c)
}

// Int64ObservableCounterOption is an alias for
// [metric.Int64ObservableCounterOption].
//
// Deprecated: Use [metric.Int64ObservableCounterOption] instead.
type Int64ObservableCounterOption metric.Int64ObservableCounterOption

// Int64ObservableUpDownCounter is an alias for
// [metric.Int64ObservableUpDownCounter].
//
// Deprecated: Use [metric.Int64ObservableUpDownCounter] instead.
type Int64ObservableUpDownCounter metric.Int64ObservableUpDownCounter

// Int64ObservableUpDownCounterConfig is an alias for
// [metric.Int64ObservableUpDownCounterConfig].
//
// Deprecated: Use [metric.Int64ObservableUpDownCounterConfig] instead.
type Int64ObservableUpDownCounterConfig metric.Int64ObservableUpDownCounterConfig

// NewInt64ObservableUpDownCounterConfig wraps
// [metric.NewInt64ObservableUpDownCounterConfig].
//
// Deprecated: Use [metric.NewInt64ObservableUpDownCounterConfig] instead.
func NewInt64ObservableUpDownCounterConfig(opts ...Int64ObservableUpDownCounterOption) Int64ObservableUpDownCounterConfig {
	o := make([]metric.Int64ObservableUpDownCounterOption, len(opts))
	for i := range opts {
		o[i] = metric.Int64ObservableUpDownCounterOption(opts[i])
	}
	c := metric.NewInt64ObservableUpDownCounterConfig(o...)
	return Int64ObservableUpDownCounterConfig(c)
}

// Int64ObservableUpDownCounterOption is an alias for
// [metric.Int64ObservableUpDownCounterOption].
//
// Deprecated: Use [metric.Int64ObservableUpDownCounterOption] instead.
type Int64ObservableUpDownCounterOption metric.Int64ObservableUpDownCounterOption

// Int64ObservableGauge is an alias for [metric.Int64ObservableGauge].
//
// Deprecated: Use [metric.Int64ObservableGauge] instead.
type Int64ObservableGauge metric.Int64ObservableGauge

// Int64ObservableGaugeConfig is an alias for
// [metric.Int64ObservableGaugeConfig].
//
// Deprecated: Use [metric.Int64ObservableGaugeConfig] instead.
type Int64ObservableGaugeConfig metric.Int64ObservableGaugeConfig

// NewInt64ObservableGaugeConfig wraps [metric.NewInt64ObservableGaugeConfig].
//
// Deprecated: Use [metric.NewInt64ObservableGaugeConfig] instead.
func NewInt64ObservableGaugeConfig(opts ...Int64ObservableGaugeOption) Int64ObservableGaugeConfig {
	o := make([]metric.Int64ObservableGaugeOption, len(opts))
	for i := range opts {
		o[i] = metric.Int64ObservableGaugeOption(opts[i])
	}
	c := metric.NewInt64ObservableGaugeConfig(o...)
	return Int64ObservableGaugeConfig(c)
}

// Int64ObservableGaugeOption is an alias for
// [metric.Int64ObservableGaugeOption].
//
// Deprecated: Use [metric.Int64ObservableGaugeOption] instead.
type Int64ObservableGaugeOption metric.Int64ObservableGaugeOption

// Int64Observer is an alias for [metric.Int64Observer].
//
// Deprecated: Use [metric.Int64Observer] instead.
type Int64Observer metric.Int64Observer

// Int64Callback is an alias for [metric.Int64Callback].
//
// Deprecated: Use [metric.Int64Callback] instead.
type Int64Callback metric.Int64Callback

// Int64ObservableOption is an alias for [metric.Int64ObservableOption].
//
// Deprecated: Use [metric.Int64ObservableOption] instead.
type Int64ObservableOption metric.Int64ObservableOption

// WithInt64Callback wraps [metric.WithInt64Callback].
//
// Deprecated: Use [metric.WithInt64Callback] instead.
func WithInt64Callback(callback Int64Callback) Int64ObservableOption {
	cback := metric.Int64Callback(callback)
	opt := metric.WithInt64Callback(cback)
	return Int64ObservableOption(opt)
}

// Float64Counter is an alias for [metric.Float64Counter].
//
// Deprecated: Use [metric.Float64Counter] instead.
type Float64Counter metric.Float64Counter

// Float64CounterConfig is an alias for [metric.Float64CounterConfig].
//
// Deprecated: Use [metric.Float64CounterConfig] instead.
type Float64CounterConfig metric.Float64CounterConfig

// NewFloat64CounterConfig wraps [metric.NewFloat64CounterConfig].
//
// Deprecated: Use [metric.NewFloat64CounterConfig] instead.
func NewFloat64CounterConfig(opts ...Float64CounterOption) Float64CounterConfig {
	o := make([]metric.Float64CounterOption, len(opts))
	for i := range opts {
		o[i] = metric.Float64CounterOption(opts[i])
	}
	c := metric.NewFloat64CounterConfig(o...)
	return Float64CounterConfig(c)
}

// Float64CounterOption is an alias for [metric.Float64CounterOption].
//
// Deprecated: Use [metric.Float64CounterOption] instead.
type Float64CounterOption metric.Float64CounterOption

// Float64UpDownCounter is an alias for [metric.Float64UpDownCounter].
//
// Deprecated: Use [metric.Float64UpDownCounter] instead.
type Float64UpDownCounter metric.Float64UpDownCounter

// Float64UpDownCounterConfig is an alias for
// [metric.Float64UpDownCounterConfig].
//
// Deprecated: Use [metric.Float64UpDownCounterConfig] instead.
type Float64UpDownCounterConfig metric.Float64UpDownCounterConfig

// NewFloat64UpDownCounterConfig wraps [metric.NewFloat64UpDownCounterConfig].
//
// Deprecated: Use [metric.NewFloat64UpDownCounterConfig] instead.
func NewFloat64UpDownCounterConfig(opts ...Float64UpDownCounterOption) Float64UpDownCounterConfig {
	o := make([]metric.Float64UpDownCounterOption, len(opts))
	for i := range opts {
		o[i] = metric.Float64UpDownCounterOption(opts[i])
	}
	c := metric.NewFloat64UpDownCounterConfig(o...)
	return Float64UpDownCounterConfig(c)
}

// Float64UpDownCounterOption is an alias for
// [metric.Float64UpDownCounterOption].
//
// Deprecated: Use [metric.Float64UpDownCounterOption] instead.
type Float64UpDownCounterOption metric.Float64UpDownCounterOption

// Float64Histogram is an alias for [metric.Float64Histogram].
//
// Deprecated: Use [metric.Float64Histogram] instead.
type Float64Histogram metric.Float64Histogram

// Float64HistogramConfig is an alias for [metric.Float64HistogramConfig].
//
// Deprecated: Use [metric.Float64HistogramConfig] instead.
type Float64HistogramConfig metric.Float64HistogramConfig

// NewFloat64HistogramConfig wraps [metric.NewFloat64HistogramConfig].
//
// Deprecated: Use [metric.NewFloat64HistogramConfig] instead.
func NewFloat64HistogramConfig(opts ...Float64HistogramOption) Float64HistogramConfig {
	o := make([]metric.Float64HistogramOption, len(opts))
	for i := range opts {
		o[i] = metric.Float64HistogramOption(opts[i])
	}
	c := metric.NewFloat64HistogramConfig(o...)
	return Float64HistogramConfig(c)
}

// Float64HistogramOption is an alias for [metric.Float64HistogramOption].
//
// Deprecated: Use [metric.Float64HistogramOption] instead.
type Float64HistogramOption metric.Float64HistogramOption

// Int64Counter is an alias for [metric.Int64Counter].
//
// Deprecated: Use [metric.Int64Counter] instead.
type Int64Counter metric.Int64Counter

// Int64CounterConfig is an alias for [metric.Int64CounterConfig].
//
// Deprecated: Use [metric.Int64CounterConfig] instead.
type Int64CounterConfig metric.Int64CounterConfig

// NewInt64CounterConfig wraps [metric.NewInt64CounterConfig].
//
// Deprecated: Use [metric.NewInt64CounterConfig] instead.
func NewInt64CounterConfig(opts ...Int64CounterOption) Int64CounterConfig {
	o := make([]metric.Int64CounterOption, len(opts))
	for i := range opts {
		o[i] = metric.Int64CounterOption(opts[i])
	}
	c := metric.NewInt64CounterConfig(o...)
	return Int64CounterConfig(c)
}

// Int64CounterOption is an alias for [metric.Int64CounterOption].
//
// Deprecated: Use [metric.Int64CounterOption] instead.
type Int64CounterOption metric.Int64CounterOption

// Int64UpDownCounter is an alias for [metric.Int64UpDownCounter].
//
// Deprecated: Use [metric.Int64UpDownCounter] instead.
type Int64UpDownCounter metric.Int64UpDownCounter

// Int64UpDownCounterConfig is an alias for [metric.Int64UpDownCounterConfig].
//
// Deprecated: Use [metric.Int64UpDownCounterConfig] instead.
type Int64UpDownCounterConfig metric.Int64UpDownCounterConfig

// NewInt64UpDownCounterConfig wraps [metric.NewInt64UpDownCounterConfig].
//
// Deprecated: Use [metric.NewInt64UpDownCounterConfig] instead.
func NewInt64UpDownCounterConfig(opts ...Int64UpDownCounterOption) Int64UpDownCounterConfig {
	o := make([]metric.Int64UpDownCounterOption, len(opts))
	for i := range opts {
		o[i] = metric.Int64UpDownCounterOption(opts[i])
	}
	c := metric.NewInt64UpDownCounterConfig(o...)
	return Int64UpDownCounterConfig(c)
}

// Int64UpDownCounterOption is an alias for [metric.Int64UpDownCounterOption].
//
// Deprecated: Use [metric.Int64UpDownCounterOption] instead.
type Int64UpDownCounterOption metric.Int64UpDownCounterOption

// Int64Histogram is an alias for [metric.Int64Histogram].
//
// Deprecated: Use [metric.Int64Histogram] instead.
type Int64Histogram metric.Int64Histogram

// Int64HistogramConfig is an alias for [metric.Int64HistogramConfig].
//
// Deprecated: Use [metric.Int64HistogramConfig] instead.
type Int64HistogramConfig metric.Int64HistogramConfig

// NewInt64HistogramConfig wraps [metric.NewInt64HistogramConfig].
//
// Deprecated: Use [metric.NewInt64HistogramConfig] instead.
func NewInt64HistogramConfig(opts ...Int64HistogramOption) Int64HistogramConfig {
	o := make([]metric.Int64HistogramOption, len(opts))
	for i := range opts {
		o[i] = metric.Int64HistogramOption(opts[i])
	}
	c := metric.NewInt64HistogramConfig(o...)
	return Int64HistogramConfig(c)
}

// Int64HistogramOption is an alias for [metric.Int64HistogramOption].
//
// Deprecated: Use [metric.Int64HistogramOption] instead.
type Int64HistogramOption metric.Int64HistogramOption

// Observable is an alias for [metric.Observable].
//
// Deprecated: Use [metric.Observable] instead.
type Observable metric.Observable

// Option is an alias for [metric.InstrumentOption].
//
// Deprecated: Use [metric.InstrumentOption] instead.
type Option metric.InstrumentOption

// WithDescription is an alias for [metric.WithDescription].
//
// Deprecated: Use [metric.WithDescription] instead.
func WithDescription(desc string) Option {
	o := metric.WithDescription(desc)
	return Option(o)
}

// WithUnit is an alias for [metric.WithUnit].
//
// Deprecated: Use [metric.WithUnit] instead.
func WithUnit(u string) Option {
	o := metric.WithUnit(u)
	return Option(o)
}
