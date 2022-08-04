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
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	rpb "go.opentelemetry.io/proto/otlp/resource/v1"
)

// ResourceMetrics transforms a metricdata.ResourceMetrics type into an OTLP
// ReResourceMetrics type.
func ResourceMetrics(rm metricdata.ResourceMetrics) (*mpb.ResourceMetrics, error) {
	sms, err := ScopeMetrics(rm.ScopeMetrics)
	return &mpb.ResourceMetrics{
		Resource: &rpb.Resource{
			Attributes: AttrIter(rm.Resource.Iter()),
		},
		ScopeMetrics: sms,
		SchemaUrl:    rm.Resource.SchemaURL(),
	}, err
}
