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

// Aggregation is the aggregation used to summarize recorded measurements.
type Aggregation struct {
	// Operation is the kind of operation performed by the aggregation and the
	// configuration for that operation. This can be Drop, Sum, LastValue, or
	// ExplicitBucketHistogram.
	Operation operation
}

// Err returns an error if Aggregation a is invalid, nil otherwise.
func (a Aggregation) Err() error {
	if a.Operation == nil {
		return errors.New("aggregation: unset operation")
	}
	switch a.Operation.(type) {
	case Drop, Sum, LastValue, ExplicitBucketHistogram:
		return nil
	}
	return fmt.Errorf("aggregation: unknown %T", a.Operation)
}

// operation is an aggregation operation. The OTel specification does not
// allow user-defined aggregations, therefore, this is not exported.
type operation interface {
	isOperation()
}

// Drop drops all data recorded.
type Drop struct{} // The Drop operation has no parameters.

func (Drop) isOperation() {}

// Sum summarizes a set of measurements as their arithmetic sum.
type Sum struct{} // The Sum operation has no parameters.

func (Sum) isOperation() {}

// LastValues summarizes a set of measurements as the last one made.
type LastValue struct{} // The LastValue operation has no parameters.

func (LastValue) isOperation() {}

// ExplicitBucketHistogram summarizes a set of measurements as an histogram
// with explicitly defined buckets.
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

func (ExplicitBucketHistogram) isOperation() {}
