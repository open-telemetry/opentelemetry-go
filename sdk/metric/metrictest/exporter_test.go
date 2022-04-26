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
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

func TestSyncInstruments(t *testing.T) {
	ctx := context.Background()
	mp, exp := metrictest.NewTestMeterProvider()
	meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestSyncInstruments")

	t.Run("Float Counter", func(t *testing.T) {
		fcnt, err := meter.SyncFloat64().Counter("fCount")
		require.NoError(t, err)

		fcnt.Add(ctx, 2)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("fCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewMonotonicFloat64(2), out.Aggregation)
	})

	t.Run("Float UpDownCounter", func(t *testing.T) {
		fudcnt, err := meter.SyncFloat64().UpDownCounter("fUDCount")
		require.NoError(t, err)

		fudcnt.Add(ctx, 3)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("fUDCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewNonMonotonicFloat64(3), out.Aggregation)

	})

	t.Run("Float Histogram", func(t *testing.T) {
		fhis, err := meter.SyncFloat64().Histogram("fHist")
		require.NoError(t, err)

		fhis.Record(ctx, 4)
		fhis.Record(ctx, 5)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("fHist")
		assert.NoError(t, err)
		assert.EqualValues(t, histogram.NewFloat64(nil, 4, 5), out.Aggregation)
	})

	t.Run("Int Counter", func(t *testing.T) {
		icnt, err := meter.SyncInt64().Counter("iCount")
		require.NoError(t, err)

		icnt.Add(ctx, 22)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("iCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewMonotonicInt64(22), out.Aggregation)

	})
	t.Run("Int UpDownCounter", func(t *testing.T) {
		iudcnt, err := meter.SyncInt64().UpDownCounter("iUDCount")
		require.NoError(t, err)

		iudcnt.Add(ctx, 23)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("iUDCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewNonMonotonicInt64(23), out.Aggregation)
	})
	t.Run("Int Histogram", func(t *testing.T) {

		ihis, err := meter.SyncInt64().Histogram("iHist")
		require.NoError(t, err)

		ihis.Record(ctx, 24)
		ihis.Record(ctx, 25)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("iHist")
		assert.NoError(t, err)
		assert.EqualValues(t, histogram.NewInt64(nil, 24, 25), out.Aggregation)
	})
}

func TestSyncDeltaInstruments(t *testing.T) {
	ctx := context.Background()
	mp, exp := metrictest.NewTestMeterProvider(metrictest.WithTemporalitySelector(view.DeltaPreferredTemporality))
	meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestSyncDeltaInstruments")

	t.Run("Float Counter", func(t *testing.T) {
		fcnt, err := meter.SyncFloat64().Counter("fCount")
		require.NoError(t, err)

		fcnt.Add(ctx, 2)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("fCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewMonotonicFloat64(2), out.Aggregation)
	})

	t.Run("Float UpDownCounter", func(t *testing.T) {
		fudcnt, err := meter.SyncFloat64().UpDownCounter("fUDCount")
		require.NoError(t, err)

		fudcnt.Add(ctx, 3)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("fUDCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewNonMonotonicFloat64(3), out.Aggregation)

	})

	t.Run("Float Histogram", func(t *testing.T) {
		fhis, err := meter.SyncFloat64().Histogram("fHist")
		require.NoError(t, err)

		fhis.Record(ctx, 4)
		fhis.Record(ctx, 5)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("fHist")
		assert.NoError(t, err)
		assert.EqualValues(t, histogram.NewFloat64(nil, 4, 5), out.Aggregation)
	})

	t.Run("Int Counter", func(t *testing.T) {
		icnt, err := meter.SyncInt64().Counter("iCount")
		require.NoError(t, err)

		icnt.Add(ctx, 22)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("iCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewMonotonicInt64(22), out.Aggregation)

	})
	t.Run("Int UpDownCounter", func(t *testing.T) {
		iudcnt, err := meter.SyncInt64().UpDownCounter("iUDCount")
		require.NoError(t, err)

		iudcnt.Add(ctx, 23)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("iUDCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewNonMonotonicInt64(23), out.Aggregation)

	})
	t.Run("Int Histogram", func(t *testing.T) {

		ihis, err := meter.SyncInt64().Histogram("iHist")
		require.NoError(t, err)

		ihis.Record(ctx, 24)
		ihis.Record(ctx, 25)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("iHist")
		assert.NoError(t, err)
		assert.EqualValues(t, histogram.NewInt64(nil, 24, 25), out.Aggregation)
	})
}

func TestAsyncInstruments(t *testing.T) {
	mp, exp := metrictest.NewTestMeterProvider()

	t.Run("Float Counter", func(t *testing.T) {
		meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestAsyncCounter_FloatCounter")

		fcnt, err := meter.AsyncFloat64().Counter("fCount")
		require.NoError(t, err)

		err = meter.RegisterCallback(
			[]instrument.Asynchronous{
				fcnt,
			}, func(ctx context.Context) {
				fcnt.Observe(ctx, 2)
			})
		require.NoError(t, err)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("fCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewMonotonicFloat64(2), out.Aggregation)
	})

	t.Run("Float UpDownCounter", func(t *testing.T) {
		meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestAsyncCounter_FloatUpDownCounter")

		fudcnt, err := meter.AsyncFloat64().UpDownCounter("fUDCount")
		require.NoError(t, err)

		err = meter.RegisterCallback(
			[]instrument.Asynchronous{
				fudcnt,
			}, func(ctx context.Context) {
				fudcnt.Observe(ctx, 3)
			})
		require.NoError(t, err)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("fUDCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewNonMonotonicFloat64(3), out.Aggregation)
	})

	t.Run("Float Gauge", func(t *testing.T) {
		meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestAsyncCounter_FloatGauge")

		fgauge, err := meter.AsyncFloat64().Gauge("fGauge")
		require.NoError(t, err)

		err = meter.RegisterCallback(
			[]instrument.Asynchronous{
				fgauge,
			}, func(ctx context.Context) {
				fgauge.Observe(ctx, 4)
			})
		require.NoError(t, err)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("fGauge")
		assert.NoError(t, err)
		assert.EqualValues(t, gauge.NewFloat64(4), out.Aggregation)
	})

	t.Run("Int Counter", func(t *testing.T) {
		meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestAsyncCounter_IntCounter")

		icnt, err := meter.AsyncInt64().Counter("iCount")
		require.NoError(t, err)

		err = meter.RegisterCallback(
			[]instrument.Asynchronous{
				icnt,
			}, func(ctx context.Context) {
				icnt.Observe(ctx, 22)
			})
		require.NoError(t, err)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("iCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewMonotonicInt64(22), out.Aggregation)
	})

	t.Run("Int UpDownCounter", func(t *testing.T) {
		meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestAsyncCounter_IntUpDownCounter")

		iudcnt, err := meter.AsyncInt64().UpDownCounter("iUDCount")
		require.NoError(t, err)

		err = meter.RegisterCallback(
			[]instrument.Asynchronous{
				iudcnt,
			}, func(ctx context.Context) {
				iudcnt.Observe(ctx, 23)
			})
		require.NoError(t, err)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("iUDCount")
		assert.NoError(t, err)
		assert.EqualValues(t, sum.NewNonMonotonicInt64(23), out.Aggregation)

	})
	t.Run("Int Gauge", func(t *testing.T) {
		meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_TestAsyncCounter_IntGauge")

		igauge, err := meter.AsyncInt64().Gauge("iGauge")
		require.NoError(t, err)

		err = meter.RegisterCallback(
			[]instrument.Asynchronous{
				igauge,
			}, func(ctx context.Context) {
				igauge.Observe(ctx, 25)
			})
		require.NoError(t, err)

		err = exp.Collect(context.Background())
		assert.NoError(t, err)

		out, err := exp.GetByName("iGauge")
		assert.NoError(t, err)
		assert.EqualValues(t, gauge.NewInt64(25), out.Aggregation)
	})
}

func ExampleExporter_GetByName() {
	mp, exp := metrictest.NewTestMeterProvider()
	meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_ExampleExporter_GetByName")

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

	fmt.Println(out.Aggregation.(aggregation.Sum).Sum().AsFloat64())
	// Output: 2.5
}

func ExampleExporter_GetByNameAndAttributes() {
	mp, exp := metrictest.NewTestMeterProvider()
	meter := mp.Meter("go.opentelemetry.io/otel/sdk/metric/metrictest/exporter_ExampleExporter_GetByNameAndAttributes")

	cnt, err := meter.SyncFloat64().Counter("fCount")
	if err != nil {
		panic("could not acquire counter")
	}

	cnt.Add(context.Background(), 4, attribute.String("foo", "bar"), attribute.Bool("found", false))

	err = exp.Collect(context.Background())
	if err != nil {
		panic("collection failed")
	}

	out, err := exp.GetByNameAndAttributes("fCount", []attribute.KeyValue{attribute.String("foo", "bar")})
	if err != nil {
		println(err.Error())
	}

	fmt.Println(out.Aggregation.(aggregation.Sum).Sum().AsFloat64())
	// Output: 4

}
