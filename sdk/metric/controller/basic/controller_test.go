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

package basic_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	ottest "go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/controller/controllertest"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/resource"
)

const envVar = "OTEL_RESOURCE_ATTRIBUTES"

func getMap(t *testing.T, cont *controller.Controller) map[string]float64 {
	out := processortest.NewOutput(attribute.DefaultEncoder())

	require.NoError(t, cont.ForEach(
		export.CumulativeExportKindSelector(),
		func(record export.Record) error {
			return out.AddRecord(record)
		},
	))
	return out.Map()
}

type testContextKey string

func testContext() context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, testContextKey("A"), "B")
}

func checkTestContext(t *testing.T, ctx context.Context) {
	require.Equal(t, "B", ctx.Value(testContextKey("A")))
}

func TestControllerUsesResource(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		envVar: "key=value,T=U",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	cases := []struct {
		name    string
		options []controller.Option
		wanted  string
	}{
		{
			name:    "explicitly empty resource",
			options: []controller.Option{controller.WithResource(resource.Empty())},
			wanted:  resource.Environment().Encoded(attribute.DefaultEncoder())},
		{
			name:    "uses default if no resource option",
			options: nil,
			wanted:  resource.Default().Encoded(attribute.DefaultEncoder())},
		{
			name:    "explicit resource",
			options: []controller.Option{controller.WithResource(resource.NewWithAttributes(attribute.String("R", "S")))},
			wanted:  "R=S,T=U,key=value"},
		{
			name: "last resource wins",
			options: []controller.Option{
				controller.WithResource(resource.Default()),
				controller.WithResource(resource.NewWithAttributes(attribute.String("R", "S"))),
			},
			wanted: "R=S,T=U,key=value"},
		{
			name:    "overlapping attributes with environment resource",
			options: []controller.Option{controller.WithResource(resource.NewWithAttributes(attribute.String("T", "V")))},
			wanted:  "T=V,key=value"},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			cont := controller.New(
				processor.New(
					processortest.AggregatorSelector(),
					export.CumulativeExportKindSelector(),
				),
				c.options...,
			)
			prov := cont.MeterProvider()

			ctr := metric.Must(prov.Meter("named")).NewFloat64Counter("calls.sum")
			ctr.Add(context.Background(), 1.)

			// Collect once
			require.NoError(t, cont.Collect(context.Background()))

			expect := map[string]float64{
				"calls.sum//" + c.wanted: 1.,
			}
			require.EqualValues(t, expect, getMap(t, cont))
		})
	}
}

func TestStartNoExporter(t *testing.T) {
	cont := controller.New(
		processor.New(
			processortest.AggregatorSelector(),
			export.CumulativeExportKindSelector(),
		),
		controller.WithCollectPeriod(time.Second),
		controller.WithResource(resource.Empty()),
	)
	mock := controllertest.NewMockClock()
	cont.SetClock(mock)

	prov := cont.MeterProvider()
	calls := int64(0)

	_ = metric.Must(prov.Meter("named")).NewInt64SumObserver("calls.lastvalue",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			calls++
			checkTestContext(t, ctx)
			result.Observe(calls, attribute.String("A", "B"))
		},
	)

	// Collect() has not been called.  The controller is unstarted.
	expect := map[string]float64{}

	// The time advances, but doesn't change the result (not collected).
	require.EqualValues(t, expect, getMap(t, cont))
	mock.Add(time.Second)
	require.EqualValues(t, expect, getMap(t, cont))
	mock.Add(time.Second)

	expect = map[string]float64{
		"calls.lastvalue/A=B/": 1,
	}

	// Collect once
	ctx := testContext()

	require.NoError(t, cont.Collect(ctx))

	require.EqualValues(t, expect, getMap(t, cont))
	mock.Add(time.Second)
	require.EqualValues(t, expect, getMap(t, cont))
	mock.Add(time.Second)

	// Again
	expect = map[string]float64{
		"calls.lastvalue/A=B/": 2,
	}

	require.NoError(t, cont.Collect(ctx))

	require.EqualValues(t, expect, getMap(t, cont))
	mock.Add(time.Second)
	require.EqualValues(t, expect, getMap(t, cont))

	// Start the controller
	require.NoError(t, cont.Start(ctx))

	for i := 1; i <= 3; i++ {
		expect = map[string]float64{
			"calls.lastvalue/A=B/": 2 + float64(i),
		}

		mock.Add(time.Second)
		require.EqualValues(t, expect, getMap(t, cont))
	}
}

func TestObserverCanceled(t *testing.T) {
	cont := controller.New(
		processor.New(
			processortest.AggregatorSelector(),
			export.CumulativeExportKindSelector(),
		),
		controller.WithCollectPeriod(0),
		controller.WithCollectTimeout(time.Millisecond),
		controller.WithResource(resource.Empty()),
	)

	prov := cont.MeterProvider()
	calls := int64(0)

	_ = metric.Must(prov.Meter("named")).NewInt64SumObserver("done.lastvalue",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			<-ctx.Done()
			calls++
			result.Observe(calls)
		},
	)
	// This relies on the context timing out
	err := cont.Collect(context.Background())
	require.Error(t, err)
	require.True(t, errors.Is(err, context.DeadlineExceeded))

	expect := map[string]float64{
		"done.lastvalue//": 1,
	}

	require.EqualValues(t, expect, getMap(t, cont))
}

func TestObserverContext(t *testing.T) {
	cont := controller.New(
		processor.New(
			processortest.AggregatorSelector(),
			export.CumulativeExportKindSelector(),
		),
		controller.WithCollectTimeout(0),
		controller.WithResource(resource.Empty()),
	)

	prov := cont.MeterProvider()

	_ = metric.Must(prov.Meter("named")).NewInt64SumObserver("done.lastvalue",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			time.Sleep(10 * time.Millisecond)
			checkTestContext(t, ctx)
			result.Observe(1)
		},
	)
	ctx := testContext()

	require.NoError(t, cont.Collect(ctx))

	expect := map[string]float64{
		"done.lastvalue//": 1,
	}

	require.EqualValues(t, expect, getMap(t, cont))
}

type blockingExporter struct {
	calls    int
	exporter *processortest.Exporter
}

func newBlockingExporter() *blockingExporter {
	return &blockingExporter{
		exporter: processortest.NewExporter(
			export.CumulativeExportKindSelector(),
			attribute.DefaultEncoder(),
		),
	}
}

func (b *blockingExporter) Export(ctx context.Context, output export.CheckpointSet) error {
	var err error
	_ = b.exporter.Export(ctx, output)
	if b.calls == 0 {
		// timeout once
		<-ctx.Done()
		err = ctx.Err()
	}
	b.calls++
	return err
}

func (*blockingExporter) ExportKindFor(
	*metric.Descriptor,
	aggregation.Kind,
) export.ExportKind {
	return export.CumulativeExportKind
}

func TestExportTimeout(t *testing.T) {
	exporter := newBlockingExporter()
	cont := controller.New(
		processor.New(
			processortest.AggregatorSelector(),
			export.CumulativeExportKindSelector(),
		),
		controller.WithCollectPeriod(time.Second),
		controller.WithPushTimeout(time.Millisecond),
		controller.WithExporter(exporter),
		controller.WithResource(resource.Empty()),
	)
	mock := controllertest.NewMockClock()
	cont.SetClock(mock)

	prov := cont.MeterProvider()

	calls := int64(0)
	_ = metric.Must(prov.Meter("named")).NewInt64SumObserver("one.lastvalue",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			calls++
			result.Observe(calls)
		},
	)

	require.NoError(t, cont.Start(context.Background()))

	// Initial empty state
	expect := map[string]float64{}
	require.EqualValues(t, expect, exporter.exporter.Values())

	// Collect after 1s, timeout
	mock.Add(time.Second)

	err := testHandler.Flush()
	require.Error(t, err)
	require.True(t, errors.Is(err, context.DeadlineExceeded))

	expect = map[string]float64{
		"one.lastvalue//": 1,
	}
	require.EqualValues(t, expect, exporter.exporter.Values())

	// Collect again
	mock.Add(time.Second)
	expect = map[string]float64{
		"one.lastvalue//": 2,
	}
	require.EqualValues(t, expect, exporter.exporter.Values())

	err = testHandler.Flush()
	require.NoError(t, err)
}

func TestCollectAfterStopThenStartAgain(t *testing.T) {
	exp := processortest.NewExporter(
		export.CumulativeExportKindSelector(),
		attribute.DefaultEncoder(),
	)
	cont := controller.New(
		processor.New(
			processortest.AggregatorSelector(),
			exp,
		),
		controller.WithCollectPeriod(time.Second),
		controller.WithExporter(exp),
		controller.WithResource(resource.Empty()),
	)
	mock := controllertest.NewMockClock()
	cont.SetClock(mock)

	prov := cont.MeterProvider()

	calls := 0
	_ = metric.Must(prov.Meter("named")).NewInt64SumObserver("one.lastvalue",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			calls++
			result.Observe(int64(calls))
		},
	)

	// No collections happen (because mock clock does not advance):
	require.NoError(t, cont.Start(context.Background()))
	require.True(t, cont.IsRunning())

	// There's one collection run by Stop():
	require.NoError(t, cont.Stop(context.Background()))

	require.EqualValues(t, map[string]float64{
		"one.lastvalue//": 1,
	}, exp.Values())
	require.NoError(t, testHandler.Flush())

	// Manual collect after Stop still works, subject to
	// CollectPeriod.
	require.NoError(t, cont.Collect(context.Background()))
	require.EqualValues(t, map[string]float64{
		"one.lastvalue//": 2,
	}, getMap(t, cont))

	require.NoError(t, testHandler.Flush())
	require.False(t, cont.IsRunning())

	// Start again, see that collection proceeds.  However,
	// explicit collection should still fail.
	require.NoError(t, cont.Start(context.Background()))
	require.True(t, cont.IsRunning())
	err := cont.Collect(context.Background())
	require.Error(t, err)
	require.Equal(t, controller.ErrControllerStarted, err)

	require.NoError(t, cont.Stop(context.Background()))
	require.EqualValues(t, map[string]float64{
		"one.lastvalue//": 3,
	}, exp.Values())
	require.False(t, cont.IsRunning())

	// Time has not advanced yet. Now let the ticker perform
	// collection:
	require.NoError(t, cont.Start(context.Background()))
	mock.Add(time.Second)
	require.EqualValues(t, map[string]float64{
		"one.lastvalue//": 4,
	}, exp.Values())

	mock.Add(time.Second)
	require.EqualValues(t, map[string]float64{
		"one.lastvalue//": 5,
	}, exp.Values())
	require.NoError(t, cont.Stop(context.Background()))
	require.EqualValues(t, map[string]float64{
		"one.lastvalue//": 6,
	}, exp.Values())
}
