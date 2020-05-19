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
}

func newSDK(t *testing.T) (metric.Meter, *metricsdk.Accumulator, *correctnessIntegrator) {
	integrator := &correctnessIntegrator{
		t: t,
	}
	accum := metricsdk.NewAccumulator(integrator, metricsdk.WithResource(testResource))
	meter := metric.WrapMeterImpl(accum, "test")
	return meter, accum, integrator
}

func (cb *correctnessIntegrator) AggregatorFor(descriptor *metric.Descriptor) (agg export.Aggregator) {
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
		atomic.AddInt64(&cb.newAggCount, 1)
	}
	return
}

func (cb *correctnessIntegrator) CheckpointSet() export.CheckpointSet {
	cb.t.Fatal("Should not be called")
	return nil
}

func (*correctnessIntegrator) FinishedCollection() {
}

func (cb *correctnessIntegrator) Process(_ context.Context, record export.Record) error {
	cb.records = append(cb.records, record)
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

	_ = Must(meter).RegisterFloat64ValueObserver("float.valueobserver", func(result metric.Float64ObserverResult) {
		result.Observe(1, kv.String("A", "B"))
		// last value wins
		result.Observe(-1, kv.String("A", "B"))
		result.Observe(-1, kv.String("C", "D"))
	})
	_ = Must(meter).RegisterInt64ValueObserver("int.valueobserver", func(result metric.Int64ObserverResult) {
		result.Observe(-1, kv.String("A", "B"))
		result.Observe(1)
		// last value wins
		result.Observe(1, kv.String("A", "B"))
		result.Observe(1)
	})
	_ = Must(meter).RegisterInt64ValueObserver("empty.valueobserver", func(result metric.Int64ObserverResult) {
	})

	collected := sdk.Collect(ctx)

	require.Equal(t, 4, collected)
	require.Equal(t, 4, len(integrator.records))

	out := batchTest.NewOutput(label.DefaultEncoder())
	for _, rec := range integrator.records {
		_ = out.AddTo(rec)
	}
	require.EqualValues(t, map[string]float64{
		"float.valueobserver/A=B/R=V": -1,
		"float.valueobserver/C=D/R=V": -1,
		"int.valueobserver//R=V":      1,
		"int.valueobserver/A=B/R=V":   1,
	}, out.Map)
}

func TestObserverBatch(t *testing.T) {
	ctx := context.Background()
	meter, sdk, integrator := newSDK(t)

	var floatObs metric.Float64ValueObserver
	var intObs metric.Int64ValueObserver
	var batch = Must(meter).NewBatchObserver(
		func(result metric.BatchObserverResult) {
			result.Observe(
				[]kv.KeyValue{
					kv.String("A", "B"),
				},
				floatObs.Observation(1),
				floatObs.Observation(-1),
				intObs.Observation(-1),
				intObs.Observation(1),
			)
			result.Observe(
				[]kv.KeyValue{
					kv.String("C", "D"),
				},
				floatObs.Observation(-1),
			)
			result.Observe(
				nil,
				intObs.Observation(1),
				intObs.Observation(1),
			)
		})
	floatObs = batch.RegisterFloat64ValueObserver("float.valueobserver")
	intObs = batch.RegisterInt64ValueObserver("int.valueobserver")

	collected := sdk.Collect(ctx)

	require.Equal(t, 4, collected)
	require.Equal(t, 4, len(integrator.records))

	out := batchTest.NewOutput(label.DefaultEncoder())
	for _, rec := range integrator.records {
		_ = out.AddTo(rec)
	}
	require.EqualValues(t, map[string]float64{
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
