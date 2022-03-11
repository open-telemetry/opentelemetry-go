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

package aggregator // import "go.opentelemetry.io/otel/sdk/metric/aggregator"

import (
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

// RangeTest is a common routine for testing for valid input values.
// This rejects NaN values.  This rejects negative values when the
// metric instrument does not support negative values, including
// monotonic counter metrics and absolute Histogram metrics.
func RangeTest[N number.Any, Traits traits.Any[N]](num N, desc *sdkapi.Descriptor) error {
	var traits Traits

	// @@@ Should we have an Inf check?

	if traits.IsNaN(num) {
		return aggregation.ErrNaNInput
	}

	switch desc.InstrumentKind() {
	case sdkapi.CounterInstrumentKind,
		sdkapi.CounterObserverInstrumentKind,
		sdkapi.HistogramInstrumentKind: // @@@ right?
		if num < 0 {
			return aggregation.ErrNegativeInput
		}
	}
	return nil
}

// Methods implements a specific aggregation behavior, e.g., a
// behavior to track a sequence of updates to an instrument.  Counter
// instruments commonly use a simple Sum aggregator, but for the
// distribution instruments (Histogram, GaugeObserver) there are a
// number of possible aggregators with different cost and accuracy
// tradeoffs.
//
// Note that any Aggregator may be attached to any instrument--this is
// the result of the OpenTelemetry API/SDK separation.  It is possible
// to attach a Sum aggregator to a Histogram instrument.
type Methods[N number.Any, Storage, Config any] interface {
	Init(ptr *Storage, cfg Config)
	Update(ptr *Storage, number N)
	SynchronizedMove(inputIsReset, output *Storage)
	Merge(output, input *Storage)
	Aggregation(ptr *Storage) aggregation.Aggregation
	Storage(aggr aggregation.Aggregation) *Storage
}
