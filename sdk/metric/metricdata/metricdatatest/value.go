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

package exporttest

import (
	"fmt"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/export"
)

// CompareValues returns true when Values are equivalent. It returns false
// when they differ, along with a message describing the difference.
func CompareValues(a, b export.Value) (equal bool, explination []string) {
	if a == nil || b == nil {
		if a != b {
			equal, explination = false, []string{notEqualStr("Values", a, b)}
		}
		return equal, explination
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, []string{fmt.Sprintf(
			"Value types not equal:\nexpected: %T\nactual: %T", a, b,
		)}
	}

	switch v := a.(type) {
	case export.Int64:
		var exp []string
		equal, exp = CompareInt64(v, b.(export.Int64))
		if !equal {
			explination = append(explination, "Int64 not equal:")
			explination = append(explination, exp...)
		}
	case export.Float64:
		var exp []string
		equal, exp = CompareFloat64(v, b.(export.Float64))
		if !equal {
			explination = append(explination, "Int64 not equal:")
			explination = append(explination, exp...)
		}
	default:
		equal = false
		explination = append(explination, fmt.Sprintf("Value of unknown types %T", a))
	}

	return equal, explination
}

// AssertValuesEqual asserts that two Values are equal.
func AssertValuesEqual(t *testing.T, expected, actual export.Value) bool {
	return assertCompare(CompareValues(expected, actual))(t)
}
