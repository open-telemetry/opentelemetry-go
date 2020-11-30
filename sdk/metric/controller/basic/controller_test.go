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
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/controller/controllertest"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/resource"
)

var testResource = resource.NewWithAttributes(label.String("R", "V"))

type handler struct {
	sync.Mutex
	err error
}

func (h *handler) Handle(err error) {
	h.Lock()
	h.err = err
	h.Unlock()
}

func (h *handler) Flush() error {
	h.Lock()
	err := h.err
	h.err = nil
	h.Unlock()
	return err
}

var testHandler *handler

func init() {
	testHandler = new(handler)
	otel.SetErrorHandler(testHandler)
}

func newExporter() *processortest.Exporter {
	return processortest.NewExporter(
		export.StatelessExportKindSelector(),
		label.DefaultEncoder(),
	)
}

func newCheckpointer() export.Checkpointer {
	return processortest.Checkpointer(
		processortest.NewProcessor(
			processortest.AggregatorSelector(),
			label.DefaultEncoder(),
		),
	)
}

func TestPushDoubleStop(t *testing.T) {
	ctx := context.Background()
	exporter := newExporter()
	checkpointer := newCheckpointer()
	p := controller.New(checkpointer, controller.WithExporter(exporter))
	p.Start(ctx)
	p.Stop(ctx)
	p.Stop(ctx)
}

func TestPushDoubleStart(t *testing.T) {
	ctx := context.Background()
	exporter := newExporter()
	checkpointer := newCheckpointer()
	p := controller.New(checkpointer, controller.WithExporter(exporter))
	p.Start(ctx)
	p.Start(ctx)
	p.Stop(ctx)
}

func TestPushTicker(t *testing.T) {
	exporter := newExporter()
	checkpointer := newCheckpointer()
	p := controller.New(
		checkpointer,
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(time.Second),
		controller.WithResource(testResource),
	)
	meter := p.MeterProvider().Meter("name")

	mock := controllertest.NewMockClock()
	p.SetClock(mock)

	ctx := context.Background()

	counter := metric.Must(meter).NewInt64Counter("counter.sum")

	p.Start(ctx)

	counter.Add(ctx, 3)

	require.EqualValues(t, map[string]float64{}, exporter.Values())

	mock.Add(time.Second)
	runtime.Gosched()

	require.EqualValues(t, map[string]float64{
		"counter.sum//R=V": 3,
	}, exporter.Values())

	require.Equal(t, 1, exporter.ExportCount())
	exporter.Reset()

	counter.Add(ctx, 7)

	mock.Add(time.Second)
	runtime.Gosched()

	require.EqualValues(t, map[string]float64{
		"counter.sum//R=V": 10,
	}, exporter.Values())

	require.Equal(t, 1, exporter.ExportCount())
	exporter.Reset()

	p.Stop(ctx)
}

func TestPushExportError(t *testing.T) {
	injector := func(name string, e error) func(r export.Record) error {
		return func(r export.Record) error {
			if r.Descriptor().Name() == name {
				return e
			}
			return nil
		}
	}
	var errAggregator = fmt.Errorf("unexpected error")
	var tests = []struct {
		name          string
		injectedError error
		expected      map[string]float64
		expectedError error
	}{
		{"errNone", nil, map[string]float64{
			"counter1.sum/X=Y/R=V": 3,
			"counter2.sum//R=V":    5,
		}, nil},
		{"errNoData", aggregation.ErrNoData, map[string]float64{
			"counter2.sum//R=V": 5,
		}, nil},
		{"errUnexpected", errAggregator, map[string]float64{}, errAggregator},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exporter := newExporter()
			exporter.InjectErr = injector("counter1.sum", tt.injectedError)

			// This test validates the error handling
			// behavior of the basic Processor is honored
			// by the push processor.
			checkpointer := processor.New(processortest.AggregatorSelector(), exporter)
			p := controller.New(
				checkpointer,
				controller.WithExporter(exporter),
				controller.WithCollectPeriod(time.Second),
				controller.WithResource(testResource),
			)

			mock := controllertest.NewMockClock()
			p.SetClock(mock)

			ctx := context.Background()

			meter := p.MeterProvider().Meter("name")
			counter1 := metric.Must(meter).NewInt64Counter("counter1.sum")
			counter2 := metric.Must(meter).NewInt64Counter("counter2.sum")

			p.Start(ctx)
			runtime.Gosched()

			counter1.Add(ctx, 3, label.String("X", "Y"))
			counter2.Add(ctx, 5)

			require.Equal(t, 0, exporter.ExportCount())
			require.Nil(t, testHandler.Flush())

			mock.Add(time.Second)
			runtime.Gosched()

			require.Equal(t, 1, exporter.ExportCount())
			if tt.expectedError == nil {
				require.EqualValues(t, tt.expected, exporter.Values())
				require.NoError(t, testHandler.Flush())
			} else {
				err := testHandler.Flush()
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err)
			}

			p.Stop(ctx)
		})
	}
}

func TestPullNoCollect(t *testing.T) {
	puller := controller.New(
		processor.New(
			processortest.AggregatorSelector(),
			export.CumulativeExportKindSelector(),
			processor.WithMemory(true),
		),
		controller.WithCollectPeriod(0),
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

func TestPullWithCollect(t *testing.T) {
	puller := controller.New(
		processor.New(
			processortest.AggregatorSelector(),
			export.CumulativeExportKindSelector(),
			processor.WithMemory(true),
		),
		controller.WithCollectPeriod(time.Second),
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

func getMap(t *testing.T, cont *controller.Controller) map[string]float64 {
	out := processortest.NewOutput(label.DefaultEncoder())

	require.NoError(t, cont.ForEach(
		export.CumulativeExportKindSelector(),
		func(record export.Record) error {
			return out.AddRecord(record)
		},
	))
	return out.Map()
}

func TestStartNoExporter(t *testing.T) {
	cont := controller.New(
		processor.New(
			processortest.AggregatorSelector(),
			export.CumulativeExportKindSelector(),
		),
		controller.WithCollectPeriod(time.Second),
	)
	mock := controllertest.NewMockClock()
	cont.SetClock(mock)

	prov := cont.MeterProvider()
	calls := int64(0)

	_ = metric.Must(prov.Meter("named")).NewInt64SumObserver("calls.lastvalue",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			calls++
			require.Equal(t, "B", ctx.Value("A"))
			result.Observe(calls, label.String("A", "B"))
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
	ctx := context.Background()
	ctx = context.WithValue(ctx, "A", "B")

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
	cont.Start(ctx)

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
	)

	prov := cont.MeterProvider()

	_ = metric.Must(prov.Meter("named")).NewInt64SumObserver("done.lastvalue",
		func(ctx context.Context, result metric.Int64ObserverResult) {
			time.Sleep(10 * time.Millisecond)
			require.Equal(t, "B", ctx.Value("A"))
			result.Observe(1)
		},
	)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "A", "B")

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
			label.DefaultEncoder(),
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

func (_ *blockingExporter) ExportKindFor(
	_ *metric.Descriptor,
	_ aggregation.Kind,
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
		controller.WithExportTimeout(time.Millisecond),
		controller.WithExporter(exporter),
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
