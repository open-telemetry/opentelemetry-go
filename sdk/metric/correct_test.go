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
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
)

type correctnessBatcher struct {
	t   *testing.T
	agg export.Aggregator
}

func (cb *correctnessBatcher) AggregatorFor(*export.Descriptor) export.Aggregator {
	return cb.agg
}

func (cb *correctnessBatcher) ReadCheckpoint() export.Producer {
	cb.t.Fatal("Should not be called")
	return nil
}

func (cb *correctnessBatcher) Process(_ context.Context, desc *export.Descriptor, labels export.Labels, agg export.Aggregator) error {
	return nil
}

func TestInputRangeTestCounter(t *testing.T) {
	ctx := context.Background()
	cagg := counter.New()
	batcher := &correctnessBatcher{
		t:   t,
		agg: cagg,
	}
	sdk := sdk.New(batcher, sdk.DefaultLabelEncoder())
	counter := sdk.NewInt64Counter("counter.name", metric.WithMonotonic(true))

	counter.Add(ctx, -1, sdk.Labels())
	sdk.Collect(ctx)
	require.Equal(t, int64(0), cagg.Sum().AsInt64())

	counter.Add(ctx, 1, sdk.Labels())
	sdk.Collect(ctx)
	require.Equal(t, int64(1), cagg.Sum().AsInt64())
}

func TestInputRangeTestMeasure(t *testing.T) {
	ctx := context.Background()
	magg := array.New()
	batcher := &correctnessBatcher{
		t:   t,
		agg: magg,
	}
	sdk := sdk.New(batcher, sdk.DefaultLabelEncoder())
	measure := sdk.NewFloat64Measure("measure.name", metric.WithAbsolute(true))

	measure.Record(ctx, -1, sdk.Labels())
	sdk.Collect(ctx)
	require.Equal(t, int64(0), magg.Count())

	measure.Record(ctx, 1, sdk.Labels())
	measure.Record(ctx, 2, sdk.Labels())
	sdk.Collect(ctx)
	require.Equal(t, int64(2), magg.Count())
}
