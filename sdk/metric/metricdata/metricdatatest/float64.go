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
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// CompareFloat64 returns true when Float64s are equivalent. It returns false
// when they differ, along with messages describing the difference.
func CompareFloat64(a, b metricdata.Float64) (equal bool, explanation []string) {
	equal = true
	if a != b {
		equal, explanation = false, append(
			explanation, notEqualStr("Float64 value", a, b),
		)
	}
	return equal, explanation
}

// AssertFloat64sEqual asserts that two Float64 are equal.
func AssertFloat64sEqual(t *testing.T, expected, actual metricdata.Float64) bool {
	return assertCompare(CompareFloat64(expected, actual))(t)
}
