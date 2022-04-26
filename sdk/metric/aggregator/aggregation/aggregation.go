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

//go:generate stringer -type=Category,Kind

package aggregation // import "go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"

import (
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

// These interfaces describe the various ways to access state from an
// Aggregation.

type (
	// Aggregation is an interface returned by the Aggregator
	// containing an interval of metric data.
	Aggregation interface {
		Kind() Kind
	}

	// HasASum
	HasASum interface {
		Sum() number.Number
	}

	// Sum returns an aggregated sum.
	Sum interface {
		// Review NOTE: Should this be Total() or Value()?
		HasASum
		IsMonotonic() bool
	}

	// Gauge returns the latest value that was aggregated.
	Gauge interface {
		Aggregation

		// Review NOTE: Should this be LastValue() or Value()?
		Gauge() number.Number
	}

	// Buckets represents histogram buckets boundaries and counts.
	//
	// For a Histogram with N defined boundaries, e.g, [x, y, z].
	// There are N+1 counts: [-inf, x), [x, y), [y, z), [z, +inf]
	Buckets struct {
		// Boundaries are floating point numbers, even when
		// aggregating integers.
		Boundaries []float64

		// Counts holds the count in each bucket.
		Counts []uint64
	}

	// Histogram returns the count of events in pre-determined buckets.
	Histogram interface {
		Aggregation
		Count() uint64
		HasASum
		Buckets() Buckets
	}
)

// Category constants describe semantic kind.  For the histogram
// category there are multiple implementations, for those distinctions
// as well as Drop, use Kind.
type Category int

const (
	UndefinedCategory Category = iota
	MonotonicSumCategory
	NonMonotonicSumCategory
	GaugeCategory
	HistogramCategory
)

type Kind int

const (
	UndefinedKind Kind = iota
	DropKind
	AnySumKind
	MonotonicSumKind
	NonMonotonicSumKind
	GaugeKind
	HistogramKind
)

func (k Kind) Category(ik sdkinstrument.Kind) Category {
	switch k {
	case AnySumKind:
		switch ik {
		case sdkinstrument.HistogramKind, sdkinstrument.CounterKind, sdkinstrument.CounterObserverKind:
			return MonotonicSumCategory
		case sdkinstrument.UpDownCounterKind, sdkinstrument.UpDownCounterObserverKind:
			return NonMonotonicSumCategory
		}
		return UndefinedCategory
	case MonotonicSumKind:
		return MonotonicSumCategory
	case NonMonotonicSumKind:
		return NonMonotonicSumCategory
	case GaugeKind:
		return GaugeCategory
	case HistogramKind:
		return HistogramCategory
	default:
		return UndefinedCategory
	}
}

// KindSelector is a per-instrument-kind Kind choice.
type KindSelector func(sdkinstrument.Kind) Kind
