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

package metrictest_test // import "go.opentelemetry.io/otel/sdk/metric/metrictest"

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
)

func TestSyncCounter(t *testing.T) {
	ctx := context.Background()
	mp, exp := metrictest.NewTestMeterProvider()
	meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestSyncCounter")

	fcnt, err := meter.SyncFloat64().Counter("fCount")
	require.NoError(t, err)
	fudcnt, err := meter.SyncFloat64().UpDownCounter("fUDCount")
	require.NoError(t, err)
	fhis, err := meter.SyncFloat64().Histogram("fHist")
	require.NoError(t, err)

	icnt, err := meter.SyncInt64().Counter("iCount")
	require.NoError(t, err)
	iudcnt, err := meter.SyncInt64().UpDownCounter("iUDCount")
	require.NoError(t, err)
	ihis, err := meter.SyncInt64().Histogram("iHist")
	require.NoError(t, err)

	fcnt.Add(ctx, 2)
	fudcnt.Add(ctx, 3)
	fhis.Record(ctx, 4)
	fhis.Record(ctx, 5)

	icnt.Add(ctx, 22)
	iudcnt.Add(ctx, 23)
	ihis.Record(ctx, 24)
	ihis.Record(ctx, 25)

	err = exp.Collect(context.Background())
	assert.NoError(t, err)

	out, err := exp.GetByName("fCount")
	assert.NoError(t, err)
	assert.InDelta(t, 2.0, out.Sum.AsFloat64(), 0.0001)
	assert.Equal(t, aggregation.SumKind, out.AggregationKind)

	out, err = exp.GetByName("fUDCount")
	assert.NoError(t, err)
	assert.InDelta(t, 3.0, out.Sum.AsFloat64(), 0.0001)
	assert.Equal(t, aggregation.SumKind, out.AggregationKind)

	out, err = exp.GetByName("fHist")
	assert.NoError(t, err)
	assert.InDelta(t, 9.0, out.Sum.AsFloat64(), 0.0001)
	assert.EqualValues(t, 2, out.Count)
	assert.Equal(t, aggregation.HistogramKind, out.AggregationKind)

	out, err = exp.GetByName("iCount")
	assert.NoError(t, err)
	assert.EqualValues(t, 22, out.Sum.AsInt64())
	assert.Equal(t, aggregation.SumKind, out.AggregationKind)

	out, err = exp.GetByName("iUDCount")
	assert.NoError(t, err)
	assert.EqualValues(t, 23, out.Sum.AsInt64())
	assert.Equal(t, aggregation.SumKind, out.AggregationKind)

	out, err = exp.GetByName("iHist")
	assert.NoError(t, err)
	assert.EqualValues(t, 49, out.Sum.AsInt64())
	assert.EqualValues(t, 2, out.Count)
	assert.Equal(t, aggregation.HistogramKind, out.AggregationKind)
}

func TestAsyncCounter(t *testing.T) {
	ctx := context.Background()
	mp, exp := metrictest.NewTestMeterProvider()
	meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestAsyncCounter")

	fcnt, err := meter.AsyncFloat64().Counter("fCount")
	require.NoError(t, err)
	fudcnt, err := meter.AsyncFloat64().UpDownCounter("fUDCount")
	require.NoError(t, err)
	fgauge, err := meter.AsyncFloat64().Gauge("fGauge")
	require.NoError(t, err)

	icnt, err := meter.AsyncInt64().Counter("iCount")
	require.NoError(t, err)
	iudcnt, err := meter.AsyncInt64().UpDownCounter("iUDCount")
	require.NoError(t, err)
	igauge, err := meter.AsyncInt64().Gauge("iGauge")
	require.NoError(t, err)

	meter.RegisterCallback(
		[]instrument.Asynchronous{
			fcnt,
			fudcnt,
			fgauge,
			icnt,
			iudcnt,
			igauge,
		}, func(context.Context) {
			fcnt.Observe(ctx, 2)
			fudcnt.Observe(ctx, 3)
			fgauge.Observe(ctx, 4)
			icnt.Observe(ctx, 22)
			iudcnt.Observe(ctx, 23)
			igauge.Observe(ctx, 25)
		})

	err = exp.Collect(context.Background())
	assert.NoError(t, err)

	out, err := exp.GetByName("fCount")
	assert.NoError(t, err)
	assert.InDelta(t, 2.0, out.Sum.AsFloat64(), 0.0001)
	assert.Equal(t, aggregation.SumKind, out.AggregationKind)

	out, err = exp.GetByName("fUDCount")
	assert.NoError(t, err)
	assert.InDelta(t, 3.0, out.Sum.AsFloat64(), 0.0001)
	assert.Equal(t, aggregation.SumKind, out.AggregationKind)

	out, err = exp.GetByName("fGauge")
	assert.NoError(t, err)
	assert.InDelta(t, 4.0, out.LastValue.AsFloat64(), 0.0001)
	assert.Equal(t, aggregation.LastValueKind, out.AggregationKind)

	out, err = exp.GetByName("iCount")
	assert.NoError(t, err)
	assert.EqualValues(t, 22, out.Sum.AsInt64())
	assert.Equal(t, aggregation.SumKind, out.AggregationKind)

	out, err = exp.GetByName("iUDCount")
	assert.NoError(t, err)
	assert.EqualValues(t, 23, out.Sum.AsInt64())
	assert.Equal(t, aggregation.SumKind, out.AggregationKind)

	out, err = exp.GetByName("iGauge")
	assert.NoError(t, err)
	assert.EqualValues(t, 25, out.LastValue.AsInt64())
	assert.Equal(t, aggregation.LastValueKind, out.AggregationKind)
}

func ExampleExporter_GetByName() {
	mp, exp := metrictest.NewTestMeterProvider()
	meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestSyncCounter")

	cnt, err := meter.SyncFloat64().Counter("fCount")
	if err != nil {
		panic("could not acquire counter")
	}

	cnt.Add(context.Background(), 2.5)

	err = exp.Collect(context.Background())
	if err != nil {
		panic("collection failed")
	}

	out, _ := exp.GetByName("fCount")

	fmt.Println(out.Sum.AsFloat64())
	// Output: 2.5
}

func ExampleExporter_GetByNameAndLabels() {
	mp, exp := metrictest.NewTestMeterProvider()
	meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestSyncCounter")

	cnt, err := meter.SyncFloat64().Counter("fCount")
	if err != nil {
		panic("could not acquire counter")
	}

	cnt.Add(context.Background(), 4, attribute.String("foo", "bar"))

	err = exp.Collect(context.Background())
	if err != nil {
		panic("collection failed")
	}

	out, err := exp.GetByNameAndLabels("fCount", []attribute.KeyValue{attribute.String("foo", "bar")})
	if err != nil {
		println(err.Error())
	}

	fmt.Println(out.Sum.AsFloat64())
	// Output: 4

}
