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
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// equalMetrics returns true when Metrics are equal. It returns false when
// they differ, along with the reasons why they differ.
func equalMetrics(a, b metricdata.Metrics) (equal bool, reasons []string) {
	equal = true
	if a.Name != b.Name {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Name", a.Name, b.Name),
		)
	}
	if a.Description != b.Description {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Description", a.Description, b.Description),
		)
	}
	if a.Unit != b.Unit {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Unit", a.Unit, b.Unit),
		)
	}

	var exp []string
	equal, exp = equalAggregations(a.Data, b.Data)
	if !equal {
		reasons = append(reasons, "Metrics Data not equal:")
		reasons = append(reasons, exp...)
	}
	return equal, reasons
}
