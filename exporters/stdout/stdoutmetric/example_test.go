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

package stdoutmetric_test

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

// Sat Jan 01 2000 00:00:00 GMT+0000.
var now = time.Unix(946684800, 0)

var mockData = metricdata.ResourceMetrics{
	Resource: resource.NewSchemaless(semconv.ServiceNameKey.String("stdoutmetric-example")),
	ScopeMetrics: []metricdata.ScopeMetrics{
		{
			Scope: instrumentation.Scope{
				Name:    "example",
				Version: "v0.0.1",
			},
			Metrics: []metricdata.Metrics{
				{
					Name:        "requests",
					Description: "Number of requests received",
					Unit:        unit.Dimensionless,
					Data: metricdata.Sum[int64]{
						IsMonotonic: true,
						Temporality: metricdata.DeltaTemporality,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Attributes: attribute.NewSet(attribute.String("server", "central")),
								StartTime:  now,
								Time:       now.Add(1 * time.Second),
								Value:      5,
							},
						},
					},
				},
				{
					Name:        "latency",
					Description: "Time spend processing received requests",
					Unit:        unit.Milliseconds,
					Data: metricdata.Histogram{
						Temporality: metricdata.DeltaTemporality,
						DataPoints: []metricdata.HistogramDataPoint{
							{
								Attributes:   attribute.NewSet(attribute.String("server", "central")),
								StartTime:    now,
								Time:         now.Add(1 * time.Second),
								Count:        10,
								Bounds:       []float64{1, 5, 10},
								BucketCounts: []uint64{1, 3, 6, 0},
								Sum:          57,
							},
						},
					},
				},
				{
					Name:        "temperature",
					Description: "CPU global temperature",
					Unit:        unit.Unit("cel(1 K)"),
					Data: metricdata.Gauge[float64]{
						DataPoints: []metricdata.DataPoint[float64]{
							{
								Attributes: attribute.NewSet(attribute.String("server", "central")),
								Time:       now.Add(1 * time.Second),
								Value:      32.4,
							},
						},
					},
				},
			},
		},
	},
}

func Example() {
	// Print with a JSON encoder that indents with two spaces.
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	exp, err := stdoutmetric.New(stdoutmetric.WithEncoder(enc))
	if err != nil {
		panic(err)
	}
	// FIXME: The exporter should be registered with an SDK and the output
	// should come from that.
	exp.Export(context.Background(), mockData)

	// Output:
	// {
	//   "Resource": [
	//     {
	//       "Key": "service.name",
	//       "Value": {
	//         "Type": "STRING",
	//         "Value": "stdoutmetric-example"
	//       }
	//     }
	//   ],
	//   "ScopeMetrics": [
	//     {
	//       "Scope": {
	//         "Name": "example",
	//         "Version": "v0.0.1"
	//       },
	//       "Metrics": [
	//         {
	//           "Name": "requests",
	//           "Description": "Number of requests received",
	//           "Unit": "1",
	//           "Data": {
	//             "DataPoints": [
	//               {
	//                 "Attributes": [
	//                   {
	//                     "Key": "server",
	//                     "Value": {
	//                       "Type": "STRING",
	//                       "Value": "central"
	//                     }
	//                   }
	//                 ],
	//                 "StartTime": "1999-12-31T16:00:00-08:00",
	//                 "Time": "1999-12-31T16:00:01-08:00",
	//                 "Value": 5
	//               }
	//             ],
	//             "Temporality": "DeltaTemporality",
	//             "IsMonotonic": true
	//           }
	//         },
	//         {
	//           "Name": "latency",
	//           "Description": "Time spend processing received requests",
	//           "Unit": "ms",
	//           "Data": {
	//             "DataPoints": [
	//               {
	//                 "Attributes": [
	//                   {
	//                     "Key": "server",
	//                     "Value": {
	//                       "Type": "STRING",
	//                       "Value": "central"
	//                     }
	//                   }
	//                 ],
	//                 "StartTime": "1999-12-31T16:00:00-08:00",
	//                 "Time": "1999-12-31T16:00:01-08:00",
	//                 "Count": 10,
	//                 "Bounds": [
	//                   1,
	//                   5,
	//                   10
	//                 ],
	//                 "BucketCounts": [
	//                   1,
	//                   3,
	//                   6,
	//                   0
	//                 ],
	//                 "Sum": 57
	//               }
	//             ],
	//             "Temporality": "DeltaTemporality"
	//           }
	//         },
	//         {
	//           "Name": "temperature",
	//           "Description": "CPU global temperature",
	//           "Unit": "cel(1 K)",
	//           "Data": {
	//             "DataPoints": [
	//               {
	//                 "Attributes": [
	//                   {
	//                     "Key": "server",
	//                     "Value": {
	//                       "Type": "STRING",
	//                       "Value": "central"
	//                     }
	//                   }
	//                 ],
	//                 "StartTime": "0001-01-01T00:00:00Z",
	//                 "Time": "1999-12-31T16:00:01-08:00",
	//                 "Value": 32.4
	//               }
	//             ]
	//           }
	//         }
	//       ]
	//     }
	//   ]
	// }
}
