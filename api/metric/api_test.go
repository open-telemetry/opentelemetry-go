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

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/metric/metrictest"
	mockTest "go.opentelemetry.io/otel/api/metric/metrictest"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/unit"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Must = metric.Must

func checkSyncBatches(ctx context.Context, t *testing.T, labels []label.KeyValue, mock *mockTest.MeterImpl, nkind metric.NumberKind, mkind metric.InstrumentKind, instrument metric.InstrumentImpl, expected ...float64) {
	t.Helper()

	batchesCount := len(mock.MeasurementBatches)
	if len(mock.MeasurementBatches) != len(expected) {
		t.Errorf("Expected %d recorded measurement batches, got %d", batchesCount, len(mock.MeasurementBatches))
	}
	recorded := metrictest.AsStructs(mock.MeasurementBatches)

	for i, batch := range mock.MeasurementBatches {
		if len(batch.Measurements) != 1 {
			t.Errorf("Expected 1 measurement in batch %d, got %d", i, len(batch.Measurements))
		}

		measurement := batch.Measurements[0]
		descriptor := measurement.Instrument.Descriptor()

		expected := metrictest.Measured{
			Name:                descriptor.Name(),
			InstrumentationName: descriptor.InstrumentationName(),
			Labels:              metrictest.LabelsToMap(labels...),
			Number:              metrictest.ResolveNumberByKind(t, nkind, expected[i]),
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
			name: "unit override",
			opts: []metric.InstrumentOption{
				metric.WithUnit("s"),
				metric.WithUnit("h"),
			},
			desc: "",
			unit: "h",
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing counter case %s (%d)", tt.name, idx)
		if diff := cmp.Diff(metric.NewInstrumentConfig(tt.opts...), metric.InstrumentConfig{
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
		mockSDK, meter := mockTest.NewMeter()
		c := Must(meter).NewFloat64Counter("test.counter.float")
		ctx := context.Background()
		labels := []label.KeyValue{label.String("A", "B")}
		c.Add(ctx, 1994.1, labels...)
		boundInstrument := c.Bind(labels...)
		boundInstrument.Add(ctx, -742)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		checkSyncBatches(ctx, t, labels, mockSDK, metric.Float64NumberKind, metric.CounterInstrumentKind, c.SyncImpl(),
			1994.1, -742, 42,
		)
	})
	t.Run("int64 counter", func(t *testing.T) {
		mockSDK, meter := mockTest.NewMeter()
		c := Must(meter).NewInt64Counter("test.counter.int")
		ctx := context.Background()
		labels := []label.KeyValue{label.String("A", "B"), label.String("C", "D")}
		c.Add(ctx, 42, labels...)
		boundInstrument := c.Bind(labels...)
		boundInstrument.Add(ctx, 4200)
		meter.RecordBatch(ctx, labels, c.Measurement(420000))
		checkSyncBatches(ctx, t, labels, mockSDK, metric.Int64NumberKind, metric.CounterInstrumentKind, c.SyncImpl(),
			42, 4200, 420000,
		)

	})
	t.Run("int64 updowncounter", func(t *testing.T) {
		mockSDK, meter := mockTest.NewMeter()
		c := Must(meter).NewInt64UpDownCounter("test.updowncounter.int")
		ctx := context.Background()
		labels := []label.KeyValue{label.String("A", "B"), label.String("C", "D")}
		c.Add(ctx, 100, labels...)
		boundInstrument := c.Bind(labels...)
		boundInstrument.Add(ctx, -100)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		checkSyncBatches(ctx, t, labels, mockSDK, metric.Int64NumberKind, metric.UpDownCounterInstrumentKind, c.SyncImpl(),
			100, -100, 42,
		)
	})
	t.Run("float64 updowncounter", func(t *testing.T) {
		mockSDK, meter := mockTest.NewMeter()
		c := Must(meter).NewFloat64UpDownCounter("test.updowncounter.float")
		ctx := context.Background()
		labels := []label.KeyValue{label.String("A", "B"), label.String("C", "D")}
		c.Add(ctx, 100.1, labels...)
		boundInstrument := c.Bind(labels...)
		boundInstrument.Add(ctx, -76)
		meter.RecordBatch(ctx, labels, c.Measurement(-100.1))
		checkSyncBatches(ctx, t, labels, mockSDK, metric.Float64NumberKind, metric.UpDownCounterInstrumentKind, c.SyncImpl(),
			100.1, -76, -100.1,
		)
	})
}

func TestValueRecorder(t *testing.T) {
	t.Run("float64 valuerecorder", func(t *testing.T) {
		mockSDK, meter := mockTest.NewMeter()
		m := Must(meter).NewFloat64ValueRecorder("test.valuerecorder.float")
		ctx := context.Background()
		labels := []label.KeyValue{}
		m.Record(ctx, 42, labels...)
		boundInstrument := m.Bind(labels...)
		boundInstrument.Record(ctx, 0)
		meter.RecordBatch(ctx, labels, m.Measurement(-100.5))
		checkSyncBatches(ctx, t, labels, mockSDK, metric.Float64NumberKind, metric.ValueRecorderInstrumentKind, m.SyncImpl(),
			42, 0, -100.5,
		)
	})
	t.Run("int64 valuerecorder", func(t *testing.T) {
		mockSDK, meter := mockTest.NewMeter()
		m := Must(meter).NewInt64ValueRecorder("test.valuerecorder.int")
		ctx := context.Background()
		labels := []label.KeyValue{label.Int("I", 1)}
		m.Record(ctx, 173, labels...)
		boundInstrument := m.Bind(labels...)
		boundInstrument.Record(ctx, 80)
		meter.RecordBatch(ctx, labels, m.Measurement(0))
		checkSyncBatches(ctx, t, labels, mockSDK, metric.Int64NumberKind, metric.ValueRecorderInstrumentKind, m.SyncImpl(),
			173, 80, 0,
		)
	})
}

func TestObserverInstruments(t *testing.T) {
	t.Run("float valueobserver", func(t *testing.T) {
		labels := []label.KeyValue{label.String("O", "P")}
		mockSDK, meter := mockTest.NewMeter()
		o := Must(meter).NewFloat64ValueObserver("test.valueobserver.float", func(_ context.Context, result metric.Float64ObserverResult) {
			result.Observe(42.1, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, metric.Float64NumberKind, metric.ValueObserverInstrumentKind, o.AsyncImpl(),
			42.1,
		)
	})
	t.Run("int valueobserver", func(t *testing.T) {
		labels := []label.KeyValue{}
		mockSDK, meter := mockTest.NewMeter()
		o := Must(meter).NewInt64ValueObserver("test.observer.int", func(_ context.Context, result metric.Int64ObserverResult) {
			result.Observe(-142, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, metric.Int64NumberKind, metric.ValueObserverInstrumentKind, o.AsyncImpl(),
			-142,
		)
	})
	t.Run("float sumobserver", func(t *testing.T) {
		labels := []label.KeyValue{label.String("O", "P")}
		mockSDK, meter := mockTest.NewMeter()
		o := Must(meter).NewFloat64SumObserver("test.sumobserver.float", func(_ context.Context, result metric.Float64ObserverResult) {
			result.Observe(42.1, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, metric.Float64NumberKind, metric.SumObserverInstrumentKind, o.AsyncImpl(),
			42.1,
		)
	})
	t.Run("int sumobserver", func(t *testing.T) {
		labels := []label.KeyValue{}
		mockSDK, meter := mockTest.NewMeter()
		o := Must(meter).NewInt64SumObserver("test.observer.int", func(_ context.Context, result metric.Int64ObserverResult) {
			result.Observe(-142, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, metric.Int64NumberKind, metric.SumObserverInstrumentKind, o.AsyncImpl(),
			-142,
		)
	})
	t.Run("float updownsumobserver", func(t *testing.T) {
		labels := []label.KeyValue{label.String("O", "P")}
		mockSDK, meter := mockTest.NewMeter()
		o := Must(meter).NewFloat64UpDownSumObserver("test.updownsumobserver.float", func(_ context.Context, result metric.Float64ObserverResult) {
			result.Observe(42.1, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, metric.Float64NumberKind, metric.UpDownSumObserverInstrumentKind, o.AsyncImpl(),
			42.1,
		)
	})
	t.Run("int updownsumobserver", func(t *testing.T) {
		labels := []label.KeyValue{}
		mockSDK, meter := mockTest.NewMeter()
		o := Must(meter).NewInt64UpDownSumObserver("test.observer.int", func(_ context.Context, result metric.Int64ObserverResult) {
			result.Observe(-142, labels...)
		})
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, metric.Int64NumberKind, metric.UpDownSumObserverInstrumentKind, o.AsyncImpl(),
			-142,
		)
	})
}

func TestBatchObserverInstruments(t *testing.T) {
	mockSDK, meter := mockTest.NewMeter()

	var obs1 metric.Int64ValueObserver
	var obs2 metric.Float64ValueObserver

	labels := []label.KeyValue{
		label.String("A", "B"),
		label.String("C", "D"),
	}

	cb := Must(meter).NewBatchObserver(
		func(_ context.Context, result metric.BatchObserverResult) {
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

	impl1 := obs1.AsyncImpl().Implementation().(*mockTest.Async)
	impl2 := obs2.AsyncImpl().Implementation().(*mockTest.Async)

	require.NotNil(t, impl1)
	require.NotNil(t, impl2)

	got := mockSDK.MeasurementBatches[0]
	require.Equal(t, labels, got.Labels)
	require.Len(t, got.Measurements, 2)

	m1 := got.Measurements[0]
	require.Equal(t, impl1, m1.Instrument.Implementation().(*mockTest.Async))
	require.Equal(t, 0, m1.Number.CompareNumber(metric.Int64NumberKind, mockTest.ResolveNumberByKind(t, metric.Int64NumberKind, 42)))

	m2 := got.Measurements[1]
	require.Equal(t, impl2, m2.Instrument.Implementation().(*mockTest.Async))
	require.Equal(t, 0, m2.Number.CompareNumber(metric.Float64NumberKind, mockTest.ResolveNumberByKind(t, metric.Float64NumberKind, 42)))
}

func checkObserverBatch(t *testing.T, labels []label.KeyValue, mock *mockTest.MeterImpl, nkind metric.NumberKind, mkind metric.InstrumentKind, observer metric.AsyncImpl, expected float64) {
	t.Helper()
	assert.Len(t, mock.MeasurementBatches, 1)
	if len(mock.MeasurementBatches) < 1 {
		return
	}
	o := observer.Implementation().(*mockTest.Async)
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
	assert.Equal(t, o, measurement.Instrument.Implementation().(*mockTest.Async))
	ft := mockTest.ResolveNumberByKind(t, nkind, expected)
	assert.Equal(t, 0, measurement.Number.CompareNumber(nkind, ft))
}

type testWrappedMeter struct {
}

var _ metric.MeterImpl = testWrappedMeter{}

func (testWrappedMeter) RecordBatch(context.Context, []label.KeyValue, ...metric.Measurement) {
}

func (testWrappedMeter) NewSyncInstrument(_ metric.Descriptor) (metric.SyncImpl, error) {
	return nil, nil
}

func (testWrappedMeter) NewAsyncInstrument(_ metric.Descriptor, _ metric.AsyncRunner) (metric.AsyncImpl, error) {
	return nil, errors.New("Test wrap error")
}

func TestWrappedInstrumentError(t *testing.T) {
	impl := &testWrappedMeter{}
	meter := metric.WrapMeterImpl(impl, "test")

	valuerecorder, err := meter.NewInt64ValueRecorder("test.valuerecorder")

	require.Equal(t, err, metric.ErrSDKReturnedNilImpl)
	require.NotNil(t, valuerecorder.SyncImpl())

	observer, err := meter.NewInt64ValueObserver("test.observer", func(_ context.Context, result metric.Int64ObserverResult) {})

	require.NotNil(t, err)
	require.NotNil(t, observer.AsyncImpl())
}

func TestNilCallbackObserverNoop(t *testing.T) {
	// Tests that a nil callback yields a no-op observer without error.
	_, meter := mockTest.NewMeter()

	observer := Must(meter).NewInt64ValueObserver("test.observer", nil)

	_, ok := observer.AsyncImpl().(metric.NoopAsync)
	require.True(t, ok)
}
