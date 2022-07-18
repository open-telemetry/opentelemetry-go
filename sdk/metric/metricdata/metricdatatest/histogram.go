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

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// CompareHistogram returns true when Histograms are equivalent. It returns
// false when they differ, along with messages describing the difference.
//
// The DataPoints each Histogram contains are compared based on containing the
// same HistogramDataPoint, not the order they are stored in.
func CompareHistogram(a, b metricdata.Histogram) (equal bool, explanation []string) {
	equal = true
	if a.Temporality != b.Temporality {
		equal, explanation = false, append(
			explanation,
			notEqualStr("Temporality", a.Temporality, b.Temporality),
		)
	}

	var exp string
	equal, exp = compareDiff(diffSlices(
		a.DataPoints,
		b.DataPoints,
		func(a, b metricdata.HistogramDataPoint) bool {
			equal, _ := CompareHistogramDataPoint(a, b)
			return equal
		},
	))
	if !equal {
		explanation = append(explanation, fmt.Sprintf(
			"Histogram DataPoints not equal:\n%s", exp,
		))
	}
	return equal, explanation
}

// AssertHistogramsEqual asserts that two Histogram are equal.
func AssertHistogramsEqual(t *testing.T, expected, actual metricdata.Histogram) bool {
	return assertCompare(CompareHistogram(expected, actual))(t)
}
