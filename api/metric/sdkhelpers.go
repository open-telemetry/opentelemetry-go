// Copyright 2019, OpenTelemetry Authors
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

package metric

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type MeterSDK interface {
	// NewInt64Counter creates a new integral counter with a given
	// name and customized with passed options.
	NewInt64Counter(name core.Name, cos ...CounterOptionApplier) Int64Counter
	// NewFloat64Counter creates a new floating point counter with
	// a given name and customized with passed options.
	NewFloat64Counter(name core.Name, cos ...CounterOptionApplier) Float64Counter
	// NewInt64Gauge creates a new integral gauge with a given
	// name and customized with passed options.
	NewInt64Gauge(name core.Name, gos ...GaugeOptionApplier) Int64Gauge
	// NewFloat64Gauge creates a new floating point gauge with a
	// given name and customized with passed options.
	NewFloat64Gauge(name core.Name, gos ...GaugeOptionApplier) Float64Gauge
	// NewInt64Measure creates a new integral measure with a given
	// name and customized with passed options.
	NewInt64Measure(name core.Name, mos ...MeasureOptionApplier) Int64Measure
	// NewFloat64Measure creates a new floating point measure with
	// a given name and customized with passed options.
	NewFloat64Measure(name core.Name, mos ...MeasureOptionApplier) Float64Measure

	// RecordBatch atomically records a batch of measurements.
	RecordBatch(context.Context, []core.KeyValue, ...Measurement)
}

// InstrumentImpl is the implementation-level interface Set/Add/Record
// individual metrics without precomputed labels.
type InstrumentImpl interface {
	// Bind creates a Bound Instrument to record metrics with
	// precomputed labels.
	Bind(ctx context.Context, labels []core.KeyValue) BoundInstrumentImpl

	// RecordOne allows the SDK to observe a single metric event.
	RecordOne(ctx context.Context, number core.Number, labels []core.KeyValue)
}

// BoundInstrumentImpl is the implementation-level interface to Set/Add/Record
// individual metrics with precomputed labels.
type BoundInstrumentImpl interface {
	// RecordOne allows the SDK to observe a single metric event.
	RecordOne(ctx context.Context, number core.Number)

	// Unbind frees the resources associated with this bound instrument. It
	// does not affect the metric this bound instrument was created through.
	Unbind()
}

// WrapInt64CounterInstrument wraps the instrument in the type-safe
// wrapper as an integral counter.
//
// It is mostly intended for SDKs.
func WrapInt64CounterInstrument(instrument InstrumentImpl) Int64Counter {
	return Int64Counter{commonMetric: newCommonMetric(instrument)}
}

// WrapFloat64CounterInstrument wraps the instrument in the type-safe
// wrapper as an floating point counter.
//
// It is mostly intended for SDKs.
func WrapFloat64CounterInstrument(instrument InstrumentImpl) Float64Counter {
	return Float64Counter{commonMetric: newCommonMetric(instrument)}
}

// WrapInt64GaugeInstrument wraps the instrument in the type-safe
// wrapper as an integral gauge.
//
// It is mostly intended for SDKs.
func WrapInt64GaugeInstrument(instrument InstrumentImpl) Int64Gauge {
	return Int64Gauge{commonMetric: newCommonMetric(instrument)}
}

// WrapFloat64GaugeInstrument wraps the instrument in the type-safe
// wrapper as an floating point gauge.
//
// It is mostly intended for SDKs.
func WrapFloat64GaugeInstrument(instrument InstrumentImpl) Float64Gauge {
	return Float64Gauge{commonMetric: newCommonMetric(instrument)}
}

// WrapInt64MeasureInstrument wraps the instrument in the type-safe
// wrapper as an integral measure.
//
// It is mostly intended for SDKs.
func WrapInt64MeasureInstrument(instrument InstrumentImpl) Int64Measure {
	return Int64Measure{commonMetric: newCommonMetric(instrument)}
}

// WrapFloat64MeasureInstrument wraps the instrument in the type-safe
// wrapper as an floating point measure.
//
// It is mostly intended for SDKs.
func WrapFloat64MeasureInstrument(instrument InstrumentImpl) Float64Measure {
	return Float64Measure{commonMetric: newCommonMetric(instrument)}
}

// ApplyCounterOptions is a helper that applies all the counter
// options to passed opts.
func ApplyCounterOptions(opts *Options, cos ...CounterOptionApplier) {
	for _, o := range cos {
		o.ApplyCounterOption(opts)
	}
}

// ApplyGaugeOptions is a helper that applies all the gauge options to
// passed opts.
func ApplyGaugeOptions(opts *Options, gos ...GaugeOptionApplier) {
	for _, o := range gos {
		o.ApplyGaugeOption(opts)
	}
}

// ApplyMeasureOptions is a helper that applies all the measure
// options to passed opts.
func ApplyMeasureOptions(opts *Options, mos ...MeasureOptionApplier) {
	for _, o := range mos {
		o.ApplyMeasureOption(opts)
	}
}
