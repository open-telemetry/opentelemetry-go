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

package stdoutmetric

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

func ExampleExporter_Export() {
	exampleData := metricdata.ResourceMetrics{
		Resource: resource.NewWithAttributes("example", attribute.String("resource-foo", "bar")),
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{
					Name:      "Scope1",
					Version:   "v1",
					SchemaURL: "anotherExample",
				},
				Metrics: []metricdata.Metrics{
					{
						Name:        "",
						Description: "",
						Unit:        "",
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(),
									StartTime:  time.Now().Add(-1 * time.Minute),
									Time:       time.Now(),
									Value:      3,
								},
							},
						},
					},
					{},
					{},
				},
			},
			{
				Scope: instrumentation.Scope{
					Name:      "Scope2",
					Version:   "v2",
					SchemaURL: "aDifferentExample",
				},
				Metrics: []metricdata.Metrics{
					{},
					{},
					{},
				},
			},
		},
	}
	exp, _ := New()
	exp.Export(context.Background(), exampleData)
	// Output: Hello
}
