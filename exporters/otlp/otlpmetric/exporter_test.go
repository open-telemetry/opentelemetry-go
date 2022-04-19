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

package otlpmetric_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/metrictransform"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/number"

	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/resource"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

var (
	// Timestamps used in this test:

	intervalStart = time.Now()
	intervalEnd   = intervalStart.Add(time.Hour)
)

type stubClient struct {
	rm []*metricpb.ResourceMetrics
}

func (m *stubClient) Start(ctx context.Context) error {
	return nil
}

func (m *stubClient) Stop(ctx context.Context) error {
	return nil
}

func (m *stubClient) UploadMetrics(ctx context.Context, protoMetrics *metricpb.ResourceMetrics) error {
	m.rm = append(m.rm, protoMetrics)
	return nil
}

var _ otlpmetric.Client = (*stubClient)(nil)

func (m *stubClient) Reset() {
	m.rm = nil
}

func newExporter(t *testing.T, opts ...otlpmetric.Option) (*otlpmetric.Exporter, *stubClient) {
	client := &stubClient{}
	exp, _ := otlpmetric.New(context.Background(), client, opts...)
	return exp, client
}

func startTime() uint64 {
	return uint64(intervalStart.UnixNano())
}

func pointTime() uint64 {
	return uint64(intervalEnd.UnixNano())
}

var (
	baseKeyValues = []attribute.KeyValue{attribute.String("host", "test.com")}
	cpuKey        = attribute.Key("CPU")

	testHistogramBoundaries = []float64{2.0, 4.0, 8.0}

	cpu1Labels = []*commonpb.KeyValue{
		{
			Key: "CPU",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_IntValue{
					IntValue: 1,
				},
			},
		},
		{
			Key: "host",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_StringValue{
					StringValue: "test.com",
				},
			},
		},
	}
	cpu2Labels = []*commonpb.KeyValue{
		{
			Key: "CPU",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_IntValue{
					IntValue: 2,
				},
			},
		},
		{
			Key: "host",
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_StringValue{
					StringValue: "test.com",
				},
			},
		},
	}

	testerAResource   = resource.NewSchemaless(attribute.String("instance", "tester-a"))
	testerAResourcePb = metrictransform.Resource(testerAResource)
)

func TestNoGroupingExport(t *testing.T) {
	runMetricExportTests(
		t,
		nil,
		reader.Metrics{
			Resource: resource.Empty(),
			Scopes: []reader.Scope{
				{
					Instruments: []reader.Instrument{
						{
							Descriptor:  sdkinstrument.NewDescriptor("int64-count", sdkinstrument.CounterKind, number.Int64Kind, "", ""),
							Temporality: aggregation.CumulativeTemporality,

							Points: []reader.Point{
								{
									Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
									Aggregation: sum.NewInt64Monotonic(11),
									Start:       intervalStart,
									End:         intervalEnd,
								},
								{
									Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(2))...),
									Aggregation: sum.NewInt64Monotonic(11),
									Start:       intervalStart,
									End:         intervalEnd,
								},
							},
						},
					},
				},
			},
		},
		[]*metricpb.ResourceMetrics{
			{
				Resource: nil,
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Name: "int64-count",
								Data: &metricpb.Metric_Sum{
									Sum: &metricpb.Sum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.NumberDataPoint{
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu2Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)
}

func TestHistogramInt64MetricGroupingExport(t *testing.T) {
	metrics := reader.Metrics{
		Resource: resource.Empty(),
		Scopes: []reader.Scope{
			{
				Instruments: []reader.Instrument{
					{
						Descriptor:  sdkinstrument.NewDescriptor("int64-histogram", sdkinstrument.HistogramKind, number.Int64Kind, "", ""),
						Temporality: aggregation.CumulativeTemporality,

						Points: []reader.Point{
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: histogram.NewInt64(testHistogramBoundaries, int64(1), int64(10)),
								Start:       intervalStart,
								End:         intervalEnd,
							},
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: histogram.NewInt64(testHistogramBoundaries, int64(1), int64(10)),
								Start:       intervalStart,
								End:         intervalEnd,
							},
						},
					},
				},
			},
		},
	}

	sum := 11.0
	expected := []*metricpb.ResourceMetrics{
		{
			Resource: nil,
			ScopeMetrics: []*metricpb.ScopeMetrics{
				{
					Metrics: []*metricpb.Metric{
						{
							Name: "int64-histogram",
							Data: &metricpb.Metric_Histogram{
								Histogram: &metricpb.Histogram{
									AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
									DataPoints: []*metricpb.HistogramDataPoint{
										{
											Attributes:        cpu1Labels,
											StartTimeUnixNano: startTime(),
											TimeUnixNano:      pointTime(),
											Count:             2,
											Sum:               &sum,
											ExplicitBounds:    testHistogramBoundaries,
											BucketCounts:      []uint64{1, 0, 0, 1},
										},
										{
											Attributes:        cpu1Labels,
											Count:             2,
											Sum:               &sum,
											ExplicitBounds:    testHistogramBoundaries,
											BucketCounts:      []uint64{1, 0, 0, 1},
											StartTimeUnixNano: startTime(),
											TimeUnixNano:      pointTime(),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	runMetricExportTests(t, nil, metrics, expected)
}

func TestHistogramFloat64MetricGroupingExport(t *testing.T) {

	metrics := reader.Metrics{
		Resource: resource.Empty(),
		Scopes: []reader.Scope{
			{
				Instruments: []reader.Instrument{
					{
						Descriptor:  sdkinstrument.NewDescriptor("float64-histogram", sdkinstrument.HistogramKind, number.Float64Kind, "", ""),
						Temporality: aggregation.CumulativeTemporality,

						Points: []reader.Point{
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: histogram.NewFloat64(testHistogramBoundaries, 1, 10),
								Start:       intervalStart,
								End:         intervalEnd,
							},
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: histogram.NewFloat64(testHistogramBoundaries, 1, 10),
								Start:       intervalStart,
								End:         intervalEnd,
							},
						},
					},
				},
			},
		},
	}
	sum := 11.0
	expected := []*metricpb.ResourceMetrics{
		{
			Resource: nil,
			ScopeMetrics: []*metricpb.ScopeMetrics{
				{
					Metrics: []*metricpb.Metric{
						{
							Name: "float64-histogram",
							Data: &metricpb.Metric_Histogram{
								Histogram: &metricpb.Histogram{
									AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
									DataPoints: []*metricpb.HistogramDataPoint{
										{
											Attributes:        cpu1Labels,
											StartTimeUnixNano: startTime(),
											TimeUnixNano:      pointTime(),
											Count:             2,
											Sum:               &sum,
											ExplicitBounds:    testHistogramBoundaries,
											BucketCounts:      []uint64{1, 0, 0, 1},
										},
										{
											Attributes:        cpu1Labels,
											Count:             2,
											Sum:               &sum,
											ExplicitBounds:    testHistogramBoundaries,
											BucketCounts:      []uint64{1, 0, 0, 1},
											StartTimeUnixNano: startTime(),
											TimeUnixNano:      pointTime(),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	runMetricExportTests(t, nil, metrics, expected)
}

func TestCountInt64MetricGroupingExport(t *testing.T) {

	metrics := reader.Metrics{
		Resource: resource.Empty(),
		Scopes: []reader.Scope{
			{
				Instruments: []reader.Instrument{
					{
						Descriptor:  sdkinstrument.NewDescriptor("int64-count", sdkinstrument.CounterKind, number.Int64Kind, "", ""),
						Temporality: aggregation.CumulativeTemporality,

						Points: []reader.Point{
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
						},
					},
				},
			},
		},
	}

	runMetricExportTests(
		t,
		nil,
		metrics,
		[]*metricpb.ResourceMetrics{
			{
				Resource: nil,
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Name: "int64-count",
								Data: &metricpb.Metric_Sum{
									Sum: &metricpb.Sum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.NumberDataPoint{
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)
}

func TestCountFloat64MetricGroupingExport(t *testing.T) {

	metrics := reader.Metrics{
		Resource: resource.Empty(),
		Scopes: []reader.Scope{
			{
				Instruments: []reader.Instrument{
					{
						Descriptor:  sdkinstrument.NewDescriptor("float64-count", sdkinstrument.CounterKind, number.Float64Kind, "", ""),
						Temporality: aggregation.CumulativeTemporality,

						Points: []reader.Point{
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewFloat64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewFloat64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
						},
					},
				},
			},
		},
	}

	runMetricExportTests(
		t,
		nil,
		metrics,
		[]*metricpb.ResourceMetrics{
			{
				Resource: nil,
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Name: "float64-count",
								Data: &metricpb.Metric_Sum{
									Sum: &metricpb.Sum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.NumberDataPoint{
											{
												Value:             &metricpb.NumberDataPoint_AsDouble{AsDouble: 11.0},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsDouble{AsDouble: 11.0},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)
}

func TestResourceMetricGroupingExport(t *testing.T) {
	metrics := reader.Metrics{
		Resource: testerAResource,
		Scopes: []reader.Scope{
			{
				Instruments: []reader.Instrument{
					{
						Descriptor:  sdkinstrument.NewDescriptor("int64-count", sdkinstrument.CounterKind, number.Int64Kind, "", ""),
						Temporality: aggregation.CumulativeTemporality,

						Points: []reader.Point{
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(2))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
						},
					},
				},
			},
		},
	}

	runMetricExportTests(
		t,
		nil,
		metrics,
		[]*metricpb.ResourceMetrics{
			{
				Resource: testerAResourcePb,
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Name: "int64-count",
								Data: &metricpb.Metric_Sum{
									Sum: &metricpb.Sum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.NumberDataPoint{
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu2Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)
}

func TestResourceInstLibMetricGroupingExport(t *testing.T) {
	// version1 := metric.WithInstrumentationVersion("v1")
	// version2 := metric.WithInstrumentationVersion("v2")
	// specialSchema := metric.WithSchemaURL("schurl")
	// summingLib := "summing-lib"
	// countingLib := "counting-lib"
	//testerAResource,
	metrics := reader.Metrics{
		Resource: testerAResource,
		Scopes: []reader.Scope{
			{
				Library: instrumentation.Library{
					Name:    "counting-lib",
					Version: "v1",
				},
				Instruments: []reader.Instrument{
					{
						Descriptor:  sdkinstrument.NewDescriptor("int64-count", sdkinstrument.CounterKind, number.Int64Kind, "", ""),
						Temporality: aggregation.CumulativeTemporality,

						Points: []reader.Point{
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(2))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
						},
					},
				},
			},
			{
				Library: instrumentation.Library{
					Name:    "counting-lib",
					Version: "v2",
				},
				Instruments: []reader.Instrument{
					{
						Descriptor:  sdkinstrument.NewDescriptor("int64-count", sdkinstrument.CounterKind, number.Int64Kind, "", ""),
						Temporality: aggregation.CumulativeTemporality,

						Points: []reader.Point{
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
						},
					},
				},
			},
			{
				Library: instrumentation.Library{
					Name:      "summing-lib",
					SchemaURL: "schurl",
				},
				Instruments: []reader.Instrument{
					{
						Descriptor:  sdkinstrument.NewDescriptor("int64-count", sdkinstrument.CounterKind, number.Int64Kind, "", ""),
						Temporality: aggregation.CumulativeTemporality,

						Points: []reader.Point{
							{
								Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
								Aggregation: sum.NewInt64Monotonic(11),
								Start:       intervalStart,
								End:         intervalEnd,
							},
						},
					},
				},
			},
		},
	}

	runMetricExportTests(
		t,
		nil,
		metrics,
		[]*metricpb.ResourceMetrics{
			{
				Resource: testerAResourcePb,
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Scope: &commonpb.InstrumentationScope{
							Name:    "counting-lib",
							Version: "v1",
						},
						Metrics: []*metricpb.Metric{
							{
								Name: "int64-count",
								Data: &metricpb.Metric_Sum{
									Sum: &metricpb.Sum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.NumberDataPoint{
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu2Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
										},
									},
								},
							},
						},
					},
					{
						Scope: &commonpb.InstrumentationScope{
							Name:    "counting-lib",
							Version: "v2",
						},
						Metrics: []*metricpb.Metric{
							{
								Name: "int64-count",
								Data: &metricpb.Metric_Sum{
									Sum: &metricpb.Sum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.NumberDataPoint{
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
										},
									},
								},
							},
						},
					},
					{
						Scope: &commonpb.InstrumentationScope{
							Name: "summing-lib",
						},
						SchemaUrl: "schurl",
						Metrics: []*metricpb.Metric{
							{
								Name: "int64-count",
								Data: &metricpb.Metric_Sum{
									Sum: &metricpb.Sum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.NumberDataPoint{
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)
}

func TestStatelessAggregationTemporality(t *testing.T) {
	type testcase struct {
		name           string
		kind           sdkinstrument.Kind
		temporality    aggregation.Temporality
		aggregation    aggregation.Aggregation
		aggTemporality metricpb.AggregationTemporality
		monotonic      bool
	}

	for _, k := range []testcase{
		{
			name:           "counter",
			kind:           sdkinstrument.CounterKind,
			temporality:    aggregation.DeltaTemporality,
			aggregation:    sum.NewInt64Monotonic(11),
			aggTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
			monotonic:      true},
		{
			name:           "updowncounter",
			kind:           sdkinstrument.UpDownCounterKind,
			temporality:    aggregation.DeltaTemporality,
			aggregation:    sum.NewInt64NonMonotonic(11),
			aggTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
			monotonic:      false},
		{
			name:           "counterobserver",
			kind:           sdkinstrument.CounterObserverKind,
			temporality:    aggregation.CumulativeTemporality,
			aggregation:    sum.NewInt64Monotonic(11),
			aggTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
			monotonic:      true},
		{
			name:           "updowncounterobserver",
			kind:           sdkinstrument.UpDownCounterObserverKind,
			temporality:    aggregation.CumulativeTemporality,
			aggregation:    sum.NewInt64NonMonotonic(11),
			aggTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
			monotonic:      false},
	} {
		t.Run(k.name, func(t *testing.T) {

			metrics := reader.Metrics{
				Resource: testerAResource,
				Scopes: []reader.Scope{
					{
						Instruments: []reader.Instrument{
							{
								Descriptor:  sdkinstrument.NewDescriptor("instrument", k.kind, number.Int64Kind, "", ""),
								Temporality: k.temporality,

								Points: []reader.Point{
									{
										Attributes:  attribute.NewSet(append(baseKeyValues, cpuKey.Int(1))...),
										Aggregation: k.aggregation,
										Start:       intervalStart,
										End:         intervalEnd,
									},
								},
							},
						},
					},
				},
			}

			runMetricExportTests(
				t,
				[]otlpmetric.Option{
					otlpmetric.WithMetricAggregationTemporalitySelector(
						aggregation.UndefinedTemporality,
					),
				},
				metrics,
				[]*metricpb.ResourceMetrics{
					{
						Resource: testerAResourcePb,
						ScopeMetrics: []*metricpb.ScopeMetrics{
							{
								Metrics: []*metricpb.Metric{
									{
										Name: "instrument",
										Data: &metricpb.Metric_Sum{
											Sum: &metricpb.Sum{
												IsMonotonic:            k.monotonic,
												AggregationTemporality: k.aggTemporality,
												DataPoints: []*metricpb.NumberDataPoint{
													{
														Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
														Attributes:        cpu1Labels,
														StartTimeUnixNano: startTime(),
														TimeUnixNano:      pointTime(),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			)
		})
	}
}

func runMetricExportTests(t *testing.T, opts []otlpmetric.Option, metrics reader.Metrics, expected []*metricpb.ResourceMetrics) {
	t.Helper()
	exp, driver := newExporter(t, opts...)

	err := exp.Export(context.Background(), metrics)
	assert.NoError(t, err)

	// assert.ElementsMatch does not equate nested slices of different order,
	// therefore this requires the top level slice to be broken down.
	// Build a map of Resource/Scope pairs to Metrics, from that validate the
	// metric elements match for all expected pairs. Finally, make we saw all
	// expected pairs.
	keyFor := func(sm *metricpb.ScopeMetrics) string {
		return fmt.Sprintf("%s/%s/%s", sm.GetScope().GetName(), sm.GetScope().GetVersion(), sm.GetSchemaUrl())
	}
	got := map[string][]*metricpb.Metric{}
	for _, rm := range driver.rm {
		for _, sm := range rm.ScopeMetrics {
			k := keyFor(sm)
			got[k] = append(got[k], sm.GetMetrics()...)
		}
	}

	seen := map[string]struct{}{}
	for _, rm := range expected {
		for _, sm := range rm.ScopeMetrics {
			k := keyFor(sm)
			seen[k] = struct{}{}
			g, ok := got[k]
			if !ok {
				t.Errorf("missing metrics for:\n\tInstrumentationScope: %q\n", k)
				continue
			}
			if !assert.Len(t, g, len(sm.GetMetrics())) {
				continue
			}
			for i, expected := range sm.GetMetrics() {
				t.Log("Expected: ", expected)
				t.Log("Actual:   ", g[i])
				assert.Equal(t, "", cmp.Diff(expected, g[i], protocmp.Transform()))
			}
		}
	}
	for k := range got {
		if _, ok := seen[k]; !ok {
			t.Errorf("did not expect metrics for:\n\tInstrumentationScope: %s\n", k)
		}
	}
}

func TestEmptyMetricExport(t *testing.T) {
	exp, driver := newExporter(t)

	// {
	records := reader.Metrics{}
	want := []*metricpb.ResourceMetrics{
		{},
	}

	driver.Reset()
	require.NoError(t, exp.Export(context.Background(), records))
	assert.Equal(t, want, driver.rm)

}
