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
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// CompareDataPoint returns true when DataPoints are equivalent. It returns
// false when they differ, along with messages describing the difference.
func CompareDataPoint(a, b metricdata.DataPoint) (equal bool, explination []string) {
	equal = true
	if !a.Attributes.Equals(&b.Attributes) {
		equal, explination = false, append(explination, notEqualStr(
			"Attributes",
			a.Attributes.Encoded(attribute.DefaultEncoder()),
			b.Attributes.Encoded(attribute.DefaultEncoder()),
		))
	}
	if !a.StartTime.Equal(b.StartTime) {
		equal, explination = false, append(explination, notEqualStr(
			"StartTime",
			a.StartTime.UnixNano(),
			b.StartTime.UnixNano(),
		))
	}
	if !a.Time.Equal(b.Time) {
		equal, explination = false, append(explination, notEqualStr(
			"Time",
			a.Time.UnixNano(),
			b.Time.UnixNano(),
		))
	}

	var exp []string
	equal, exp = CompareValues(a.Value, b.Value)
	if !equal {
		explination = append(explination, "DataPoint Value not equal:")
		explination = append(explination, exp...)
	}
	return equal, explination
}

// AssertDataPointsEqual asserts that two DataPoint are equal.
func AssertDataPointsEqual(t *testing.T, expected, actual metricdata.DataPoint) bool {
	return assertCompare(CompareDataPoint(expected, actual))(t)
}
