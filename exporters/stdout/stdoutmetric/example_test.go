// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutmetric_test

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	// Sat Jan 01 2000 00:00:00 GMT+0000.
	now = time.Date(2000, time.January, 0o1, 0, 0, 0, 0, time.FixedZone("GMT", 0))

	res = resource.NewSchemaless(
		semconv.ServiceName("stdoutmetric-example"),
	)

	mockData = metricdata.ResourceMetrics{
		Resource: res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{Name: "example", Version: "0.0.1"},
				Metrics: []metricdata.Metrics{
					{
						Name:        "requests",
						Description: "Number of requests received",
						Unit:        "1",
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
						Name:        "system.cpu.time",
						Description: "Accumulated CPU time spent",
						Unit:        "s",
						Data: metricdata.Sum[float64]{
							IsMonotonic: true,
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.DataPoint[float64]{
								{
									Attributes: attribute.NewSet(attribute.String("state", "user")),
									StartTime:  now,
									Time:       now.Add(1 * time.Second),
									Value:      0.5,
								},
							},
						},
					},
					{
						Name:        "requests.size",
						Description: "Size of received requests",
						Unit:        "kb",
						Data: metricdata.Histogram[int64]{
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.HistogramDataPoint[int64]{
								{
									Attributes:   attribute.NewSet(attribute.String("server", "central")),
									StartTime:    now,
									Time:         now.Add(1 * time.Second),
									Count:        10,
									Bounds:       []float64{1, 5, 10},
									BucketCounts: []uint64{1, 3, 6, 0},
									Sum:          128,
									Min:          metricdata.NewExtrema[int64](3),
									Max:          metricdata.NewExtrema[int64](30),
								},
							},
						},
					},
					{
						Name:        "latency",
						Description: "Time spend processing received requests",
						Unit:        "ms",
						Data: metricdata.Histogram[float64]{
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.HistogramDataPoint[float64]{
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
						Name:        "system.memory.usage",
						Description: "Memory usage",
						Unit:        "By",
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(attribute.String("state", "used")),
									Time:       now.Add(1 * time.Second),
									Value:      100,
								},
							},
						},
					},
					{
						Name:        "temperature",
						Description: "CPU global temperature",
						Unit:        "cel(1 K)",
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
)

func Example() {
	// Print with a JSON encoder that indents with two spaces.
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	exp, err := stdoutmetric.New(
		stdoutmetric.WithEncoder(enc),
		stdoutmetric.WithoutTimestamps(),
	)
	if err != nil {
		panic(err)
	}

	// Register the exporter with an SDK via a periodic reader.
	sdk := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exp)),
	)

	ctx := context.Background()
	// This is where the sdk would be used to create a Meter and from that
	// instruments that would make measurements of your code. To simulate that
	// behavior, call export directly with mocked data.
	_ = exp.Export(ctx, &mockData)

	// Ensure the periodic reader is cleaned up by shutting down the sdk.
	_ = sdk.Shutdown(ctx)

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
	//         "Version": "0.0.1",
	//         "SchemaURL": "",
	//         "Attributes": null
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
	//                 "StartTime": "0001-01-01T00:00:00Z",
	//                 "Time": "0001-01-01T00:00:00Z",
	//                 "Value": 5
	//               }
	//             ],
	//             "Temporality": "DeltaTemporality",
	//             "IsMonotonic": true
	//           }
	//         },
	//         {
	//           "Name": "system.cpu.time",
	//           "Description": "Accumulated CPU time spent",
	//           "Unit": "s",
	//           "Data": {
	//             "DataPoints": [
	//               {
	//                 "Attributes": [
	//                   {
	//                     "Key": "state",
	//                     "Value": {
	//                       "Type": "STRING",
	//                       "Value": "user"
	//                     }
	//                   }
	//                 ],
	//                 "StartTime": "0001-01-01T00:00:00Z",
	//                 "Time": "0001-01-01T00:00:00Z",
	//                 "Value": 0.5
	//               }
	//             ],
	//             "Temporality": "CumulativeTemporality",
	//             "IsMonotonic": true
	//           }
	//         },
	//         {
	//           "Name": "requests.size",
	//           "Description": "Size of received requests",
	//           "Unit": "kb",
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
	//                 "Time": "0001-01-01T00:00:00Z",
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
	//                 "Min": 3,
	//                 "Max": 30,
	//                 "Sum": 128
	//               }
	//             ],
	//             "Temporality": "DeltaTemporality"
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
	//                 "StartTime": "0001-01-01T00:00:00Z",
	//                 "Time": "0001-01-01T00:00:00Z",
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
	//                 "Min": null,
	//                 "Max": null,
	//                 "Sum": 57
	//               }
	//             ],
	//             "Temporality": "DeltaTemporality"
	//           }
	//         },
	//         {
	//           "Name": "system.memory.usage",
	//           "Description": "Memory usage",
	//           "Unit": "By",
	//           "Data": {
	//             "DataPoints": [
	//               {
	//                 "Attributes": [
	//                   {
	//                     "Key": "state",
	//                     "Value": {
	//                       "Type": "STRING",
	//                       "Value": "used"
	//                     }
	//                   }
	//                 ],
	//                 "StartTime": "0001-01-01T00:00:00Z",
	//                 "Time": "0001-01-01T00:00:00Z",
	//                 "Value": 100
	//               }
	//             ]
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
	//                 "Time": "0001-01-01T00:00:00Z",
	//                 "Value": 32.4
	//               }
	//             ]
	//           }
	//         }
	//       ]
	//     }
	//   ]
	// }
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
	//   "ScopeMetrics": null
	// }
}
