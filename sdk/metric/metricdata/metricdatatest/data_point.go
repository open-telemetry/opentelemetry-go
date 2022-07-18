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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// equalDataPoints returns true when DataPoints are equal. It returns false
// when they differ, along with the reasons why they differ.
func equalDataPoints(a, b metricdata.DataPoint) (equal bool, reasons []string) {
	equal = true
	if !a.Attributes.Equals(&b.Attributes) {
		equal, reasons = false, append(reasons, notEqualStr(
			"Attributes",
			a.Attributes.Encoded(attribute.DefaultEncoder()),
			b.Attributes.Encoded(attribute.DefaultEncoder()),
		))
	}
	if !a.StartTime.Equal(b.StartTime) {
		equal, reasons = false, append(reasons, notEqualStr(
			"StartTime",
			a.StartTime.UnixNano(),
			b.StartTime.UnixNano(),
		))
	}
	if !a.Time.Equal(b.Time) {
		equal, reasons = false, append(reasons, notEqualStr(
			"Time",
			a.Time.UnixNano(),
			b.Time.UnixNano(),
		))
	}

	var r []string
	equal, r = equalValues(a.Value, b.Value)
	if !equal {
		reasons = append(reasons, "DataPoint Value not equal:")
		reasons = append(reasons, r...)
	}
	return equal, reasons
}
