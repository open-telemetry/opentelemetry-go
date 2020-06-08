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

package aggregation // import "go.opentelemetry.io/otel/sdk/export/metric/aggregation"

import (
	"fmt"
	"time"

	"go.opentelemetry.io/otel/api/metric"
)

// These interfaces describe the various ways to access state from an
// Aggregation.

type (
	// Aggregation is one of the interfaces below.
	Aggregation interface {
		// Kind returns the kind of aggregation used.
		Kind() Kind
	}

	// Sum returns an aggregated sum.
	Sum interface {
		Sum() (metric.Number, error)
	}

	// Sum returns the number of values that were aggregated.
	Count interface {
		Count() (int64, error)
	}

	// Min returns the minimum value over the set of values that were aggregated.
	Min interface {
		Min() (metric.Number, error)
	}

	// Max returns the maximum value over the set of values that were aggregated.
	Max interface {
		Max() (metric.Number, error)
	}

	// Quantile returns an exact or estimated quantile over the
	// set of values that were aggregated.
	Quantile interface {
		Quantile(float64) (metric.Number, error)
	}

	// LastValue returns the latest value that was aggregated.
	LastValue interface {
		LastValue() (metric.Number, time.Time, error)
	}

	// Points returns the raw set of values that were aggregated.
	Points interface {
		Points() ([]metric.Number, error)
	}

	// Buckets represents histogram buckets boundaries and counts.
	//
	// For a Histogram with N defined boundaries, e.g, [x, y, z].
	// There are N+1 counts: [-inf, x), [x, y), [y, z), [z, +inf]
	Buckets struct {
		// Boundaries are floating point numbers, even when
		// aggregating integers.
		Boundaries []float64

		// Counts are floating point numbers to account for
		// the possibility of sampling which allows for
		// non-integer count values.
		Counts []float64
	}

	// Histogram returns the count of events in pre-determined buckets.
	Histogram interface {
		Sum
		Histogram() (Buckets, error)
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

type (
	Kind string
)

const (
	SumKind            Kind = "sum"
	MinMaxSumCountKind Kind = "minmaxsumcount"
	HistogramKind      Kind = "histogram"
	LastValueKind      Kind = "lastvalue"
	SketchKind         Kind = "sketch"
	ExactKind          Kind = "exact"
)

var (
	ErrInvalidQuantile  = fmt.Errorf("the requested quantile is out of range")
	ErrNegativeInput    = fmt.Errorf("negative value is out of range for this instrument")
	ErrNaNInput         = fmt.Errorf("NaN value is an invalid input")
	ErrInconsistentType = fmt.Errorf("inconsistent aggregator types")
	ErrNoSubtraction    = fmt.Errorf("aggregator does not subtract")

	// ErrNoData is returned when (due to a race with collection)
	// the Aggregator is check-pointed before the first value is set.
	// The aggregator should simply be skipped in this case.
	ErrNoData = fmt.Errorf("no data collected by this aggregator")
)

// String returns a string representation of the aggregation kind.
func (k Kind) String() string {
	return string(k)
}
