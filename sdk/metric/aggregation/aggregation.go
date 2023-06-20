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

// Package aggregation contains configuration types that define the
// aggregation operation used to summarizes recorded measurements.
package aggregation // import "go.opentelemetry.io/otel/sdk/metric/aggregation"

import (
	"errors"
	"fmt"
)

// errAgg is wrapped by misconfigured aggregations.
var errAgg = errors.New("aggregation")

// Aggregation is the aggregation used to summarize recorded measurements.
type Aggregation interface {
	// private attempts to ensure no user-defined Aggregation are allowed. The
	// OTel specification does not allow user-defined Aggregation currently.
	private()

	// Copy returns a deep copy of the Aggregation.
	Copy() Aggregation

	// Err returns an error for any misconfigured Aggregation.
	Err() error
}

// Drop is an aggregation that drops all recorded data.
type Drop struct{} // Drop has no parameters.

var _ Aggregation = Drop{}

func (Drop) private() {}

// Copy returns a deep copy of d.
func (d Drop) Copy() Aggregation { return d }

// Err returns an error for any misconfiguration. A Drop aggregation has no
// parameters and cannot be misconfigured, therefore this always returns nil.
func (Drop) Err() error { return nil }

// Default is an aggregation that uses the default instrument kind selection
// mapping to select another aggregation. A metric reader can be configured to
// make an aggregation selection based on instrument kind that differs from
// the default. This aggregation ensures the default is used.
//
// See the "go.opentelemetry.io/otel/sdk/metric".DefaultAggregationSelector
// for information about the default instrument kind selection mapping.
type Default struct{} // Default has no parameters.

var _ Aggregation = Default{}

func (Default) private() {}

// Copy returns a deep copy of d.
func (d Default) Copy() Aggregation { return d }

// Err returns an error for any misconfiguration. A Default aggregation has no
// parameters and cannot be misconfigured, therefore this always returns nil.
func (Default) Err() error { return nil }

// Sum is an aggregation that summarizes a set of measurements as their
// arithmetic sum.
type Sum struct{} // Sum has no parameters.

var _ Aggregation = Sum{}

func (Sum) private() {}

// Copy returns a deep copy of s.
func (s Sum) Copy() Aggregation { return s }

// Err returns an error for any misconfiguration. A Sum aggregation has no
// parameters and cannot be misconfigured, therefore this always returns nil.
func (Sum) Err() error { return nil }

// LastValue is an aggregation that summarizes a set of measurements as the
// last one made.
type LastValue struct{} // LastValue has no parameters.

var _ Aggregation = LastValue{}

func (LastValue) private() {}

// Copy returns a deep copy of l.
func (l LastValue) Copy() Aggregation { return l }

// Err returns an error for any misconfiguration. A LastValue aggregation has
// no parameters and cannot be misconfigured, therefore this always returns
// nil.
func (LastValue) Err() error { return nil }

// ExplicitBucketHistogram is an aggregation that summarizes a set of
// measurements as an histogram with explicitly defined buckets.
type ExplicitBucketHistogram struct {
	// Boundaries are the increasing bucket boundary values. Boundary values
	// define bucket upper bounds. Buckets are exclusive of their lower
	// boundary and inclusive of their upper bound (except at positive
	// infinity). A measurement is defined to fall into the greatest-numbered
	// bucket with a boundary that is greater than or equal to the
	// measurement. As an example, boundaries defined as:
	//
	// []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 1000}
	//
	// Will define these buckets:
	//
	// (-∞, 0], (0, 5.0], (5.0, 10.0], (10.0, 25.0], (25.0, 50.0],
	// (50.0, 75.0], (75.0, 100.0], (100.0, 250.0], (250.0, 500.0],
	// (500.0, 1000.0], (1000.0, +∞)
	Boundaries []float64
	// NoMinMax indicates whether to not record the min and max of the
	// distribution. By default, these extrema are recorded.
	//
	// Recording these extrema for cumulative data is expected to have little
	// value, they will represent the entire life of the instrument instead of
	// just the current collection cycle. It is recommended to set this to true
	// for that type of data to avoid computing the low-value extrema.
	NoMinMax bool
}

var _ Aggregation = ExplicitBucketHistogram{}

func (ExplicitBucketHistogram) private() {}

// errHist is returned by misconfigured ExplicitBucketHistograms.
var errHist = fmt.Errorf("%w: explicit bucket histogram", errAgg)

// Err returns an error for any misconfiguration.
func (h ExplicitBucketHistogram) Err() error {
	if len(h.Boundaries) <= 1 {
		return nil
	}

	// Check boundaries are monotonic.
	i := h.Boundaries[0]
	for _, j := range h.Boundaries[1:] {
		if i >= j {
			return fmt.Errorf("%w: non-monotonic boundaries: %v", errHist, h.Boundaries)
		}
		i = j
	}

	return nil
}

// Copy returns a deep copy of h.
func (h ExplicitBucketHistogram) Copy() Aggregation {
	b := make([]float64, len(h.Boundaries))
	copy(b, h.Boundaries)
	return ExplicitBucketHistogram{
		Boundaries: b,
		NoMinMax:   h.NoMinMax,
	}
}

// ExponentialHistogram is an aggregation that summarizes a set of
// measurements as an histogram with buckets widths that grow exponentially.
type ExponentialHistogram struct {
	// MaxSize is the maximum number of buckets to use for the histogram.
	MaxSize int
	// MaxScale is the maximum resolution scale to use for the histogram.
	MaxScale int
	// ZeroThreshold sets the minimum value blow which will be recorded as zero.
	ZeroThreshold float64

	// NoMinMax indicates whether to not record the min and max of the
	// distribution. By default, these extrema are recorded.
	//
	// Recording these extrema for cumulative data is expected to have little
	// value, they will represent the entire life of the instrument instead of
	// just the current collection cycle. It is recommended to set this to true
	// for that type of data to avoid computing the low-value extrema.
	NoMinMax bool
}

// DefaultExponentialHistogram returns the default Exponential Histogram aggregation used by the SDK.
func DefaultExponentialHistogram() ExponentialHistogram {
	return ExponentialHistogram{
		MaxSize:       160,
		MaxScale:      20,
		ZeroThreshold: 0.0,
	}
}

var _ Aggregation = ExponentialHistogram{}

// private attempts to ensure no user-defined Aggregation are allowed. The
// OTel specification does not allow user-defined Aggregation currently.
func (e ExponentialHistogram) private() {}

// Copy returns a deep copy of the Aggregation.
func (e ExponentialHistogram) Copy() Aggregation {
	return e
}

const (
	expoMaxScale = 20
	expoMinScale = -10
)

// errExpoHist is returned by misconfigured ExplicitBucketHistograms.
var errExpoHist = fmt.Errorf("%w: explicit bucket histogram", errAgg)

// Err returns an error for any misconfigured Aggregation.
func (e ExponentialHistogram) Err() error {
	if e.MaxScale < expoMinScale {
		return fmt.Errorf("%w: max size %d is less than minimum scale %d", errExpoHist, e.MaxSize, expoMinScale)
	}
	if e.MaxScale > expoMaxScale {
		return fmt.Errorf("%w: max size %d is greater than maximum scale %d", errExpoHist, e.MaxSize, expoMaxScale)
	}
	if e.MaxSize <= 0 {
		return fmt.Errorf("%w: max size %d is less than or equal to zero", errExpoHist, e.MaxSize)
	}
	return nil
}
