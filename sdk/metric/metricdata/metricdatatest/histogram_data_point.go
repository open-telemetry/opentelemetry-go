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

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// CompareHistogramDataPoint returns true when HistogramDataPoints are
// equivalent. It returns false when they differ, along with messages
// describing the difference.
func CompareHistogramDataPoint(a, b metricdata.HistogramDataPoint) (equal bool, explanation []string) {
	equal = true
	if !a.Attributes.Equals(&b.Attributes) {
		equal, explanation = false, append(explanation, fmt.Sprintf(
			"Attributes not equal:\nexpected: %s\nactual: %s",
			a.Attributes.Encoded(attribute.DefaultEncoder()),
			b.Attributes.Encoded(attribute.DefaultEncoder()),
		))
	}
	if !a.StartTime.Equal(b.StartTime) {
		equal, explanation = false, append(explanation, fmt.Sprintf(
			"StartTime not equal:\nexpected: %d\nactual: %d",
			a.StartTime.UnixNano(),
			b.StartTime.UnixNano(),
		))
	}
	if !a.Time.Equal(b.Time) {
		equal, explanation = false, append(explanation, fmt.Sprintf(
			"Time not equal:\nexpected: %d\nactual: %d",
			a.Time.UnixNano(),
			b.Time.UnixNano(),
		))
	}
	if a.Count != b.Count {
		equal, explanation = false, append(explanation, fmt.Sprintf(
			"Count not equal:\nexpected: %d\nactual: %d",
			a.Count,
			b.Count,
		))
	}
	if !equalSlices(a.Bounds, b.Bounds) {
		equal, explanation = false, append(explanation, fmt.Sprintf(
			"Bounds not equal:\nexpected: %v\nactual: %v",
			a.Bounds,
			b.Bounds,
		))
	}
	if !equalSlices(a.BucketCounts, b.BucketCounts) {
		equal, explanation = false, append(explanation, fmt.Sprintf(
			"BucketCounts not equal:\nexpected: %v\nactual: %v",
			a.BucketCounts,
			b.BucketCounts,
		))
	}
	if !equalPtrValues(a.Min, b.Min) {
		equal, explanation = false, append(explanation, fmt.Sprintf(
			"Min not equal:\nexpected: %v\nactual: %v",
			a.Min,
			b.Min,
		))
	}
	if !equalPtrValues(a.Max, b.Max) {
		equal, explanation = false, append(explanation, fmt.Sprintf(
			"Max not equal:\nexpected: %v\nactual: %v",
			a.Max,
			b.Max,
		))
	}
	if a.Sum != b.Sum {
		equal, explanation = false, append(explanation, fmt.Sprintf(
			"Sum not equal:\nexpected: %g\nactual: %g",
			a.Sum,
			b.Sum,
		))
	}
	return equal, explanation
}

// AssertHistogramDataPointsEqual asserts that two HistogramDataPoint are equal.
func AssertHistogramDataPointsEqual(t *testing.T, expected, actual metricdata.HistogramDataPoint) bool {
	t.Helper()
	return assertCompare(CompareHistogramDataPoint(expected, actual))(t)
}
