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
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// CompareValues returns true when Values are equivalent. It returns false
// when they differ, along with a message describing the difference.
func CompareValues(a, b metricdata.Value) (equal bool, explanation []string) {
	if a == nil || b == nil {
		if a != b {
			equal, explanation = false, []string{notEqualStr("Values", a, b)}
		}
		return equal, explanation
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, []string{fmt.Sprintf(
			"Value types not equal:\nexpected: %T\nactual: %T", a, b,
		)}
	}

	switch v := a.(type) {
	case metricdata.Int64:
		var exp []string
		equal, exp = CompareInt64(v, b.(metricdata.Int64))
		if !equal {
			explanation = append(explanation, "Int64 not equal:")
			explanation = append(explanation, exp...)
		}
	case metricdata.Float64:
		var exp []string
		equal, exp = CompareFloat64(v, b.(metricdata.Float64))
		if !equal {
			explanation = append(explanation, "Int64 not equal:")
			explanation = append(explanation, exp...)
		}
	default:
		equal = false
		explanation = append(explanation, fmt.Sprintf("Value of unknown types %T", a))
	}

	return equal, explanation
}

// AssertValuesEqual asserts that two Values are equal.
func AssertValuesEqual(t *testing.T, expected, actual metricdata.Value) bool {
	return assertCompare(CompareValues(expected, actual))(t)
}
