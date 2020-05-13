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
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/unit"
	mockTest "go.opentelemetry.io/otel/internal/metric"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Must = metric.Must

func TestOptions(t *testing.T) {
	type testcase struct {
		name string
		opts []metric.Option
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
			opts: []metric.Option{
				metric.WithDescription("stuff"),
			},
			desc: "stuff",
			unit: "",
		},
		{
			name: "description override",
			opts: []metric.Option{
				metric.WithDescription("stuff"),
				metric.WithDescription("things"),
			},
			desc: "things",
			unit: "",
		},
		{
			name: "unit",
			opts: []metric.Option{
				metric.WithUnit("s"),
			},
			desc: "",
			unit: "s",
		},
		{
			name: "unit override",
			opts: []metric.Option{
				metric.WithUnit("s"),
				metric.WithUnit("h"),
			},
			desc: "",
			unit: "h",
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing counter case %s (%d)", tt.name, idx)
		if diff := cmp.Diff(metric.Configure(tt.opts), metric.Config{
			Description: tt.desc,
			Unit:        tt.unit,
		}); diff != "" {
			t.Errorf("Compare options: -got +want %s", diff)
		}
	}
}

func TestCounter(t *testing.T) {
	{
		mockSDK, meter := mockTest.NewMeter()
		c := Must(meter).NewFloat64Counter("test.counter.float")
		ctx := context.Background()
		labels := []kv.KeyValue{kv.String("A", "B")}
		c.Add(ctx, 42, labels...)
		boundInstrument := c.Bind(labels...)
		boundInstrument.Add(ctx, 42)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		t.Log("Testing float counter")
		checkBatches(t, ctx, labels, mockSDK, metric.Float64NumberKind, c.SyncImpl())
	}
	{
		mockSDK, meter := mockTest.NewMeter()
		c := Must(meter).NewInt64Counter("test.counter.int")
		ctx := context.Background()
		labels := []kv.KeyValue{kv.String("A", "B"), kv.String("C", "D")}
		c.Add(ctx, 42, labels...)
		boundInstrument := c.Bind(labels...)
		boundInstrument.Add(ctx, 42)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		t.Log("Testing int counter")
		checkBatches(t, ctx, labels, mockSDK, metric.Int64NumberKind, c.SyncImpl())
	}
}

func TestMeasure(t *testing.T) {
	{
		mockSDK, meter := mockTest.NewMeter()
		m := Must(meter).NewFloat64Measure("test.measure.float")
		ctx := context.Background()
		labels := []kv.KeyValue{}
		m.Record(ctx, 42, labels...)
		boundInstrument := m.Bind(labels...)
		boundInstrument.Record(ctx, 42)
		meter.RecordBatch(ctx, labels, m.Measurement(42))
		t.Log("Testing float measure")
		checkBatches(t, ctx, labels, mockSDK, metric.Float64NumberKind, m.SyncImpl())
	}
	{
		mockSDK, meter := mockTest.NewMeter()
		m := Must(meter).NewInt64Measure("test.measure.int")
		ctx := context.Background()
		labels := []kv.KeyValue{kv.Int("I", 1)}
		m.Record(ctx, 42, labels...)
		boundInstrument := m.Bind(labels...)
		boundInstrument.Record(ctx, 42)
		meter.RecordBatch(ctx, labels, m.Measurement(42))
		t.Log("Testing int measure")
		checkBatches(t, ctx, labels, mockSDK, metric.Int64NumberKind, m.SyncImpl())
	}
}

func TestObserver(t *testing.T) {
	{
		labels := []kv.KeyValue{kv.String("O", "P")}
		mockSDK, meter := mockTest.NewMeter()
		o := Must(meter).RegisterFloat64Observer("test.observer.float", func(result metric.Float64ObserverResult) {
			result.Observe(42, labels...)
		})
		t.Log("Testing float observer")

		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, metric.Float64NumberKind, o.AsyncImpl())
	}
	{
		labels := []kv.KeyValue{}
		mockSDK, meter := mockTest.NewMeter()
		o := Must(meter).RegisterInt64Observer("test.observer.int", func(result metric.Int64ObserverResult) {
			result.Observe(42, labels...)
		})
		t.Log("Testing int observer")
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, metric.Int64NumberKind, o.AsyncImpl())
	}
}

func checkBatches(t *testing.T, ctx context.Context, labels []kv.KeyValue, mock *mockTest.MeterImpl, kind metric.NumberKind, instrument metric.InstrumentImpl) {
	t.Helper()
	if len(mock.MeasurementBatches) != 3 {
		t.Errorf("Expected 3 recorded measurement batches, got %d", len(mock.MeasurementBatches))
	}
	ourInstrument := instrument.Implementation().(*mockTest.Sync)
	for i, got := range mock.MeasurementBatches {
		if got.Ctx != ctx {
			d := func(c context.Context) string {
				return fmt.Sprintf("(ptr: %p, ctx %#v)", c, c)
			}
			t.Errorf("Wrong recorded context in batch %d, expected %s, got %s", i, d(ctx), d(got.Ctx))
		}
		if !assert.Equal(t, got.Labels, labels) {
			t.Errorf("Wrong recorded label set in batch %d, expected %v, got %v", i, labels, got.Labels)
		}
		if len(got.Measurements) != 1 {
			t.Errorf("Expected 1 measurement in batch %d, got %d", i, len(got.Measurements))
		}
		minMLen := 1
		if minMLen > len(got.Measurements) {
			minMLen = len(got.Measurements)
		}
		for j := 0; j < minMLen; j++ {
			measurement := got.Measurements[j]
			if measurement.Instrument.Implementation() != ourInstrument {
				d := func(iface interface{}) string {
					i := iface.(*mockTest.Instrument)
					return fmt.Sprintf("(ptr: %p, instrument %#v)", i, i)
				}
				t.Errorf("Wrong recorded instrument in measurement %d in batch %d, expected %s, got %s", j, i, d(ourInstrument), d(measurement.Instrument.Implementation()))
			}
			ft := fortyTwo(t, kind)
			if measurement.Number.CompareNumber(kind, ft) != 0 {
				t.Errorf("Wrong recorded value in measurement %d in batch %d, expected %s, got %s", j, i, ft.Emit(kind), measurement.Number.Emit(kind))
			}
		}
	}
}

func checkObserverBatch(t *testing.T, labels []kv.KeyValue, mock *mockTest.MeterImpl, kind metric.NumberKind, observer metric.AsyncImpl) {
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
	assert.Equal(t, o, measurement.Instrument.Implementation().(*mockTest.Async))
	ft := fortyTwo(t, kind)
	assert.Equal(t, 0, measurement.Number.CompareNumber(kind, ft))
}

func fortyTwo(t *testing.T, kind metric.NumberKind) metric.Number {
	t.Helper()
	switch kind {
	case metric.Int64NumberKind:
		return metric.NewInt64Number(42)
	case metric.Float64NumberKind:
		return metric.NewFloat64Number(42)
	}
	t.Errorf("Invalid value kind %q", kind)
	return metric.NewInt64Number(0)
}

type testWrappedMeter struct {
}

var _ metric.MeterImpl = testWrappedMeter{}

func (testWrappedMeter) RecordBatch(context.Context, []kv.KeyValue, ...metric.Measurement) {
}

func (testWrappedMeter) NewSyncInstrument(_ metric.Descriptor) (metric.SyncImpl, error) {
	return nil, nil
}

func (testWrappedMeter) NewAsyncInstrument(_ metric.Descriptor, _ func(func(metric.Number, []kv.KeyValue))) (metric.AsyncImpl, error) {
	return nil, errors.New("Test wrap error")
}

func TestWrappedInstrumentError(t *testing.T) {
	impl := &testWrappedMeter{}
	meter := metric.WrapMeterImpl(impl, "test")

	measure, err := meter.NewInt64Measure("test.measure")

	require.Equal(t, err, metric.ErrSDKReturnedNilImpl)
	require.NotNil(t, measure.SyncImpl())

	observer, err := meter.RegisterInt64Observer("test.observer", func(result metric.Int64ObserverResult) {})

	require.NotNil(t, err)
	require.NotNil(t, observer.AsyncImpl())
}

func TestNilCallbackObserverNoop(t *testing.T) {
	// Tests that a nil callback yields a no-op observer without error.
	_, meter := mockTest.NewMeter()

	observer := Must(meter).RegisterInt64Observer("test.observer", nil)

	_, ok := observer.AsyncImpl().(metric.NoopAsync)
	require.True(t, ok)
}
