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

package pull_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/controller/controllertest"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func TestPullNoCache(t *testing.T) {
	puller := pull.New(
		basic.New(
			selector.NewWithExactDistribution(),
			export.CumulativeExportKindSelector(),
			basic.WithMemory(true),
		),
		pull.WithCachePeriod(0),
	)

	ctx := context.Background()
	meter := puller.MeterProvider().Meter("nocache")
	counter := metric.Must(meter).NewInt64Counter("counter.sum")

	counter.Add(ctx, 10, label.String("A", "B"))

	require.NoError(t, puller.Collect(ctx))
	records := processortest.NewOutput(label.DefaultEncoder())
	require.NoError(t, puller.ForEach(export.CumulativeExportKindSelector(), records.AddRecord))

	require.EqualValues(t, map[string]float64{
		"counter.sum/A=B/": 10,
	}, records.Map())

	counter.Add(ctx, 10, label.String("A", "B"))

	require.NoError(t, puller.Collect(ctx))
	records = processortest.NewOutput(label.DefaultEncoder())
	require.NoError(t, puller.ForEach(export.CumulativeExportKindSelector(), records.AddRecord))

	require.EqualValues(t, map[string]float64{
		"counter.sum/A=B/": 20,
	}, records.Map())
}

func TestPullWithCache(t *testing.T) {
	puller := pull.New(
		basic.New(
			selector.NewWithExactDistribution(),
			export.CumulativeExportKindSelector(),
			basic.WithMemory(true),
		),
		pull.WithCachePeriod(time.Second),
	)
	mock := controllertest.NewMockClock()
	puller.SetClock(mock)

	ctx := context.Background()
	meter := puller.MeterProvider().Meter("nocache")
	counter := metric.Must(meter).NewInt64Counter("counter.sum")

	counter.Add(ctx, 10, label.String("A", "B"))

	require.NoError(t, puller.Collect(ctx))
	records := processortest.NewOutput(label.DefaultEncoder())
	require.NoError(t, puller.ForEach(export.CumulativeExportKindSelector(), records.AddRecord))

	require.EqualValues(t, map[string]float64{
		"counter.sum/A=B/": 10,
	}, records.Map())

	counter.Add(ctx, 10, label.String("A", "B"))

	// Cached value!
	require.NoError(t, puller.Collect(ctx))
	records = processortest.NewOutput(label.DefaultEncoder())
	require.NoError(t, puller.ForEach(export.CumulativeExportKindSelector(), records.AddRecord))

	require.EqualValues(t, map[string]float64{
		"counter.sum/A=B/": 10,
	}, records.Map())

	mock.Add(time.Second)
	runtime.Gosched()

	// Re-computed value!
	require.NoError(t, puller.Collect(ctx))
	records = processortest.NewOutput(label.DefaultEncoder())
	require.NoError(t, puller.ForEach(export.CumulativeExportKindSelector(), records.AddRecord))

	require.EqualValues(t, map[string]float64{
		"counter.sum/A=B/": 20,
	}, records.Map())

}
