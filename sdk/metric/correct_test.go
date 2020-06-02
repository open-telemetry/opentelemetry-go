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
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	batchTest "go.opentelemetry.io/otel/sdk/metric/integrator/test"
	"go.opentelemetry.io/otel/sdk/resource"
)

var Must = metric.Must
var testResource = resource.New(kv.String("R", "V"))

type correctnessIntegrator struct {
	newAggCount int64

	t *testing.T

	records []export.Record

	sync.Mutex
	err error
}

func newSDK(t *testing.T) (metric.Meter, *metricsdk.Accumulator, *correctnessIntegrator) {
	integrator := &correctnessIntegrator{
		t: t,
	}
	accum := metricsdk.NewAccumulator(
		integrator,
		metricsdk.WithResource(testResource),
		metricsdk.WithErrorHandler(func(err error) {
			integrator.Lock()
			defer integrator.Unlock()
			integrator.err = err
		}),
	)
	meter := metric.WrapMeterImpl(accum, "test")
	return meter, accum, integrator
}

func (ci *correctnessIntegrator) sdkErr() error {
	ci.Lock()
	defer ci.Unlock()
	t := ci.err
	ci.err = nil
	return t
}

func (ci *correctnessIntegrator) AggregatorFor(descriptor *metric.Descriptor) (agg export.Aggregator) {
	name := descriptor.Name()

	switch {
	case strings.HasSuffix(name, ".counter"):
		agg = sum.New()
	case strings.HasSuffix(name, ".disabled"):
		agg = nil
	default:
		agg = array.New()
	}
	if agg != nil {
		atomic.AddInt64(&ci.newAggCount, 1)
	}
	return
}

func (ci *correctnessIntegrator) CheckpointSet() export.CheckpointSet {
	ci.t.Fatal("Should not be called")
	return nil
}

func (*correctnessIntegrator) FinishedCollection() {
}

func (ci *correctnessIntegrator) Process(_ context.Context, record export.Record) error {
	ci.records = append(ci.records, record)
	return nil
}

func TestInputRangeCounter(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	var sdkErr error
	sdk.SetErrorHandler(func(handleErr error) {
		sdkErr = handleErr
	})

	counter := Must(meter).NewInt64Counter("name.counter")

	counter.Add(ctx, -1)
	require.Equal(t, aggregator.ErrNegativeInput, sdkErr)
	sdkErr = nil

	checkpointed := sdk.Collect(ctx)
	require.Equal(t, 0, checkpointed)

	integrator.records = nil
	counter.Add(ctx, 1)
	checkpointed = sdk.Collect(ctx)
	sum, err := integrator.records[0].Aggregator().(aggregator.Sum).Sum()
	require.Equal(t, int64(1), sum.AsInt64())
	require.Equal(t, 1, checkpointed)
	require.Nil(t, err)
	require.Nil(t, sdkErr)
}

func TestInputRangeUpDownCounter(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	var sdkErr error
	sdk.SetErrorHandler(func(handleErr error) {
		sdkErr = handleErr
	})

	counter := Must(meter).NewInt64UpDownCounter("name.updowncounter")

	counter.Add(ctx, -1)
	counter.Add(ctx, -1)
	counter.Add(ctx, 2)
	counter.Add(ctx, 1)

	checkpointed := sdk.Collect(ctx)
	sum, err := integrator.records[0].Aggregator().(aggregator.Sum).Sum()
	require.Equal(t, int64(1), sum.AsInt64())
	require.Equal(t, 1, checkpointed)
	require.Nil(t, err)
	require.Nil(t, sdkErr)
}

func TestInputRangeValueRecorder(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	var sdkErr error
	sdk.SetErrorHandler(func(handleErr error) {
		sdkErr = handleErr
	})

	valuerecorder := Must(meter).NewFloat64ValueRecorder("name.valuerecorder")

	valuerecorder.Record(ctx, math.NaN())
	require.Equal(t, aggregator.ErrNaNInput, sdkErr)
	sdkErr = nil

	checkpointed := sdk.Collect(ctx)
	require.Equal(t, 0, checkpointed)

	valuerecorder.Record(ctx, 1)
	valuerecorder.Record(ctx, 2)

	integrator.records = nil
	checkpointed = sdk.Collect(ctx)

	count, err := integrator.records[0].Aggregator().(aggregator.Distribution).Count()
	require.Equal(t, int64(2), count)
	require.Equal(t, 1, checkpointed)
	require.Nil(t, sdkErr)
	require.Nil(t, err)
}

func TestDisabledInstrument(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	valuerecorder := Must(meter).NewFloat64ValueRecorder("name.disabled")

	valuerecorder.Record(ctx, -1)
	checkpointed := sdk.Collect(ctx)

	require.Equal(t, 0, checkpointed)
	require.Equal(t, 0, len(integrator.records))
}

func TestRecordNaN(t *testing.T) {
	ctx := context.Background()
	meter, sdk, _ := newSDK(t)

	var sdkErr error
	sdk.SetErrorHandler(func(handleErr error) {
		sdkErr = handleErr
	})
	c := Must(meter).NewFloat64Counter("sum.name")

	require.Nil(t, sdkErr)
	c.Add(ctx, math.NaN())
	require.Error(t, sdkErr)
}

func TestSDKLabelsDeduplication(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	counter := Must(meter).NewInt64Counter("counter")

	const (
		maxKeys = 21
		keySets = 2
		repeats = 3
	)
	var keysA []kv.Key
	var keysB []kv.Key

	for i := 0; i < maxKeys; i++ {
		keysA = append(keysA, kv.Key(fmt.Sprintf("A%03d", i)))
		keysB = append(keysB, kv.Key(fmt.Sprintf("B%03d", i)))
	}

	var allExpect [][]kv.KeyValue
	for numKeys := 0; numKeys < maxKeys; numKeys++ {

		var kvsA []kv.KeyValue
		var kvsB []kv.KeyValue
		for r := 0; r < repeats; r++ {
			for i := 0; i < numKeys; i++ {
				kvsA = append(kvsA, keysA[i].Int(r))
				kvsB = append(kvsB, keysB[i].Int(r))
			}
		}

		var expectA []kv.KeyValue
		var expectB []kv.KeyValue
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

	var actual [][]kv.KeyValue
	for _, rec := range integrator.records {
		sum, _ := rec.Aggregator().(aggregator.Sum).Sum()
		require.Equal(t, sum, metric.NewInt64Number(2))

		kvs := rec.Labels().ToSlice()
		actual = append(actual, kvs)
	}

	require.ElementsMatch(t, allExpect, actual)
}

func newSetIter(kvs ...kv.KeyValue) label.Iterator {
	labels := label.NewSet(kvs...)
	return labels.Iter()
}

func TestDefaultLabelEncoder(t *testing.T) {
	encoder := label.DefaultEncoder()

	encoded := encoder.Encode(newSetIter(kv.String("A", "B"), kv.String("C", "D")))
	require.Equal(t, `A=B,C=D`, encoded)

	encoded = encoder.Encode(newSetIter(kv.String("A", "B,c=d"), kv.String(`C\`, "D")))
	require.Equal(t, `A=B\,c\=d,C\\=D`, encoded)

	encoded = encoder.Encode(newSetIter(kv.String(`\`, `=`), kv.String(`,`, `\`)))
	require.Equal(t, `\,=\\,\\=\=`, encoded)

	// Note: the label encoder does not sort or de-dup values,
	// that is done in Labels(...).
	encoded = encoder.Encode(newSetIter(
		kv.Int("I", 1),
		kv.Uint("U", 1),
		kv.Int32("I32", 1),
		kv.Uint32("U32", 1),
		kv.Int64("I64", 1),
		kv.Uint64("U64", 1),
		kv.Float64("F64", 1),
		kv.Float64("F64", 1),
		kv.String("S", "1"),
		kv.Bool("B", true),
	))
	require.Equal(t, "B=true,F64=1,I=1,I32=1,I64=1,S=1,U=1,U32=1,U64=1", encoded)
}

func TestObserverCollection(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	_ = Must(meter).NewFloat64ValueObserver("float.valueobserver", func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1, kv.String("A", "B"))
		// last value wins
		result.Observe(-1, kv.String("A", "B"))
		result.Observe(-1, kv.String("C", "D"))
	})
	_ = Must(meter).NewInt64ValueObserver("int.valueobserver", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(-1, kv.String("A", "B"))
		result.Observe(1)
		// last value wins
		result.Observe(1, kv.String("A", "B"))
		result.Observe(1)
	})

	_ = Must(meter).NewFloat64SumObserver("float.sumobserver", func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1, kv.String("A", "B"))
		result.Observe(2, kv.String("A", "B"))
		result.Observe(1, kv.String("C", "D"))
	})
	_ = Must(meter).NewInt64SumObserver("int.sumobserver", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(2, kv.String("A", "B"))
		result.Observe(1)
		// last value wins
		result.Observe(1, kv.String("A", "B"))
		result.Observe(1)
	})

	_ = Must(meter).NewFloat64UpDownSumObserver("float.updownsumobserver", func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1, kv.String("A", "B"))
		result.Observe(-2, kv.String("A", "B"))
		result.Observe(1, kv.String("C", "D"))
	})
	_ = Must(meter).NewInt64UpDownSumObserver("int.updownsumobserver", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(2, kv.String("A", "B"))
		result.Observe(1)
		// last value wins
		result.Observe(1, kv.String("A", "B"))
		result.Observe(-1)
	})

	_ = Must(meter).NewInt64ValueObserver("empty.valueobserver", func(_ context.Context, result metric.Int64ObserverResult) {
	})

	collected := sdk.Collect(ctx)

	require.Equal(t, collected, len(integrator.records))

	out := batchTest.NewOutput(label.DefaultEncoder())
	for _, rec := range integrator.records {
		_ = out.AddTo(rec)
	}
	require.EqualValues(t, map[string]float64{
		"float.valueobserver/A=B/R=V": -1,
		"float.valueobserver/C=D/R=V": -1,
		"int.valueobserver//R=V":      1,
		"int.valueobserver/A=B/R=V":   1,

		"float.sumobserver/A=B/R=V": 2,
		"float.sumobserver/C=D/R=V": 1,
		"int.sumobserver//R=V":      1,
		"int.sumobserver/A=B/R=V":   1,

		"float.updownsumobserver/A=B/R=V": -2,
		"float.updownsumobserver/C=D/R=V": 1,
		"int.updownsumobserver//R=V":      -1,
		"int.updownsumobserver/A=B/R=V":   1,
	}, out.Map)
}

func TestSumObserverInputRange(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	_ = Must(meter).NewFloat64SumObserver("float.sumobserver", func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(-2, kv.String("A", "B"))
		require.Equal(t, aggregator.ErrNegativeInput, integrator.sdkErr())
		result.Observe(-1, kv.String("C", "D"))
		require.Equal(t, aggregator.ErrNegativeInput, integrator.sdkErr())
	})
	_ = Must(meter).NewInt64SumObserver("int.sumobserver", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(-1, kv.String("A", "B"))
		require.Equal(t, aggregator.ErrNegativeInput, integrator.sdkErr())
		result.Observe(-1)
		require.Equal(t, aggregator.ErrNegativeInput, integrator.sdkErr())
	})

	collected := sdk.Collect(ctx)

	require.Equal(t, 0, collected)
	require.Equal(t, 0, len(integrator.records))

	// check that the error condition was reset
	require.NoError(t, integrator.sdkErr())
}

func TestObserverBatch(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	var floatValueObs metric.Float64ValueObserver
	var intValueObs metric.Int64ValueObserver
	var floatSumObs metric.Float64SumObserver
	var intSumObs metric.Int64SumObserver
	var floatUpDownSumObs metric.Float64UpDownSumObserver
	var intUpDownSumObs metric.Int64UpDownSumObserver

	var batch = Must(meter).NewBatchObserver(
		func(_ context.Context, result metric.BatchObserverResult) {
			result.Observe(
				[]kv.KeyValue{
					kv.String("A", "B"),
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
				[]kv.KeyValue{
					kv.String("C", "D"),
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
	floatValueObs = batch.NewFloat64ValueObserver("float.valueobserver")
	intValueObs = batch.NewInt64ValueObserver("int.valueobserver")
	floatSumObs = batch.NewFloat64SumObserver("float.sumobserver")
	intSumObs = batch.NewInt64SumObserver("int.sumobserver")
	floatUpDownSumObs = batch.NewFloat64UpDownSumObserver("float.updownsumobserver")
	intUpDownSumObs = batch.NewInt64UpDownSumObserver("int.updownsumobserver")

	collected := sdk.Collect(ctx)

	require.Equal(t, collected, len(integrator.records))

	out := batchTest.NewOutput(label.DefaultEncoder())
	for _, rec := range integrator.records {
		_ = out.AddTo(rec)
	}
	require.EqualValues(t, map[string]float64{
		"float.sumobserver//R=V":    1.1,
		"float.sumobserver/A=B/R=V": 1000,
		"int.sumobserver//R=V":      10,
		"int.sumobserver/A=B/R=V":   100,

		"int.updownsumobserver/A=B/R=V":   -100,
		"float.updownsumobserver/A=B/R=V": -1000,
		"int.updownsumobserver//R=V":      10,
		"float.updownsumobserver/C=D/R=V": -1,

		"float.valueobserver/A=B/R=V": -1,
		"float.valueobserver/C=D/R=V": -1,
		"int.valueobserver//R=V":      1,
		"int.valueobserver/A=B/R=V":   1,
	}, out.Map)
}

func TestRecordBatch(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	counter1 := Must(meter).NewInt64Counter("int64.counter")
	counter2 := Must(meter).NewFloat64Counter("float64.counter")
	valuerecorder1 := Must(meter).NewInt64ValueRecorder("int64.valuerecorder")
	valuerecorder2 := Must(meter).NewFloat64ValueRecorder("float64.valuerecorder")

	sdk.RecordBatch(
		ctx,
		[]kv.KeyValue{
			kv.String("A", "B"),
			kv.String("C", "D"),
		},
		counter1.Measurement(1),
		counter2.Measurement(2),
		valuerecorder1.Measurement(3),
		valuerecorder2.Measurement(4),
	)

	sdk.Collect(ctx)

	out := batchTest.NewOutput(label.DefaultEncoder())
	for _, rec := range integrator.records {
		_ = out.AddTo(rec)
	}
	require.EqualValues(t, map[string]float64{
		"int64.counter/A=B,C=D/R=V":         1,
		"float64.counter/A=B,C=D/R=V":       2,
		"int64.valuerecorder/A=B,C=D/R=V":   3,
		"float64.valuerecorder/A=B,C=D/R=V": 4,
	}, out.Map)
}

// TestRecordPersistence ensures that a direct-called instrument that
// is repeatedly used each interval results in a persistent record, so
// that its encoded labels will be cached across collection intervals.
func TestRecordPersistence(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	c := Must(meter).NewFloat64Counter("sum.name")
	b := c.Bind(kv.String("bound", "true"))
	uk := kv.String("bound", "false")

	for i := 0; i < 100; i++ {
		c.Add(ctx, 1, uk)
		b.Add(ctx, 1)
		sdk.Collect(ctx)
	}

	require.Equal(t, int64(2), integrator.newAggCount)
}

func TestIncorrectInstruments(t *testing.T) {
	// The Batch observe/record APIs are susceptible to
	// uninitialized instruments.
	var counter metric.Int64Counter
	var observer metric.Int64ValueObserver

	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	// Now try with uninitialized instruments.
	meter.RecordBatch(ctx, nil, counter.Measurement(1))
	meter.NewBatchObserver(func(_ context.Context, result metric.BatchObserverResult) {
		result.Observe(nil, observer.Observation(1))
	})

	collected := sdk.Collect(ctx)
	require.Equal(t, metricsdk.ErrUninitializedInstrument, integrator.sdkErr())
	require.Equal(t, 0, collected)

	// Now try with instruments from another SDK.
	var noopMeter metric.Meter
	counter = metric.Must(noopMeter).NewInt64Counter("counter")
	observer = metric.Must(noopMeter).NewBatchObserver(
		func(context.Context, metric.BatchObserverResult) {},
	).NewInt64ValueObserver("observer")

	meter.RecordBatch(ctx, nil, counter.Measurement(1))
	meter.NewBatchObserver(func(_ context.Context, result metric.BatchObserverResult) {
		result.Observe(nil, observer.Observation(1))
	})

	collected = sdk.Collect(ctx)
	require.Equal(t, 0, collected)
	require.Equal(t, metricsdk.ErrUninitializedInstrument, integrator.sdkErr())
}

func TestSyncInAsync(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	counter := Must(meter).NewFloat64Counter("counter")
	_ = Must(meter).NewInt64ValueObserver("observer",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			result.Observe(10)
			counter.Add(ctx, 100)
		},
	)

	sdk.Collect(ctx)

	out := batchTest.NewOutput(label.DefaultEncoder())
	for _, rec := range integrator.records {
		_ = out.AddTo(rec)
	}
	require.EqualValues(t, map[string]float64{
		"counter//R=V":  100,
		"observer//R=V": 10,
	}, out.Map)
}
