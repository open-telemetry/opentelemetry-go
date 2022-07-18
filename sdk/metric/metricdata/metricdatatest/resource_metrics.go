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

// CompareResourceMetrics returns true when ResourceMetrics are equivalent. It
// returns false when they differ, along with messages describing the
// difference.
//
// The ScopeMetrics each ResourceMetrics contains are compared based on
// containing the same ScopeMetrics, not the order they are stored in.
func CompareResourceMetrics(a, b metricdata.ResourceMetrics) (equal bool, explanation []string) {
	equal = true
	if !a.Resource.Equal(b.Resource) {
		equal, explanation = false, append(
			explanation, notEqualStr("Resources", a.Resource, b.Resource),
		)
	}

	var exp string
	equal, exp = compareDiff(diffSlices(
		a.ScopeMetrics,
		b.ScopeMetrics,
		func(a, b metricdata.ScopeMetrics) bool {
			equal, _ := CompareScopeMetrics(a, b)
			return equal
		},
	))
	if !equal {
		explanation = append(explanation, fmt.Sprintf(
			"ResourceMetrics ScopeMetrics not equal:\n%s", exp,
		))
	}
	return equal, explanation
}

// AssertResourceMetricsEqual asserts that two ResourceMetrics are equal.
func AssertResourceMetricsEqual(t *testing.T, expected, actual metricdata.ResourceMetrics) bool {
	t.Helper()
	return assertCompare(CompareResourceMetrics(expected, actual))(t)
}
