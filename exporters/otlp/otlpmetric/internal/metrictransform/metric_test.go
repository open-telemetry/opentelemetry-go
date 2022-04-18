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
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
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
		attrs := attribute.NewSet(test.kvs...)
		assert.Equal(t, test.expected, Iterator(attrs.Iter()))
	}
}

func TestSumIntDataPoints(t *testing.T) {
	desc := metrictest.NewDescriptor("", sdkapi.HistogramInstrumentKind, number.Int64Kind)
	attrs := attribute.NewSet(attribute.String("one", "1"))
	sums := sum.New(2)
	s, ckpt := &sums[0], &sums[1]

	assert.NoError(t, s.Update(context.Background(), number.Number(1), &desc))
	require.NoError(t, s.SynchronizedMove(ckpt, &desc))
	record := export.NewRecord(&desc, &attrs, ckpt.Aggregation(), intervalStart, intervalEnd)

	value, err := ckpt.Sum()
	require.NoError(t, err)

	if m, err := sumPoint(record, value, record.StartTime(), record.EndTime(), aggregation.CumulativeTemporality, true); assert.NoError(t, err) {
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
	}
}

func TestSumFloatDataPoints(t *testing.T) {
	desc := metrictest.NewDescriptor("", sdkapi.HistogramInstrumentKind, number.Float64Kind)
	attrs := attribute.NewSet(attribute.String("one", "1"))
	sums := sum.New(2)
	s, ckpt := &sums[0], &sums[1]

	assert.NoError(t, s.Update(context.Background(), number.NewFloat64Number(1), &desc))
	require.NoError(t, s.SynchronizedMove(ckpt, &desc))
	record := export.NewRecord(&desc, &attrs, ckpt.Aggregation(), intervalStart, intervalEnd)
	value, err := ckpt.Sum()
	require.NoError(t, err)

	if m, err := sumPoint(record, value, record.StartTime(), record.EndTime(), aggregation.DeltaTemporality, false); assert.NoError(t, err) {
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

	}
}

func TestLastValueIntDataPoints(t *testing.T) {
	desc := metrictest.NewDescriptor("", sdkapi.HistogramInstrumentKind, number.Int64Kind)
	attrs := attribute.NewSet(attribute.String("one", "1"))
	lvs := lastvalue.New(2)
	lv, ckpt := &lvs[0], &lvs[1]

	assert.NoError(t, lv.Update(context.Background(), number.Number(100), &desc))
	require.NoError(t, lv.SynchronizedMove(ckpt, &desc))
	record := export.NewRecord(&desc, &attrs, ckpt.Aggregation(), intervalStart, intervalEnd)
	value, timestamp, err := ckpt.LastValue()
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
	}
}

func TestSumErrUnknownValueType(t *testing.T) {
	desc := metrictest.NewDescriptor("", sdkapi.HistogramInstrumentKind, number.Kind(-1))
	attrs := attribute.NewSet()
	s := &sum.New(1)[0]
	record := export.NewRecord(&desc, &attrs, s, intervalStart, intervalEnd)
	value, err := s.Sum()
	require.NoError(t, err)

	_, err = sumPoint(record, value, record.StartTime(), record.EndTime(), aggregation.CumulativeTemporality, true)
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

func (t *testAgg) Update(ctx context.Context, number number.Number, descriptor *sdkapi.Descriptor) error {
	return nil
}
func (t *testAgg) SynchronizedMove(destination aggregator.Aggregator, descriptor *sdkapi.Descriptor) error {
	return nil
}
func (t *testAgg) Merge(aggregator aggregator.Aggregator, descriptor *sdkapi.Descriptor) error {
	return nil
}

type testErrSum struct {
	err error
}

type testErrLastValue struct {
	err error
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

var _ aggregator.Aggregator = &testAgg{}
var _ aggregation.Aggregation = &testAgg{}
var _ aggregation.Sum = &testErrSum{}
var _ aggregation.LastValue = &testErrLastValue{}

func TestRecordAggregatorIncompatibleErrors(t *testing.T) {
	makeMpb := func(kind aggregation.Kind, agg aggregation.Aggregation) (*metricpb.Metric, error) {
		desc := metrictest.NewDescriptor("things", sdkapi.CounterInstrumentKind, number.Int64Kind)
		attrs := attribute.NewSet()
		test := &testAgg{
			kind: kind,
			agg:  agg,
		}
		return Record(aggregation.CumulativeTemporalitySelector(), export.NewRecord(&desc, &attrs, test, intervalStart, intervalEnd))
	}

	mpb, err := makeMpb(aggregation.SumKind, &lastvalue.New(1)[0])

	require.Error(t, err)
	require.Nil(t, mpb)
	require.True(t, errors.Is(err, ErrIncompatibleAgg))

	mpb, err = makeMpb(aggregation.LastValueKind, &sum.New(1)[0])

	require.Error(t, err)
	require.Nil(t, mpb)
	require.True(t, errors.Is(err, ErrIncompatibleAgg))
}

func TestRecordAggregatorUnexpectedErrors(t *testing.T) {
	makeMpb := func(kind aggregation.Kind, agg aggregation.Aggregation) (*metricpb.Metric, error) {
		desc := metrictest.NewDescriptor("things", sdkapi.CounterInstrumentKind, number.Int64Kind)
		attrs := attribute.NewSet()
		return Record(aggregation.CumulativeTemporalitySelector(), export.NewRecord(&desc, &attrs, agg, intervalStart, intervalEnd))
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
}
