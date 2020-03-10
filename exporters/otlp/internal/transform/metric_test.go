// Copyright 2020, OpenTelemetry Authors //
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

package transform

import (
	"context"
	"testing"

	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/unit"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
)

func TestKeysValues(t *testing.T) {
	tests := []struct {
		kvs          []core.KeyValue
		expectedKeys []string
		expectedVals []string
	}{
		{
			[]core.KeyValue{},
			[]string{},
			[]string{},
		},
		{
			[]core.KeyValue{
				core.Key("true").Bool(true),
				core.Key("one").Int64(1),
				core.Key("two").Uint64(2),
				core.Key("three").Float64(3),
				core.Key("four").Int32(4),
				core.Key("five").Uint32(5),
				core.Key("six").Float32(6),
				core.Key("seven").Int(7),
				core.Key("eight").Uint(8),
				core.Key("the").String("final word"),
			},
			[]string{"true", "one", "two", "three", "four", "five", "six", "seven", "eight", "the"},
			[]string{"true", "1", "2", "3", "4", "5", "6", "7", "8", "final word"},
		},
	}

	for _, test := range tests {
		assert.Equal(t, keys(test.kvs), test.expectedKeys)
		assert.Equal(t, values(test.kvs), test.expectedVals)
	}
}

func TestMinMaxSumCountValue(t *testing.T) {
	mmsc := minmaxsumcount.New(&metricsdk.Descriptor{})
	mmsc.Update(context.Background(), 1, &metricsdk.Descriptor{})
	mmsc.Update(context.Background(), 10, &metricsdk.Descriptor{})

	// Prior to checkpointing everything should be zero.
	min, max, sum, count, err := minMaxSumCountValues(mmsc)
	assert.Nil(t, err)
	assert.Equal(t, min, core.NewInt64Number(0))
	assert.Equal(t, max, core.NewInt64Number(0))
	assert.Equal(t, sum, core.NewInt64Number(0))
	assert.Equal(t, count, int64(0))

	// Checkpoint to set non-zero values
	mmsc.Checkpoint(context.Background(), &metricsdk.Descriptor{})
	min, max, sum, count, err = minMaxSumCountValues(mmsc)
	assert.Nil(t, err)
	assert.Equal(t, min, core.NewInt64Number(1))
	assert.Equal(t, max, core.NewInt64Number(10))
	assert.Equal(t, sum, core.NewInt64Number(11))
	assert.Equal(t, count, int64(2))
}

func TestMinMaxSumCountMetricDescriptor(t *testing.T) {
	tests := []struct {
		name        string
		metricKind  metricsdk.Kind
		keys        []core.Key
		description string
		unit        unit.Unit
		numberKind  core.NumberKind
		labels      []core.KeyValue
		expected    *metricpb.MetricDescriptor
	}{
		{
			"mmsc-test-a",
			metricsdk.MeasureKind,
			[]core.Key{},
			"test-a-description",
			unit.Dimensionless,
			core.Int64NumberKind,
			[]core.KeyValue{},
			&metricpb.MetricDescriptor{
				Name:        "mmsc-test-a",
				Description: "test-a-description",
				Unit:        "1",
				Type:        metricpb.MetricDescriptor_SUMMARY,
				LabelKeys:   []string{},
			},
		},
		{
			"mmsc-test-b",
			metricsdk.CounterKind, // This shouldn't change anything.
			[]core.Key{"test"},    // This shouldn't change anything.
			"test-b-description",
			unit.Bytes,
			core.Float64NumberKind, // This shouldn't change anything.
			[]core.KeyValue{core.Key("A").String("1")},
			&metricpb.MetricDescriptor{
				Name:        "mmsc-test-b",
				Description: "test-b-description",
				Unit:        "By",
				Type:        metricpb.MetricDescriptor_SUMMARY,
				LabelKeys:   []string{"A"},
			},
		},
	}

	for _, test := range tests {
		desc := metricsdk.NewDescriptor(test.name, test.metricKind, test.keys, test.description, test.unit, test.numberKind, false)
		labels := metricsdk.NewLabels(test.labels, "", nil)
		mmsc := minmaxsumcount.New(&metricsdk.Descriptor{})
		got, err := minMaxSumCount(desc, labels, mmsc)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, got.MetricDescriptor)
	}
}

func TestMinMaxSumCountDatapoints(t *testing.T) {
	desc := metricsdk.NewDescriptor("", metricsdk.MeasureKind, []core.Key{}, "", unit.Dimensionless, core.Int64NumberKind, false)
	labels := metricsdk.NewLabels([]core.KeyValue{}, "", nil)
	mmsc := minmaxsumcount.New(&metricsdk.Descriptor{})

	// test zero values.
	m, err := minMaxSumCount(desc, labels, mmsc)
	assert.Nil(t, err)
	expected := []*metricpb.SummaryDataPoint{
		{
			LabelValues: []string{},
			Value: &metricpb.SummaryValue{
				Count: 0,
				Sum:   0,
				PercentileValues: []*metricpb.SummaryValue_ValueAtPercentile{
					{
						Percentile: 0.0,
						Value:      0,
					},
					{
						Percentile: 100.0,
						Value:      0,
					},
				},
			},
		},
	}
	assert.Equal(t, []*metricpb.Int64DataPoint(nil), m.Int64Datapoints)
	assert.Equal(t, []*metricpb.DoubleDataPoint(nil), m.DoubleDatapoints)
	assert.Equal(t, []*metricpb.HistogramDataPoint(nil), m.HistogramDatapoints)
	assert.Equal(t, expected, m.SummaryDatapoints)

	// test with non-zero values.
	mmsc.Update(context.Background(), 1, &metricsdk.Descriptor{})
	mmsc.Update(context.Background(), 10, &metricsdk.Descriptor{})
	mmsc.Checkpoint(context.Background(), &metricsdk.Descriptor{})
	m, err = minMaxSumCount(desc, labels, mmsc)
	assert.Nil(t, err)
	expected = []*metricpb.SummaryDataPoint{
		{
			LabelValues: []string{},
			Value: &metricpb.SummaryValue{
				Count: 2,
				Sum:   11,
				PercentileValues: []*metricpb.SummaryValue_ValueAtPercentile{
					{
						Percentile: 0.0,
						Value:      1,
					},
					{
						Percentile: 100.0,
						Value:      10,
					},
				},
			},
		},
	}
	assert.Equal(t, []*metricpb.Int64DataPoint(nil), m.Int64Datapoints)
	assert.Equal(t, []*metricpb.DoubleDataPoint(nil), m.DoubleDatapoints)
	assert.Equal(t, []*metricpb.HistogramDataPoint(nil), m.HistogramDatapoints)
	assert.Equal(t, expected, m.SummaryDatapoints)
}

func TestSumMetricDescriptor(t *testing.T) {
	tests := []struct {
		name        string
		metricKind  metricsdk.Kind
		keys        []core.Key
		description string
		unit        unit.Unit
		numberKind  core.NumberKind
		labels      []core.KeyValue
		expected    *metricpb.MetricDescriptor
	}{
		{
			"sum-test-a",
			metricsdk.CounterKind,
			[]core.Key{},
			"test-a-description",
			unit.Dimensionless,
			core.Int64NumberKind,
			[]core.KeyValue{},
			&metricpb.MetricDescriptor{
				Name:        "sum-test-a",
				Description: "test-a-description",
				Unit:        "1",
				Type:        metricpb.MetricDescriptor_COUNTER_INT64,
				LabelKeys:   []string{},
			},
		},
		{
			"sum-test-b",
			metricsdk.MeasureKind, // This shouldn't change anything.
			[]core.Key{"test"},    // This shouldn't change anything.
			"test-b-description",
			unit.Milliseconds,
			core.Float64NumberKind,
			[]core.KeyValue{core.Key("A").String("1")},
			&metricpb.MetricDescriptor{
				Name:        "sum-test-b",
				Description: "test-b-description",
				Unit:        "ms",
				Type:        metricpb.MetricDescriptor_COUNTER_DOUBLE,
				LabelKeys:   []string{"A"},
			},
		},
	}

	for _, test := range tests {
		desc := metricsdk.NewDescriptor(test.name, test.metricKind, test.keys, test.description, test.unit, test.numberKind, false)
		labels := metricsdk.NewLabels(test.labels, "", nil)
		got, err := sum(desc, labels, counter.New())
		assert.Nil(t, err)
		assert.Equal(t, test.expected, got.MetricDescriptor)
	}
}

func TestSumInt64Datapoints(t *testing.T) {
	desc := metricsdk.NewDescriptor("", metricsdk.MeasureKind, []core.Key{}, "", unit.Dimensionless, core.Int64NumberKind, false)
	labels := metricsdk.NewLabels([]core.KeyValue{}, "", nil)
	s := counter.New()

	// test zero values.
	m, err := sum(desc, labels, s)
	assert.Nil(t, err)
	expected := []*metricpb.Int64DataPoint{
		{
			LabelValues: []string{},
			Value:       &metricpb.Int64Value{Value: 0},
		},
	}
	assert.Equal(t, expected, m.Int64Datapoints)
	assert.Equal(t, []*metricpb.DoubleDataPoint(nil), m.DoubleDatapoints)
	assert.Equal(t, []*metricpb.HistogramDataPoint(nil), m.HistogramDatapoints)
	assert.Equal(t, []*metricpb.SummaryDataPoint(nil), m.SummaryDatapoints)

	// test with non-zero values.
	s.Update(context.Background(), core.Number(1), &metricsdk.Descriptor{})
	s.Checkpoint(context.Background(), &metricsdk.Descriptor{})
	m, err = sum(desc, labels, s)
	assert.Nil(t, err)
	expected = []*metricpb.Int64DataPoint{
		{
			LabelValues: []string{},
			Value:       &metricpb.Int64Value{Value: 1},
		},
	}
	assert.Equal(t, expected, m.Int64Datapoints)
	assert.Equal(t, []*metricpb.DoubleDataPoint(nil), m.DoubleDatapoints)
	assert.Equal(t, []*metricpb.HistogramDataPoint(nil), m.HistogramDatapoints)
	assert.Equal(t, []*metricpb.SummaryDataPoint(nil), m.SummaryDatapoints)
}
