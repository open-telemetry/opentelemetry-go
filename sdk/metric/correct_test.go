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

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	batchTest "go.opentelemetry.io/otel/sdk/metric/batcher/test"
)

var Must = metric.Must

type correctnessBatcher struct {
	newAggCount int64

	t *testing.T

	records []export.Record
}

func (cb *correctnessBatcher) AggregatorFor(descriptor *metric.Descriptor) (agg export.Aggregator) {
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

func (cb *correctnessBatcher) CheckpointSet() export.CheckpointSet {
	cb.t.Fatal("Should not be called")
	return nil
}

func (*correctnessBatcher) FinishedCollection() {
}

func (cb *correctnessBatcher) Process(_ context.Context, record export.Record) error {
	cb.records = append(cb.records, record)
	return nil
}

func TestInputRangeTestCounter(t *testing.T) {
	ctx := context.Background()
	batcher := &correctnessBatcher{
		t: t,
	}
	sdk := metricsdk.New(batcher)
	meter := metric.WrapMeterImpl(sdk, "test")

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

	batcher.records = nil
	counter.Add(ctx, 1)
	checkpointed = sdk.Collect(ctx)
	sum, err := batcher.records[0].Aggregator().(aggregator.Sum).Sum()
	require.Equal(t, int64(1), sum.AsInt64())
	require.Equal(t, 1, checkpointed)
	require.Nil(t, err)
	require.Nil(t, sdkErr)
}

func TestInputRangeTestMeasure(t *testing.T) {
	ctx := context.Background()
	batcher := &correctnessBatcher{
		t: t,
	}
	sdk := metricsdk.New(batcher)
	meter := metric.WrapMeterImpl(sdk, "test")

	var sdkErr error
	sdk.SetErrorHandler(func(handleErr error) {
		sdkErr = handleErr
	})

	measure := Must(meter).NewFloat64Measure("name.measure")

	measure.Record(ctx, math.NaN())
	require.Equal(t, aggregator.ErrNaNInput, sdkErr)
	sdkErr = nil

	checkpointed := sdk.Collect(ctx)
	require.Equal(t, 0, checkpointed)

	measure.Record(ctx, 1)
	measure.Record(ctx, 2)

	batcher.records = nil
	checkpointed = sdk.Collect(ctx)

	count, err := batcher.records[0].Aggregator().(aggregator.Distribution).Count()
	require.Equal(t, int64(2), count)
	require.Equal(t, 1, checkpointed)
	require.Nil(t, sdkErr)
	require.Nil(t, err)
}

func TestDisabledInstrument(t *testing.T) {
	ctx := context.Background()
	batcher := &correctnessBatcher{
		t: t,
	}
	sdk := metricsdk.New(batcher)
	meter := metric.WrapMeterImpl(sdk, "test")

	measure := Must(meter).NewFloat64Measure("name.disabled")

	measure.Record(ctx, -1)
	checkpointed := sdk.Collect(ctx)

	require.Equal(t, 0, checkpointed)
	require.Equal(t, 0, len(batcher.records))
}

func TestRecordNaN(t *testing.T) {
	ctx := context.Background()
	batcher := &correctnessBatcher{
		t: t,
	}

	sdk := metricsdk.New(batcher)
	meter := metric.WrapMeterImpl(sdk, "test")

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
	batcher := &correctnessBatcher{
		t: t,
	}
	sdk := metricsdk.New(batcher)
	meter := metric.WrapMeterImpl(sdk, "test")

	counter := Must(meter).NewInt64Counter("counter")

	const (
		maxKeys = 21
		keySets = 2
		repeats = 3
	)
	var keysA []core.Key
	var keysB []core.Key

	for i := 0; i < maxKeys; i++ {
		keysA = append(keysA, core.Key(fmt.Sprintf("A%03d", i)))
		keysB = append(keysB, core.Key(fmt.Sprintf("B%03d", i)))
	}

	var allExpect [][]core.KeyValue
	for numKeys := 0; numKeys < maxKeys; numKeys++ {

		var kvsA []core.KeyValue
		var kvsB []core.KeyValue
		for r := 0; r < repeats; r++ {
			for i := 0; i < numKeys; i++ {
				kvsA = append(kvsA, keysA[i].Int(r))
				kvsB = append(kvsB, keysB[i].Int(r))
			}
		}

		var expectA []core.KeyValue
		var expectB []core.KeyValue
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

	var actual [][]core.KeyValue
	for _, rec := range batcher.records {
		sum, _ := rec.Aggregator().(aggregator.Sum).Sum()
		require.Equal(t, sum, metric.NewInt64Number(2))

		kvs := rec.Labels().ToSlice()
		actual = append(actual, kvs)
	}

	require.ElementsMatch(t, allExpect, actual)
}

func newSetIter(kvs ...core.KeyValue) label.Iterator {
	labels := label.NewSet(kvs...)
	return labels.Iter()
}

func TestDefaultLabelEncoder(t *testing.T) {
	encoder := label.DefaultEncoder()

	encoded := encoder.Encode(newSetIter(key.String("A", "B"), key.String("C", "D")))
	require.Equal(t, `A=B,C=D`, encoded)

	encoded = encoder.Encode(newSetIter(key.String("A", "B,c=d"), key.String(`C\`, "D")))
	require.Equal(t, `A=B\,c\=d,C\\=D`, encoded)

	encoded = encoder.Encode(newSetIter(key.String(`\`, `=`), key.String(`,`, `\`)))
	require.Equal(t, `\,=\\,\\=\=`, encoded)

	// Note: the label encoder does not sort or de-dup values,
	// that is done in Labels(...).
	encoded = encoder.Encode(newSetIter(
		key.Int("I", 1),
		key.Uint("U", 1),
		key.Int32("I32", 1),
		key.Uint32("U32", 1),
		key.Int64("I64", 1),
		key.Uint64("U64", 1),
		key.Float64("F64", 1),
		key.Float64("F64", 1),
		key.String("S", "1"),
		key.Bool("B", true),
	))
	require.Equal(t, "B=true,F64=1,I=1,I32=1,I64=1,S=1,U=1,U32=1,U64=1", encoded)
}

func TestObserverCollection(t *testing.T) {
	ctx := context.Background()
	batcher := &correctnessBatcher{
		t: t,
	}

	sdk := metricsdk.New(batcher)
	meter := metric.WrapMeterImpl(sdk, "test")

	_ = Must(meter).RegisterFloat64Observer("float.observer", func(result metric.Float64ObserverResult) {
		result.Observe(1, key.String("A", "B"))
		// last value wins
		result.Observe(-1, key.String("A", "B"))
		result.Observe(-1, key.String("C", "D"))
	})
	_ = Must(meter).RegisterInt64Observer("int.observer", func(result metric.Int64ObserverResult) {
		result.Observe(-1, key.String("A", "B"))
		result.Observe(1)
		// last value wins
		result.Observe(1, key.String("A", "B"))
		result.Observe(1)
	})
	_ = Must(meter).RegisterInt64Observer("empty.observer", func(result metric.Int64ObserverResult) {
	})

	collected := sdk.Collect(ctx)

	require.Equal(t, 4, collected)
	require.Equal(t, 4, len(batcher.records))

	out := batchTest.NewOutput(label.DefaultEncoder())
	for _, rec := range batcher.records {
		_ = out.AddTo(rec)
	}
	require.EqualValues(t, map[string]float64{
		"float.observer/A=B": -1,
		"float.observer/C=D": -1,
		"int.observer/":      1,
		"int.observer/A=B":   1,
	}, out.Map)
}

func TestRecordBatch(t *testing.T) {
	ctx := context.Background()
	batcher := &correctnessBatcher{
		t: t,
	}

	sdk := metricsdk.New(batcher)
	meter := metric.WrapMeterImpl(sdk, "test")

	counter1 := Must(meter).NewInt64Counter("int64.counter")
	counter2 := Must(meter).NewFloat64Counter("float64.counter")
	measure1 := Must(meter).NewInt64Measure("int64.measure")
	measure2 := Must(meter).NewFloat64Measure("float64.measure")

	sdk.RecordBatch(
		ctx,
		[]core.KeyValue{
			key.String("A", "B"),
			key.String("C", "D"),
		},
		counter1.Measurement(1),
		counter2.Measurement(2),
		measure1.Measurement(3),
		measure2.Measurement(4),
	)

	sdk.Collect(ctx)

	out := batchTest.NewOutput(label.DefaultEncoder())
	for _, rec := range batcher.records {
		_ = out.AddTo(rec)
	}
	require.EqualValues(t, map[string]float64{
		"int64.counter/A=B,C=D":   1,
		"float64.counter/A=B,C=D": 2,
		"int64.measure/A=B,C=D":   3,
		"float64.measure/A=B,C=D": 4,
	}, out.Map)
}

// TestRecordPersistence ensures that a direct-called instrument that
// is repeatedly used each interval results in a persistent record, so
// that its encoded labels will be cached across collection intervals.
func TestRecordPersistence(t *testing.T) {
	ctx := context.Background()
	batcher := &correctnessBatcher{
		t: t,
	}

	sdk := metricsdk.New(batcher)
	meter := metric.WrapMeterImpl(sdk, "test")

	c := Must(meter).NewFloat64Counter("sum.name")
	b := c.Bind(key.String("bound", "true"))
	uk := key.String("bound", "false")

	for i := 0; i < 100; i++ {
		c.Add(ctx, 1, uk)
		b.Add(ctx, 1)
		sdk.Collect(ctx)
	}

	require.Equal(t, int64(2), batcher.newAggCount)
}
