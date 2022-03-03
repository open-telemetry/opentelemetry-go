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
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/controller/controllertest"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/resource"
)

const envVar = "OTEL_RESOURCE_ATTRIBUTES"

func getMap(t *testing.T, cont *controller.Controller) map[string]float64 {
	out := processortest.NewOutput(attribute.DefaultEncoder())

	require.NoError(t, cont.ForEach(
		func(_ instrumentation.Library, reader export.Reader) error {
			return reader.ForEach(
				aggregation.CumulativeTemporalitySelector(),
				func(record export.Record) error {
					return out.AddRecord(record)
				},
			)
		}))
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
	const envVal = "T=U,key=value"
	store, err := ottest.SetEnvVariables(map[string]string{
		envVar: envVal,
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
			wanted:  envVal,
		},
		{
			name:    "uses default if no resource option",
			options: nil,
			wanted:  resource.Default().Encoded(attribute.DefaultEncoder()),
		},
		{
			name:    "explicit resource",
			options: []controller.Option{controller.WithResource(resource.NewSchemaless(attribute.String("R", "S")))},
			wanted:  "R=S," + envVal,
		},
		{
			name: "multi resource",
			options: []controller.Option{
				controller.WithResource(resource.NewSchemaless(attribute.String("R", "WRONG"))),
				controller.WithResource(resource.NewSchemaless(attribute.String("R", "S"))),
				controller.WithResource(resource.NewSchemaless(attribute.String("W", "X"))),
				controller.WithResource(resource.NewSchemaless(attribute.String("T", "V"))),
			},
			wanted: "R=S,T=V,W=X,key=value",
		},
		{
			name: "user override environment",
			options: []controller.Option{
				controller.WithResource(resource.NewSchemaless(attribute.String("T", "V"))),
				controller.WithResource(resource.NewSchemaless(attribute.String("key", "I win"))),
			},
			wanted: "T=V,key=I win",
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			sel := aggregation.CumulativeTemporalitySelector()
			exp := processortest.New(sel, attribute.DefaultEncoder())
			cont := controller.New(
				processor.NewFactory(
					processortest.AggregatorSelector(),
					exp,
				),
				append(c.options, controller.WithExporter(exp))...,
			)
			ctx := context.Background()
			require.NoError(t, cont.Start(ctx))

			ctr, _ := cont.Meter("named").SyncFloat64().Counter("calls.sum")
			ctr.Add(context.Background(), 1.)

			// Collect once
			require.NoError(t, cont.Stop(ctx))

			expect := map[string]float64{
				"calls.sum//" + c.wanted: 1.,
			}
			require.EqualValues(t, expect, exp.Values())
		})
	}
}

func TestStartNoExporter(t *testing.T) {
	cont := controller.New(
		processor.NewFactory(
			processortest.AggregatorSelector(),
			aggregation.CumulativeTemporalitySelector(),
		),
		controller.WithCollectPeriod(time.Second),
		controller.WithResource(resource.Empty()),
	)
	mock := controllertest.NewMockClock()
	cont.SetClock(mock)
	meter := cont.Meter("go.opentelemetry.io/otel/sdk/metric/controller/basic_test#StartNoExporter")

	calls := int64(0)

	counterObserver, err := meter.AsyncInt64().Counter("calls.lastvalue")
	require.NoError(t, err)

	err = meter.RegisterCallback([]instrument.Asynchronous{counterObserver}, func(ctx context.Context) {
		calls++
		checkTestContext(t, ctx)
		counterObserver.Observe(ctx, calls, attribute.String("A", "B"))
	})
	require.NoError(t, err)

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
		processor.NewFactory(
			processortest.AggregatorSelector(),
			aggregation.CumulativeTemporalitySelector(),
		),
		controller.WithCollectPeriod(0),
		controller.WithCollectTimeout(time.Millisecond),
		controller.WithResource(resource.Empty()),
	)
	meter := cont.Meter("go.opentelemetry.io/otel/sdk/metric/controller/basic_test#ObserverCanceled")

	calls := int64(0)

	counterObserver, err := meter.AsyncInt64().Counter("done.lastvalue")
	require.NoError(t, err)

	err = meter.RegisterCallback([]instrument.Asynchronous{counterObserver}, func(ctx context.Context) {
		<-ctx.Done()
		calls++
		counterObserver.Observe(ctx, calls)
	})
	require.NoError(t, err)

	// This relies on the context timing out
	err = cont.Collect(context.Background())
	require.Error(t, err)
	require.True(t, errors.Is(err, context.DeadlineExceeded))

	expect := map[string]float64{
		"done.lastvalue//": 1,
	}

	require.EqualValues(t, expect, getMap(t, cont))
}

func TestObserverContext(t *testing.T) {
	cont := controller.New(
		processor.NewFactory(
			processortest.AggregatorSelector(),
			aggregation.CumulativeTemporalitySelector(),
		),
		controller.WithCollectTimeout(0),
		controller.WithResource(resource.Empty()),
	)
	meter := cont.Meter("go.opentelemetry.io/otel/sdk/metric/controller/basic_test#ObserverContext")

	counterObserver, err := meter.AsyncInt64().Counter("done.lastvalue")
	require.NoError(t, err)

	err = meter.RegisterCallback([]instrument.Asynchronous{counterObserver}, func(ctx context.Context) {
		time.Sleep(10 * time.Millisecond)
		checkTestContext(t, ctx)
		counterObserver.Observe(ctx, 1)
	})
	require.NoError(t, err)

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
		exporter: processortest.New(
			aggregation.CumulativeTemporalitySelector(),
			attribute.DefaultEncoder(),
		),
	}
}

func (b *blockingExporter) Export(ctx context.Context, res *resource.Resource, output export.InstrumentationLibraryReader) error {
	var err error
	_ = b.exporter.Export(ctx, res, output)
	if b.calls == 0 {
		// timeout once
		<-ctx.Done()
		err = ctx.Err()
	}
	b.calls++
	return err
}

func (*blockingExporter) TemporalityFor(*sdkapi.Descriptor, aggregation.Kind) aggregation.Temporality {
	return aggregation.CumulativeTemporality
}

func TestExportTimeout(t *testing.T) {
	exporter := newBlockingExporter()
	cont := controller.New(
		processor.NewFactory(
			processortest.AggregatorSelector(),
			aggregation.CumulativeTemporalitySelector(),
		),
		controller.WithCollectPeriod(time.Second),
		controller.WithPushTimeout(time.Millisecond),
		controller.WithExporter(exporter),
		controller.WithResource(resource.Empty()),
	)
	mock := controllertest.NewMockClock()
	cont.SetClock(mock)
	meter := cont.Meter("go.opentelemetry.io/otel/sdk/metric/controller/basic_test#ExportTimeout")

	calls := int64(0)
	counterObserver, err := meter.AsyncInt64().Counter("one.lastvalue")
	require.NoError(t, err)

	err = meter.RegisterCallback([]instrument.Asynchronous{counterObserver}, func(ctx context.Context) {
		calls++
		counterObserver.Observe(ctx, calls)
	})
	require.NoError(t, err)

	require.NoError(t, cont.Start(context.Background()))

	// Initial empty state
	expect := map[string]float64{}
	require.EqualValues(t, expect, exporter.exporter.Values())

	// Collect after 1s, timeout
	mock.Add(time.Second)

	err = testHandler.Flush()
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
	exp := processortest.New(
		aggregation.CumulativeTemporalitySelector(),
		attribute.DefaultEncoder(),
	)
	cont := controller.New(
		processor.NewFactory(
			processortest.AggregatorSelector(),
			exp,
		),
		controller.WithCollectPeriod(time.Second),
		controller.WithExporter(exp),
		controller.WithResource(resource.Empty()),
	)
	mock := controllertest.NewMockClock()
	cont.SetClock(mock)

	meter := cont.Meter("go.opentelemetry.io/otel/sdk/metric/controller/basic_test#CollectAfterStopThenStartAgain")

	calls := 0
	counterObserver, err := meter.AsyncInt64().Counter("one.lastvalue")
	require.NoError(t, err)

	err = meter.RegisterCallback([]instrument.Asynchronous{counterObserver}, func(ctx context.Context) {
		calls++
		counterObserver.Observe(ctx, int64(calls))
	})
	require.NoError(t, err)

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
	err = cont.Collect(context.Background())
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

func TestRegistryFunction(t *testing.T) {
	exp := processortest.New(
		aggregation.CumulativeTemporalitySelector(),
		attribute.DefaultEncoder(),
	)
	cont := controller.New(
		processor.NewFactory(
			processortest.AggregatorSelector(),
			exp,
		),
		controller.WithCollectPeriod(time.Second),
		controller.WithExporter(exp),
		controller.WithResource(resource.Empty()),
	)

	m1 := cont.Meter("test")
	m2 := cont.Meter("test")

	require.NotNil(t, m1)
	require.Equal(t, m1, m2)

	c1, err := m1.SyncInt64().Counter("counter.sum")
	require.NoError(t, err)

	c2, err := m1.SyncInt64().Counter("counter.sum")
	require.NoError(t, err)

	require.Equal(t, c1, c2)

	ctx := context.Background()

	require.NoError(t, cont.Start(ctx))

	c1.Add(ctx, 10)
	c2.Add(ctx, 10)

	require.NoError(t, cont.Stop(ctx))

	require.EqualValues(t, map[string]float64{
		"counter.sum//": 20,
	}, exp.Values())
}
