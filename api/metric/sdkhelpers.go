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

// LabelSetDelegate is a general-purpose delegating implementation of
// the LabelSet interface.  This is implemented by the default
// Provider returned by api/global.SetMeterProvider(), and should be
// tested for by implementations before converting a LabelSet to their
// private concrete type.
type LabelSetDelegate interface {
	Delegate() LabelSet
}

// InstrumentImpl is the implementation-level interface Set/Add/Record
// individual metrics without precomputed labels.
type InstrumentImpl interface {
	// Bind creates a Bound Instrument to record metrics with
	// precomputed labels.
	Bind(labels LabelSet) BoundInstrumentImpl

	// RecordOne allows the SDK to observe a single metric event.
	RecordOne(ctx context.Context, number core.Number, labels LabelSet)
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
func WrapInt64CounterInstrument(instrument InstrumentImpl, err error) (Int64Counter, error) {
	common, err := newCommonMetric(instrument, err)
	return Int64Counter{commonMetric: common}, err
}

// WrapFloat64CounterInstrument wraps the instrument in the type-safe
// wrapper as an floating point counter.
//
// It is mostly intended for SDKs.
func WrapFloat64CounterInstrument(instrument InstrumentImpl, err error) (Float64Counter, error) {
	common, err := newCommonMetric(instrument, err)
	return Float64Counter{commonMetric: common}, err
}

// WrapInt64MeasureInstrument wraps the instrument in the type-safe
// wrapper as an integral measure.
//
// It is mostly intended for SDKs.
func WrapInt64MeasureInstrument(instrument InstrumentImpl, err error) (Int64Measure, error) {
	common, err := newCommonMetric(instrument, err)
	return Int64Measure{commonMetric: common}, err
}

// WrapFloat64MeasureInstrument wraps the instrument in the type-safe
// wrapper as an floating point measure.
//
// It is mostly intended for SDKs.
func WrapFloat64MeasureInstrument(instrument InstrumentImpl, err error) (Float64Measure, error) {
	common, err := newCommonMetric(instrument, err)
	return Float64Measure{commonMetric: common}, err
}

// Configure is a helper that applies all the options to a Config.
func Configure(opts []Option) Config {
	var config Config
	for _, o := range opts {
		o.Apply(&config)
	}
	return config
}
