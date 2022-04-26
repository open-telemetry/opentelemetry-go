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
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

// Sentinel errors for Aggregator interface.
var (
	ErrNegativeInput = fmt.Errorf("negative value is out of range for this instrument")
	ErrNaNInput      = fmt.Errorf("NaN value is an invalid input")
	ErrInfInput      = fmt.Errorf("±Inf value is an invalid input")
)

// RangeTest is a common routine for testing for valid input values.
// This rejects NaN and Inf values.  This rejects negative values when the
// aggregation does not support negative values, including
// monotonic counter metrics and Histogram metrics.
func RangeTest[N number.Any, Traits number.Traits[N]](num N, kind aggregation.Category) bool {
	var traits Traits

	if traits.IsInf(num) {
		otel.Handle(ErrInfInput)
		return false
	}

	if traits.IsNaN(num) {
		otel.Handle(ErrNaNInput)
		return false
	}

	// Check for negative values
	switch kind {
	case aggregation.MonotonicSumCategory,
		aggregation.HistogramCategory:
		if num < 0 {
			otel.Handle(ErrNegativeInput)
			return false
		}
	}
	return true
}

type HistogramConfig struct {
	ExplicitBoundaries []float64
}

type Config struct {
	Histogram HistogramConfig
}

// Boundaries implements Histogram.Defaults.
func (hc HistogramConfig) Boundaries() []float64 {
	return hc.ExplicitBoundaries
}

// Methods implements a specific aggregation behavior.  Methods
// are parameterized by the type of the number (int64, flot64),
// the Storage (generally an `Storage` struct in the same package),
// and the Config (generally a `Config` struct in the same package).
type Methods[N number.Any, Storage any] interface {
	// Init initializes the storage with its configuration.
	Init(ptr *Storage, cfg Config)

	// Update modifies the aggregator concurrently with respect to
	// SynchronizedMove() for single new measurement.
	Update(ptr *Storage, number N)

	// SynchronizedMove concurrently copies and resets the
	// `inputIsReset` aggregator.
	SynchronizedMove(inputIsReset, output *Storage)

	// Simply reset the storage.
	Reset(ptr *Storage)

	// Merge adds the contents of `input` to `output`.
	Merge(output, input *Storage)

	// SubtractSwap removes the contents of `operand` from `valueToModify`
	SubtractSwap(valueToModify, operandToModify *Storage)

	// ToAggregation returns an exporter-ready value.
	ToAggregation(ptr *Storage) aggregation.Aggregation

	// ToStorage returns the underlying storage of an existing Aggregation.
	ToStorage(aggregation.Aggregation) (*Storage, bool)

	// Kind returns the Kind of aggregator.
	Kind() aggregation.Kind

	// HasChange returns true if there have been any (discernible)
	// Updates.  This tests whether an aggregation has zero sum,
	// zero count, or zero difference, depending on the
	// aggregation.  If the instrument is asynchronous, this will
	// be used after subtraction.
	HasChange(ptr *Storage) bool
}

// ConfigSelector is a per-instrument-kind, per-number-kind Config choice.
type ConfigSelector func(sdkinstrument.Kind) (int64Config, float64Config Config)
