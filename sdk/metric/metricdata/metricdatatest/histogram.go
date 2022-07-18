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

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// equalHistograms returns true when Histograms are equal. It returns false
// when they differ, along with the reasons why they differ.
//
// The DataPoints each Histogram contains are compared based on containing the
// same HistogramDataPoint, not the order they are stored in.
func equalHistograms(a, b metricdata.Histogram) (equal bool, reasons []string) {
	equal = true
	if a.Temporality != b.Temporality {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Temporality", a.Temporality, b.Temporality),
		)
	}

	var r string
	equal, r = compareDiff(diffSlices(
		a.DataPoints,
		b.DataPoints,
		func(a, b metricdata.HistogramDataPoint) bool {
			equal, _ := equalHistogramDataPoints(a, b)
			return equal
		},
	))
	if !equal {
		reasons = append(reasons, fmt.Sprintf(
			"Histogram DataPoints not equal:\n%s", r,
		))
	}
	return equal, reasons
}
