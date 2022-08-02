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
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestExporterExport(t *testing.T) {
	now := time.Unix(1558033527, 0)

	exampleData := metricdata.ResourceMetrics{
		Resource: resource.NewWithAttributes("example", attribute.String("resource-foo", "bar")),
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{
					Name:      "scope1",
					Version:   "v1",
					SchemaURL: "anotherExample",
				},
				Metrics: []metricdata.Metrics{
					{
						Name:        "gauge",
						Description: "Gauge's description",
						Unit:        unit.Dimensionless,
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(),
									StartTime:  now.Add(-1 * time.Minute),
									Time:       now,
									Value:      3,
								},
							},
						},
					},
					{
						Name:        "sum",
						Description: "Sum's description",
						Unit:        unit.Bytes,
						Data: metricdata.Sum[float64]{
							Temporality: metricdata.DeltaTemporality,
							IsMonotonic: false,
							DataPoints: []metricdata.DataPoint[float64]{
								{
									Attributes: attribute.NewSet(),
									StartTime:  now.Add(-1 * time.Minute),
									Time:       now,
									Value:      3,
								},
								{
									Attributes: attribute.NewSet(attribute.String("foo", "bar")),
									StartTime:  now.Add(-1 * time.Minute),
									Time:       now,
									Value:      3,
								},
							},
						},
					},
					{
						Name:        "histogram",
						Description: "Histogram's description",
						Unit:        unit.Bytes,
						Data: metricdata.Histogram{
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.HistogramDataPoint{
								{
									Attributes: attribute.NewSet(),
									StartTime:  now.Add(-1 * time.Minute),
									Time:       now,
									Count:      6,
									Bounds: []float64{
										1.0,
										5.0,
									},
									BucketCounts: []uint64{
										1, 2, 3,
									},
									Min: nil,
									Max: nil,
									Sum: 200,
								},
							},
						},
					},
				},
			},
			{
				Scope: instrumentation.Scope{
					Name:      "Scope2",
					Version:   "v2",
					SchemaURL: "aDifferentExample",
				},
				Metrics: []metricdata.Metrics{
					{
						Name:        "gauge",
						Description: "Gauge's description",
						Unit:        unit.Dimensionless,
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(),
									StartTime:  now.Add(-1 * time.Minute),
									Time:       now,
									Value:      3,
								},
							},
						},
					},
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	exp, err := New(
		WithWriter(buf),
		WithPrettyPrint(),
	)
	require.NoError(t, err)

	err = exp.Export(context.Background(), exampleData)

	require.NoError(t, err)
	assert.Equal(t, expectedOutput, buf.String())
}

var expectedOutput = `{
	"Resource": [
		{
			"Key": "resource-foo",
			"Value": {
				"Type": "STRING",
				"Value": "bar"
			}
		}
	],
	"ScopeMetrics": [
		{
			"Scope": {
				"Name": "scope1",
				"Version": "v1",
				"SchemaURL": "anotherExample"
			},
			"Metrics": [
				{
					"Name": "gauge",
					"Description": "Gauge's description",
					"Unit": "1",
					"Data": {
						"DataPoints": [
							{
								"Attributes": [],
								"StartTime": "2019-05-16T19:04:27Z",
								"Time": "2019-05-16T19:05:27Z",
								"Value": 3
							}
						]
					}
				},
				{
					"Name": "sum",
					"Description": "Sum's description",
					"Unit": "By",
					"Data": {
						"DataPoints": [
							{
								"Attributes": [],
								"StartTime": "2019-05-16T19:04:27Z",
								"Time": "2019-05-16T19:05:27Z",
								"Value": 3
							},
							{
								"Attributes": [
									{
										"Key": "foo",
										"Value": {
											"Type": "STRING",
											"Value": "bar"
										}
									}
								],
								"StartTime": "2019-05-16T19:04:27Z",
								"Time": "2019-05-16T19:05:27Z",
								"Value": 3
							}
						],
						"Temporality": 2,
						"IsMonotonic": false
					}
				},
				{
					"Name": "histogram",
					"Description": "Histogram's description",
					"Unit": "By",
					"Data": {
						"DataPoints": [
							{
								"Attributes": [],
								"StartTime": "2019-05-16T19:04:27Z",
								"Time": "2019-05-16T19:05:27Z",
								"Count": 6,
								"Bounds": [
									1,
									5
								],
								"BucketCounts": [
									1,
									2,
									3
								],
								"Min": null,
								"Max": null,
								"Sum": 200
							}
						],
						"Temporality": 1
					}
				}
			]
		},
		{
			"Scope": {
				"Name": "Scope2",
				"Version": "v2",
				"SchemaURL": "aDifferentExample"
			},
			"Metrics": [
				{
					"Name": "gauge",
					"Description": "Gauge's description",
					"Unit": "1",
					"Data": {
						"DataPoints": [
							{
								"Attributes": [],
								"StartTime": "2019-05-16T19:04:27Z",
								"Time": "2019-05-16T19:05:27Z",
								"Value": 3
							}
						]
					}
				}
			]
		}
	]
}
`