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
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/export/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
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

func (m *stubClient) UploadMetrics(ctx context.Context, protoMetrics []*metricpb.ResourceMetrics) error {
	m.rm = append(m.rm, protoMetrics...)
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

type checkpointSet struct { //nolint:govet // ignore 'fieldalignment' error
	sync.RWMutex
	records []metricsdk.Record
}

func (m *checkpointSet) ForEach(_ metricsdk.ExportKindSelector, fn func(metricsdk.Record) error) error {
	for _, r := range m.records {
		if err := fn(r); err != nil && err != aggregation.ErrNoData {
			return err
		}
	}
	return nil
}

type record struct { //nolint:govet // ignore 'fieldalignment' error
	name     string
	iKind    metric.InstrumentKind
	nKind    number.Kind
	resource *resource.Resource
	opts     []metric.InstrumentOption
	labels   []attribute.KeyValue
}

var (
	baseKeyValues = []attribute.KeyValue{attribute.String("host", "test.com")}
	cpuKey        = attribute.Key("CPU")

	testInstA = resource.NewSchemaless(attribute.String("instance", "tester-a"))
	testInstB = resource.NewSchemaless(attribute.String("instance", "tester-b"))

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

	testerAResource = &resourcepb.Resource{
		Attributes: []*commonpb.KeyValue{
			{
				Key: "instance",
				Value: &commonpb.AnyValue{
					Value: &commonpb.AnyValue_StringValue{
						StringValue: "tester-a",
					},
				},
			},
		},
	}
	testerBResource = &resourcepb.Resource{
		Attributes: []*commonpb.KeyValue{
			{
				Key: "instance",
				Value: &commonpb.AnyValue{
					Value: &commonpb.AnyValue_StringValue{
						StringValue: "tester-b",
					},
				},
			},
		},
	}
)

func TestNoGroupingExport(t *testing.T) {
	runMetricExportTests(
		t,
		nil,
		[]record{
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				nil,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				nil,
				nil,
				append(baseKeyValues, cpuKey.Int(2)),
			},
		},
		[]*metricpb.ResourceMetrics{
			{
				Resource: nil,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
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

func TestValuerecorderMetricGroupingExport(t *testing.T) {
	r := record{
		"valuerecorder",
		metric.ValueRecorderInstrumentKind,
		number.Int64Kind,
		nil,
		nil,
		append(baseKeyValues, cpuKey.Int(1)),
	}
	expected := []*metricpb.ResourceMetrics{
		{
			Resource: nil,
			InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
				{
					Metrics: []*metricpb.Metric{
						{
							Name: "valuerecorder",
							Data: &metricpb.Metric_Histogram{
								Histogram: &metricpb.Histogram{
									AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
									DataPoints: []*metricpb.HistogramDataPoint{
										{
											Attributes:        cpu1Labels,
											StartTimeUnixNano: startTime(),
											TimeUnixNano:      pointTime(),
											Count:             2,
											Sum:               11,
											ExplicitBounds:    testHistogramBoundaries,
											BucketCounts:      []uint64{1, 0, 0, 1},
										},
										{
											Attributes:        cpu1Labels,
											Count:             2,
											Sum:               11,
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
	runMetricExportTests(t, nil, []record{r, r}, expected)
}

func TestCountInt64MetricGroupingExport(t *testing.T) {
	r := record{
		"int64-count",
		metric.CounterInstrumentKind,
		number.Int64Kind,
		nil,
		nil,
		append(baseKeyValues, cpuKey.Int(1)),
	}
	runMetricExportTests(
		t,
		nil,
		[]record{r, r},
		[]*metricpb.ResourceMetrics{
			{
				Resource: nil,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
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
	r := record{
		"float64-count",
		metric.CounterInstrumentKind,
		number.Float64Kind,
		nil,
		nil,
		append(baseKeyValues, cpuKey.Int(1)),
	}
	runMetricExportTests(
		t,
		nil,
		[]record{r, r},
		[]*metricpb.ResourceMetrics{
			{
				Resource: nil,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
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
	runMetricExportTests(
		t,
		nil,
		[]record{
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				nil,
				append(baseKeyValues, cpuKey.Int(2)),
			},
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				testInstB,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
		},
		[]*metricpb.ResourceMetrics{
			{
				Resource: testerAResource,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
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
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Resource: testerBResource,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
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
	countingLib1 := []metric.InstrumentOption{
		metric.WithInstrumentationName("counting-lib"),
		metric.WithInstrumentationVersion("v1"),
	}
	countingLib2 := []metric.InstrumentOption{
		metric.WithInstrumentationName("counting-lib"),
		metric.WithInstrumentationVersion("v2"),
	}
	summingLib := []metric.InstrumentOption{
		metric.WithInstrumentationName("summing-lib"),
	}
	runMetricExportTests(
		t,
		nil,
		[]record{
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				countingLib2,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(2)),
			},
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				summingLib,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterInstrumentKind,
				number.Int64Kind,
				testInstB,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(1)),
			},
		},
		[]*metricpb.ResourceMetrics{
			{
				Resource: testerAResource,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
					{
						InstrumentationLibrary: &commonpb.InstrumentationLibrary{
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
						InstrumentationLibrary: &commonpb.InstrumentationLibrary{
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
						InstrumentationLibrary: &commonpb.InstrumentationLibrary{
							Name: "summing-lib",
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
				},
			},
			{
				Resource: testerBResource,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
					{
						InstrumentationLibrary: &commonpb.InstrumentationLibrary{
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

func TestStatelessExportKind(t *testing.T) {
	type testcase struct { //nolint:govet // ignore 'fieldalignment' error
		name           string
		instrumentKind metric.InstrumentKind
		aggTemporality metricpb.AggregationTemporality
		monotonic      bool
	}

	for _, k := range []testcase{
		{"counter", metric.CounterInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA, true},
		{"updowncounter", metric.UpDownCounterInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA, false},
		{"sumobserver", metric.SumObserverInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE, true},
		{"updownsumobserver", metric.UpDownSumObserverInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE, false},
	} {
		t.Run(k.name, func(t *testing.T) {
			runMetricExportTests(
				t,
				[]otlpmetric.Option{
					otlpmetric.WithMetricExportKindSelector(
						metricsdk.StatelessExportKindSelector(),
					),
				},
				[]record{
					{
						"instrument",
						k.instrumentKind,
						number.Int64Kind,
						testInstA,
						nil,
						append(baseKeyValues, cpuKey.Int(1)),
					},
				},
				[]*metricpb.ResourceMetrics{
					{
						Resource: testerAResource,
						InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
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

func runMetricExportTests(t *testing.T, opts []otlpmetric.Option, rs []record, expected []*metricpb.ResourceMetrics) {
	exp, driver := newExporter(t, opts...)

	recs := map[attribute.Distinct][]metricsdk.Record{}
	resources := map[attribute.Distinct]*resource.Resource{}
	for _, r := range rs {
		lcopy := make([]attribute.KeyValue, len(r.labels))
		copy(lcopy, r.labels)
		desc := metric.NewDescriptor(r.name, r.iKind, r.nKind, r.opts...)
		labs := attribute.NewSet(lcopy...)

		var agg, ckpt metricsdk.Aggregator
		if r.iKind.Adding() {
			agg, ckpt = metrictest.Unslice2(sum.New(2))
		} else {
			agg, ckpt = metrictest.Unslice2(histogram.New(2, &desc, histogram.WithExplicitBoundaries(testHistogramBoundaries)))
		}

		ctx := context.Background()
		if r.iKind.Synchronous() {
			// For synchronous instruments, perform two updates: 1 and 10
			switch r.nKind {
			case number.Int64Kind:
				require.NoError(t, agg.Update(ctx, number.NewInt64Number(1), &desc))
				require.NoError(t, agg.Update(ctx, number.NewInt64Number(10), &desc))
			case number.Float64Kind:
				require.NoError(t, agg.Update(ctx, number.NewFloat64Number(1), &desc))
				require.NoError(t, agg.Update(ctx, number.NewFloat64Number(10), &desc))
			default:
				t.Fatalf("invalid number kind: %v", r.nKind)
			}
		} else {
			// For asynchronous instruments, perform a single update: 11
			switch r.nKind {
			case number.Int64Kind:
				require.NoError(t, agg.Update(ctx, number.NewInt64Number(11), &desc))
			case number.Float64Kind:
				require.NoError(t, agg.Update(ctx, number.NewFloat64Number(11), &desc))
			default:
				t.Fatalf("invalid number kind: %v", r.nKind)
			}
		}
		require.NoError(t, agg.SynchronizedMove(ckpt, &desc))

		equiv := r.resource.Equivalent()
		resources[equiv] = r.resource
		recs[equiv] = append(recs[equiv], metricsdk.NewRecord(&desc, &labs, r.resource, ckpt.Aggregation(), intervalStart, intervalEnd))
	}
	for _, records := range recs {
		assert.NoError(t, exp.Export(context.Background(), &checkpointSet{records: records}))
	}

	// assert.ElementsMatch does not equate nested slices of different order,
	// therefore this requires the top level slice to be broken down.
	// Build a map of Resource/InstrumentationLibrary pairs to Metrics, from
	// that validate the metric elements match for all expected pairs. Finally,
	// make we saw all expected pairs.
	type key struct {
		resource, instrumentationLibrary string
	}
	got := map[key][]*metricpb.Metric{}
	for _, rm := range driver.rm {
		for _, ilm := range rm.InstrumentationLibraryMetrics {
			k := key{
				resource:               rm.GetResource().String(),
				instrumentationLibrary: ilm.GetInstrumentationLibrary().String(),
			}
			got[k] = ilm.GetMetrics()
		}
	}
	seen := map[key]struct{}{}
	for _, rm := range expected {
		for _, ilm := range rm.InstrumentationLibraryMetrics {
			k := key{
				resource:               rm.GetResource().String(),
				instrumentationLibrary: ilm.GetInstrumentationLibrary().String(),
			}
			seen[k] = struct{}{}
			g, ok := got[k]
			if !ok {
				t.Errorf("missing metrics for:\n\tResource: %s\n\tInstrumentationLibrary: %s\n", k.resource, k.instrumentationLibrary)
				continue
			}
			if !assert.Len(t, g, len(ilm.GetMetrics())) {
				continue
			}
			for i, expected := range ilm.GetMetrics() {
				assert.Equal(t, expected.Name, g[i].Name)
				assert.Equal(t, expected.Unit, g[i].Unit)
				assert.Equal(t, expected.Description, g[i].Description)
				switch g[i].Data.(type) {
				case *metricpb.Metric_Gauge:
					assert.ElementsMatch(t, expected.GetGauge().GetDataPoints(), g[i].GetGauge().GetDataPoints())
				case *metricpb.Metric_Sum:
					assert.Equal(t,
						expected.GetSum().GetAggregationTemporality(),
						g[i].GetSum().GetAggregationTemporality(),
					)
					assert.Equal(t,
						expected.GetSum().GetIsMonotonic(),
						g[i].GetSum().GetIsMonotonic(),
					)
					assert.ElementsMatch(t, expected.GetSum().GetDataPoints(), g[i].GetSum().GetDataPoints())
				case *metricpb.Metric_Histogram:
					assert.Equal(
						t,
						expected.GetHistogram().GetAggregationTemporality(),
						g[i].GetHistogram().GetAggregationTemporality(),
					)
					assert.ElementsMatch(t, expected.GetHistogram().GetDataPoints(), g[i].GetHistogram().GetDataPoints())
				case *metricpb.Metric_Summary:
					assert.ElementsMatch(t, expected.GetSummary().GetDataPoints(), g[i].GetSummary().GetDataPoints())
				default:
					assert.Failf(t, "unknown data type", g[i].Name)
				}
			}
		}
	}
	for k := range got {
		if _, ok := seen[k]; !ok {
			t.Errorf("did not expect metrics for:\n\tResource: %s\n\tInstrumentationLibrary: %s\n", k.resource, k.instrumentationLibrary)
		}
	}
}

func TestEmptyMetricExport(t *testing.T) {
	exp, driver := newExporter(t)

	for _, test := range []struct {
		records []metricsdk.Record
		want    []*metricpb.ResourceMetrics
	}{
		{
			[]metricsdk.Record(nil),
			[]*metricpb.ResourceMetrics(nil),
		},
		{
			[]metricsdk.Record{},
			[]*metricpb.ResourceMetrics(nil),
		},
	} {
		driver.Reset()
		require.NoError(t, exp.Export(context.Background(), &checkpointSet{records: test.records}))
		assert.Equal(t, test.want, driver.rm)
	}
}
