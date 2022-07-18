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

// equalScopeMetrics returns true when ScopeMetrics are equal. It returns
// false when they differ, along with the reasons why they differ.
//
// The Metrics each ScopeMetrics contains are compared based on containing the
// same Metrics, not the order they are stored in.
func equalScopeMetrics(a, b metricdata.ScopeMetrics) (equal bool, reasons []string) {
	equal = true
	if a.Scope != b.Scope {
		equal, reasons = false, append(
			reasons,
			notEqualStr("Scope", a.Scope, b.Scope),
		)
	}

	var r string
	equal, r = compareDiff(diffSlices(
		a.Metrics,
		b.Metrics,
		func(a, b metricdata.Metrics) bool {
			equal, _ := equalMetrics(a, b)
			return equal
		},
	))
	if !equal {
		reasons = append(reasons, fmt.Sprintf(
			"ScopeMetrics Metrics not equal:\n%s", r,
		))
	}
	return equal, reasons
}
