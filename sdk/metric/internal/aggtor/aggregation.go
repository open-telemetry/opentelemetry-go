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

//go:build go1.18
// +build go1.18

package aggtor // import "go.opentelemetry.io/otel/sdk/metric/internal/aggtor"

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

// Aggregation is a single data point in a timeseries that summarizes
// measurements made during a time span.
type Aggregation struct {
	// Timestamp defines the time the last measurement was made. If zero, no
	// measurements were made for this time span. The time is represented as a
	// unix timestamp with nanosecond precision.
	Timestamp uint64

	// Attributes are the unique dimensions Value describes.
	Attributes *attribute.Set

	// Value is the summarization of the measurements made.
	Value value
}

var errIncompatible = errors.New("incompatible aggregation")

// Fold combines other into a.
func (a Aggregation) Fold(other Aggregation) error {
	if other.Timestamp > a.Timestamp {
		a.Timestamp = other.Timestamp
	}
	if !a.Attributes.Equals(other.Attributes) {
		return fmt.Errorf("%w: attributes not equal", errIncompatible)
	}
	return a.Value.fold(other.Value)
}

type value interface {
	// fold combines other into the value. It will return an errIncompatible
	// if other is not a compatible type with value.
	fold(other value) error
}

// SingleValue summarizes a set of measurements as a single numeric value.
type SingleValue[N int64 | float64] struct {
	Value N
}

func (v SingleValue[N]) fold(other value) error {
	o, ok := other.(SingleValue[N])
	if !ok {
		return fmt.Errorf("%w: value types %T and %T", errIncompatible, v, other)
	}
	v.Value += o.Value
	return nil
}

// HistogramValue summarizes a set of measurements as a histogram.
type HistogramValue struct {
	Bounds   []float64
	Counts   []uint64
	Sum      float64
	Min, Max float64
}

func (v HistogramValue) fold(other value) error {
	o, ok := other.(HistogramValue)
	if !ok {
		return fmt.Errorf("%w: value types %T and %T", errIncompatible, v, other)
	}
	if !sliceEqual[float64](v.Bounds, o.Bounds) || len(o.Counts) != len(v.Counts) {
		return fmt.Errorf("%w: different histogram binning", errIncompatible)
	}
	v.Sum += o.Sum
	for i, c := range o.Counts {
		v.Counts[i] += c
	}
	if o.Min < v.Min {
		v.Min = o.Min
	}
	if o.Max > v.Max {
		v.Max = o.Max
	}
	return nil
}

func sliceEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
