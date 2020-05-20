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

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
	controllerTest "go.opentelemetry.io/otel/sdk/metric/controller/test"
	"go.opentelemetry.io/otel/sdk/metric/integrator/test"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func TestPullNoCache(t *testing.T) {
	puller := pull.New(
		selector.NewWithExactDistribution(),
		pull.WithCachePeriod(0),
		pull.WithStateful(true),
	)

	ctx := context.Background()
	meter := puller.Provider().Meter("nocache")
	counter := metric.Must(meter).NewInt64Counter("counter")

	counter.Add(ctx, 10, kv.String("A", "B"))

	puller.Collect(ctx)
	records := test.NewOutput(label.DefaultEncoder())
	_ = puller.ForEach(records.AddTo)

	require.EqualValues(t, map[string]float64{
		"counter/A=B/": 10,
	}, records.Map)

	counter.Add(ctx, 10, kv.String("A", "B"))

	puller.Collect(ctx)
	records = test.NewOutput(label.DefaultEncoder())
	_ = puller.ForEach(records.AddTo)

	require.EqualValues(t, map[string]float64{
		"counter/A=B/": 20,
	}, records.Map)
}

func TestPullWithCache(t *testing.T) {
	puller := pull.New(
		selector.NewWithExactDistribution(),
		pull.WithCachePeriod(time.Second),
		pull.WithStateful(true),
	)
	mock := controllerTest.NewMockClock()
	puller.SetClock(mock)

	ctx := context.Background()
	meter := puller.Provider().Meter("nocache")
	counter := metric.Must(meter).NewInt64Counter("counter")

	counter.Add(ctx, 10, kv.String("A", "B"))

	puller.Collect(ctx)
	records := test.NewOutput(label.DefaultEncoder())
	_ = puller.ForEach(records.AddTo)

	require.EqualValues(t, map[string]float64{
		"counter/A=B/": 10,
	}, records.Map)

	counter.Add(ctx, 10, kv.String("A", "B"))

	// Cached value!
	puller.Collect(ctx)
	records = test.NewOutput(label.DefaultEncoder())
	_ = puller.ForEach(records.AddTo)

	require.EqualValues(t, map[string]float64{
		"counter/A=B/": 10,
	}, records.Map)

	mock.Add(time.Second)
	runtime.Gosched()

	// Re-computed value!
	puller.Collect(ctx)
	records = test.NewOutput(label.DefaultEncoder())
	_ = puller.ForEach(records.AddTo)

	require.EqualValues(t, map[string]float64{
		"counter/A=B/": 20,
	}, records.Map)

}
