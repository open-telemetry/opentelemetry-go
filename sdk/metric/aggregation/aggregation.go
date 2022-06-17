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

//go:build go1.17
// +build go1.17

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

	// Err returns an error for any misconfigured Aggregation.
	Err() error
}

// Drop is an aggregation that drops all recorded data.
type Drop struct{} // Drop has no parameters.

var _ Aggregation = Drop{}

func (Drop) private() {}

// Err returns an error for any misconfiguration. A Drop aggregation has no
// parameters and cannot be misconfigured, therefore this always returns nil.
func (Drop) Err() error { return nil }

// Default is an aggregation that uses the default instrument kind selection
// mapping to select another aggregation. A metric reader can be configured to
// make an aggregation selection based on instrument kind that differs from
// the default selection mapping. This aggregation ensures that the default
// instrument kind mapping. The default instrument kind mapping uses the
// following selection: Counter ⇨ Sum, Asynchronous Counter ⇨ Sum,
// UpDownCounter ⇨ Sum, Asynchronous UpDownCounter ⇨ Sum, Asynchronous Gauge ⇨
// LastValue, Histogram ⇨ ExplicitBucketHistogram.
type Default struct{} // Default has no parameters.

var _ Aggregation = Default{}

func (Default) private() {}

// Err returns an error for any misconfiguration. A Default aggregation has no
// parameters and cannot be misconfigured, therefore this always returns nil.
func (Default) Err() error { return nil }

// Sum is an aggregation that summarizes a set of measurements as their
// arithmetic sum.
type Sum struct{} // Sum has no parameters.

var _ Aggregation = Sum{}

func (Sum) private() {}

// Err returns an error for any misconfiguration. A Sum aggregation has no
// parameters and cannot be misconfigured, therefore this always returns nil.
func (Sum) Err() error { return nil }

// LastValue is an aggregation that summarizes a set of measurements as the
// last one made.
type LastValue struct{} // LastValue has no parameters.

var _ Aggregation = LastValue{}

func (LastValue) private() {}

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
	// RecordMinMax indicates whether to record the min and max of the
	// distribution.
	RecordMinMax bool
}

var _ Aggregation = ExplicitBucketHistogram{}

func (ExplicitBucketHistogram) private() {}

// errHist is returned by misconfigured ExplicitBucketHistograms.
var errHist = fmt.Errorf("%w: explicit bucket histogram", errAgg)

// Err returns an error for any misconfiguration.
func (h ExplicitBucketHistogram) Err() error {
	// Check boundaries are monotonic.
	if len(h.Boundaries) <= 1 {
		return nil
	}

	i := h.Boundaries[0]
	for _, j := range h.Boundaries[1:] {
		if i >= j {
			return fmt.Errorf("%w: non-monotonic boundaries: %v", errHist, h.Boundaries)
		}
	}

	return nil
}
