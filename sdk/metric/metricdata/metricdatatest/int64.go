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

// CompareInt64 returns true when Int64s are equivalent. It returns false when
// they differ, along with messages describing the difference.
func CompareInt64(a, b metricdata.Int64) (equal bool, explanation []string) {
	equal = true
	if a != b {
		equal, explanation = false, append(
			explanation, notEqualStr("Int64 value", a, b),
		)
	}
	return equal, explanation
}

// AssertInt64sEqual asserts that two Int64 are equal.
func AssertInt64sEqual(t *testing.T, expected, actual metricdata.Int64) bool {
	return assertCompare(CompareInt64(expected, actual))(t)
}
