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

// CompareSum returns true when Sums are equivalent. It returns false when
// they differ, along with messages describing the difference.
//
// The DataPoints each Sum contains are compared based on containing the same
// DataPoints, not the order they are stored in.
func CompareSum(a, b metricdata.Sum) (equal bool, explanation []string) {
	equal = true
	if a.Temporality != b.Temporality {
		equal, explanation = false, append(
			explanation,
			notEqualStr("Temporality", a.Temporality, b.Temporality),
		)
	}
	if a.IsMonotonic != b.IsMonotonic {
		equal, explanation = false, append(
			explanation,
			notEqualStr("IsMonotonic", a.IsMonotonic, b.IsMonotonic),
		)
	}

	var exp string
	equal, exp = compareDiff(diffSlices(
		a.DataPoints,
		b.DataPoints,
		func(a, b metricdata.DataPoint) bool {
			equal, _ := CompareDataPoint(a, b)
			return equal
		},
	))
	if !equal {
		explanation = append(explanation, fmt.Sprintf(
			"Sum DataPoints not equal:\n%s", exp,
		))
	}
	return equal, explanation
}

// AssertSumsEqual asserts that two Sum are equal.
func AssertSumsEqual(t *testing.T, expected, actual metricdata.Sum) bool {
	return assertCompare(CompareSum(expected, actual))(t)
}
