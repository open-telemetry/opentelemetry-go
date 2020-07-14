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

	colmetricpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/collector/metrics/v1"
	commonpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/common/v1"
	metricpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/metrics/v1"
	resourcepb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/resource/v1"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/test"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
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
	mKind    metric.Kind
	nKind    metric.NumberKind
	resource *resource.Resource
	opts     []metric.InstrumentOption
	labels   []kv.KeyValue
}

var (
	baseKeyValues = []kv.KeyValue{kv.String("host", "test.com")}
	cpuKey        = kv.Key("CPU")

	testInstA = resource.New(kv.String("instance", "tester-a"))
	testInstB = resource.New(kv.String("instance", "tester-b"))

	md = &metricpb.MetricDescriptor{
		Name: "int64-count",
		Type: metricpb.MetricDescriptor_INT64,
	}

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
		[]record{
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
				nil,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
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
								MetricDescriptor: md,
								Int64DataPoints: []*metricpb.Int64DataPoint{
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
	)
}

func TestValuerecorderMetricGroupingExport(t *testing.T) {
	r := record{
		"valuerecorder",
		metric.ValueRecorderKind,
		metric.Int64NumberKind,
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
							MetricDescriptor: &metricpb.MetricDescriptor{
								Name: "valuerecorder",
								Type: metricpb.MetricDescriptor_SUMMARY,
							},
							SummaryDataPoints: []*metricpb.SummaryDataPoint{
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
									Count: 2,
									Sum:   11,
									PercentileValues: []*metricpb.SummaryDataPoint_ValueAtPercentile{
										{
											Percentile: 0.0,
											Value:      1.0,
										},
										{
											Percentile: 100.0,
											Value:      10.0,
										},
									},
									StartTimeUnixNano: startTime(),
									TimeUnixNano:      pointTime(),
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
									Count: 2,
									Sum:   11,
									PercentileValues: []*metricpb.SummaryDataPoint_ValueAtPercentile{
										{
											Percentile: 0.0,
											Value:      1.0,
										},
										{
											Percentile: 100.0,
											Value:      10.0,
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
	}
	runMetricExportTests(t, []record{r, r}, expected)
	//changing the number kind should make no difference.
	r.nKind = metric.Float64NumberKind
	runMetricExportTests(t, []record{r, r}, expected)
}

func TestCountInt64MetricGroupingExport(t *testing.T) {
	r := record{
		"int64-count",
		metric.CounterKind,
		metric.Int64NumberKind,
		nil,
		nil,
		append(baseKeyValues, cpuKey.Int(1)),
	}
	runMetricExportTests(
		t,
		[]record{r, r},
		[]metricpb.ResourceMetrics{
			{
				Resource: nil,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								MetricDescriptor: md,
								Int64DataPoints: []*metricpb.Int64DataPoint{
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
	)
}

func TestCountFloat64MetricGroupingExport(t *testing.T) {
	r := record{
		"float64-count",
		metric.CounterKind,
		metric.Float64NumberKind,
		nil,
		nil,
		append(baseKeyValues, cpuKey.Int(1)),
	}
	runMetricExportTests(
		t,
		[]record{r, r},
		[]metricpb.ResourceMetrics{
			{
				Resource: nil,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								MetricDescriptor: &metricpb.MetricDescriptor{
									Name: "float64-count",
									Type: metricpb.MetricDescriptor_DOUBLE,
								},
								DoubleDataPoints: []*metricpb.DoubleDataPoint{
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
	)
}

func TestResourceMetricGroupingExport(t *testing.T) {
	runMetricExportTests(
		t,
		[]record{
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
				testInstA,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
				testInstA,
				nil,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
				testInstA,
				nil,
				append(baseKeyValues, cpuKey.Int(2)),
			},
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
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
								MetricDescriptor: md,
								Int64DataPoints: []*metricpb.Int64DataPoint{
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
				Resource: testerBResource,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								MetricDescriptor: md,
								Int64DataPoints: []*metricpb.Int64DataPoint{
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
		[]record{
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
				testInstA,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
				testInstA,
				countingLib2,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
				testInstA,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
				testInstA,
				countingLib1,
				append(baseKeyValues, cpuKey.Int(2)),
			},
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
				testInstA,
				summingLib,
				append(baseKeyValues, cpuKey.Int(1)),
			},
			{
				"int64-count",
				metric.CounterKind,
				metric.Int64NumberKind,
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
								MetricDescriptor: md,
								Int64DataPoints: []*metricpb.Int64DataPoint{
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
					{
						InstrumentationLibrary: &commonpb.InstrumentationLibrary{
							Name:    "counting-lib",
							Version: "v2",
						},
						Metrics: []*metricpb.Metric{
							{
								MetricDescriptor: md,
								Int64DataPoints: []*metricpb.Int64DataPoint{
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
					{
						InstrumentationLibrary: &commonpb.InstrumentationLibrary{
							Name: "summing-lib",
						},
						Metrics: []*metricpb.Metric{
							{
								MetricDescriptor: md,
								Int64DataPoints: []*metricpb.Int64DataPoint{
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
				Resource: testerBResource,
				InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
					{
						InstrumentationLibrary: &commonpb.InstrumentationLibrary{
							Name:    "counting-lib",
							Version: "v1",
						},
						Metrics: []*metricpb.Metric{
							{
								MetricDescriptor: md,
								Int64DataPoints: []*metricpb.Int64DataPoint{
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
	)
}

// What works single-threaded should work multi-threaded
func runMetricExportTests(t *testing.T, rs []record, expected []metricpb.ResourceMetrics) {
	t.Run("1 goroutine", func(t *testing.T) {
		runMetricExportTest(t, NewUnstartedExporter(WorkerCount(1)), rs, expected)
	})
	t.Run("20 goroutines", func(t *testing.T) {
		runMetricExportTest(t, NewUnstartedExporter(WorkerCount(20)), rs, expected)
	})
}

func runMetricExportTest(t *testing.T, exp *Exporter, rs []record, expected []metricpb.ResourceMetrics) {
	msc := &metricsServiceClientStub{}
	exp.metricExporter = msc
	exp.started = true

	recs := map[label.Distinct][]metricsdk.Record{}
	resources := map[label.Distinct]*resource.Resource{}
	for _, r := range rs {
		desc := metric.NewDescriptor(r.name, r.mKind, r.nKind, r.opts...)
		labs := label.NewSet(r.labels...)

		var agg, ckpt metricsdk.Aggregator
		switch r.mKind {
		case metric.CounterKind:
			agg, ckpt = test.Unslice2(sum.New(2))
		default:
			agg, ckpt = test.Unslice2(minmaxsumcount.New(2, &desc))
		}

		ctx := context.Background()
		switch r.nKind {
		case metric.Int64NumberKind:
			require.NoError(t, agg.Update(ctx, metric.NewInt64Number(1), &desc))
			require.NoError(t, agg.Update(ctx, metric.NewInt64Number(10), &desc))
		case metric.Float64NumberKind:
			require.NoError(t, agg.Update(ctx, metric.NewFloat64Number(1), &desc))
			require.NoError(t, agg.Update(ctx, metric.NewFloat64Number(10), &desc))
		default:
			t.Fatalf("invalid number kind: %v", r.nKind)
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
				assert.Equal(t, expected.GetMetricDescriptor(), g[i].GetMetricDescriptor())
				// Compare each list directly because there is no order
				// guarantee with the concurrent processing design of the exporter
				// and ElementsMatch does not apply to contained slices.
				assert.ElementsMatch(t, expected.GetInt64DataPoints(), g[i].GetInt64DataPoints())
				assert.ElementsMatch(t, expected.GetDoubleDataPoints(), g[i].GetDoubleDataPoints())
				assert.ElementsMatch(t, expected.GetHistogramDataPoints(), g[i].GetHistogramDataPoints())
				assert.ElementsMatch(t, expected.GetSummaryDataPoints(), g[i].GetSummaryDataPoints())
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
