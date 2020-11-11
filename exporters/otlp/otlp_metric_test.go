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

package otlp

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	colmetricpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/metrics/v1"
	commonpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/common/v1"
	metricpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/metrics/v1"
	resourcepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/resource/v1"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric/number"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/export/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/resource"

	"google.golang.org/grpc"
)

var (
	// Timestamps used in this test:

	intervalStart = time.Now()
	intervalEnd   = intervalStart.Add(time.Hour)
)

func startTime() uint64 {
	return uint64(intervalStart.UnixNano())
}

func pointTime() uint64 {
	return uint64(intervalEnd.UnixNano())
}

type metricsServiceClientStub struct {
	rm []metricpb.ResourceMetrics
}

func (m *metricsServiceClientStub) Export(ctx context.Context, in *colmetricpb.ExportMetricsServiceRequest, opts ...grpc.CallOption) (*colmetricpb.ExportMetricsServiceResponse, error) {
	for _, rm := range in.GetResourceMetrics() {
		if rm == nil {
			continue
		}
		m.rm = append(m.rm, *rm)
	}
	return &colmetricpb.ExportMetricsServiceResponse{}, nil
}

func (m *metricsServiceClientStub) ResourceMetrics() []metricpb.ResourceMetrics {
	return m.rm
}

func (m *metricsServiceClientStub) Reset() {
	m.rm = nil
}

type checkpointSet struct {
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

type record struct {
	name     string
	iKind    otel.InstrumentKind
	nKind    number.Kind
	resource *resource.Resource
	opts     []otel.InstrumentOption
	labels   []label.KeyValue
}

var (
	baseKeyValues = []label.KeyValue{label.String("host", "test.com")}
	cpuKey        = label.Key("CPU")

	testInstA = resource.NewWithAttributes(label.String("instance", "tester-a"))
	testInstB = resource.NewWithAttributes(label.String("instance", "tester-b"))

	testHistogramBoundaries = []float64{2.0, 4.0, 8.0}

	cpu1Labels = []*commonpb.StringKeyValue{
		{
			Key:   "CPU",
			Value: "1",
		},
		{
			Key:   "host",
			Value: "test.com",
		},
	}
	cpu2Labels = []*commonpb.StringKeyValue{
		{
			Key:   "CPU",
			Value: "2",
		},
		{
			Key:   "host",
			Value: "test.com",
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
				otel.CounterInstrumentKind,
				number.Int64Kind,
				nil,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				otel.CounterInstrumentKind,
				number.Int64Kind,
				nil,
				nil,
				append(baseKeyValues, cpuKey.Int(2)),
			},
		},
		[]metricpb.ResourceMetrics{
			{
				Resource: nil,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Name: "int64-count",
								Data: &metricpb.Metric_IntSum{
									IntSum: &metricpb.IntSum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.IntDataPoint{
											{
												Value:             11,
												Labels:            cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             11,
												Labels:            cpu2Labels,
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
		otel.ValueRecorderInstrumentKind,
		number.Int64Kind,
		nil,
		nil,
		append(baseKeyValues, cpuKey.Int(1)),
	}
	expected := []metricpb.ResourceMetrics{
		{
			Resource: nil,
			InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
				{
					Metrics: []*metricpb.Metric{
						{
							Name: "valuerecorder",
							Data: &metricpb.Metric_IntHistogram{
								IntHistogram: &metricpb.IntHistogram{
									AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
									DataPoints: []*metricpb.IntHistogramDataPoint{
										{
											Labels: []*commonpb.StringKeyValue{
												{
													Key:   "CPU",
													Value: "1",
												},
												{
													Key:   "host",
													Value: "test.com",
												},
											},
											StartTimeUnixNano: startTime(),
											TimeUnixNano:      pointTime(),
											Count:             2,
											Sum:               11,
											ExplicitBounds:    testHistogramBoundaries,
											BucketCounts:      []uint64{1, 0, 0, 1},
										},
										{
											Labels: []*commonpb.StringKeyValue{
												{
													Key:   "CPU",
													Value: "1",
												},
												{
													Key:   "host",
													Value: "test.com",
												},
											},
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
		otel.CounterInstrumentKind,
		number.Int64Kind,
		nil,
		nil,
		append(baseKeyValues, cpuKey.Int(1)),
	}
	runMetricExportTests(
		t,
		nil,
		[]record{r, r},
		[]metricpb.ResourceMetrics{
			{
				Resource: nil,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Name: "int64-count",
								Data: &metricpb.Metric_IntSum{
									IntSum: &metricpb.IntSum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.IntDataPoint{
											{
												Value:             11,
												Labels:            cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             11,
												Labels:            cpu1Labels,
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
		otel.CounterInstrumentKind,
		number.Float64Kind,
		nil,
		nil,
		append(baseKeyValues, cpuKey.Int(1)),
	}
	runMetricExportTests(
		t,
		nil,
		[]record{r, r},
		[]metricpb.ResourceMetrics{
			{
				Resource: nil,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Name: "float64-count",
								Data: &metricpb.Metric_DoubleSum{
									DoubleSum: &metricpb.DoubleSum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.DoubleDataPoint{
											{
												Value: 11,
												Labels: []*commonpb.StringKeyValue{
													{
														Key:   "CPU",
														Value: "1",
													},
													{
														Key:   "host",
														Value: "test.com",
													},
												},
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value: 11,
												Labels: []*commonpb.StringKeyValue{
													{
														Key:   "CPU",
														Value: "1",
													},
													{
														Key:   "host",
														Value: "test.com",
													},
												},
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
				otel.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				otel.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				otel.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				nil,
				append(baseKeyValues, cpuKey.Int(2)),
			},
			{
				"int64-count",
				otel.CounterInstrumentKind,
				number.Int64Kind,
				testInstB,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
		},
		[]metricpb.ResourceMetrics{
			{
				Resource: testerAResource,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Name: "int64-count",
								Data: &metricpb.Metric_IntSum{
									IntSum: &metricpb.IntSum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.IntDataPoint{
											{
												Value:             11,
												Labels:            cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             11,
												Labels:            cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             11,
												Labels:            cpu2Labels,
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
								Data: &metricpb.Metric_IntSum{
									IntSum: &metricpb.IntSum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.IntDataPoint{
											{
												Value:             11,
												Labels:            cpu1Labels,
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
	countingLib1 := []otel.InstrumentOption{
		otel.WithInstrumentationName("counting-lib"),
		otel.WithInstrumentationVersion("v1"),
	}
	countingLib2 := []otel.InstrumentOption{
		otel.WithInstrumentationName("counting-lib"),
		otel.WithInstrumentationVersion("v2"),
	}
	summingLib := []otel.InstrumentOption{
		otel.WithInstrumentationName("summing-lib"),
	}
	runMetricExportTests(
		t,
		nil,
		[]record{
			{
				"int64-count",
				otel.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				otel.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				countingLib2,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				otel.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				otel.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(2)),
			},
			{
				"int64-count",
				otel.CounterInstrumentKind,
				number.Int64Kind,
				testInstA,
				summingLib,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				otel.CounterInstrumentKind,
				number.Int64Kind,
				testInstB,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(1)),
			},
		},
		[]metricpb.ResourceMetrics{
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
								Data: &metricpb.Metric_IntSum{
									IntSum: &metricpb.IntSum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.IntDataPoint{
											{
												Value:             11,
												Labels:            cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             11,
												Labels:            cpu1Labels,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             11,
												Labels:            cpu2Labels,
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
								Data: &metricpb.Metric_IntSum{
									IntSum: &metricpb.IntSum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.IntDataPoint{
											{
												Value:             11,
												Labels:            cpu1Labels,
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
								Data: &metricpb.Metric_IntSum{
									IntSum: &metricpb.IntSum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.IntDataPoint{
											{
												Value:             11,
												Labels:            cpu1Labels,
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
								Data: &metricpb.Metric_IntSum{
									IntSum: &metricpb.IntSum{
										IsMonotonic:            true,
										AggregationTemporality: metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
										DataPoints: []*metricpb.IntDataPoint{
											{
												Value:             11,
												Labels:            cpu1Labels,
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
	type testcase struct {
		name           string
		instrumentKind otel.InstrumentKind
		aggTemporality metricpb.AggregationTemporality
		monotonic      bool
	}

	for _, k := range []testcase{
		{"counter", otel.CounterInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA, true},
		{"updowncounter", otel.UpDownCounterInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA, false},
		{"sumobserver", otel.SumObserverInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE, true},
		{"updownsumobserver", otel.UpDownSumObserverInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE, false},
	} {
		t.Run(k.name, func(t *testing.T) {
			runMetricExportTests(
				t,
				[]ExporterOption{
					WithMetricExportKindSelector(
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
				[]metricpb.ResourceMetrics{
					{
						Resource: testerAResource,
						InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
							{
								Metrics: []*metricpb.Metric{
									{
										Name: "instrument",
										Data: &metricpb.Metric_IntSum{
											IntSum: &metricpb.IntSum{
												IsMonotonic:            k.monotonic,
												AggregationTemporality: k.aggTemporality,
												DataPoints: []*metricpb.IntDataPoint{
													{
														Value:             11,
														Labels:            cpu1Labels,
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

// What works single-threaded should work multi-threaded
func runMetricExportTests(t *testing.T, opts []ExporterOption, rs []record, expected []metricpb.ResourceMetrics) {
	t.Run("1 goroutine", func(t *testing.T) {
		runMetricExportTest(t, NewUnstartedExporter(append(opts[:len(opts):len(opts)], WorkerCount(1))...), rs, expected)
	})
	t.Run("20 goroutines", func(t *testing.T) {
		runMetricExportTest(t, NewUnstartedExporter(append(opts[:len(opts):len(opts)], WorkerCount(20))...), rs, expected)
	})
}

func runMetricExportTest(t *testing.T, exp *Exporter, rs []record, expected []metricpb.ResourceMetrics) {
	msc := &metricsServiceClientStub{}
	exp.metricExporter = msc
	exp.started = true

	recs := map[label.Distinct][]metricsdk.Record{}
	resources := map[label.Distinct]*resource.Resource{}
	for _, r := range rs {
		lcopy := make([]label.KeyValue, len(r.labels))
		copy(lcopy, r.labels)
		desc := otel.NewDescriptor(r.name, r.iKind, r.nKind, r.opts...)
		labs := label.NewSet(lcopy...)

		var agg, ckpt metricsdk.Aggregator
		if r.iKind.Adding() {
			agg, ckpt = metrictest.Unslice2(sum.New(2))
		} else {
			agg, ckpt = metrictest.Unslice2(histogram.New(2, &desc, testHistogramBoundaries))
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
	for _, rm := range msc.ResourceMetrics() {
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
				case *metricpb.Metric_IntGauge:
					assert.ElementsMatch(t, expected.GetIntGauge().DataPoints, g[i].GetIntGauge().DataPoints)
				case *metricpb.Metric_IntHistogram:
					assert.Equal(t,
						expected.GetIntHistogram().GetAggregationTemporality(),
						g[i].GetIntHistogram().GetAggregationTemporality(),
					)
					assert.ElementsMatch(t, expected.GetIntHistogram().DataPoints, g[i].GetIntHistogram().DataPoints)
				case *metricpb.Metric_IntSum:
					assert.Equal(t,
						expected.GetIntSum().GetAggregationTemporality(),
						g[i].GetIntSum().GetAggregationTemporality(),
					)
					assert.Equal(t,
						expected.GetIntSum().GetIsMonotonic(),
						g[i].GetIntSum().GetIsMonotonic(),
					)
					assert.ElementsMatch(t, expected.GetIntSum().DataPoints, g[i].GetIntSum().DataPoints)
				case *metricpb.Metric_DoubleGauge:
					assert.ElementsMatch(t, expected.GetDoubleGauge().DataPoints, g[i].GetDoubleGauge().DataPoints)
				case *metricpb.Metric_DoubleHistogram:
					assert.Equal(t,
						expected.GetDoubleHistogram().GetAggregationTemporality(),
						g[i].GetDoubleHistogram().GetAggregationTemporality(),
					)
					assert.ElementsMatch(t, expected.GetDoubleHistogram().DataPoints, g[i].GetDoubleHistogram().DataPoints)
				case *metricpb.Metric_DoubleSum:
					assert.Equal(t,
						expected.GetDoubleSum().GetAggregationTemporality(),
						g[i].GetDoubleSum().GetAggregationTemporality(),
					)
					assert.Equal(t,
						expected.GetDoubleSum().GetIsMonotonic(),
						g[i].GetDoubleSum().GetIsMonotonic(),
					)
					assert.ElementsMatch(t, expected.GetDoubleSum().DataPoints, g[i].GetDoubleSum().DataPoints)
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
	msc := &metricsServiceClientStub{}
	exp := NewUnstartedExporter()
	exp.metricExporter = msc
	exp.started = true

	for _, test := range []struct {
		records []metricsdk.Record
		want    []metricpb.ResourceMetrics
	}{
		{
			[]metricsdk.Record(nil),
			[]metricpb.ResourceMetrics(nil),
		},
		{
			[]metricsdk.Record{},
			[]metricpb.ResourceMetrics(nil),
		},
	} {
		msc.Reset()
		require.NoError(t, exp.Export(context.Background(), &checkpointSet{records: test.records}))
		assert.Equal(t, test.want, msc.ResourceMetrics())
	}
}
