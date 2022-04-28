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
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
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

type testRecord struct {
	name  string
	iKind sdkapi.InstrumentKind
	nKind number.Kind
	attrs []attribute.KeyValue

	meterName string
	meterOpts []metric.MeterOption
}

func record(
	name string,
	iKind sdkapi.InstrumentKind,
	nKind number.Kind,
	attrs []attribute.KeyValue,
	meterName string,
	meterOpts ...metric.MeterOption) testRecord {
	return testRecord{
		name:      name,
		iKind:     iKind,
		nKind:     nKind,
		attrs:     attrs,
		meterName: meterName,
		meterOpts: meterOpts,
	}
}

var (
	baseKeyValues = []attribute.KeyValue{attribute.String("host", "test.com")}
	cpuKey        = attribute.Key("CPU")

	testHistogramBoundaries = []float64{2.0, 4.0, 8.0}

	cpu1Attrs = []*commonpb.KeyValue{
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
	cpu2Attrs = []*commonpb.KeyValue{
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

const (
	// Most of this test uses an empty instrumentation library name.
	testLibName = ""
)

func TestNoGroupingExport(t *testing.T) {
	runMetricExportTests(
		t,
		nil,
		resource.Empty(),
		[]testRecord{
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(1)),
				testLibName,
			),
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(2)),
				testLibName,
			),
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
												Attributes:        cpu1Attrs,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu2Attrs,
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
	r := record(
		"int64-histogram",
		sdkapi.HistogramInstrumentKind,
		number.Int64Kind,
		append(baseKeyValues, cpuKey.Int(1)),
		testLibName,
	)
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
											Attributes:        cpu1Attrs,
											StartTimeUnixNano: startTime(),
											TimeUnixNano:      pointTime(),
											Count:             2,
											Sum:               &sum,
											ExplicitBounds:    testHistogramBoundaries,
											BucketCounts:      []uint64{1, 0, 0, 1},
										},
										{
											Attributes:        cpu1Attrs,
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
	runMetricExportTests(t, nil, resource.Empty(), []testRecord{r, r}, expected)
}

func TestHistogramFloat64MetricGroupingExport(t *testing.T) {
	r := record(
		"float64-histogram",
		sdkapi.HistogramInstrumentKind,
		number.Float64Kind,
		append(baseKeyValues, cpuKey.Int(1)),
		testLibName,
	)
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
											Attributes:        cpu1Attrs,
											StartTimeUnixNano: startTime(),
											TimeUnixNano:      pointTime(),
											Count:             2,
											Sum:               &sum,
											ExplicitBounds:    testHistogramBoundaries,
											BucketCounts:      []uint64{1, 0, 0, 1},
										},
										{
											Attributes:        cpu1Attrs,
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
	runMetricExportTests(t, nil, resource.Empty(), []testRecord{r, r}, expected)
}

func TestCountInt64MetricGroupingExport(t *testing.T) {
	r := record(
		"int64-count",
		sdkapi.CounterInstrumentKind,
		number.Int64Kind,
		append(baseKeyValues, cpuKey.Int(1)),
		testLibName,
	)
	runMetricExportTests(
		t,
		nil,
		resource.Empty(),
		[]testRecord{r, r},
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
												Attributes:        cpu1Attrs,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Attrs,
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
	r := record(
		"float64-count",
		sdkapi.CounterInstrumentKind,
		number.Float64Kind,
		append(baseKeyValues, cpuKey.Int(1)),
		testLibName,
	)
	runMetricExportTests(
		t,
		nil,
		resource.Empty(),
		[]testRecord{r, r},
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
												Attributes:        cpu1Attrs,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsDouble{AsDouble: 11.0},
												Attributes:        cpu1Attrs,
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
		testerAResource,
		[]testRecord{
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(1)),
				testLibName,
			),
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(1)),
				testLibName,
			),
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(2)),
				testLibName,
			),
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(1)),
				testLibName,
			),
		},
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
												Attributes:        cpu1Attrs,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Attrs,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu2Attrs,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Attrs,
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
	version1 := metric.WithInstrumentationVersion("v1")
	version2 := metric.WithInstrumentationVersion("v2")
	specialSchema := metric.WithSchemaURL("schurl")
	summingLib := "summing-lib"
	countingLib := "counting-lib"
	runMetricExportTests(
		t,
		nil,
		testerAResource,
		[]testRecord{
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(1)),
				countingLib,
				version1,
			),
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(1)),
				countingLib,
				version2,
			),
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(1)),
				countingLib,
				version1,
			),
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(2)),
				countingLib,
				version1,
			),
			record(
				"int64-count",
				sdkapi.CounterInstrumentKind,
				number.Int64Kind,
				append(baseKeyValues, cpuKey.Int(1)),
				summingLib,
				specialSchema,
			),
		},
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
												Attributes:        cpu1Attrs,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu1Attrs,
												StartTimeUnixNano: startTime(),
												TimeUnixNano:      pointTime(),
											},
											{
												Value:             &metricpb.NumberDataPoint_AsInt{AsInt: 11},
												Attributes:        cpu2Attrs,
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
												Attributes:        cpu1Attrs,
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
												Attributes:        cpu1Attrs,
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
		instrumentKind sdkapi.InstrumentKind
		aggTemporality metricpb.AggregationTemporality
		monotonic      bool
	}

	for _, k := range []testcase{
		{"counter", sdkapi.CounterInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA, true},
		{"updowncounter", sdkapi.UpDownCounterInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA, false},
		{"counterobserver", sdkapi.CounterObserverInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE, true},
		{"updowncounterobserver", sdkapi.UpDownCounterObserverInstrumentKind, metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE, false},
	} {
		t.Run(k.name, func(t *testing.T) {
			runMetricExportTests(
				t,
				[]otlpmetric.Option{
					otlpmetric.WithMetricAggregationTemporalitySelector(
						aggregation.StatelessTemporalitySelector(),
					),
				},
				testerAResource,
				[]testRecord{
					record(
						"instrument",
						k.instrumentKind,
						number.Int64Kind,
						append(baseKeyValues, cpuKey.Int(1)),
						testLibName,
					),
				},
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
														Attributes:        cpu1Attrs,
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

func runMetricExportTests(t *testing.T, opts []otlpmetric.Option, res *resource.Resource, records []testRecord, expected []*metricpb.ResourceMetrics) {
	exp, driver := newExporter(t, opts...)

	libraryRecs := map[instrumentation.Library][]export.Record{}
	for _, r := range records {
		lcopy := make([]attribute.KeyValue, len(r.attrs))
		copy(lcopy, r.attrs)
		desc := metrictest.NewDescriptor(r.name, r.iKind, r.nKind)
		labs := attribute.NewSet(lcopy...)

		var agg, ckpt aggregator.Aggregator
		if r.iKind.Adding() {
			sums := sum.New(2)
			agg, ckpt = &sums[0], &sums[1]
		} else {
			histos := histogram.New(2, &desc, histogram.WithExplicitBoundaries(testHistogramBoundaries))
			agg, ckpt = &histos[0], &histos[1]
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

		meterCfg := metric.NewMeterConfig(r.meterOpts...)
		lib := instrumentation.Library{
			Name:      r.meterName,
			Version:   meterCfg.InstrumentationVersion(),
			SchemaURL: meterCfg.SchemaURL(),
		}
		libraryRecs[lib] = append(libraryRecs[lib], export.NewRecord(&desc, &labs, ckpt.Aggregation(), intervalStart, intervalEnd))
	}
	assert.NoError(t, exp.Export(context.Background(), res, processortest.MultiInstrumentationLibraryReader(libraryRecs)))

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

	for _, test := range []struct {
		records []export.Record
		want    []*metricpb.ResourceMetrics
	}{
		{
			[]export.Record(nil),
			[]*metricpb.ResourceMetrics(nil),
		},
		{
			[]export.Record{},
			[]*metricpb.ResourceMetrics(nil),
		},
	} {
		driver.Reset()
		require.NoError(t, exp.Export(context.Background(), resource.Empty(), processortest.MultiInstrumentationLibraryReader(map[instrumentation.Library][]export.Record{
			{
				Name: testLibName,
			}: test.records,
		})))
		assert.Equal(t, test.want, driver.rm)
	}
}
