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

package transform

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/common/v1"
	metricpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/metrics/v1"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/unit"
	"go.opentelemetry.io/otel/exporters/metric/test"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	sumAgg "go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

var (
	// Timestamps used in this test:

	intervalStart = time.Now()
	intervalEnd   = intervalStart.Add(time.Hour)
)

func TestStringKeyValues(t *testing.T) {
	tests := []struct {
		kvs      []kv.KeyValue
		expected []*commonpb.StringKeyValue
	}{
		{
			nil,
			nil,
		},
		{
			[]kv.KeyValue{},
			nil,
		},
		{
			[]kv.KeyValue{
				kv.Bool("true", true),
				kv.Int64("one", 1),
				kv.Uint64("two", 2),
				kv.Float64("three", 3),
				kv.Int32("four", 4),
				kv.Uint32("five", 5),
				kv.Float32("six", 6),
				kv.Int("seven", 7),
				kv.Uint("eight", 8),
				kv.String("the", "final word"),
			},
			[]*commonpb.StringKeyValue{
				{Key: "eight", Value: "8"},
				{Key: "five", Value: "5"},
				{Key: "four", Value: "4"},
				{Key: "one", Value: "1"},
				{Key: "seven", Value: "7"},
				{Key: "six", Value: "6"},
				{Key: "the", Value: "final word"},
				{Key: "three", Value: "3"},
				{Key: "true", Value: "true"},
				{Key: "two", Value: "2"},
			},
		},
	}

	for _, test := range tests {
		labels := label.NewSet(test.kvs...)
		assert.Equal(t, test.expected, stringKeyValues(labels.Iter()))
	}
}

func TestMinMaxSumCountValue(t *testing.T) {
	mmsc, ckpt := test.Unslice2(minmaxsumcount.New(2, &metric.Descriptor{}))

	assert.NoError(t, mmsc.Update(context.Background(), 1, &metric.Descriptor{}))
	assert.NoError(t, mmsc.Update(context.Background(), 10, &metric.Descriptor{}))

	// Prior to checkpointing ErrNoData should be returned.
	_, _, _, _, err := minMaxSumCountValues(ckpt.(aggregation.MinMaxSumCount))
	assert.EqualError(t, err, aggregation.ErrNoData.Error())

	// Checkpoint to set non-zero values
	require.NoError(t, mmsc.SynchronizedMove(ckpt, &metric.Descriptor{}))
	min, max, sum, count, err := minMaxSumCountValues(ckpt.(aggregation.MinMaxSumCount))
	if assert.NoError(t, err) {
		assert.Equal(t, min, metric.NewInt64Number(1))
		assert.Equal(t, max, metric.NewInt64Number(10))
		assert.Equal(t, sum, metric.NewInt64Number(11))
		assert.Equal(t, count, int64(2))
	}
}

func TestMinMaxSumCountMetricDescriptor(t *testing.T) {
	tests := []struct {
		name        string
		metricKind  metric.Kind
		description string
		unit        unit.Unit
		numberKind  metric.NumberKind
		labels      []kv.KeyValue
		expected    *metricpb.MetricDescriptor
	}{
		{
			"mmsc-test-a",
			metric.ValueRecorderKind,
			"test-a-description",
			unit.Dimensionless,
			metric.Int64NumberKind,
			[]kv.KeyValue{},
			&metricpb.MetricDescriptor{
				Name:        "mmsc-test-a",
				Description: "test-a-description",
				Unit:        "1",
				Type:        metricpb.MetricDescriptor_SUMMARY,
			},
		},
		{
			"mmsc-test-b",
			metric.CounterKind, // This shouldn't change anything.
			"test-b-description",
			unit.Bytes,
			metric.Float64NumberKind, // This shouldn't change anything.
			[]kv.KeyValue{kv.String("A", "1")},
			&metricpb.MetricDescriptor{
				Name:        "mmsc-test-b",
				Description: "test-b-description",
				Unit:        "By",
				Type:        metricpb.MetricDescriptor_SUMMARY,
			},
		},
	}

	ctx := context.Background()
	mmsc, ckpt := test.Unslice2(minmaxsumcount.New(2, &metric.Descriptor{}))
	if !assert.NoError(t, mmsc.Update(ctx, 1, &metric.Descriptor{})) {
		return
	}
	require.NoError(t, mmsc.SynchronizedMove(ckpt, &metric.Descriptor{}))
	for _, test := range tests {
		desc := metric.NewDescriptor(test.name, test.metricKind, test.numberKind,
			metric.WithDescription(test.description),
			metric.WithUnit(test.unit))
		labels := label.NewSet(test.labels...)
		record := export.NewRecord(&desc, &labels, nil, ckpt.Aggregation(), intervalStart, intervalEnd)
		got, err := minMaxSumCount(record, ckpt.(aggregation.MinMaxSumCount))
		if assert.NoError(t, err) {
			assert.Equal(t, test.expected, got.MetricDescriptor)
		}
	}
}

func TestMinMaxSumCountDatapoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderKind, metric.Int64NumberKind)
	labels := label.NewSet()
	mmsc, ckpt := test.Unslice2(minmaxsumcount.New(2, &desc))

	assert.NoError(t, mmsc.Update(context.Background(), 1, &desc))
	assert.NoError(t, mmsc.Update(context.Background(), 10, &desc))
	require.NoError(t, mmsc.SynchronizedMove(ckpt, &desc))
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
			StartTimeUnixNano: uint64(intervalStart.UnixNano()),
			TimeUnixNano:      uint64(intervalEnd.UnixNano()),
		},
	}
	record := export.NewRecord(&desc, &labels, nil, ckpt.Aggregation(), intervalStart, intervalEnd)
	m, err := minMaxSumCount(record, ckpt.(aggregation.MinMaxSumCount))
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
	mmsc := &minmaxsumcount.New(1, &metric.Descriptor{})[0]
	_, _, _, _, err := minMaxSumCountValues(mmsc)
	assert.Error(t, err)
	assert.Equal(t, aggregation.ErrNoData, err)
}

func TestSumMetricDescriptor(t *testing.T) {
	tests := []struct {
		name        string
		metricKind  metric.Kind
		description string
		unit        unit.Unit
		numberKind  metric.NumberKind
		labels      []kv.KeyValue
		expected    *metricpb.MetricDescriptor
	}{
		{
			"sum-test-a",
			metric.CounterKind,
			"test-a-description",
			unit.Dimensionless,
			metric.Int64NumberKind,
			[]kv.KeyValue{},
			&metricpb.MetricDescriptor{
				Name:        "sum-test-a",
				Description: "test-a-description",
				Unit:        "1",
				Type:        metricpb.MetricDescriptor_INT64,
			},
		},
		{
			"sum-test-b",
			metric.ValueRecorderKind, // This shouldn't change anything.
			"test-b-description",
			unit.Milliseconds,
			metric.Float64NumberKind,
			[]kv.KeyValue{kv.String("A", "1")},
			&metricpb.MetricDescriptor{
				Name:        "sum-test-b",
				Description: "test-b-description",
				Unit:        "ms",
				Type:        metricpb.MetricDescriptor_DOUBLE,
			},
		},
	}

	for _, test := range tests {
		desc := metric.NewDescriptor(test.name, test.metricKind, test.numberKind,
			metric.WithDescription(test.description),
			metric.WithUnit(test.unit),
		)
		labels := label.NewSet(test.labels...)
		emptyAgg := &sumAgg.New(1)[0]
		record := export.NewRecord(&desc, &labels, nil, emptyAgg, intervalStart, intervalEnd)
		got, err := sum(record, emptyAgg)
		if assert.NoError(t, err) {
			assert.Equal(t, test.expected, got.MetricDescriptor)
		}
	}
}

func TestSumInt64DataPoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderKind, metric.Int64NumberKind)
	labels := label.NewSet()
	s, ckpt := test.Unslice2(sumAgg.New(2))
	assert.NoError(t, s.Update(context.Background(), metric.Number(1), &desc))
	require.NoError(t, s.SynchronizedMove(ckpt, &desc))
	record := export.NewRecord(&desc, &labels, nil, ckpt.Aggregation(), intervalStart, intervalEnd)
	if m, err := sum(record, ckpt.(aggregation.Sum)); assert.NoError(t, err) {
		assert.Equal(t, []*metricpb.Int64DataPoint{{
			Value:             1,
			StartTimeUnixNano: uint64(intervalStart.UnixNano()),
			TimeUnixNano:      uint64(intervalEnd.UnixNano()),
		}}, m.Int64DataPoints)
		assert.Equal(t, []*metricpb.DoubleDataPoint(nil), m.DoubleDataPoints)
		assert.Equal(t, []*metricpb.HistogramDataPoint(nil), m.HistogramDataPoints)
		assert.Equal(t, []*metricpb.SummaryDataPoint(nil), m.SummaryDataPoints)
	}
}

func TestSumFloat64DataPoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderKind, metric.Float64NumberKind)
	labels := label.NewSet()
	s, ckpt := test.Unslice2(sumAgg.New(2))
	assert.NoError(t, s.Update(context.Background(), metric.NewFloat64Number(1), &desc))
	require.NoError(t, s.SynchronizedMove(ckpt, &desc))
	record := export.NewRecord(&desc, &labels, nil, ckpt.Aggregation(), intervalStart, intervalEnd)
	if m, err := sum(record, ckpt.(aggregation.Sum)); assert.NoError(t, err) {
		assert.Equal(t, []*metricpb.Int64DataPoint(nil), m.Int64DataPoints)
		assert.Equal(t, []*metricpb.DoubleDataPoint{{
			Value:             1,
			StartTimeUnixNano: uint64(intervalStart.UnixNano()),
			TimeUnixNano:      uint64(intervalEnd.UnixNano()),
		}}, m.DoubleDataPoints)
		assert.Equal(t, []*metricpb.HistogramDataPoint(nil), m.HistogramDataPoints)
		assert.Equal(t, []*metricpb.SummaryDataPoint(nil), m.SummaryDataPoints)
	}
}

func TestSumErrUnknownValueType(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderKind, metric.NumberKind(-1))
	labels := label.NewSet()
	s := &sumAgg.New(1)[0]
	record := export.NewRecord(&desc, &labels, nil, s, intervalStart, intervalEnd)
	_, err := sum(record, s)
	assert.Error(t, err)
	if !errors.Is(err, ErrUnknownValueType) {
		t.Errorf("expected ErrUnknownValueType, got %v", err)
	}
}
