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

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

import (
	"go.opentelemetry.io/otel/attribute"
)

// Aggregation is a single data point in a timeseries that summarizes
// measurements made during a time span.
type Aggregation struct {
	// TODO(#2968): Replace this with the export.Aggregation type once #2961
	// is merged.

	// Timestamp defines the time the last measurement was made. If zero, no
	// measurements were made for this time span. The time is represented as a
	// unix timestamp with nanosecond precision.
	Timestamp uint64

	// Attributes are the unique dimensions Value describes.
	Attributes *attribute.Set

	// Value is the summarization of the measurements made.
	Value value
}

type value interface {
	private()
}

// SingleValue summarizes a set of measurements as a single value.
type SingleValue[N int64 | float64] struct {
	Value N
}

func (SingleValue[N]) private() {}

// HistogramValue summarizes a set of measurements as a histogram.
type HistogramValue struct {
	Bounds   []float64
	Counts   []uint64
	Sum      float64
	Min, Max float64
}

func (HistogramValue) private() {}
