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

// Counter is an instrument that records increasing incremental values.
//
// Warning: Methods may be added to this interface in minor releases. See
// [go.opentelemetry.io/otel/metric] package documentation on API
// implementation for information on how to set default behavior for
// unimplemented methods.
type Counter[N int64 | float64] interface {
	embedded.Counter[N]

	// Add records a change to the counter.
	Add(ctx context.Context, incr N, attrs ...attribute.KeyValue)
}

// CounterConfig contains options for synchronous counter instruments that
// records increasing incremental values.
type CounterConfig[N int64 | float64] struct {
	description string
	unit        string
}

// NewCounterConfig returns a new [CounterConfig] with all opts applied.
func NewCounterConfig[N int64 | float64](opts ...CounterOption[N]) CounterConfig[N] {
	var config CounterConfig[N]
	for _, o := range opts {
		config = o.applyCounter(config)
	}
	return config
}

// Description returns the configured description.
func (c CounterConfig[N]) Description() string {
	return c.description
}

// Unit returns the configured unit.
func (c CounterConfig[N]) Unit() string {
	return c.unit
}

// CounterOption applies options to a [CounterConfig]. See [Option] for other
// options that can be used as an CounterOption.
type CounterOption[N int64 | float64] interface {
	applyCounter(CounterConfig[N]) CounterConfig[N]
}

// UpDownCounter is an instrument that records increasing or decreasing
// incremental values.
//
// Warning: Methods may be added to this interface in minor releases. See
// [go.opentelemetry.io/otel/metric] package documentation on API
// implementation for information on how to set default behavior for
// unimplemented methods.
type UpDownCounter[N int64 | float64] interface {
	embedded.UpDownCounter[N]

	// Add records a change to the counter.
	Add(ctx context.Context, incr N, attrs ...attribute.KeyValue)
}

// UpDownCounterConfig contains options for synchronous counter
// instruments that record int64 values.
type UpDownCounterConfig[N int64 | float64] struct {
	description string
	unit        string
}

// NewUpDownCounterConfig returns a new [UpDownCounterConfig] with
// all opts applied.
func NewUpDownCounterConfig[N int64 | float64](opts ...UpDownCounterOption[N]) UpDownCounterConfig[N] {
	var config UpDownCounterConfig[N]
	for _, o := range opts {
		config = o.applyUpDownCounter(config)
	}
	return config
}

// Description returns the configured description.
func (c UpDownCounterConfig[N]) Description() string {
	return c.description
}

// Unit returns the configured unit.
func (c UpDownCounterConfig[N]) Unit() string {
	return c.unit
}

// UpDownCounterOption applies options to a [UpDownCounterConfig].
// See [Option] for other options that can be used as an
// UpDownCounterOption.
type UpDownCounterOption[N int64 | float64] interface {
	applyUpDownCounter(UpDownCounterConfig[N]) UpDownCounterConfig[N]
}

// Histogram is an instrument that records a distribution of values.
//
// Warning: Methods may be added to this interface in minor releases. See
// [go.opentelemetry.io/otel/metric] package documentation on API
// implementation for information on how to set default behavior for
// unimplemented methods.
type Histogram[N int64 | float64] interface {
	embedded.Histogram[N]

	// Record adds an additional value to the distribution.
	Record(ctx context.Context, incr N, attrs ...attribute.KeyValue)
}

// HistogramConfig contains options for synchronous Histogram instruments.
type HistogramConfig[N int64 | float64] struct {
	description string
	unit        string
}

// NewHistogramConfig returns a new [HistogramConfig] with all opts applied.
func NewHistogramConfig[N int64 | float64](opts ...HistogramOption[N]) HistogramConfig[N] {
	var config HistogramConfig[N]
	for _, o := range opts {
		config = o.applyHistogram(config)
	}
	return config
}

// Description returns the configured description.
func (c HistogramConfig[N]) Description() string {
	return c.description
}

// Unit returns the configured unit.
func (c HistogramConfig[N]) Unit() string {
	return c.unit
}

// HistogramOption applies options to a [HistogramConfig]. See [Option] for
// other options that can be used as an HistogramOption.
type HistogramOption[N int64 | float64] interface {
	applyHistogram(HistogramConfig[N]) HistogramConfig[N]
}
