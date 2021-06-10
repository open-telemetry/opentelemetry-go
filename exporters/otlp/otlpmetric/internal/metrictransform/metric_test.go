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

package metrictransform

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/export/metric/metrictest"
	arrAgg "go.opentelemetry.io/otel/sdk/metric/aggregator/exact"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	lvAgg "go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	sumAgg "go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/resource"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

var (
	// Timestamps used in this test:

	intervalStart = time.Now()
	intervalEnd   = intervalStart.Add(time.Hour)
)

const (
	otelCumulative = metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE
	otelDelta      = metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA
)

func TestStringKeyValues(t *testing.T) {
	tests := []struct {
		kvs      []attribute.KeyValue
		expected []*commonpb.KeyValue
	}{
		{
			nil,
			nil,
		},
		{
			[]attribute.KeyValue{},
			nil,
		},
		{
			[]attribute.KeyValue{
				attribute.Bool("true", true),
				attribute.Int64("one", 1),
				attribute.Int64("two", 2),
				attribute.Float64("three", 3),
				attribute.Int("four", 4),
				attribute.Int("five", 5),
				attribute.Float64("six", 6),
				attribute.Int("seven", 7),
				attribute.Int("eight", 8),
				attribute.String("the", "final word"),
			},
			[]*commonpb.KeyValue{
				{Key: "eight", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 8}}},
				{Key: "five", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 5}}},
				{Key: "four", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 4}}},
				{Key: "one", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 1}}},
				{Key: "seven", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 7}}},
				{Key: "six", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_DoubleValue{DoubleValue: 6.0}}},
				{Key: "the", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "final word"}}},
				{Key: "three", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_DoubleValue{DoubleValue: 3.0}}},
				{Key: "true", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_BoolValue{BoolValue: true}}},
				{Key: "two", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_IntValue{IntValue: 2}}},
			},
		},
	}

	for _, test := range tests {
		labels := attribute.NewSet(test.kvs...)
		assert.Equal(t, test.expected, keyValues(labels.Iter()))
	}
}

func TestMinMaxSumCountValue(t *testing.T) {
	mmsc, ckpt := metrictest.Unslice2(minmaxsumcount.New(2, &metric.Descriptor{}))

	assert.NoError(t, mmsc.Update(context.Background(), 1, &metric.Descriptor{}))
	assert.NoError(t, mmsc.Update(context.Background(), 10, &metric.Descriptor{}))

	// Prior to checkpointing ErrNoData should be returned.
	_, _, _, _, err := minMaxSumCountValues(ckpt.(aggregation.MinMaxSumCount))
	assert.EqualError(t, err, aggregation.ErrNoData.Error())

	// Checkpoint to set non-zero values
	require.NoError(t, mmsc.SynchronizedMove(ckpt, &metric.Descriptor{}))
	min, max, sum, count, err := minMaxSumCountValues(ckpt.(aggregation.MinMaxSumCount))
	if assert.NoError(t, err) {
		assert.Equal(t, min, number.NewInt64Number(1))
		assert.Equal(t, max, number.NewInt64Number(10))
		assert.Equal(t, sum, number.NewInt64Number(11))
		assert.Equal(t, count, uint64(2))
	}
}

func TestMinMaxSumCountDatapoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderInstrumentKind, number.Int64Kind)
	labels := attribute.NewSet(attribute.String("one", "1"))
	mmsc, ckpt := metrictest.Unslice2(minmaxsumcount.New(2, &desc))

	assert.NoError(t, mmsc.Update(context.Background(), 1, &desc))
	assert.NoError(t, mmsc.Update(context.Background(), 10, &desc))
	require.NoError(t, mmsc.SynchronizedMove(ckpt, &desc))
	expected := []*metricpb.SummaryDataPoint{
		{
			Count:             2,
			Sum:               11,
			StartTimeUnixNano: uint64(intervalStart.UnixNano()),
			TimeUnixNano:      uint64(intervalEnd.UnixNano()),
			Attributes: []*commonpb.KeyValue{
				{
					Key:   "one",
					Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "1"}},
				},
			},
			QuantileValues: []*metricpb.SummaryDataPoint_ValueAtQuantile{
				{
					Quantile: 0.0,
					Value:    1.0,
				},
				{
					Quantile: 1.0,
					Value:    10.0,
				},
			},
		},
	}
	record := export.NewRecord(&desc, &labels, nil, ckpt.Aggregation(), intervalStart, intervalEnd)
	m, err := minMaxSumCount(record, ckpt.(aggregation.MinMaxSumCount))
	if assert.NoError(t, err) {
		assert.Nil(t, m.GetGauge())
		assert.Nil(t, m.GetSum())
		assert.Nil(t, m.GetHistogram())
		assert.Equal(t, expected, m.GetSummary().DataPoints)
		assert.Nil(t, m.GetIntGauge())     // nolint
		assert.Nil(t, m.GetIntSum())       // nolint
		assert.Nil(t, m.GetIntHistogram()) // nolint
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

func TestSumIntDataPoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderInstrumentKind, number.Int64Kind)
	labels := attribute.NewSet(attribute.String("one", "1"))
	s, ckpt := metrictest.Unslice2(sumAgg.New(2))
	assert.NoError(t, s.Update(context.Background(), number.Number(1), &desc))
	require.NoError(t, s.SynchronizedMove(ckpt, &desc))
	record := export.NewRecord(&desc, &labels, nil, ckpt.Aggregation(), intervalStart, intervalEnd)
	sum, ok := ckpt.(aggregation.Sum)
	require.True(t, ok, "ckpt is not an aggregation.Sum: %T", ckpt)
	value, err := sum.Sum()
	require.NoError(t, err)

	if m, err := sumPoint(record, value, record.StartTime(), record.EndTime(), export.CumulativeExportKind, true); assert.NoError(t, err) {
		assert.Nil(t, m.GetGauge())
		assert.Equal(t, &metricpb.Sum{
			AggregationTemporality: otelCumulative,
			IsMonotonic:            true,
			DataPoints: []*metricpb.NumberDataPoint{{
				StartTimeUnixNano: uint64(intervalStart.UnixNano()),
				TimeUnixNano:      uint64(intervalEnd.UnixNano()),
				Attributes: []*commonpb.KeyValue{
					{
						Key:   "one",
						Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "1"}},
					},
				},
				Value: &metricpb.NumberDataPoint_AsInt{
					AsInt: 1,
				},
			}},
		}, m.GetSum())
		assert.Nil(t, m.GetHistogram())
		assert.Nil(t, m.GetSummary())
		assert.Nil(t, m.GetIntGauge())     // nolint
		assert.Nil(t, m.GetIntSum())       // nolint
		assert.Nil(t, m.GetIntHistogram()) // nolint
	}
}

func TestSumFloatDataPoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderInstrumentKind, number.Float64Kind)
	labels := attribute.NewSet(attribute.String("one", "1"))
	s, ckpt := metrictest.Unslice2(sumAgg.New(2))
	assert.NoError(t, s.Update(context.Background(), number.NewFloat64Number(1), &desc))
	require.NoError(t, s.SynchronizedMove(ckpt, &desc))
	record := export.NewRecord(&desc, &labels, nil, ckpt.Aggregation(), intervalStart, intervalEnd)
	sum, ok := ckpt.(aggregation.Sum)
	require.True(t, ok, "ckpt is not an aggregation.Sum: %T", ckpt)
	value, err := sum.Sum()
	require.NoError(t, err)

	if m, err := sumPoint(record, value, record.StartTime(), record.EndTime(), export.DeltaExportKind, false); assert.NoError(t, err) {
		assert.Nil(t, m.GetGauge())
		assert.Equal(t, &metricpb.Sum{
			IsMonotonic:            false,
			AggregationTemporality: otelDelta,
			DataPoints: []*metricpb.NumberDataPoint{{
				Value: &metricpb.NumberDataPoint_AsDouble{
					AsDouble: 1.0,
				},
				StartTimeUnixNano: uint64(intervalStart.UnixNano()),
				TimeUnixNano:      uint64(intervalEnd.UnixNano()),
				Attributes: []*commonpb.KeyValue{
					{
						Key:   "one",
						Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "1"}},
					},
				},
			}}}, m.GetSum())
		assert.Nil(t, m.GetHistogram())
		assert.Nil(t, m.GetSummary())
		assert.Nil(t, m.GetIntGauge())     // nolint
		assert.Nil(t, m.GetIntSum())       // nolint
		assert.Nil(t, m.GetIntHistogram()) // nolint
	}
}

func TestLastValueIntDataPoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderInstrumentKind, number.Int64Kind)
	labels := attribute.NewSet(attribute.String("one", "1"))
	s, ckpt := metrictest.Unslice2(lvAgg.New(2))
	assert.NoError(t, s.Update(context.Background(), number.Number(100), &desc))
	require.NoError(t, s.SynchronizedMove(ckpt, &desc))
	record := export.NewRecord(&desc, &labels, nil, ckpt.Aggregation(), intervalStart, intervalEnd)
	sum, ok := ckpt.(aggregation.LastValue)
	require.True(t, ok, "ckpt is not an aggregation.LastValue: %T", ckpt)
	value, timestamp, err := sum.LastValue()
	require.NoError(t, err)

	if m, err := gaugePoint(record, value, time.Time{}, timestamp); assert.NoError(t, err) {
		assert.Equal(t, []*metricpb.NumberDataPoint{{
			StartTimeUnixNano: 0,
			TimeUnixNano:      uint64(timestamp.UnixNano()),
			Attributes: []*commonpb.KeyValue{
				{
					Key:   "one",
					Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "1"}},
				},
			},
			Value: &metricpb.NumberDataPoint_AsInt{
				AsInt: 100,
			},
		}}, m.GetGauge().DataPoints)
		assert.Nil(t, m.GetSum())
		assert.Nil(t, m.GetHistogram())
		assert.Nil(t, m.GetSummary())
		assert.Nil(t, m.GetIntGauge())     // nolint
		assert.Nil(t, m.GetIntSum())       // nolint
		assert.Nil(t, m.GetIntHistogram()) // nolint
	}
}

func TestExactIntDataPoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderInstrumentKind, number.Int64Kind)
	labels := attribute.NewSet(attribute.String("one", "1"))
	e, ckpt := metrictest.Unslice2(arrAgg.New(2))
	assert.NoError(t, e.Update(context.Background(), number.Number(100), &desc))
	require.NoError(t, e.SynchronizedMove(ckpt, &desc))
	record := export.NewRecord(&desc, &labels, nil, ckpt.Aggregation(), intervalStart, intervalEnd)
	p, ok := ckpt.(aggregation.Points)
	require.True(t, ok, "ckpt is not an aggregation.Points: %T", ckpt)
	pts, err := p.Points()
	require.NoError(t, err)

	if m, err := gaugeArray(record, pts); assert.NoError(t, err) {
		assert.Equal(t, []*metricpb.NumberDataPoint{{
			StartTimeUnixNano: toNanos(intervalStart),
			TimeUnixNano:      toNanos(intervalEnd),
			Attributes: []*commonpb.KeyValue{
				{
					Key:   "one",
					Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "1"}},
				},
			},
			Value: &metricpb.NumberDataPoint_AsInt{
				AsInt: 100,
			},
		}}, m.GetGauge().DataPoints)
		assert.Nil(t, m.GetSum())
		assert.Nil(t, m.GetHistogram())
		assert.Nil(t, m.GetSummary())
		assert.Nil(t, m.GetIntGauge())     // nolint
		assert.Nil(t, m.GetIntSum())       // nolint
		assert.Nil(t, m.GetIntHistogram()) // nolint
	}
}

func TestExactFloatDataPoints(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderInstrumentKind, number.Float64Kind)
	labels := attribute.NewSet(attribute.String("one", "1"))
	e, ckpt := metrictest.Unslice2(arrAgg.New(2))
	assert.NoError(t, e.Update(context.Background(), number.NewFloat64Number(100), &desc))
	require.NoError(t, e.SynchronizedMove(ckpt, &desc))
	record := export.NewRecord(&desc, &labels, nil, ckpt.Aggregation(), intervalStart, intervalEnd)
	p, ok := ckpt.(aggregation.Points)
	require.True(t, ok, "ckpt is not an aggregation.Points: %T", ckpt)
	pts, err := p.Points()
	require.NoError(t, err)

	if m, err := gaugeArray(record, pts); assert.NoError(t, err) {
		assert.Equal(t, []*metricpb.NumberDataPoint{{
			Value: &metricpb.NumberDataPoint_AsDouble{
				AsDouble: 100,
			},
			StartTimeUnixNano: toNanos(intervalStart),
			TimeUnixNano:      toNanos(intervalEnd),
			Attributes: []*commonpb.KeyValue{
				{
					Key:   "one",
					Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "1"}},
				},
			},
		}}, m.GetGauge().DataPoints)
		assert.Nil(t, m.GetSum())
		assert.Nil(t, m.GetHistogram())
		assert.Nil(t, m.GetSummary())
		assert.Nil(t, m.GetIntGauge())     // nolint
		assert.Nil(t, m.GetIntSum())       // nolint
		assert.Nil(t, m.GetIntHistogram()) // nolint
	}
}

func TestSumErrUnknownValueType(t *testing.T) {
	desc := metric.NewDescriptor("", metric.ValueRecorderInstrumentKind, number.Kind(-1))
	labels := attribute.NewSet()
	s := &sumAgg.New(1)[0]
	record := export.NewRecord(&desc, &labels, nil, s, intervalStart, intervalEnd)
	value, err := s.Sum()
	require.NoError(t, err)

	_, err = sumPoint(record, value, record.StartTime(), record.EndTime(), export.CumulativeExportKind, true)
	assert.Error(t, err)
	if !errors.Is(err, ErrUnknownValueType) {
		t.Errorf("expected ErrUnknownValueType, got %v", err)
	}
}

type testAgg struct {
	kind aggregation.Kind
	agg  aggregation.Aggregation
}

func (t *testAgg) Kind() aggregation.Kind {
	return t.kind
}

func (t *testAgg) Aggregation() aggregation.Aggregation {
	return t.agg
}

// None of these three are used:

func (t *testAgg) Update(ctx context.Context, number number.Number, descriptor *metric.Descriptor) error {
	return nil
}
func (t *testAgg) SynchronizedMove(destination export.Aggregator, descriptor *metric.Descriptor) error {
	return nil
}
func (t *testAgg) Merge(aggregator export.Aggregator, descriptor *metric.Descriptor) error {
	return nil
}

type testErrSum struct {
	err error
}

type testErrLastValue struct {
	err error
}

type testErrMinMaxSumCount struct {
	testErrSum
}

func (te *testErrLastValue) LastValue() (number.Number, time.Time, error) {
	return 0, time.Time{}, te.err
}
func (te *testErrLastValue) Kind() aggregation.Kind {
	return aggregation.LastValueKind
}

func (te *testErrSum) Sum() (number.Number, error) {
	return 0, te.err
}
func (te *testErrSum) Kind() aggregation.Kind {
	return aggregation.SumKind
}

func (te *testErrMinMaxSumCount) Min() (number.Number, error) {
	return 0, te.err
}

func (te *testErrMinMaxSumCount) Max() (number.Number, error) {
	return 0, te.err
}

func (te *testErrMinMaxSumCount) Count() (uint64, error) {
	return 0, te.err
}

var _ export.Aggregator = &testAgg{}
var _ aggregation.Aggregation = &testAgg{}
var _ aggregation.Sum = &testErrSum{}
var _ aggregation.LastValue = &testErrLastValue{}
var _ aggregation.MinMaxSumCount = &testErrMinMaxSumCount{}

func TestRecordAggregatorIncompatibleErrors(t *testing.T) {
	makeMpb := func(kind aggregation.Kind, agg aggregation.Aggregation) (*metricpb.Metric, error) {
		desc := metric.NewDescriptor("things", metric.CounterInstrumentKind, number.Int64Kind)
		labels := attribute.NewSet()
		res := resource.Empty()
		test := &testAgg{
			kind: kind,
			agg:  agg,
		}
		return Record(export.CumulativeExportKindSelector(), export.NewRecord(&desc, &labels, res, test, intervalStart, intervalEnd))
	}

	mpb, err := makeMpb(aggregation.SumKind, &lastvalue.New(1)[0])

	require.Error(t, err)
	require.Nil(t, mpb)
	require.True(t, errors.Is(err, ErrIncompatibleAgg))

	mpb, err = makeMpb(aggregation.LastValueKind, &sum.New(1)[0])

	require.Error(t, err)
	require.Nil(t, mpb)
	require.True(t, errors.Is(err, ErrIncompatibleAgg))

	mpb, err = makeMpb(aggregation.MinMaxSumCountKind, &lastvalue.New(1)[0])

	require.Error(t, err)
	require.Nil(t, mpb)
	require.True(t, errors.Is(err, ErrIncompatibleAgg))

	mpb, err = makeMpb(aggregation.ExactKind, &lastvalue.New(1)[0])

	require.Error(t, err)
	require.Nil(t, mpb)
	require.True(t, errors.Is(err, ErrIncompatibleAgg))
}

func TestRecordAggregatorUnexpectedErrors(t *testing.T) {
	makeMpb := func(kind aggregation.Kind, agg aggregation.Aggregation) (*metricpb.Metric, error) {
		desc := metric.NewDescriptor("things", metric.CounterInstrumentKind, number.Int64Kind)
		labels := attribute.NewSet()
		res := resource.Empty()
		return Record(export.CumulativeExportKindSelector(), export.NewRecord(&desc, &labels, res, agg, intervalStart, intervalEnd))
	}

	errEx := fmt.Errorf("timeout")

	mpb, err := makeMpb(aggregation.SumKind, &testErrSum{errEx})

	require.Error(t, err)
	require.Nil(t, mpb)
	require.True(t, errors.Is(err, errEx))

	mpb, err = makeMpb(aggregation.LastValueKind, &testErrLastValue{errEx})

	require.Error(t, err)
	require.Nil(t, mpb)
	require.True(t, errors.Is(err, errEx))

	mpb, err = makeMpb(aggregation.MinMaxSumCountKind, &testErrMinMaxSumCount{testErrSum{errEx}})

	require.Error(t, err)
	require.Nil(t, mpb)
	require.True(t, errors.Is(err, errEx))
}
