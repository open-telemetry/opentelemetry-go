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

package metric_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

// This example can be found: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/supplementary-guidelines.md#asynchronous-example
func TestCumulativeAsynchronousExample(t *testing.T) {
	ctx := context.Background()
	filter := attribute.Filter(func(kv attribute.KeyValue) bool {
		return kv.Key != "tid"
	})
	reader := metric.NewManualReader()

	defaultView := metric.NewView(metric.Instrument{Name: "pageFaults"}, metric.Stream{Name: "pageFaults"})
	filteredView := metric.NewView(metric.Instrument{Name: "pageFaults"}, metric.Stream{Name: "filteredPageFaults", AttributeFilter: filter})

	meter := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithView(defaultView),
		metric.WithView(filteredView),
	).Meter("AsynchronousExample")

	ctr, err := meter.Int64ObservableCounter("pageFaults")
	assert.NoError(t, err)

	tid1Attrs := []attribute.KeyValue{attribute.String("pid", "1001"), attribute.Int("tid", 1)}
	tid2Attrs := []attribute.KeyValue{attribute.String("pid", "1001"), attribute.Int("tid", 2)}
	tid3Attrs := []attribute.KeyValue{attribute.String("pid", "1001"), attribute.Int("tid", 3)}

	attrs := [][]attribute.KeyValue{tid1Attrs, tid2Attrs, tid3Attrs}

	pfValues := []int64{0, 0, 0}

	_, err = meter.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
		for i := range pfValues {
			if pfValues[i] != 0 {
				ctr.Observe(ctx, pfValues[i], attrs[i]...)
			}
		}
	})
	assert.NoError(t, err)

	filteredAttributeSet := attribute.NewSet(attribute.KeyValue{Key: "pid", Value: attribute.StringValue("1001")})

	// During the time range (T0, T1]:
	//     pid = 1001, tid = 1, #PF = 50
	//     pid = 1001, tid = 2, #PF = 30
	atomic.StoreInt64(&pfValues[0], 50)
	atomic.StoreInt64(&pfValues[1], 30)

	wantScopeMetrics := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{Name: "AsynchronousExample"},
		Metrics: []metricdata.Metrics{
			{
				Name: "filteredPageFaults",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: filteredAttributeSet,
							Value:      80,
						},
					},
				},
			},
			{
				Name: "pageFaults",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(tid1Attrs...),
							Value:      50,
						},
						{
							Attributes: attribute.NewSet(tid2Attrs...),
							Value:      30,
						},
					},
				},
			},
		},
	}

	metrics, err := reader.Collect(ctx)
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, wantScopeMetrics, metrics.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())

	wantFilterValue := &wantScopeMetrics.Metrics[0].Data.(metricdata.Sum[int64]).DataPoints[0].Value
	wantDataPoint1Value := &wantScopeMetrics.Metrics[1].Data.(metricdata.Sum[int64]).DataPoints[0].Value
	wantDataPoint2Value := &wantScopeMetrics.Metrics[1].Data.(metricdata.Sum[int64]).DataPoints[1].Value

	// During the time range (T1, T2]:
	//     pid = 1001, tid = 1, #PF = 53
	//     pid = 1001, tid = 2, #PF = 38

	atomic.StoreInt64(&pfValues[0], 53)
	atomic.StoreInt64(&pfValues[1], 38)

	*wantFilterValue = 91
	*wantDataPoint1Value = 53
	*wantDataPoint2Value = 38

	metrics, err = reader.Collect(ctx)
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, wantScopeMetrics, metrics.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())

	// During the time range (T2, T3]
	//     pid = 1001, tid = 1, #PF = 56
	//     pid = 1001, tid = 2, #PF = 42

	atomic.StoreInt64(&pfValues[0], 56)
	atomic.StoreInt64(&pfValues[1], 42)

	*wantFilterValue = 98
	*wantDataPoint1Value = 56
	*wantDataPoint2Value = 42

	metrics, err = reader.Collect(ctx)
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, wantScopeMetrics, metrics.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())

	// During the time range (T3, T4]:
	//     pid = 1001, tid = 1, #PF = 60
	//     pid = 1001, tid = 2, #PF = 47

	atomic.StoreInt64(&pfValues[0], 60)
	atomic.StoreInt64(&pfValues[1], 47)

	*wantFilterValue = 107
	*wantDataPoint1Value = 60
	*wantDataPoint2Value = 47

	metrics, err = reader.Collect(ctx)
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, wantScopeMetrics, metrics.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())

	// During the time range (T4, T5]:
	//     thread 1 died, thread 3 started
	//     pid = 1001, tid = 2, #PF = 53
	//     pid = 1001, tid = 3, #PF = 5

	atomic.StoreInt64(&pfValues[0], 0)
	atomic.StoreInt64(&pfValues[1], 53)
	atomic.StoreInt64(&pfValues[2], 5)

	*wantFilterValue = 58
	wantAgg := metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints: []metricdata.DataPoint[int64]{
			{
				Attributes: attribute.NewSet(tid1Attrs...),
				Value:      60,
			},
			{
				Attributes: attribute.NewSet(tid2Attrs...),
				Value:      53,
			},
			{
				Attributes: attribute.NewSet(tid3Attrs...),
				Value:      5,
			},
		},
	}
	wantScopeMetrics.Metrics[1].Data = wantAgg

	metrics, err = reader.Collect(ctx)
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, wantScopeMetrics, metrics.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())
}

// This example can be found: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/supplementary-guidelines.md#asynchronous-example

func TestDeltaAsynchronousExample(t *testing.T) {
	ctx := context.Background()
	filter := attribute.Filter(func(kv attribute.KeyValue) bool {
		return kv.Key != "tid"
	})
	reader := metric.NewManualReader(metric.WithTemporalitySelector(func(ik metric.InstrumentKind) metricdata.Temporality { return metricdata.DeltaTemporality }))

	defaultView := metric.NewView(metric.Instrument{Name: "pageFaults"}, metric.Stream{Name: "pageFaults"})
	filteredView := metric.NewView(metric.Instrument{Name: "pageFaults"}, metric.Stream{Name: "filteredPageFaults", AttributeFilter: filter})

	meter := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithView(defaultView),
		metric.WithView(filteredView),
	).Meter("AsynchronousExample")

	ctr, err := meter.Int64ObservableCounter("pageFaults")
	assert.NoError(t, err)

	tid1Attrs := []attribute.KeyValue{attribute.String("pid", "1001"), attribute.Int("tid", 1)}
	tid2Attrs := []attribute.KeyValue{attribute.String("pid", "1001"), attribute.Int("tid", 2)}
	tid3Attrs := []attribute.KeyValue{attribute.String("pid", "1001"), attribute.Int("tid", 3)}

	attrs := [][]attribute.KeyValue{tid1Attrs, tid2Attrs, tid3Attrs}

	pfValues := []int64{0, 0, 0}

	_, err = meter.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
		for i := range pfValues {
			if pfValues[i] != 0 {
				ctr.Observe(ctx, pfValues[i], attrs[i]...)
			}
		}
	})
	assert.NoError(t, err)

	filteredAttributeSet := attribute.NewSet(attribute.KeyValue{Key: "pid", Value: attribute.StringValue("1001")})

	// During the time range (T0, T1]:
	//     pid = 1001, tid = 1, #PF = 50
	//     pid = 1001, tid = 2, #PF = 30
	atomic.StoreInt64(&pfValues[0], 50)
	atomic.StoreInt64(&pfValues[1], 30)

	wantScopeMetrics := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{Name: "AsynchronousExample"},
		Metrics: []metricdata.Metrics{
			{
				Name: "filteredPageFaults",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.DeltaTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: filteredAttributeSet,
							Value:      80,
						},
					},
				},
			},
			{
				Name: "pageFaults",
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.DeltaTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(tid1Attrs...),
							Value:      50,
						},
						{
							Attributes: attribute.NewSet(tid2Attrs...),
							Value:      30,
						},
					},
				},
			},
		},
	}

	metrics, err := reader.Collect(ctx)
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, wantScopeMetrics, metrics.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())

	wantFilterValue := &wantScopeMetrics.Metrics[0].Data.(metricdata.Sum[int64]).DataPoints[0].Value
	wantDataPoint1Value := &wantScopeMetrics.Metrics[1].Data.(metricdata.Sum[int64]).DataPoints[0].Value
	wantDataPoint2Value := &wantScopeMetrics.Metrics[1].Data.(metricdata.Sum[int64]).DataPoints[1].Value

	// During the time range (T1, T2]:
	//     pid = 1001, tid = 1, #PF = 53
	//     pid = 1001, tid = 2, #PF = 38

	atomic.StoreInt64(&pfValues[0], 53)
	atomic.StoreInt64(&pfValues[1], 38)

	*wantFilterValue = 11
	*wantDataPoint1Value = 3
	*wantDataPoint2Value = 8

	metrics, err = reader.Collect(ctx)
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, wantScopeMetrics, metrics.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())

	// During the time range (T2, T3]
	//     pid = 1001, tid = 1, #PF = 56
	//     pid = 1001, tid = 2, #PF = 42

	atomic.StoreInt64(&pfValues[0], 56)
	atomic.StoreInt64(&pfValues[1], 42)

	*wantFilterValue = 7
	*wantDataPoint1Value = 3
	*wantDataPoint2Value = 4

	metrics, err = reader.Collect(ctx)
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, wantScopeMetrics, metrics.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())

	// During the time range (T3, T4]:
	//     pid = 1001, tid = 1, #PF = 60
	//     pid = 1001, tid = 2, #PF = 47

	atomic.StoreInt64(&pfValues[0], 60)
	atomic.StoreInt64(&pfValues[1], 47)

	*wantFilterValue = 9
	*wantDataPoint1Value = 4
	*wantDataPoint2Value = 5

	metrics, err = reader.Collect(ctx)
	assert.NoError(t, err)
	metricdatatest.AssertEqual(t, wantScopeMetrics, metrics.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())

	// During the time range (T4, T5]:
	//     thread 1 died, thread 3 started
	//     pid = 1001, tid = 2, #PF = 53
	//     pid = 1001, tid = 3, #PF = 5

	atomic.StoreInt64(&pfValues[0], 0)
	atomic.StoreInt64(&pfValues[1], 53)
	atomic.StoreInt64(&pfValues[2], 5)

	*wantFilterValue = -49

	wantAgg := metricdata.Sum[int64]{
		Temporality: metricdata.DeltaTemporality,
		IsMonotonic: true,
		DataPoints: []metricdata.DataPoint[int64]{
			{
				Attributes: attribute.NewSet(tid1Attrs...),
				Value:      0,
			},
			{
				Attributes: attribute.NewSet(tid2Attrs...),
				Value:      6,
			},
			{
				Attributes: attribute.NewSet(tid3Attrs...),
				Value:      5,
			},
		},
	}
	wantScopeMetrics.Metrics[1].Data = wantAgg

	metrics, err = reader.Collect(ctx)
	assert.NoError(t, err)

	metricdatatest.AssertEqual(t, wantScopeMetrics, metrics.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())
}
