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

package aggregator

import (
	"fmt"
	"math"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

var (
	ErrEmptyDataSet     = fmt.Errorf("The result is not defined on an empty data set")
	ErrInvalidQuantile  = fmt.Errorf("The requested quantile is out of range")
	ErrNegativeInput    = fmt.Errorf("Negative value is out of range for this instrument")
	ErrNaNInput         = fmt.Errorf("NaN value is an invalid input")
	ErrNonMonotoneInput = fmt.Errorf("The new value is not monotone")
	ErrInconsistentType = fmt.Errorf("Cannot merge different aggregator types")

	// ErrNoLastValue is returned by the gauge.Aggregator when
	// (due to a race with collection) the Aggregator is
	// checkpointed before the first value is set.  The aggregator
	// should simply be skipped in this case.
	ErrNoLastValue = fmt.Errorf("No value has been set")
)

func RangeTest(number core.Number, descriptor *export.Descriptor) error {
	numberKind := descriptor.NumberKind()

	if numberKind == core.Float64NumberKind && math.IsNaN(number.AsFloat64()) {
		// NOTE: add this to the specification.
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
