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

// Package aggregation contains types that define the aggregation operation
// used to summarizes recorded measurements.

package aggregation // import "go.opentelemetry.io/otel/sdk/metric/aggregation"

import (
	"errors" // Aggregation defines the aggregation operation to use.
	"fmt"
)

type Aggregation struct {
	Operation operation
}

// Err returns an error if a is invalid, nil otherwise.
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

type operation interface {
	isOperation()
}

// Drop aggregation drops all data recorded.
type Drop struct {
	// Drop aggregation has no parameters.
}

func (Drop) isOperation() {}

// Sum aggregation summarizes a set of measurements as their arithmetic sum.
type Sum struct {
	// Sum aggregation has no parameters.
}

func (Sum) isOperation() {}

// LastValues summarizes a set of measurements as the last one made.
type LastValue struct {
	// LastValue aggregation has no parameters.
}

func (LastValue) isOperation() {}

// ExplicitBucketHistogram summarizes a set of measurements as an histogram
// with explicitly defined buckets.
type ExplicitBucketHistogram struct {
	Boundaries   []float64
	RecordMinMax bool
}

func (ExplicitBucketHistogram) isOperation() {}
