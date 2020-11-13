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

	"go.opentelemetry.io/otel/global"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/resource"
)

var Must = metric.Must
var testResource = resource.NewWithAttributes(label.String("R", "V"))

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
	global.SetErrorHandler(testHandler)
}

// correctnessProcessor could be replaced with processortest.Processor
// with a non-default aggregator selector.  TODO(#872) use the
// processortest code here.
type correctnessProcessor struct {
	t *testing.T
	*testSelector

	accumulations []export.Accumulation
}

type testSelector struct {
	selector    export.AggregatorSelector
	newAggCount int
}

func (ts *testSelector) AggregatorFor(desc *metric.Descriptor, aggPtrs ...*export.Aggregator) {
	ts.newAggCount += len(aggPtrs)
	processortest.AggregatorSelector().AggregatorFor(desc, aggPtrs...)
}

func newSDK(t *testing.T) (metric.Meter, *metricsdk.Accumulator, *correctnessProcessor) {
	testHandler.Reset()
	processor := &correctnessProcessor{
		t:            t,
		testSelector: &testSelector{selector: processortest.AggregatorSelector()},
	}
	accum := metricsdk.NewAccumulator(
		processor,
		testResource,
	)
	meter := metric.WrapMeterImpl(accum, "test")
	return meter, accum, processor
}

func (ci *correctnessProcessor) Process(accumulation export.Accumulation) error {
	ci.accumulations = append(ci.accumulations, accumulation)
	return nil
}

func TestInputRangeCounter(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	counter := Must(meter).NewInt64Counter("name.sum")

	counter.Add(ctx, -1)
	require.Equal(t, aggregation.ErrNegativeInput, testHandler.Flush())

	checkpointed := sdk.Collect(ctx)
	require.Equal(t, 0, checkpointed)

	processor.accumulations = nil
	counter.Add(ctx, 1)
	checkpointed = sdk.Collect(ctx)
	sum, err := processor.accumulations[0].Aggregator().(aggregation.Sum).Sum()
	require.Equal(t, int64(1), sum.AsInt64())
	require.Equal(t, 1, checkpointed)
	require.Nil(t, err)
	require.Nil(t, testHandler.Flush())
}

func TestInputRangeUpDownCounter(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	counter := Must(meter).NewInt64UpDownCounter("name.sum")

	counter.Add(ctx, -1)
	counter.Add(ctx, -1)
	counter.Add(ctx, 2)
	counter.Add(ctx, 1)

	checkpointed := sdk.Collect(ctx)
	sum, err := processor.accumulations[0].Aggregator().(aggregation.Sum).Sum()
	require.Equal(t, int64(1), sum.AsInt64())
	require.Equal(t, 1, checkpointed)
	require.Nil(t, err)
	require.Nil(t, testHandler.Flush())
}

func TestInputRangeValueRecorder(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	valuerecorder := Must(meter).NewFloat64ValueRecorder("name.exact")

	valuerecorder.Record(ctx, math.NaN())
	require.Equal(t, aggregation.ErrNaNInput, testHandler.Flush())

	checkpointed := sdk.Collect(ctx)
	require.Equal(t, 0, checkpointed)

	valuerecorder.Record(ctx, 1)
	valuerecorder.Record(ctx, 2)

	processor.accumulations = nil
	checkpointed = sdk.Collect(ctx)

	count, err := processor.accumulations[0].Aggregator().(aggregation.Distribution).Count()
	require.Equal(t, int64(2), count)
	require.Equal(t, 1, checkpointed)
	require.Nil(t, testHandler.Flush())
	require.Nil(t, err)
}

func TestDisabledInstrument(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	valuerecorder := Must(meter).NewFloat64ValueRecorder("name.disabled")

	valuerecorder.Record(ctx, -1)
	checkpointed := sdk.Collect(ctx)

	require.Equal(t, 0, checkpointed)
	require.Equal(t, 0, len(processor.accumulations))
}

func TestRecordNaN(t *testing.T) {
	ctx := context.Background()
	meter, _, _ := newSDK(t)

	c := Must(meter).NewFloat64Counter("name.sum")

	require.Nil(t, testHandler.Flush())
	c.Add(ctx, math.NaN())
	require.Error(t, testHandler.Flush())
}

func TestSDKLabelsDeduplication(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	counter := Must(meter).NewInt64Counter("name.sum")

	const (
		maxKeys = 21
		keySets = 2
		repeats = 3
	)
	var keysA []label.Key
	var keysB []label.Key

	for i := 0; i < maxKeys; i++ {
		keysA = append(keysA, label.Key(fmt.Sprintf("A%03d", i)))
		keysB = append(keysB, label.Key(fmt.Sprintf("B%03d", i)))
	}

	var allExpect [][]label.KeyValue
	for numKeys := 0; numKeys < maxKeys; numKeys++ {

		var kvsA []label.KeyValue
		var kvsB []label.KeyValue
		for r := 0; r < repeats; r++ {
			for i := 0; i < numKeys; i++ {
				kvsA = append(kvsA, keysA[i].Int(r))
				kvsB = append(kvsB, keysB[i].Int(r))
			}
		}

		var expectA []label.KeyValue
		var expectB []label.KeyValue
		for i := 0; i < numKeys; i++ {
			expectA = append(expectA, keysA[i].Int(repeats-1))
			expectB = append(expectB, keysB[i].Int(repeats-1))
		}

		counter.Add(ctx, 1, kvsA...)
		counter.Add(ctx, 1, kvsA...)
		allExpect = append(allExpect, expectA)

		if numKeys != 0 {
			// In this case A and B sets are the same.
			counter.Add(ctx, 1, kvsB...)
			counter.Add(ctx, 1, kvsB...)
			allExpect = append(allExpect, expectB)
		}

	}

	sdk.Collect(ctx)

	var actual [][]label.KeyValue
	for _, rec := range processor.accumulations {
		sum, _ := rec.Aggregator().(aggregation.Sum).Sum()
		require.Equal(t, sum, number.NewInt64Number(2))

		kvs := rec.Labels().ToSlice()
		actual = append(actual, kvs)
	}

	require.ElementsMatch(t, allExpect, actual)
}

func newSetIter(kvs ...label.KeyValue) label.Iterator {
	labels := label.NewSet(kvs...)
	return labels.Iter()
}

func TestDefaultLabelEncoder(t *testing.T) {
	encoder := label.DefaultEncoder()

	encoded := encoder.Encode(newSetIter(label.String("A", "B"), label.String("C", "D")))
	require.Equal(t, `A=B,C=D`, encoded)

	encoded = encoder.Encode(newSetIter(label.String("A", "B,c=d"), label.String(`C\`, "D")))
	require.Equal(t, `A=B\,c\=d,C\\=D`, encoded)

	encoded = encoder.Encode(newSetIter(label.String(`\`, `=`), label.String(`,`, `\`)))
	require.Equal(t, `\,=\\,\\=\=`, encoded)

	// Note: the label encoder does not sort or de-dup values,
	// that is done in Labels(...).
	encoded = encoder.Encode(newSetIter(
		label.Int("I", 1),
		label.Uint("U", 1),
		label.Int32("I32", 1),
		label.Uint32("U32", 1),
		label.Int64("I64", 1),
		label.Uint64("U64", 1),
		label.Float64("F64", 1),
		label.Float64("F64", 1),
		label.String("S", "1"),
		label.Bool("B", true),
	))
	require.Equal(t, "B=true,F64=1,I=1,I32=1,I64=1,S=1,U=1,U32=1,U64=1", encoded)
}

func TestObserverCollection(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	_ = Must(meter).NewFloat64ValueObserver("float.valueobserver.lastvalue", func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1, label.String("A", "B"))
		// last value wins
		result.Observe(-1, label.String("A", "B"))
		result.Observe(-1, label.String("C", "D"))
	})
	_ = Must(meter).NewInt64ValueObserver("int.valueobserver.lastvalue", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(-1, label.String("A", "B"))
		result.Observe(1)
		// last value wins
		result.Observe(1, label.String("A", "B"))
		result.Observe(1)
	})

	_ = Must(meter).NewFloat64SumObserver("float.sumobserver.sum", func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1, label.String("A", "B"))
		result.Observe(2, label.String("A", "B"))
		result.Observe(1, label.String("C", "D"))
	})
	_ = Must(meter).NewInt64SumObserver("int.sumobserver.sum", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(2, label.String("A", "B"))
		result.Observe(1)
		// last value wins
		result.Observe(1, label.String("A", "B"))
		result.Observe(1)
	})

	_ = Must(meter).NewFloat64UpDownSumObserver("float.updownsumobserver.sum", func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1, label.String("A", "B"))
		result.Observe(-2, label.String("A", "B"))
		result.Observe(1, label.String("C", "D"))
	})
	_ = Must(meter).NewInt64UpDownSumObserver("int.updownsumobserver.sum", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(2, label.String("A", "B"))
		result.Observe(1)
		// last value wins
		result.Observe(1, label.String("A", "B"))
		result.Observe(-1)
	})

	_ = Must(meter).NewInt64ValueObserver("empty.valueobserver.sum", func(_ context.Context, result metric.Int64ObserverResult) {
	})

	collected := sdk.Collect(ctx)

	require.Equal(t, collected, len(processor.accumulations))

	out := processortest.NewOutput(label.DefaultEncoder())
	for _, rec := range processor.accumulations {
		require.NoError(t, out.AddAccumulation(rec))
	}
	require.EqualValues(t, map[string]float64{
		"float.valueobserver.lastvalue/A=B/R=V": -1,
		"float.valueobserver.lastvalue/C=D/R=V": -1,
		"int.valueobserver.lastvalue//R=V":      1,
		"int.valueobserver.lastvalue/A=B/R=V":   1,

		"float.sumobserver.sum/A=B/R=V": 2,
		"float.sumobserver.sum/C=D/R=V": 1,
		"int.sumobserver.sum//R=V":      1,
		"int.sumobserver.sum/A=B/R=V":   1,

		"float.updownsumobserver.sum/A=B/R=V": -2,
		"float.updownsumobserver.sum/C=D/R=V": 1,
		"int.updownsumobserver.sum//R=V":      -1,
		"int.updownsumobserver.sum/A=B/R=V":   1,
	}, out.Map())
}

func TestSumObserverInputRange(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	// TODO: these tests are testing for negative values, not for _descending values_. Fix.
	_ = Must(meter).NewFloat64SumObserver("float.sumobserver.sum", func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(-2, label.String("A", "B"))
		require.Equal(t, aggregation.ErrNegativeInput, testHandler.Flush())
		result.Observe(-1, label.String("C", "D"))
		require.Equal(t, aggregation.ErrNegativeInput, testHandler.Flush())
	})
	_ = Must(meter).NewInt64SumObserver("int.sumobserver.sum", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(-1, label.String("A", "B"))
		require.Equal(t, aggregation.ErrNegativeInput, testHandler.Flush())
		result.Observe(-1)
		require.Equal(t, aggregation.ErrNegativeInput, testHandler.Flush())
	})

	collected := sdk.Collect(ctx)

	require.Equal(t, 0, collected)
	require.Equal(t, 0, len(processor.accumulations))

	// check that the error condition was reset
	require.NoError(t, testHandler.Flush())
}

func TestObserverBatch(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	var floatValueObs metric.Float64ValueObserver
	var intValueObs metric.Int64ValueObserver
	var floatSumObs metric.Float64SumObserver
	var intSumObs metric.Int64SumObserver
	var floatUpDownSumObs metric.Float64UpDownSumObserver
	var intUpDownSumObs metric.Int64UpDownSumObserver

	var batch = Must(meter).NewBatchObserver(
		func(_ context.Context, result metric.BatchObserverResult) {
			result.Observe(
				[]label.KeyValue{
					label.String("A", "B"),
				},
				floatValueObs.Observation(1),
				floatValueObs.Observation(-1),
				intValueObs.Observation(-1),
				intValueObs.Observation(1),
				floatSumObs.Observation(1000),
				intSumObs.Observation(100),
				floatUpDownSumObs.Observation(-1000),
				intUpDownSumObs.Observation(-100),
			)
			result.Observe(
				[]label.KeyValue{
					label.String("C", "D"),
				},
				floatValueObs.Observation(-1),
				floatSumObs.Observation(-1),
				floatUpDownSumObs.Observation(-1),
			)
			result.Observe(
				nil,
				intValueObs.Observation(1),
				intValueObs.Observation(1),
				intSumObs.Observation(10),
				floatSumObs.Observation(1.1),
				intUpDownSumObs.Observation(10),
			)
		})
	floatValueObs = batch.NewFloat64ValueObserver("float.valueobserver.lastvalue")
	intValueObs = batch.NewInt64ValueObserver("int.valueobserver.lastvalue")
	floatSumObs = batch.NewFloat64SumObserver("float.sumobserver.sum")
	intSumObs = batch.NewInt64SumObserver("int.sumobserver.sum")
	floatUpDownSumObs = batch.NewFloat64UpDownSumObserver("float.updownsumobserver.sum")
	intUpDownSumObs = batch.NewInt64UpDownSumObserver("int.updownsumobserver.sum")

	collected := sdk.Collect(ctx)

	require.Equal(t, collected, len(processor.accumulations))

	out := processortest.NewOutput(label.DefaultEncoder())
	for _, rec := range processor.accumulations {
		require.NoError(t, out.AddAccumulation(rec))
	}
	require.EqualValues(t, map[string]float64{
		"float.sumobserver.sum//R=V":    1.1,
		"float.sumobserver.sum/A=B/R=V": 1000,
		"int.sumobserver.sum//R=V":      10,
		"int.sumobserver.sum/A=B/R=V":   100,

		"int.updownsumobserver.sum/A=B/R=V":   -100,
		"float.updownsumobserver.sum/A=B/R=V": -1000,
		"int.updownsumobserver.sum//R=V":      10,
		"float.updownsumobserver.sum/C=D/R=V": -1,

		"float.valueobserver.lastvalue/A=B/R=V": -1,
		"float.valueobserver.lastvalue/C=D/R=V": -1,
		"int.valueobserver.lastvalue//R=V":      1,
		"int.valueobserver.lastvalue/A=B/R=V":   1,
	}, out.Map())
}

func TestRecordBatch(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	counter1 := Must(meter).NewInt64Counter("int64.sum")
	counter2 := Must(meter).NewFloat64Counter("float64.sum")
	valuerecorder1 := Must(meter).NewInt64ValueRecorder("int64.exact")
	valuerecorder2 := Must(meter).NewFloat64ValueRecorder("float64.exact")

	sdk.RecordBatch(
		ctx,
		[]label.KeyValue{
			label.String("A", "B"),
			label.String("C", "D"),
		},
		counter1.Measurement(1),
		counter2.Measurement(2),
		valuerecorder1.Measurement(3),
		valuerecorder2.Measurement(4),
	)

	sdk.Collect(ctx)

	out := processortest.NewOutput(label.DefaultEncoder())
	for _, rec := range processor.accumulations {
		require.NoError(t, out.AddAccumulation(rec))
	}
	require.EqualValues(t, map[string]float64{
		"int64.sum/A=B,C=D/R=V":     1,
		"float64.sum/A=B,C=D/R=V":   2,
		"int64.exact/A=B,C=D/R=V":   3,
		"float64.exact/A=B,C=D/R=V": 4,
	}, out.Map())
}

// TestRecordPersistence ensures that a direct-called instrument that
// is repeatedly used each interval results in a persistent record, so
// that its encoded labels will be cached across collection intervals.
func TestRecordPersistence(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	c := Must(meter).NewFloat64Counter("name.sum")
	b := c.Bind(label.String("bound", "true"))
	uk := label.String("bound", "false")

	for i := 0; i < 100; i++ {
		c.Add(ctx, 1, uk)
		b.Add(ctx, 1)
		sdk.Collect(ctx)
	}

	require.Equal(t, 4, processor.newAggCount)
}

func TestIncorrectInstruments(t *testing.T) {
	// The Batch observe/record APIs are susceptible to
	// uninitialized instruments.
	var counter metric.Int64Counter
	var observer metric.Int64ValueObserver

	ctx := context.Background()
	meter, sdk, _ := newSDK(t)

	// Now try with uninitialized instruments.
	meter.RecordBatch(ctx, nil, counter.Measurement(1))
	meter.NewBatchObserver(func(_ context.Context, result metric.BatchObserverResult) {
		result.Observe(nil, observer.Observation(1))
	})

	collected := sdk.Collect(ctx)
	require.Equal(t, metricsdk.ErrUninitializedInstrument, testHandler.Flush())
	require.Equal(t, 0, collected)

	// Now try with instruments from another SDK.
	var noopMeter metric.Meter
	counter = metric.Must(noopMeter).NewInt64Counter("name.sum")
	observer = metric.Must(noopMeter).NewBatchObserver(
		func(context.Context, metric.BatchObserverResult) {},
	).NewInt64ValueObserver("observer")

	meter.RecordBatch(ctx, nil, counter.Measurement(1))
	meter.NewBatchObserver(func(_ context.Context, result metric.BatchObserverResult) {
		result.Observe(nil, observer.Observation(1))
	})

	collected = sdk.Collect(ctx)
	require.Equal(t, 0, collected)
	require.Equal(t, metricsdk.ErrUninitializedInstrument, testHandler.Flush())
}

func TestSyncInAsync(t *testing.T) {
	ctx := context.Background()
	meter, sdk, processor := newSDK(t)

	counter := Must(meter).NewFloat64Counter("counter.sum")
	_ = Must(meter).NewInt64ValueObserver("observer.lastvalue",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			result.Observe(10)
			counter.Add(ctx, 100)
		},
	)

	sdk.Collect(ctx)

	out := processortest.NewOutput(label.DefaultEncoder())
	for _, rec := range processor.accumulations {
		require.NoError(t, out.AddAccumulation(rec))
	}
	require.EqualValues(t, map[string]float64{
		"counter.sum//R=V":        100,
		"observer.lastvalue//R=V": 10,
	}, out.Map())
}
