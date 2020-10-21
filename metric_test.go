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

package otel_test

import (
	"context"
	"errors"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/unit"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Must = otel.Must

var (
	syncKinds = []otel.InstrumentKind{
		otel.ValueRecorderInstrumentKind,
		otel.CounterInstrumentKind,
		otel.UpDownCounterInstrumentKind,
	}
	asyncKinds = []otel.InstrumentKind{
		otel.ValueObserverInstrumentKind,
		otel.SumObserverInstrumentKind,
		otel.UpDownSumObserverInstrumentKind,
	}
	addingKinds = []otel.InstrumentKind{
		otel.CounterInstrumentKind,
		otel.UpDownCounterInstrumentKind,
		otel.SumObserverInstrumentKind,
		otel.UpDownSumObserverInstrumentKind,
	}
	groupingKinds = []otel.InstrumentKind{
		otel.ValueRecorderInstrumentKind,
		otel.ValueObserverInstrumentKind,
	}

	monotonicKinds = []otel.InstrumentKind{
		otel.CounterInstrumentKind,
		otel.SumObserverInstrumentKind,
	}

	nonMonotonicKinds = []otel.InstrumentKind{
		otel.UpDownCounterInstrumentKind,
		otel.UpDownSumObserverInstrumentKind,
		otel.ValueRecorderInstrumentKind,
		otel.ValueObserverInstrumentKind,
	}

	precomputedSumKinds = []otel.InstrumentKind{
		otel.SumObserverInstrumentKind,
		otel.UpDownSumObserverInstrumentKind,
	}

	nonPrecomputedSumKinds = []otel.InstrumentKind{
		otel.CounterInstrumentKind,
		otel.UpDownCounterInstrumentKind,
		otel.ValueRecorderInstrumentKind,
		otel.ValueObserverInstrumentKind,
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

func checkSyncBatches(ctx context.Context, t *testing.T, labels []label.KeyValue, mock *oteltest.MeterImpl, nkind otel.NumberKind, mkind otel.InstrumentKind, instrument otel.InstrumentImpl, expected ...float64) {
	t.Helper()

	batchesCount := len(mock.MeasurementBatches)
	if len(mock.MeasurementBatches) != len(expected) {
		t.Errorf("Expected %d recorded measurement batches, got %d", batchesCount, len(mock.MeasurementBatches))
	}
	recorded := oteltest.AsStructs(mock.MeasurementBatches)

	for i, batch := range mock.MeasurementBatches {
		if len(batch.Measurements) != 1 {
			t.Errorf("Expected 1 measurement in batch %d, got %d", i, len(batch.Measurements))
		}

		measurement := batch.Measurements[0]
		descriptor := measurement.Instrument.Descriptor()

		expected := oteltest.Measured{
			Name:                descriptor.Name(),
			InstrumentationName: descriptor.InstrumentationName(),
			Labels:              oteltest.LabelsToMap(labels...),
			Number:              oteltest.ResolveNumberByKind(t, nkind, expected[i]),
		}
		require.Equal(t, expected, recorded[i])
	}
}

func TestOptions(t *testing.T) {
	type testcase struct {
		name string
		opts []otel.InstrumentOption
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
			opts: []otel.InstrumentOption{
				otel.WithDescription("stuff"),
			},
			desc: "stuff",
			unit: "",
		},
		{
			name: "description override",
			opts: []otel.InstrumentOption{
				otel.WithDescription("stuff"),
				otel.WithDescription("things"),
			},
			desc: "things",
			unit: "",
		},
		{
			name: "unit",
			opts: []otel.InstrumentOption{
				otel.WithUnit("s"),
			},
			desc: "",
			unit: "s",
		},
		{
			name: "unit override",
			opts: []otel.InstrumentOption{
				otel.WithUnit("s"),
				otel.WithUnit("h"),
			},
			desc: "",
			unit: "h",
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing counter case %s (%d)", tt.name, idx)
		if diff := cmp.Diff(otel.NewInstrumentConfig(tt.opts...), otel.InstrumentConfig{
			Description: tt.desc,
			Unit:        tt.unit,
		}); diff != "" {
			t.Errorf("Compare options: -got +want %s", diff)
		}
	}
}

func TestCounter(t *testing.T) {
	// N.B. the API does not check for negative
	// values, that's the SDK's responsibility.
	t.Run("float64 counter", func(t *testing.T) {
		mockSDK, meter := oteltest.NewMeter()
		c := Must(meter).NewFloat64Counter("test.counter.float")
		ctx := context.Background()
		labels := []label.KeyValue{label.String("A", "B")}
		c.Add(ctx, 1994.1, labels...)
		boundInstrument := c.Bind(labels...)
		boundInstrument.Add(ctx, -742)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		checkSyncBatches(ctx, t, labels, mockSDK, otel.Float64NumberKind, otel.CounterInstrumentKind, c.SyncImpl(),
			1994.1, -742, 42,
		)
	})
	t.Run("int64 counter", func(t *testing.T) {
		mockSDK, meter := oteltest.NewMeter()
		c := Must(meter).NewInt64Counter("test.counter.int")
		ctx := context.Background()
		labels := []label.KeyValue{label.String("A", "B"), label.String("C", "D")}
		c.Add(ctx, 42, labels...)
		boundInstrument := c.Bind(labels...)
		boundInstrument.Add(ctx, 4200)
		meter.RecordBatch(ctx, labels, c.Measurement(420000))
		checkSyncBatches(ctx, t, labels, mockSDK, otel.Int64NumberKind, otel.CounterInstrumentKind, c.SyncImpl(),
			42, 4200, 420000,
		)

	})
	t.Run("int64 updowncounter", func(t *testing.T) {
		mockSDK, meter := oteltest.NewMeter()
		c := Must(meter).NewInt64UpDownCounter("test.updowncounter.int")
		ctx := context.Background()
		labels := []label.KeyValue{label.String("A", "B"), label.String("C", "D")}
		c.Add(ctx, 100, labels...)
		boundInstrument := c.Bind(labels...)
		boundInstrument.Add(ctx, -100)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		checkSyncBatches(ctx, t, labels, mockSDK, otel.Int64NumberKind, otel.UpDownCounterInstrumentKind, c.SyncImpl(),
			100, -100, 42,
		)
	})
	t.Run("float64 updowncounter", func(t *testing.T) {
		mockSDK, meter := oteltest.NewMeter()
		c := Must(meter).NewFloat64UpDownCounter("test.updowncounter.float")
		ctx := context.Background()
		labels := []label.KeyValue{label.String("A", "B"), label.String("C", "D")}
		c.Add(ctx, 100.1, labels...)
		boundInstrument := c.Bind(labels...)
		boundInstrument.Add(ctx, -76)
		meter.RecordBatch(ctx, labels, c.Measurement(-100.1))
		checkSyncBatches(ctx, t, labels, mockSDK, otel.Float64NumberKind, otel.UpDownCounterInstrumentKind, c.SyncImpl(),
			100.1, -76, -100.1,
		)
	})
}

func TestValueRecorder(t *testing.T) {
	t.Run("float64 valuerecorder", func(t *testing.T) {
		mockSDK, meter := oteltest.NewMeter()
		m := Must(meter).NewFloat64ValueRecorder("test.valuerecorder.float")
		ctx := context.Background()
		labels := []label.KeyValue{}
		m.Record(ctx, 42, labels...)
		boundInstrument := m.Bind(labels...)
		boundInstrument.Record(ctx, 0)
		meter.RecordBatch(ctx, labels, m.Measurement(-100.5))
		checkSyncBatches(ctx, t, labels, mockSDK, otel.Float64NumberKind, otel.ValueRecorderInstrumentKind, m.SyncImpl(),
			42, 0, -100.5,
		)
	})
	t.Run("int64 valuerecorder", func(t *testing.T) {
		mockSDK, meter := oteltest.NewMeter()
		m := Must(meter).NewInt64ValueRecorder("test.valuerecorder.int")
		ctx := context.Background()
		labels := []label.KeyValue{label.Int("I", 1)}
		m.Record(ctx, 173, labels...)
		boundInstrument := m.Bind(labels...)
		boundInstrument.Record(ctx, 80)
		meter.RecordBatch(ctx, labels, m.Measurement(0))
		checkSyncBatches(ctx, t, labels, mockSDK, otel.Int64NumberKind, otel.ValueRecorderInstrumentKind, m.SyncImpl(),
			173, 80, 0,
		)
	})
}

func TestObserverInstruments(t *testing.T) {
	t.Run("float valueobserver", func(t *testing.T) {
		labels := []label.KeyValue{label.String("O", "P")}
		mockSDK, meter := oteltest.NewMeter()
		o := Must(meter).NewFloat64ValueObserver("test.valueobserver.float", func(_ context.Context, result otel.Float64ObserverResult) {
			result.Observe(42.1, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, otel.Float64NumberKind, otel.ValueObserverInstrumentKind, o.AsyncImpl(),
			42.1,
		)
	})
	t.Run("int valueobserver", func(t *testing.T) {
		labels := []label.KeyValue{}
		mockSDK, meter := oteltest.NewMeter()
		o := Must(meter).NewInt64ValueObserver("test.observer.int", func(_ context.Context, result otel.Int64ObserverResult) {
			result.Observe(-142, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, otel.Int64NumberKind, otel.ValueObserverInstrumentKind, o.AsyncImpl(),
			-142,
		)
	})
	t.Run("float sumobserver", func(t *testing.T) {
		labels := []label.KeyValue{label.String("O", "P")}
		mockSDK, meter := oteltest.NewMeter()
		o := Must(meter).NewFloat64SumObserver("test.sumobserver.float", func(_ context.Context, result otel.Float64ObserverResult) {
			result.Observe(42.1, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, otel.Float64NumberKind, otel.SumObserverInstrumentKind, o.AsyncImpl(),
			42.1,
		)
	})
	t.Run("int sumobserver", func(t *testing.T) {
		labels := []label.KeyValue{}
		mockSDK, meter := oteltest.NewMeter()
		o := Must(meter).NewInt64SumObserver("test.observer.int", func(_ context.Context, result otel.Int64ObserverResult) {
			result.Observe(-142, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, otel.Int64NumberKind, otel.SumObserverInstrumentKind, o.AsyncImpl(),
			-142,
		)
	})
	t.Run("float updownsumobserver", func(t *testing.T) {
		labels := []label.KeyValue{label.String("O", "P")}
		mockSDK, meter := oteltest.NewMeter()
		o := Must(meter).NewFloat64UpDownSumObserver("test.updownsumobserver.float", func(_ context.Context, result otel.Float64ObserverResult) {
			result.Observe(42.1, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, otel.Float64NumberKind, otel.UpDownSumObserverInstrumentKind, o.AsyncImpl(),
			42.1,
		)
	})
	t.Run("int updownsumobserver", func(t *testing.T) {
		labels := []label.KeyValue{}
		mockSDK, meter := oteltest.NewMeter()
		o := Must(meter).NewInt64UpDownSumObserver("test.observer.int", func(_ context.Context, result otel.Int64ObserverResult) {
			result.Observe(-142, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, otel.Int64NumberKind, otel.UpDownSumObserverInstrumentKind, o.AsyncImpl(),
			-142,
		)
	})
}

func TestBatchObserverInstruments(t *testing.T) {
	mockSDK, meter := oteltest.NewMeter()

	var obs1 otel.Int64ValueObserver
	var obs2 otel.Float64ValueObserver

	labels := []label.KeyValue{
		label.String("A", "B"),
		label.String("C", "D"),
	}

	cb := Must(meter).NewBatchObserver(
		func(_ context.Context, result otel.BatchObserverResult) {
			result.Observe(labels,
				obs1.Observation(42),
				obs2.Observation(42.0),
			)
		},
	)
	obs1 = cb.NewInt64ValueObserver("test.observer.int")
	obs2 = cb.NewFloat64ValueObserver("test.observer.float")

	mockSDK.RunAsyncInstruments()

	require.Len(t, mockSDK.MeasurementBatches, 1)

	impl1 := obs1.AsyncImpl().Implementation().(*oteltest.Async)
	impl2 := obs2.AsyncImpl().Implementation().(*oteltest.Async)

	require.NotNil(t, impl1)
	require.NotNil(t, impl2)

	got := mockSDK.MeasurementBatches[0]
	require.Equal(t, labels, got.Labels)
	require.Len(t, got.Measurements, 2)

	m1 := got.Measurements[0]
	require.Equal(t, impl1, m1.Instrument.Implementation().(*oteltest.Async))
	require.Equal(t, 0, m1.Number.CompareNumber(otel.Int64NumberKind, oteltest.ResolveNumberByKind(t, otel.Int64NumberKind, 42)))

	m2 := got.Measurements[1]
	require.Equal(t, impl2, m2.Instrument.Implementation().(*oteltest.Async))
	require.Equal(t, 0, m2.Number.CompareNumber(otel.Float64NumberKind, oteltest.ResolveNumberByKind(t, otel.Float64NumberKind, 42)))
}

func checkObserverBatch(t *testing.T, labels []label.KeyValue, mock *oteltest.MeterImpl, nkind otel.NumberKind, mkind otel.InstrumentKind, observer otel.AsyncImpl, expected float64) {
	t.Helper()
	assert.Len(t, mock.MeasurementBatches, 1)
	if len(mock.MeasurementBatches) < 1 {
		return
	}
	o := observer.Implementation().(*oteltest.Async)
	if !assert.NotNil(t, o) {
		return
	}
	got := mock.MeasurementBatches[0]
	assert.Equal(t, labels, got.Labels)
	assert.Len(t, got.Measurements, 1)
	if len(got.Measurements) < 1 {
		return
	}
	measurement := got.Measurements[0]
	require.Equal(t, mkind, measurement.Instrument.Descriptor().InstrumentKind())
	assert.Equal(t, o, measurement.Instrument.Implementation().(*oteltest.Async))
	ft := oteltest.ResolveNumberByKind(t, nkind, expected)
	assert.Equal(t, 0, measurement.Number.CompareNumber(nkind, ft))
}

type testWrappedMeter struct {
}

var _ otel.MeterImpl = testWrappedMeter{}

func (testWrappedMeter) RecordBatch(context.Context, []label.KeyValue, ...otel.Measurement) {
}

func (testWrappedMeter) NewSyncInstrument(_ otel.Descriptor) (otel.SyncImpl, error) {
	return nil, nil
}

func (testWrappedMeter) NewAsyncInstrument(_ otel.Descriptor, _ otel.AsyncRunner) (otel.AsyncImpl, error) {
	return nil, errors.New("Test wrap error")
}

func TestWrappedInstrumentError(t *testing.T) {
	impl := &testWrappedMeter{}
	meter := otel.WrapMeterImpl(impl, "test")

	valuerecorder, err := meter.NewInt64ValueRecorder("test.valuerecorder")

	require.Equal(t, err, otel.ErrSDKReturnedNilImpl)
	require.NotNil(t, valuerecorder.SyncImpl())

	observer, err := meter.NewInt64ValueObserver("test.observer", func(_ context.Context, result otel.Int64ObserverResult) {})

	require.NotNil(t, err)
	require.NotNil(t, observer.AsyncImpl())
}

func TestNilCallbackObserverNoop(t *testing.T) {
	// Tests that a nil callback yields a no-op observer without error.
	_, meter := oteltest.NewMeter()

	observer := Must(meter).NewInt64ValueObserver("test.observer", nil)

	_, ok := observer.AsyncImpl().(otel.NoopAsync)
	require.True(t, ok)
}
