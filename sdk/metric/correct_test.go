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
	"fmt"
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/nonrecording"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type handler struct {
	sync.Mutex
	err error
}

func (h *handler) Handle(err error) {
	h.Lock()
	h.err = err
	h.Unlock()
}

func (h *handler) Reset() {
	h.Lock()
	h.err = nil
	h.Unlock()
}

func (h *handler) Flush() error {
	h.Lock()
	err := h.err
	h.err = nil
	h.Unlock()
	return err
}

var testHandler *handler

func init() {
	testHandler = new(handler)
	otel.SetErrorHandler(testHandler)
}

type testSelector struct {
	selector    export.AggregatorSelector
	newAggCount int
}

func (ts *testSelector) AggregatorFor(desc *sdkapi.Descriptor, aggPtrs ...*aggregator.Aggregator) {
	ts.newAggCount += len(aggPtrs)
	processortest.AggregatorSelector().AggregatorFor(desc, aggPtrs...)
}

func newSDK(t *testing.T) (metric.Meter, *metricsdk.Accumulator, *testSelector, *processortest.Processor) {
	testHandler.Reset()
	testSelector := &testSelector{selector: processortest.AggregatorSelector()}
	processor := processortest.NewProcessor(
		testSelector,
		attribute.DefaultEncoder(),
	)
	accum := metricsdk.NewAccumulator(
		processor,
	)
	meter := sdkapi.WrapMeterImpl(accum)
	return meter, accum, testSelector, processor
}

func TestInputRangeCounter(t *testing.T) {
	ctx := context.Background()
	meter, sdk, _, processor := newSDK(t)

	counter, err := meter.SyncInt64().Counter("name.sum")
	require.NoError(t, err)

	counter.Add(ctx, -1)
	require.Equal(t, aggregation.ErrNegativeInput, testHandler.Flush())

	checkpointed := sdk.Collect(ctx)
	require.Equal(t, 0, checkpointed)

	processor.Reset()
	counter.Add(ctx, 1)
	checkpointed = sdk.Collect(ctx)
	require.Equal(t, map[string]float64{
		"name.sum//": 1,
	}, processor.Values())
	require.Equal(t, 1, checkpointed)
	require.Nil(t, testHandler.Flush())
}

func TestInputRangeUpDownCounter(t *testing.T) {
	ctx := context.Background()
	meter, sdk, _, processor := newSDK(t)

	counter, err := meter.SyncInt64().UpDownCounter("name.sum")
	require.NoError(t, err)

	counter.Add(ctx, -1)
	counter.Add(ctx, -1)
	counter.Add(ctx, 2)
	counter.Add(ctx, 1)

	checkpointed := sdk.Collect(ctx)
	require.Equal(t, map[string]float64{
		"name.sum//": 1,
	}, processor.Values())
	require.Equal(t, 1, checkpointed)
	require.Nil(t, testHandler.Flush())
}

func TestInputRangeHistogram(t *testing.T) {
	ctx := context.Background()
	meter, sdk, _, processor := newSDK(t)

	histogram, err := meter.SyncFloat64().Histogram("name.histogram")
	require.NoError(t, err)

	histogram.Record(ctx, math.NaN())
	require.Equal(t, aggregation.ErrNaNInput, testHandler.Flush())

	checkpointed := sdk.Collect(ctx)
	require.Equal(t, 0, checkpointed)

	histogram.Record(ctx, 1)
	histogram.Record(ctx, 2)

	processor.Reset()
	checkpointed = sdk.Collect(ctx)

	require.Equal(t, map[string]float64{
		"name.histogram//": 3,
	}, processor.Values())
	require.Equal(t, 1, checkpointed)
	require.Nil(t, testHandler.Flush())
}

func TestDisabledInstrument(t *testing.T) {
	ctx := context.Background()
	meter, sdk, _, processor := newSDK(t)

	histogram, err := meter.SyncFloat64().Histogram("name.disabled")
	require.NoError(t, err)

	histogram.Record(ctx, -1)
	checkpointed := sdk.Collect(ctx)

	require.Equal(t, 0, checkpointed)
	require.Equal(t, map[string]float64{}, processor.Values())
}

func TestRecordNaN(t *testing.T) {
	ctx := context.Background()
	meter, _, _, _ := newSDK(t)

	c, err := meter.SyncFloat64().Counter("name.sum")
	require.NoError(t, err)

	require.Nil(t, testHandler.Flush())
	c.Add(ctx, math.NaN())
	require.Error(t, testHandler.Flush())
}

func TestSDKAttrsDeduplication(t *testing.T) {
	ctx := context.Background()
	meter, sdk, _, processor := newSDK(t)

	counter, err := meter.SyncInt64().Counter("name.sum")
	require.NoError(t, err)

	const (
		maxKeys = 21
		keySets = 2
		repeats = 3
	)
	var keysA []attribute.Key
	var keysB []attribute.Key

	for i := 0; i < maxKeys; i++ {
		keysA = append(keysA, attribute.Key(fmt.Sprintf("A%03d", i)))
		keysB = append(keysB, attribute.Key(fmt.Sprintf("B%03d", i)))
	}

	allExpect := map[string]float64{}
	for numKeys := 0; numKeys < maxKeys; numKeys++ {

		var kvsA []attribute.KeyValue
		var kvsB []attribute.KeyValue
		for r := 0; r < repeats; r++ {
			for i := 0; i < numKeys; i++ {
				kvsA = append(kvsA, keysA[i].Int(r))
				kvsB = append(kvsB, keysB[i].Int(r))
			}
		}

		var expectA []attribute.KeyValue
		var expectB []attribute.KeyValue
		for i := 0; i < numKeys; i++ {
			expectA = append(expectA, keysA[i].Int(repeats-1))
			expectB = append(expectB, keysB[i].Int(repeats-1))
		}

		counter.Add(ctx, 1, kvsA...)
		counter.Add(ctx, 1, kvsA...)
		format := func(attrs []attribute.KeyValue) string {
			str := attribute.DefaultEncoder().Encode(newSetIter(attrs...))
			return fmt.Sprint("name.sum/", str, "/")
		}
		allExpect[format(expectA)] += 2

		if numKeys != 0 {
			// In this case A and B sets are the same.
			counter.Add(ctx, 1, kvsB...)
			counter.Add(ctx, 1, kvsB...)
			allExpect[format(expectB)] += 2
		}

	}

	sdk.Collect(ctx)

	require.EqualValues(t, allExpect, processor.Values())
}

func newSetIter(kvs ...attribute.KeyValue) attribute.Iterator {
	attrs := attribute.NewSet(kvs...)
	return attrs.Iter()
}

func TestDefaultAttributeEncoder(t *testing.T) {
	encoder := attribute.DefaultEncoder()

	encoded := encoder.Encode(newSetIter(attribute.String("A", "B"), attribute.String("C", "D")))
	require.Equal(t, `A=B,C=D`, encoded)

	encoded = encoder.Encode(newSetIter(attribute.String("A", "B,c=d"), attribute.String(`C\`, "D")))
	require.Equal(t, `A=B\,c\=d,C\\=D`, encoded)

	encoded = encoder.Encode(newSetIter(attribute.String(`\`, `=`), attribute.String(`,`, `\`)))
	require.Equal(t, `\,=\\,\\=\=`, encoded)

	// Note: the attr encoder does not sort or de-dup values,
	// that is done in Attributes(...).
	encoded = encoder.Encode(newSetIter(
		attribute.Int("I", 1),
		attribute.Int64("I64", 1),
		attribute.Float64("F64", 1),
		attribute.Float64("F64", 1),
		attribute.String("S", "1"),
		attribute.Bool("B", true),
	))
	require.Equal(t, "B=true,F64=1,I=1,I64=1,S=1", encoded)
}

func TestObserverCollection(t *testing.T) {
	ctx := context.Background()
	meter, sdk, _, processor := newSDK(t)
	mult := 1

	gaugeF, err := meter.AsyncFloat64().Gauge("float.gauge.lastvalue")
	require.NoError(t, err)
	err = meter.RegisterCallback([]instrument.Asynchronous{
		gaugeF,
	}, func(ctx context.Context) {
		gaugeF.Observe(ctx, float64(mult), attribute.String("A", "B"))
		// last value wins
		gaugeF.Observe(ctx, float64(-mult), attribute.String("A", "B"))
		gaugeF.Observe(ctx, float64(-mult), attribute.String("C", "D"))
	})
	require.NoError(t, err)

	gaugeI, err := meter.AsyncInt64().Gauge("int.gauge.lastvalue")
	require.NoError(t, err)
	err = meter.RegisterCallback([]instrument.Asynchronous{
		gaugeI,
	}, func(ctx context.Context) {
		gaugeI.Observe(ctx, int64(-mult), attribute.String("A", "B"))
		gaugeI.Observe(ctx, int64(mult))
		// last value wins
		gaugeI.Observe(ctx, int64(mult), attribute.String("A", "B"))
		gaugeI.Observe(ctx, int64(mult))
	})
	require.NoError(t, err)

	counterF, err := meter.AsyncFloat64().Counter("float.counterobserver.sum")
	require.NoError(t, err)
	err = meter.RegisterCallback([]instrument.Asynchronous{
		counterF,
	}, func(ctx context.Context) {
		counterF.Observe(ctx, float64(mult), attribute.String("A", "B"))
		counterF.Observe(ctx, float64(2*mult), attribute.String("A", "B"))
		counterF.Observe(ctx, float64(mult), attribute.String("C", "D"))
	})
	require.NoError(t, err)

	counterI, err := meter.AsyncInt64().Counter("int.counterobserver.sum")
	require.NoError(t, err)
	err = meter.RegisterCallback([]instrument.Asynchronous{
		counterI,
	}, func(ctx context.Context) {
		counterI.Observe(ctx, int64(2*mult), attribute.String("A", "B"))
		counterI.Observe(ctx, int64(mult))
		// last value wins
		counterI.Observe(ctx, int64(mult), attribute.String("A", "B"))
		counterI.Observe(ctx, int64(mult))
	})
	require.NoError(t, err)

	updowncounterF, err := meter.AsyncFloat64().UpDownCounter("float.updowncounterobserver.sum")
	require.NoError(t, err)
	err = meter.RegisterCallback([]instrument.Asynchronous{
		updowncounterF,
	}, func(ctx context.Context) {
		updowncounterF.Observe(ctx, float64(mult), attribute.String("A", "B"))
		updowncounterF.Observe(ctx, float64(-2*mult), attribute.String("A", "B"))
		updowncounterF.Observe(ctx, float64(mult), attribute.String("C", "D"))
	})
	require.NoError(t, err)

	updowncounterI, err := meter.AsyncInt64().UpDownCounter("int.updowncounterobserver.sum")
	require.NoError(t, err)
	err = meter.RegisterCallback([]instrument.Asynchronous{
		updowncounterI,
	}, func(ctx context.Context) {
		updowncounterI.Observe(ctx, int64(2*mult), attribute.String("A", "B"))
		updowncounterI.Observe(ctx, int64(mult))
		// last value wins
		updowncounterI.Observe(ctx, int64(mult), attribute.String("A", "B"))
		updowncounterI.Observe(ctx, int64(-mult))
	})
	require.NoError(t, err)

	unused, err := meter.AsyncInt64().Gauge("empty.gauge.sum")
	require.NoError(t, err)
	err = meter.RegisterCallback([]instrument.Asynchronous{
		unused,
	}, func(ctx context.Context) {
	})
	require.NoError(t, err)

	for mult = 0; mult < 3; mult++ {
		processor.Reset()

		collected := sdk.Collect(ctx)
		require.Equal(t, collected, len(processor.Values()))

		mult := float64(mult)
		require.EqualValues(t, map[string]float64{
			"float.gauge.lastvalue/A=B/": -mult,
			"float.gauge.lastvalue/C=D/": -mult,
			"int.gauge.lastvalue//":      mult,
			"int.gauge.lastvalue/A=B/":   mult,

			"float.counterobserver.sum/A=B/": 3 * mult,
			"float.counterobserver.sum/C=D/": mult,
			"int.counterobserver.sum//":      2 * mult,
			"int.counterobserver.sum/A=B/":   3 * mult,

			"float.updowncounterobserver.sum/A=B/": -mult,
			"float.updowncounterobserver.sum/C=D/": mult,
			"int.updowncounterobserver.sum//":      0,
			"int.updowncounterobserver.sum/A=B/":   3 * mult,
		}, processor.Values())
	}
}

func TestCounterObserverInputRange(t *testing.T) {
	ctx := context.Background()
	meter, sdk, _, processor := newSDK(t)

	// TODO: these tests are testing for negative values, not for _descending values_. Fix.
	counterF, _ := meter.AsyncFloat64().Counter("float.counterobserver.sum")
	err := meter.RegisterCallback([]instrument.Asynchronous{
		counterF,
	}, func(ctx context.Context) {
		counterF.Observe(ctx, -2, attribute.String("A", "B"))
		require.Equal(t, aggregation.ErrNegativeInput, testHandler.Flush())
		counterF.Observe(ctx, -1, attribute.String("C", "D"))
		require.Equal(t, aggregation.ErrNegativeInput, testHandler.Flush())
	})
	require.NoError(t, err)
	counterI, _ := meter.AsyncInt64().Counter("int.counterobserver.sum")
	err = meter.RegisterCallback([]instrument.Asynchronous{
		counterI,
	}, func(ctx context.Context) {
		counterI.Observe(ctx, -1, attribute.String("A", "B"))
		require.Equal(t, aggregation.ErrNegativeInput, testHandler.Flush())
		counterI.Observe(ctx, -1)
		require.Equal(t, aggregation.ErrNegativeInput, testHandler.Flush())
	})
	require.NoError(t, err)

	collected := sdk.Collect(ctx)

	require.Equal(t, 0, collected)
	require.EqualValues(t, map[string]float64{}, processor.Values())

	// check that the error condition was reset
	require.NoError(t, testHandler.Flush())
}

func TestObserverBatch(t *testing.T) {
	ctx := context.Background()
	meter, sdk, _, processor := newSDK(t)

	floatGaugeObs, _ := meter.AsyncFloat64().Gauge("float.gauge.lastvalue")
	intGaugeObs, _ := meter.AsyncInt64().Gauge("int.gauge.lastvalue")
	floatCounterObs, _ := meter.AsyncFloat64().Counter("float.counterobserver.sum")
	intCounterObs, _ := meter.AsyncInt64().Counter("int.counterobserver.sum")
	floatUpDownCounterObs, _ := meter.AsyncFloat64().UpDownCounter("float.updowncounterobserver.sum")
	intUpDownCounterObs, _ := meter.AsyncInt64().UpDownCounter("int.updowncounterobserver.sum")

	err := meter.RegisterCallback([]instrument.Asynchronous{
		floatGaugeObs,
		intGaugeObs,
		floatCounterObs,
		intCounterObs,
		floatUpDownCounterObs,
		intUpDownCounterObs,
	}, func(ctx context.Context) {
		ab := attribute.String("A", "B")
		floatGaugeObs.Observe(ctx, 1, ab)
		floatGaugeObs.Observe(ctx, -1, ab)
		intGaugeObs.Observe(ctx, -1, ab)
		intGaugeObs.Observe(ctx, 1, ab)
		floatCounterObs.Observe(ctx, 1000, ab)
		intCounterObs.Observe(ctx, 100, ab)
		floatUpDownCounterObs.Observe(ctx, -1000, ab)
		intUpDownCounterObs.Observe(ctx, -100, ab)

		cd := attribute.String("C", "D")
		floatGaugeObs.Observe(ctx, -1, cd)
		floatCounterObs.Observe(ctx, -1, cd)
		floatUpDownCounterObs.Observe(ctx, -1, cd)

		intGaugeObs.Observe(ctx, 1)
		intGaugeObs.Observe(ctx, 1)
		intCounterObs.Observe(ctx, 10)
		floatCounterObs.Observe(ctx, 1.1)
		intUpDownCounterObs.Observe(ctx, 10)
	})
	require.NoError(t, err)

	collected := sdk.Collect(ctx)

	require.Equal(t, collected, len(processor.Values()))

	require.EqualValues(t, map[string]float64{
		"float.counterobserver.sum//":    1.1,
		"float.counterobserver.sum/A=B/": 1000,
		"int.counterobserver.sum//":      10,
		"int.counterobserver.sum/A=B/":   100,

		"int.updowncounterobserver.sum/A=B/":   -100,
		"float.updowncounterobserver.sum/A=B/": -1000,
		"int.updowncounterobserver.sum//":      10,
		"float.updowncounterobserver.sum/C=D/": -1,

		"float.gauge.lastvalue/A=B/": -1,
		"float.gauge.lastvalue/C=D/": -1,
		"int.gauge.lastvalue//":      1,
		"int.gauge.lastvalue/A=B/":   1,
	}, processor.Values())
}

// TestRecordPersistence ensures that a direct-called instrument that is
// repeatedly used each interval results in a persistent record, so that its
// encoded attribute will be cached across collection intervals.
func TestRecordPersistence(t *testing.T) {
	ctx := context.Background()
	meter, sdk, selector, _ := newSDK(t)

	c, err := meter.SyncFloat64().Counter("name.sum")
	require.NoError(t, err)

	uk := attribute.String("bound", "false")

	for i := 0; i < 100; i++ {
		c.Add(ctx, 1, uk)
		sdk.Collect(ctx)
	}

	require.Equal(t, 2, selector.newAggCount)
}

func TestIncorrectInstruments(t *testing.T) {
	// The Batch observe/record APIs are susceptible to
	// uninitialized instruments.
	var observer asyncint64.Gauge

	ctx := context.Background()
	meter, sdk, _, processor := newSDK(t)

	// Now try with uninitialized instruments.
	err := meter.RegisterCallback([]instrument.Asynchronous{
		observer,
	}, func(ctx context.Context) {
		observer.Observe(ctx, 1)
	})
	require.ErrorIs(t, err, metricsdk.ErrBadInstrument)

	collected := sdk.Collect(ctx)
	err = testHandler.Flush()
	require.NoError(t, err)
	require.Equal(t, 0, collected)

	// Now try with instruments from another SDK.
	noopMeter := nonrecording.NewNoopMeter()
	observer, _ = noopMeter.AsyncInt64().Gauge("observer")

	err = meter.RegisterCallback(
		[]instrument.Asynchronous{observer},
		func(ctx context.Context) {
			observer.Observe(ctx, 1)
		},
	)
	require.ErrorIs(t, err, metricsdk.ErrBadInstrument)

	collected = sdk.Collect(ctx)
	require.Equal(t, 0, collected)
	require.EqualValues(t, map[string]float64{}, processor.Values())

	err = testHandler.Flush()
	require.NoError(t, err)
}

func TestSyncInAsync(t *testing.T) {
	ctx := context.Background()
	meter, sdk, _, processor := newSDK(t)

	counter, _ := meter.SyncFloat64().Counter("counter.sum")
	gauge, _ := meter.AsyncInt64().Gauge("observer.lastvalue")

	err := meter.RegisterCallback([]instrument.Asynchronous{
		gauge,
	}, func(ctx context.Context) {
		gauge.Observe(ctx, 10)
		counter.Add(ctx, 100)
	})
	require.NoError(t, err)

	sdk.Collect(ctx)

	require.EqualValues(t, map[string]float64{
		"counter.sum//":        100,
		"observer.lastvalue//": 10,
	}, processor.Values())
}
