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
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/export"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
)

type monotoneBatcher struct {
	t *testing.T

	collections  int
	currentValue *core.Number
	currentTime  *time.Time
}

func (m *monotoneBatcher) AggregatorFor(rec export.MetricRecord) export.MetricAggregator {
	return gauge.New()
}

func (m *monotoneBatcher) Export(_ context.Context, record export.MetricRecord, agg export.MetricAggregator) {
	require.Equal(m.t, "my.gauge.name", record.Descriptor().Name())
	require.Equal(m.t, 1, len(record.Labels()))
	require.Equal(m.t, "a", string(record.Labels()[0].Key))
	require.Equal(m.t, "b", record.Labels()[0].Value.Emit())

	gauge := agg.(*gauge.Aggregator)
	val := gauge.AsNumber()
	ts := gauge.Timestamp()

	m.currentValue = &val
	m.currentTime = &ts
	m.collections++
}

func TestMonotoneGauge(t *testing.T) {
	ctx := context.Background()
	batcher := &monotoneBatcher{
		t: t,
	}
	sdk := sdk.New(batcher)

	gauge := sdk.NewInt64Gauge("my.gauge.name", metric.WithMonotonic(true))

	handle := gauge.AcquireHandle(sdk.Labels(key.String("a", "b")))

	require.Nil(t, batcher.currentTime)
	require.Nil(t, batcher.currentValue)

	before := time.Now()

	handle.Set(ctx, 1)

	// Until collection, expect nil.
	require.Nil(t, batcher.currentTime)
	require.Nil(t, batcher.currentValue)

	sdk.Collect(ctx)

	require.NotNil(t, batcher.currentValue)
	require.Equal(t, core.NewInt64Number(1), *batcher.currentValue)
	require.True(t, before.Before(*batcher.currentTime))

	before = *batcher.currentTime

	// Collect would ordinarily flush the record, except we're using a handle.
	sdk.Collect(ctx)

	require.Equal(t, 2, batcher.collections)

	// Increase the value to 2.
	handle.Set(ctx, 2)

	sdk.Collect(ctx)

	require.Equal(t, 3, batcher.collections)
	require.Equal(t, core.NewInt64Number(2), *batcher.currentValue)
	require.True(t, before.Before(*batcher.currentTime))

	before = *batcher.currentTime

	sdk.Collect(ctx)
	require.Equal(t, 4, batcher.collections)

	// Try to lower the value to 1, it will fail.
	handle.Set(ctx, 1)
	sdk.Collect(ctx)

	// The value and timestamp are both unmodified
	require.Equal(t, 5, batcher.collections)
	require.Equal(t, core.NewInt64Number(2), *batcher.currentValue)
	require.Equal(t, before, *batcher.currentTime)

	// Update with the same value, update the timestamp.
	handle.Set(ctx, 2)
	sdk.Collect(ctx)

	require.Equal(t, 6, batcher.collections)
	require.Equal(t, core.NewInt64Number(2), *batcher.currentValue)
	require.True(t, before.Before(*batcher.currentTime))
}
