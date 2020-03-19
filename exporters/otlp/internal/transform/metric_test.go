// Copyright 2020, OpenTelemetry Authors
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

package transform

import (
	"context"
	"testing"

	commonpb "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/unit"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	sumAgg "go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

func TestStringKeyValues(t *testing.T) {
	tests := []struct {
		kvs      []core.KeyValue
		expected []*commonpb.StringKeyValue
	}{
		{
			[]core.KeyValue{},
			[]*commonpb.StringKeyValue{},
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
			[]*commonpb.StringKeyValue{
				{Key: "true", Value: "true"},
				{Key: "one", Value: "1"},
				{Key: "two", Value: "2"},
				{Key: "three", Value: "3"},
				{Key: "four", Value: "4"},
				{Key: "five", Value: "5"},
				{Key: "six", Value: "6"},
				{Key: "seven", Value: "7"},
				{Key: "eight", Value: "8"},
				{Key: "the", Value: "final word"},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, stringKeyValues(test.kvs))
	}
}

func TestMinMaxSumCountValue(t *testing.T) {
	mmsc := minmaxsumcount.New(&metric.Descriptor{})
	assert.NoError(t, mmsc.Update(context.Background(), 1, &metric.Descriptor{}))
	assert.NoError(t, mmsc.Update(context.Background(), 10, &metric.Descriptor{}))

	// Prior to checkpointing ErrNoData should be returned.
	_, _, _, _, err := minMaxSumCountValues(mmsc)
	assert.EqualError(t, err, aggregator.ErrNoData.Error())

	// Checkpoint to set non-zero values
	mmsc.Checkpoint(context.Background(), &metric.Descriptor{})
	min, max, sum, count, err := minMaxSumCountValues(mmsc)
	if assert.NoError(t, err) {
		assert.Equal(t, min, core.NewInt64Number(1))
		assert.Equal(t, max, core.NewInt64Number(10))
		assert.Equal(t, sum, core.NewInt64Number(11))
		assert.Equal(t, count, int64(2))
	}
}

func TestMinMaxSumCountMetricDescriptor(t *testing.T) {
	tests := []struct {
		name        string
		metricKind  metric.Kind
		keys        []core.Key
		description string
		unit        unit.Unit
		numberKind  core.NumberKind
		labels      []core.KeyValue
		expected    *metricpb.MetricDescriptor
	}{
		{
			"mmsc-test-a",
			metric.MeasureKind,
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
				Labels:      []*commonpb.StringKeyValue{},
			},
		},
		{
			"mmsc-test-b",
			metric.CounterKind, // This shouldn't change anything.
			[]core.Key{"test"}, // This shouldn't change anything.
			"test-b-description",
			unit.Bytes,
			core.Float64NumberKind, // This shouldn't change anything.
			[]core.KeyValue{core.Key("A").String("1")},
			&metricpb.MetricDescriptor{
				Name:        "mmsc-test-b",
				Description: "test-b-description",
				Unit:        "By",
				Type:        metricpb.MetricDescriptor_SUMMARY,
				Labels:      []*commonpb.StringKeyValue{{Key: "A", Value: "1"}},
			},
		},
	}

	ctx := context.Background()
	mmsc := minmaxsumcount.New(&metric.Descriptor{})
	if !assert.NoError(t, mmsc.Update(ctx, 1, &metric.Descriptor{})) {
		return
	}
	mmsc.Checkpoint(ctx, &metric.Descriptor{})
	for _, test := range tests {
		desc := metric.NewDescriptor(test.name, test.metricKind, test.numberKind,
			metric.WithKeys(test.keys...),
			metric.WithDescription(test.description),
			metric.WithUnit(test.unit))
		labels := export.NewLabels(test.labels, "", nil)
		got, err := minMaxSumCount(&desc, labels, mmsc)
		if assert.NoError(t, err) {
			assert.Equal(t, test.expected, got.MetricDescriptor)
		}
	}
}

func TestMinMaxSumCountDatapoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.MeasureKind, core.Int64NumberKind)
	labels := export.NewLabels([]core.KeyValue{}, "", nil)
	mmsc := minmaxsumcount.New(&desc)
	assert.NoError(t, mmsc.Update(context.Background(), 1, &desc))
	assert.NoError(t, mmsc.Update(context.Background(), 10, &desc))
	mmsc.Checkpoint(context.Background(), &desc)
	expected := []*metricpb.SummaryDataPoint{
		{
			Count: 2,
			Sum:   11,
			PercentileValues: []*metricpb.SummaryDataPoint_ValueAtPercentile{
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
	}
	m, err := minMaxSumCount(&desc, labels, mmsc)
	if assert.NoError(t, err) {
		assert.Equal(t, []*metricpb.Int64DataPoint(nil), m.Int64DataPoints)
		assert.Equal(t, []*metricpb.DoubleDataPoint(nil), m.DoubleDataPoints)
		assert.Equal(t, []*metricpb.HistogramDataPoint(nil), m.HistogramDataPoints)
		assert.Equal(t, expected, m.SummaryDataPoints)
	}
}

func TestMinMaxSumCountPropagatesErrors(t *testing.T) {
	// ErrNoData should be returned by both the Min and Max values of
	// a MinMaxSumCount Aggregator. Use this fact to check the error is
	// correctly returned.
	mmsc := minmaxsumcount.New(&metric.Descriptor{})
	_, _, _, _, err := minMaxSumCountValues(mmsc)
	assert.Error(t, err)
	assert.Equal(t, aggregator.ErrNoData, err)
}

func TestSumMetricDescriptor(t *testing.T) {
	tests := []struct {
		name        string
		metricKind  metric.Kind
		keys        []core.Key
		description string
		unit        unit.Unit
		numberKind  core.NumberKind
		labels      []core.KeyValue
		expected    *metricpb.MetricDescriptor
	}{
		{
			"sum-test-a",
			metric.CounterKind,
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
				Labels:      []*commonpb.StringKeyValue{},
			},
		},
		{
			"sum-test-b",
			metric.MeasureKind, // This shouldn't change anything.
			[]core.Key{"test"}, // This shouldn't change anything.
			"test-b-description",
			unit.Milliseconds,
			core.Float64NumberKind,
			[]core.KeyValue{core.Key("A").String("1")},
			&metricpb.MetricDescriptor{
				Name:        "sum-test-b",
				Description: "test-b-description",
				Unit:        "ms",
				Type:        metricpb.MetricDescriptor_COUNTER_DOUBLE,
				Labels:      []*commonpb.StringKeyValue{{Key: "A", Value: "1"}},
			},
		},
	}

	for _, test := range tests {
		desc := metric.NewDescriptor(test.name, test.metricKind, test.numberKind,
			metric.WithKeys(test.keys...),
			metric.WithDescription(test.description),
			metric.WithUnit(test.unit),
		)
		labels := export.NewLabels(test.labels, "", nil)
		got, err := sum(&desc, labels, sumAgg.New())
		if assert.NoError(t, err) {
			assert.Equal(t, test.expected, got.MetricDescriptor)
		}
	}
}

func TestSumInt64DataPoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.MeasureKind, core.Int64NumberKind)
	labels := export.NewLabels([]core.KeyValue{}, "", nil)
	s := sumAgg.New()
	assert.NoError(t, s.Update(context.Background(), core.Number(1), &desc))
	s.Checkpoint(context.Background(), &desc)
	if m, err := sum(&desc, labels, s); assert.NoError(t, err) {
		assert.Equal(t, []*metricpb.Int64DataPoint{{Value: 1}}, m.Int64DataPoints)
		assert.Equal(t, []*metricpb.DoubleDataPoint(nil), m.DoubleDataPoints)
		assert.Equal(t, []*metricpb.HistogramDataPoint(nil), m.HistogramDataPoints)
		assert.Equal(t, []*metricpb.SummaryDataPoint(nil), m.SummaryDataPoints)
	}
}

func TestSumFloat64DataPoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.MeasureKind, core.Float64NumberKind)
	labels := export.NewLabels([]core.KeyValue{}, "", nil)
	s := sumAgg.New()
	assert.NoError(t, s.Update(context.Background(), core.NewFloat64Number(1), &desc))
	s.Checkpoint(context.Background(), &desc)
	if m, err := sum(&desc, labels, s); assert.NoError(t, err) {
		assert.Equal(t, []*metricpb.Int64DataPoint(nil), m.Int64DataPoints)
		assert.Equal(t, []*metricpb.DoubleDataPoint{{Value: 1}}, m.DoubleDataPoints)
		assert.Equal(t, []*metricpb.HistogramDataPoint(nil), m.HistogramDataPoints)
		assert.Equal(t, []*metricpb.SummaryDataPoint(nil), m.SummaryDataPoints)
	}
}
