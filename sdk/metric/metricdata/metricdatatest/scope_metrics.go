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

// CompareScopeMetrics returns true when ScopeMetrics are equivalent. It
// returns false when they differ, along with messages describing the
// difference.
//
// The Metrics each ScopeMetrics contains are compared based on containing the
// same Metrics, not the order they are stored in.
func CompareScopeMetrics(a, b metricdata.ScopeMetrics) (equal bool, explanation []string) {
	equal = true
	if a.Scope != b.Scope {
		equal, explanation = false, append(
			explanation,
			notEqualStr("Scope", a.Scope, b.Scope),
		)
	}

	var exp string
	equal, exp = compareDiff(diffSlices(
		a.Metrics,
		b.Metrics,
		func(a, b metricdata.Metrics) bool {
			equal, _ := CompareMetrics(a, b)
			return equal
		},
	))
	if !equal {
		explanation = append(explanation, fmt.Sprintf(
			"ScopeMetrics Metrics not equal:\n%s", exp,
		))
	}
	return equal, explanation
}

// AssertScopeMetricsEqual asserts that two ScopeMetrics are equal.
func AssertScopeMetricsEqual(t *testing.T, expected, actual metricdata.ScopeMetrics) bool {
	t.Helper()
	return assertCompare(CompareScopeMetrics(expected, actual))(t)
}
