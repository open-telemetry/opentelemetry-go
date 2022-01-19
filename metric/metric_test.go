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
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/metrictest"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/unit"
)

var Must = metric.Must

var (
	syncKinds = []sdkapi.InstrumentKind{
		sdkapi.HistogramInstrumentKind,
		sdkapi.CounterInstrumentKind,
		sdkapi.UpDownCounterInstrumentKind,
	}
	asyncKinds = []sdkapi.InstrumentKind{
		sdkapi.GaugeObserverInstrumentKind,
		sdkapi.CounterObserverInstrumentKind,
		sdkapi.UpDownCounterObserverInstrumentKind,
	}
	addingKinds = []sdkapi.InstrumentKind{
		sdkapi.CounterInstrumentKind,
		sdkapi.UpDownCounterInstrumentKind,
		sdkapi.CounterObserverInstrumentKind,
		sdkapi.UpDownCounterObserverInstrumentKind,
	}
	groupingKinds = []sdkapi.InstrumentKind{
		sdkapi.HistogramInstrumentKind,
		sdkapi.GaugeObserverInstrumentKind,
	}

	monotonicKinds = []sdkapi.InstrumentKind{
		sdkapi.CounterInstrumentKind,
		sdkapi.CounterObserverInstrumentKind,
	}

	nonMonotonicKinds = []sdkapi.InstrumentKind{
		sdkapi.UpDownCounterInstrumentKind,
		sdkapi.UpDownCounterObserverInstrumentKind,
		sdkapi.HistogramInstrumentKind,
		sdkapi.GaugeObserverInstrumentKind,
	}

	precomputedSumKinds = []sdkapi.InstrumentKind{
		sdkapi.CounterObserverInstrumentKind,
		sdkapi.UpDownCounterObserverInstrumentKind,
	}

	nonPrecomputedSumKinds = []sdkapi.InstrumentKind{
		sdkapi.CounterInstrumentKind,
		sdkapi.UpDownCounterInstrumentKind,
		sdkapi.HistogramInstrumentKind,
		sdkapi.GaugeObserverInstrumentKind,
	}
)

func TestSynchronous(t *testing.T) {
	for _, k := range syncKinds {
		require.True(t, k.Synchronous())
		require.False(t, k.Asynchronous())
	}
	for _, k := range asyncKinds {
		require.True(t, k.Asynchronous())
		require.False(t, k.Synchronous())
	}
}

func TestGrouping(t *testing.T) {
	for _, k := range groupingKinds {
		require.True(t, k.Grouping())
		require.False(t, k.Adding())
	}
	for _, k := range addingKinds {
		require.True(t, k.Adding())
		require.False(t, k.Grouping())
	}
}

func TestMonotonic(t *testing.T) {
	for _, k := range monotonicKinds {
		require.True(t, k.Monotonic())
	}
	for _, k := range nonMonotonicKinds {
		require.False(t, k.Monotonic())
	}
}

func TestPrecomputedSum(t *testing.T) {
	for _, k := range precomputedSumKinds {
		require.True(t, k.PrecomputedSum())
	}
	for _, k := range nonPrecomputedSumKinds {
		require.False(t, k.PrecomputedSum())
	}
}

func checkSyncBatches(ctx context.Context, t *testing.T, labels []attribute.KeyValue, provider *metrictest.MeterProvider, nkind number.Kind, mkind sdkapi.InstrumentKind, instrument sdkapi.InstrumentImpl, expected ...float64) {
	t.Helper()

	batchesCount := len(provider.MeasurementBatches)
	if len(provider.MeasurementBatches) != len(expected) {
		t.Errorf("Expected %d recorded measurement batches, got %d", batchesCount, len(provider.MeasurementBatches))
	}
	recorded := metrictest.AsStructs(provider.MeasurementBatches)

	for i, batch := range provider.MeasurementBatches {
		if len(batch.Measurements) != 1 {
			t.Errorf("Expected 1 measurement in batch %d, got %d", i, len(batch.Measurements))
		}

		measurement := batch.Measurements[0]
		descriptor := measurement.Instrument.Descriptor()

		expected := metrictest.Measured{
			Name: descriptor.Name(),
			Library: metrictest.Library{
				InstrumentationName: "apitest",
			},
			Labels: metrictest.LabelsToMap(labels...),
			Number: metrictest.ResolveNumberByKind(t, nkind, expected[i]),
		}
		require.Equal(t, expected, recorded[i])
	}
}

func TestOptions(t *testing.T) {
	type testcase struct {
		name string
		opts []metric.InstrumentOption
		desc string
		unit unit.Unit
	}
	testcases := []testcase{
		{
			name: "no opts",
			opts: nil,
			desc: "",
			unit: "",
		},
		{
			name: "description",
			opts: []metric.InstrumentOption{
				metric.WithDescription("stuff"),
			},
			desc: "stuff",
			unit: "",
		},
		{
			name: "description override",
			opts: []metric.InstrumentOption{
				metric.WithDescription("stuff"),
				metric.WithDescription("things"),
			},
			desc: "things",
			unit: "",
		},
		{
			name: "unit",
			opts: []metric.InstrumentOption{
				metric.WithUnit("s"),
			},
			desc: "",
			unit: "s",
		},
		{
			name: "description override",
			opts: []metric.InstrumentOption{
				metric.WithDescription("stuff"),
				metric.WithDescription("things"),
			},
			desc: "things",
			unit: "",
		},
		{
			name: "unit",
			opts: []metric.InstrumentOption{
				metric.WithUnit("s"),
			},
			desc: "",
			unit: "s",
		},

		{
			name: "unit override",
			opts: []metric.InstrumentOption{
				metric.WithUnit("s"),
				metric.WithUnit("h"),
			},
			desc: "",
			unit: "h",
		},
		{
			name: "all",
			opts: []metric.InstrumentOption{
				metric.WithDescription("stuff"),
				metric.WithUnit("s"),
			},
			desc: "stuff",
			unit: "s",
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing counter case %s (%d)", tt.name, idx)
		cfg := metric.NewInstrumentConfig(tt.opts...)
		if diff := cmp.Diff(cfg.Description(), tt.desc); diff != "" {
			t.Errorf("Compare Description: -got +want %s", diff)
		}
		if diff := cmp.Diff(cfg.Unit(), tt.unit); diff != "" {
			t.Errorf("Compare Unit: -got +want %s", diff)
		}
	}
}
func testPair() (*metrictest.MeterProvider, metric.Meter) {
	provider := metrictest.NewMeterProvider()
	return provider, provider.Meter("apitest")
}

func TestCounter(t *testing.T) {
	// N.B. the API does not check for negative
	// values, that's the SDK's responsibility.
	t.Run("float64 counter", func(t *testing.T) {
		provider, meter := testPair()
		c := Must(meter).NewFloat64Counter("test.counter.float")
		ctx := context.Background()
		labels := []attribute.KeyValue{attribute.String("A", "B")}
		c.Add(ctx, 1994.1, labels...)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		checkSyncBatches(ctx, t, labels, provider, number.Float64Kind, sdkapi.CounterInstrumentKind, c.SyncImpl(),
			1994.1, 42,
		)
	})
	t.Run("int64 counter", func(t *testing.T) {
		provider, meter := testPair()
		c := Must(meter).NewInt64Counter("test.counter.int")
		ctx := context.Background()
		labels := []attribute.KeyValue{attribute.String("A", "B"), attribute.String("C", "D")}
		c.Add(ctx, 42, labels...)
		meter.RecordBatch(ctx, labels, c.Measurement(420000))
		checkSyncBatches(ctx, t, labels, provider, number.Int64Kind, sdkapi.CounterInstrumentKind, c.SyncImpl(),
			42, 420000,
		)

	})
	t.Run("int64 updowncounter", func(t *testing.T) {
		provider, meter := testPair()
		c := Must(meter).NewInt64UpDownCounter("test.updowncounter.int")
		ctx := context.Background()
		labels := []attribute.KeyValue{attribute.String("A", "B"), attribute.String("C", "D")}
		c.Add(ctx, 100, labels...)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		checkSyncBatches(ctx, t, labels, provider, number.Int64Kind, sdkapi.UpDownCounterInstrumentKind, c.SyncImpl(),
			100, 42,
		)
	})
	t.Run("float64 updowncounter", func(t *testing.T) {
		provider, meter := testPair()
		c := Must(meter).NewFloat64UpDownCounter("test.updowncounter.float")
		ctx := context.Background()
		labels := []attribute.KeyValue{attribute.String("A", "B"), attribute.String("C", "D")}
		c.Add(ctx, 100.1, labels...)
		meter.RecordBatch(ctx, labels, c.Measurement(-100.1))
		checkSyncBatches(ctx, t, labels, provider, number.Float64Kind, sdkapi.UpDownCounterInstrumentKind, c.SyncImpl(),
			100.1, -100.1,
		)
	})
}

func TestHistogram(t *testing.T) {
	t.Run("float64 histogram", func(t *testing.T) {
		provider, meter := testPair()
		m := Must(meter).NewFloat64Histogram("test.histogram.float")
		ctx := context.Background()
		labels := []attribute.KeyValue{}
		m.Record(ctx, 42, labels...)
		meter.RecordBatch(ctx, labels, m.Measurement(-100.5))
		checkSyncBatches(ctx, t, labels, provider, number.Float64Kind, sdkapi.HistogramInstrumentKind, m.SyncImpl(),
			42, -100.5,
		)
	})
	t.Run("int64 histogram", func(t *testing.T) {
		provider, meter := testPair()
		m := Must(meter).NewInt64Histogram("test.histogram.int")
		ctx := context.Background()
		labels := []attribute.KeyValue{attribute.Int("I", 1)}
		m.Record(ctx, 173, labels...)
		meter.RecordBatch(ctx, labels, m.Measurement(0))
		checkSyncBatches(ctx, t, labels, provider, number.Int64Kind, sdkapi.HistogramInstrumentKind, m.SyncImpl(),
			173, 0,
		)
	})
}

func TestObserverInstruments(t *testing.T) {
	t.Run("float gauge", func(t *testing.T) {
		labels := []attribute.KeyValue{attribute.String("O", "P")}
		provider, meter := testPair()
		o := Must(meter).NewFloat64GaugeObserver("test.gauge.float", func(_ context.Context, result metric.Float64ObserverResult) {
			result.Observe(42.1, labels...)
		})
		provider.RunAsyncInstruments()
		checkObserverBatch(t, labels, provider, number.Float64Kind, sdkapi.GaugeObserverInstrumentKind, o.AsyncImpl(),
			42.1,
		)
	})
	t.Run("int gauge", func(t *testing.T) {
		labels := []attribute.KeyValue{}
		provider, meter := testPair()
		o := Must(meter).NewInt64GaugeObserver("test.gauge.int", func(_ context.Context, result metric.Int64ObserverResult) {
			result.Observe(-142, labels...)
		})
		provider.RunAsyncInstruments()
		checkObserverBatch(t, labels, provider, number.Int64Kind, sdkapi.GaugeObserverInstrumentKind, o.AsyncImpl(),
			-142,
		)
	})
	t.Run("float counterobserver", func(t *testing.T) {
		labels := []attribute.KeyValue{attribute.String("O", "P")}
		provider, meter := testPair()
		o := Must(meter).NewFloat64CounterObserver("test.counter.float", func(_ context.Context, result metric.Float64ObserverResult) {
			result.Observe(42.1, labels...)
		})
		provider.RunAsyncInstruments()
		checkObserverBatch(t, labels, provider, number.Float64Kind, sdkapi.CounterObserverInstrumentKind, o.AsyncImpl(),
			42.1,
		)
	})
	t.Run("int counterobserver", func(t *testing.T) {
		labels := []attribute.KeyValue{}
		provider, meter := testPair()
		o := Must(meter).NewInt64CounterObserver("test.counter.int", func(_ context.Context, result metric.Int64ObserverResult) {
			result.Observe(-142, labels...)
		})
		provider.RunAsyncInstruments()
		checkObserverBatch(t, labels, provider, number.Int64Kind, sdkapi.CounterObserverInstrumentKind, o.AsyncImpl(),
			-142,
		)
	})
	t.Run("float updowncounterobserver", func(t *testing.T) {
		labels := []attribute.KeyValue{attribute.String("O", "P")}
		provider, meter := testPair()
		o := Must(meter).NewFloat64UpDownCounterObserver("test.updowncounter.float", func(_ context.Context, result metric.Float64ObserverResult) {
			result.Observe(42.1, labels...)
		})
		provider.RunAsyncInstruments()
		checkObserverBatch(t, labels, provider, number.Float64Kind, sdkapi.UpDownCounterObserverInstrumentKind, o.AsyncImpl(),
			42.1,
		)
	})
	t.Run("int updowncounterobserver", func(t *testing.T) {
		labels := []attribute.KeyValue{}
		provider, meter := testPair()
		o := Must(meter).NewInt64UpDownCounterObserver("test..int", func(_ context.Context, result metric.Int64ObserverResult) {
			result.Observe(-142, labels...)
		})
		provider.RunAsyncInstruments()
		checkObserverBatch(t, labels, provider, number.Int64Kind, sdkapi.UpDownCounterObserverInstrumentKind, o.AsyncImpl(),
			-142,
		)
	})
}

func TestBatchObserverInstruments(t *testing.T) {
	provider, meter := testPair()

	var obs1 metric.Int64GaugeObserver
	var obs2 metric.Float64GaugeObserver

	labels := []attribute.KeyValue{
		attribute.String("A", "B"),
		attribute.String("C", "D"),
	}

	cb := Must(meter).NewBatchObserver(
		func(_ context.Context, result metric.BatchObserverResult) {
			result.Observe(labels,
				obs1.Observation(42),
				obs2.Observation(42.0),
			)
		},
	)
	obs1 = cb.NewInt64GaugeObserver("test.gauge.int")
	obs2 = cb.NewFloat64GaugeObserver("test.gauge.float")

	provider.RunAsyncInstruments()

	require.Len(t, provider.MeasurementBatches, 1)

	impl1 := obs1.AsyncImpl().Implementation().(*metrictest.Async)
	impl2 := obs2.AsyncImpl().Implementation().(*metrictest.Async)

	require.NotNil(t, impl1)
	require.NotNil(t, impl2)

	got := provider.MeasurementBatches[0]
	require.Equal(t, labels, got.Labels)
	require.Len(t, got.Measurements, 2)

	m1 := got.Measurements[0]
	require.Equal(t, impl1, m1.Instrument.Implementation().(*metrictest.Async))
	require.Equal(t, 0, m1.Number.CompareNumber(number.Int64Kind, metrictest.ResolveNumberByKind(t, number.Int64Kind, 42)))

	m2 := got.Measurements[1]
	require.Equal(t, impl2, m2.Instrument.Implementation().(*metrictest.Async))
	require.Equal(t, 0, m2.Number.CompareNumber(number.Float64Kind, metrictest.ResolveNumberByKind(t, number.Float64Kind, 42)))
}

func checkObserverBatch(t *testing.T, labels []attribute.KeyValue, provider *metrictest.MeterProvider, nkind number.Kind, mkind sdkapi.InstrumentKind, observer sdkapi.AsyncImpl, expected float64) {
	t.Helper()
	assert.Len(t, provider.MeasurementBatches, 1)
	if len(provider.MeasurementBatches) < 1 {
		return
	}
	o := observer.Implementation().(*metrictest.Async)
	if !assert.NotNil(t, o) {
		return
	}
	got := provider.MeasurementBatches[0]
	assert.Equal(t, labels, got.Labels)
	assert.Len(t, got.Measurements, 1)
	if len(got.Measurements) < 1 {
		return
	}
	measurement := got.Measurements[0]
	require.Equal(t, mkind, measurement.Instrument.Descriptor().InstrumentKind())
	assert.Equal(t, o, measurement.Instrument.Implementation().(*metrictest.Async))
	ft := metrictest.ResolveNumberByKind(t, nkind, expected)
	assert.Equal(t, 0, measurement.Number.CompareNumber(nkind, ft))
}

type testWrappedMeter struct {
}

var _ sdkapi.MeterImpl = testWrappedMeter{}

func (testWrappedMeter) RecordBatch(context.Context, []attribute.KeyValue, ...sdkapi.Measurement) {
}

func (testWrappedMeter) NewSyncInstrument(_ sdkapi.Descriptor) (sdkapi.SyncImpl, error) {
	return nil, nil
}

func (testWrappedMeter) NewAsyncInstrument(_ sdkapi.Descriptor, _ sdkapi.AsyncRunner) (sdkapi.AsyncImpl, error) {
	return nil, errors.New("Test wrap error")
}

func TestWrappedInstrumentError(t *testing.T) {
	impl := &testWrappedMeter{}
	meter := metric.WrapMeterImpl(impl)

	histogram, err := meter.NewInt64Histogram("test.histogram")

	require.Equal(t, err, metric.ErrSDKReturnedNilImpl)
	require.NotNil(t, histogram.SyncImpl())

	observer, err := meter.NewInt64GaugeObserver("test.observer", func(_ context.Context, result metric.Int64ObserverResult) {})

	require.NotNil(t, err)
	require.NotNil(t, observer.AsyncImpl())
}

func TestNilCallbackObserverNoop(t *testing.T) {
	// Tests that a nil callback yields a no-op observer without error.
	_, meter := testPair()

	observer := Must(meter).NewInt64GaugeObserver("test.observer", nil)

	impl := observer.AsyncImpl().Implementation()
	desc := observer.AsyncImpl().Descriptor()
	require.Equal(t, nil, impl)
	require.Equal(t, "", desc.Name())
}
