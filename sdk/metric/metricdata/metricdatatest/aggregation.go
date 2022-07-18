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

// equalAggregations returns true when a and b are equal. It returns false
// when they differ, along with the reasons why they differ.
func equalAggregations(a, b metricdata.Aggregation) (equal bool, reasons []string) {
	equal = true
	if a == nil || b == nil {
		if a != b {
			equal, reasons = false, []string{notEqualStr("Aggregation", a, b)}
		}
		return equal, reasons
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, []string{fmt.Sprintf(
			"Aggregation types not equal:\nexpected: %T\nactual: %T", a, b,
		)}
	}

	switch v := a.(type) {
	case metricdata.Gauge:
		var exp []string
		equal, exp = equalGauges(v, b.(metricdata.Gauge))
		if !equal {
			reasons = append(reasons, "Gauge not equal:")
			reasons = append(reasons, exp...)
		}
	case metricdata.Sum:
		var exp []string
		equal, exp = equalSums(v, b.(metricdata.Sum))
		if !equal {
			reasons = append(reasons, "Sum not equal:")
			reasons = append(reasons, exp...)
		}
	case metricdata.Histogram:
		var exp []string
		equal, exp = equalHistograms(v, b.(metricdata.Histogram))
		if !equal {
			reasons = append(reasons, "Histogram not equal:")
			reasons = append(reasons, exp...)
		}
	default:
		equal = false
		reasons = append(reasons, fmt.Sprintf("Aggregation of unknown types %T", a))
	}
	return equal, reasons
}
