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

package exporttest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"fmt"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// CompareAggregations returns true when a and b are equivalent. It returns
// false when they differ, along with messages describing the difference.
func CompareAggregations(a, b metricdata.Aggregation) (equal bool, explination []string) {
	if a == nil || b == nil {
		if a != b {
			equal, explination = false, []string{notEqualStr("Aggregation", a, b)}
		}
		return equal, explination
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, []string{fmt.Sprintf(
			"Aggregation types not equal:\nexpected: %T\nactual: %T", a, b,
		)}
	}

	switch v := a.(type) {
	case metricdata.Gauge:
		var exp []string
		equal, exp = CompareGauge(v, b.(metricdata.Gauge))
		if !equal {
			explination = append(explination, "Gauge not equal:")
			explination = append(explination, exp...)
		}
	case metricdata.Sum:
		var exp []string
		equal, exp = CompareSum(v, b.(metricdata.Sum))
		if !equal {
			explination = append(explination, "Sum not equal:")
			explination = append(explination, exp...)
		}
	case metricdata.Histogram:
		var exp []string
		equal, exp = CompareHistogram(v, b.(metricdata.Histogram))
		if !equal {
			explination = append(explination, "Histogram not equal:")
			explination = append(explination, exp...)
		}
	default:
		equal = false
		explination = append(explination, fmt.Sprintf("Aggregation of unknown types %T", a))
	}
	return equal, explination
}

// AssertAggregationsEqual asserts that two Aggregations are equal.
func AssertAggregationsEqual(t *testing.T, expected, actual metricdata.Aggregation) bool {
	return assertCompare(CompareAggregations(expected, actual))(t)
}
