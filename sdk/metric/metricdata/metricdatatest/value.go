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

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// equalValues returns true when Values are equal. It returns false when they
// differ, along with the reasons why they differ.
func equalValues(a, b metricdata.Value) (equal bool, reasons []string) {
	equal = true
	if a == nil || b == nil {
		if a != b {
			equal, reasons = false, []string{notEqualStr("Values", a, b)}
		}
		return equal, reasons
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, []string{fmt.Sprintf(
			"Value types not equal:\nexpected: %T\nactual: %T", a, b,
		)}
	}

	switch v := a.(type) {
	case metricdata.Int64:
		var r []string
		equal, r = equalInt64(v, b.(metricdata.Int64))
		if !equal {
			reasons = append(reasons, "Int64 not equal:")
			reasons = append(reasons, r...)
		}
	case metricdata.Float64:
		var r []string
		equal, r = equalFloat64(v, b.(metricdata.Float64))
		if !equal {
			reasons = append(reasons, "Int64 not equal:")
			reasons = append(reasons, r...)
		}
	default:
		equal = false
		reasons = append(reasons, fmt.Sprintf("Value of unknown types %T", a))
	}

	return equal, reasons
}
