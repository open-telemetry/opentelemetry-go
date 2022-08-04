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

package transform // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/transform"

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	cpb "go.opentelemetry.io/proto/otlp/common/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

func ScopeMetrics(sms []metricdata.ScopeMetrics) ([]*mpb.ScopeMetrics, error) {
	var errs []string
	out := make([]*mpb.ScopeMetrics, 0, len(sms))
	for _, sm := range sms {
		ms, err := Metrics(sm.Metrics)
		if err != nil {
			errs = append(errs, err.Error())
		}

		out = append(out, &mpb.ScopeMetrics{
			Scope: &cpb.InstrumentationScope{
				Name:    sm.Scope.Name,
				Version: sm.Scope.Version,
			},
			Metrics:   ms,
			SchemaUrl: sm.Scope.SchemaURL,
		})
	}
	if len(errs) > 0 {
		return out, fmt.Errorf("transform ScopeMetrics: %s", strings.Join(errs, "; "))
	}
	return out, nil
}
