// Copyright 2019, OpenTelemetry Authors
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
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
)

type correctnessBatcher struct {
	t       *testing.T
	agg     export.Aggregator
	records []export.Record
}

type testLabelEncoder struct{}

func (cb *correctnessBatcher) AggregatorFor(*export.Descriptor) export.Aggregator {
	return cb.agg
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

func (testLabelEncoder) Encode(labels []core.KeyValue) string {
	return fmt.Sprint(labels)
}

func TestInputRangeTestCounter(t *testing.T) {
	ctx := context.Background()
	cagg := counter.New()
	batcher := &correctnessBatcher{
		t:   t,
		agg: cagg,
	}
	sdk := sdk.New(batcher, sdk.NewDefaultLabelEncoder())

	var sdkErr error
	sdk.SetErrorHandler(func(handleErr error) {
		sdkErr = handleErr
	})

	counter := sdk.NewInt64Counter("counter.name", metric.WithMonotonic(true))

	counter.Add(ctx, -1, sdk.Labels())
	require.Equal(t, aggregator.ErrNegativeInput, sdkErr)
	sdkErr = nil

	sdk.Collect(ctx)
	sum, err := cagg.Sum()
	require.Equal(t, int64(0), sum.AsInt64())
	require.Nil(t, err)

	counter.Add(ctx, 1, sdk.Labels())
	checkpointed := sdk.Collect(ctx)

	sum, err = cagg.Sum()
	require.Equal(t, int64(1), sum.AsInt64())
	require.Equal(t, 1, checkpointed)
	require.Nil(t, err)
	require.Nil(t, sdkErr)
}

func TestInputRangeTestMeasure(t *testing.T) {
	ctx := context.Background()
	magg := array.New()
	batcher := &correctnessBatcher{
		t:   t,
		agg: magg,
	}
	sdk := sdk.New(batcher, sdk.NewDefaultLabelEncoder())

	var sdkErr error
	sdk.SetErrorHandler(func(handleErr error) {
		sdkErr = handleErr
	})

	measure := sdk.NewFloat64Measure("measure.name", metric.WithAbsolute(true))

	measure.Record(ctx, -1, sdk.Labels())
	require.Equal(t, aggregator.ErrNegativeInput, sdkErr)
	sdkErr = nil

	sdk.Collect(ctx)
	count, err := magg.Count()
	require.Equal(t, int64(0), count)
	require.Nil(t, err)

	measure.Record(ctx, 1, sdk.Labels())
	measure.Record(ctx, 2, sdk.Labels())
	checkpointed := sdk.Collect(ctx)

	count, err = magg.Count()
	require.Equal(t, int64(2), count)
	require.Equal(t, 1, checkpointed)
	require.Nil(t, sdkErr)
	require.Nil(t, err)
}

func TestDisabledInstrument(t *testing.T) {
	ctx := context.Background()
	batcher := &correctnessBatcher{
		t:   t,
		agg: nil,
	}
	sdk := sdk.New(batcher, sdk.NewDefaultLabelEncoder())
	measure := sdk.NewFloat64Measure("measure.name", metric.WithAbsolute(true))

	measure.Record(ctx, -1, sdk.Labels())
	checkpointed := sdk.Collect(ctx)

	require.Equal(t, 0, checkpointed)
}

func TestRecordNaN(t *testing.T) {
	ctx := context.Background()
	batcher := &correctnessBatcher{
		t:   t,
		agg: gauge.New(),
	}
	sdk := sdk.New(batcher, sdk.NewDefaultLabelEncoder())

	var sdkErr error
	sdk.SetErrorHandler(func(handleErr error) {
		sdkErr = handleErr
	})
	g := sdk.NewFloat64Gauge("gauge.name")

	require.Nil(t, sdkErr)
	g.Set(ctx, math.NaN(), sdk.Labels())
	require.Error(t, sdkErr)
}

func TestSDKLabelEncoder(t *testing.T) {
	ctx := context.Background()
	cagg := counter.New()
	batcher := &correctnessBatcher{
		t:   t,
		agg: cagg,
	}
	sdk := sdk.New(batcher, testLabelEncoder{})

	measure := sdk.NewFloat64Measure("measure")
	measure.Record(ctx, 1, sdk.Labels(key.String("A", "B"), key.String("C", "D")))

	sdk.Collect(ctx)

	require.Equal(t, 1, len(batcher.records))

	labels := batcher.records[0].Labels()
	require.Equal(t, `[{A {8 0 B}} {C {8 0 D}}]`, labels.Encoded())
}

func TestDefaultLabelEncoder(t *testing.T) {
	encoder := sdk.NewDefaultLabelEncoder()

	encoded := encoder.Encode([]core.KeyValue{key.String("A", "B"), key.String("C", "D")})
	require.Equal(t, `A=B,C=D`, encoded)

	encoded = encoder.Encode([]core.KeyValue{key.String("A", "B,c=d"), key.String(`C\`, "D")})
	require.Equal(t, `A=B\,c\=d,C\\=D`, encoded)

	encoded = encoder.Encode([]core.KeyValue{key.String(`\`, `=`), key.String(`,`, `\`)})
	require.Equal(t, `\\=\=,\,=\\`, encoded)

	// Note: the label encoder does not sort or de-dup values,
	// that is done in Labels(...).
	encoded = encoder.Encode([]core.KeyValue{
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
	})
	require.Equal(t, "I=1,U=1,I32=1,U32=1,I64=1,U64=1,F64=1,F64=1,S=1,B=true", encoded)
}
