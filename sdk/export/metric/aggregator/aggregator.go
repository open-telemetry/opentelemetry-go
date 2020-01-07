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

package aggregator // import "go.opentelemetry.io/otel/sdk/export/metric/aggregator"

import (
	"fmt"
	"math"
	"time"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

// These interfaces describe the various ways to access state from an
// Aggregator.

type (
	// Sum returns an aggregated sum.
	Sum interface {
		Sum() (core.Number, error)
	}

	// Sum returns the number of values that were aggregated.
	Count interface {
		Count() (int64, error)
	}

	// Min returns the minimum value over the set of values that were aggregated.
	Min interface {
		Min() (core.Number, error)
	}

	// Max returns the maximum value over the set of values that were aggregated.
	Max interface {
		Max() (core.Number, error)
	}

	// Quantile returns an exact or estimated quantile over the
	// set of values that were aggregated.
	Quantile interface {
		Quantile(float64) (core.Number, error)
	}

	// LastValue returns the latest value that was aggregated.
	LastValue interface {
		LastValue() (core.Number, time.Time, error)
	}

	// Points returns the raw set of values that were aggregated.
	Points interface {
		Points() ([]core.Number, error)
	}

	// MinMaxSumCount supports the Min, Max, Sum, and Count interfaces.
	MinMaxSumCount interface {
		Min
		Max
		Sum
		Count
	}

	// Distribution supports the Min, Max, Sum, Count, and Quantile
	// interfaces.
	Distribution interface {
		MinMaxSumCount
		Quantile
	}
)

var (
	ErrInvalidQuantile  = fmt.Errorf("the requested quantile is out of range")
	ErrNegativeInput    = fmt.Errorf("negative value is out of range for this instrument")
	ErrNaNInput         = fmt.Errorf("NaN value is an invalid input")
	ErrNonMonotoneInput = fmt.Errorf("the new value is not monotone")
	ErrInconsistentType = fmt.Errorf("inconsistent aggregator types")

	// ErrNoLastValue is returned by the LastValue interface when
	// (due to a race with collection) the Aggregator is
	// checkpointed before the first value is set.  The aggregator
	// should simply be skipped in this case.
	ErrNoLastValue = fmt.Errorf("no value has been set")

	// ErrEmptyDataSet is returned by Max and Quantile interfaces
	// when (due to a race with collection) the Aggregator is
	// checkpointed before the first value is set.  The aggregator
	// should simply be skipped in this case.
	ErrEmptyDataSet = fmt.Errorf("the result is not defined on an empty data set")
)

// NewInconsistentMergeError formats an error describing an attempt to
// merge different-type aggregators.  The result can be unwrapped as
// an ErrInconsistentType.
func NewInconsistentMergeError(a1, a2 export.Aggregator) error {
	return fmt.Errorf("cannot merge %T with %T: %w", a1, a2, ErrInconsistentType)
}

// RangeTest is a commmon routine for testing for valid input values.
// This rejects NaN values.  This rejects negative values when the
// metric instrument does not support negative values, including
// monotonic counter metrics and absolute measure metrics.
func RangeTest(number core.Number, descriptor *export.Descriptor) error {
	numberKind := descriptor.NumberKind()

	if numberKind == core.Float64NumberKind && math.IsNaN(number.AsFloat64()) {
		return ErrNaNInput
	}

	switch descriptor.MetricKind() {
	case export.CounterKind, export.MeasureKind:
		if !descriptor.Alternate() && number.IsNegative(numberKind) {
			return ErrNegativeInput
		}
	}
	return nil
}
