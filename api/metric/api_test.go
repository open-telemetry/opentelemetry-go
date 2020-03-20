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
	"errors"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/unit"
	mockTest "go.opentelemetry.io/otel/internal/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Must = metric.Must

func TestOptions(t *testing.T) {
	type testcase struct {
		name     string
		opts     []metric.Option
		keys     []core.Key
		desc     string
		unit     unit.Unit
		resource resource.Resource
	}
	testcases := []testcase{
		{
			name:     "no opts",
			opts:     nil,
			keys:     nil,
			desc:     "",
			unit:     "",
			resource: resource.Resource{},
		},
		{
			name: "keys keys keys",
			opts: []metric.Option{
				metric.WithKeys(key.New("foo"), key.New("foo2")),
				metric.WithKeys(key.New("bar"), key.New("bar2")),
				metric.WithKeys(key.New("baz"), key.New("baz2")),
			},
			keys: []core.Key{
				key.New("foo"), key.New("foo2"),
				key.New("bar"), key.New("bar2"),
				key.New("baz"), key.New("baz2"),
			},
			desc:     "",
			unit:     "",
			resource: resource.Resource{},
		},
		{
			name: "description",
			opts: []metric.Option{
				metric.WithDescription("stuff"),
			},
			keys:     nil,
			desc:     "stuff",
			unit:     "",
			resource: resource.Resource{},
		},
		{
			name: "description override",
			opts: []metric.Option{
				metric.WithDescription("stuff"),
				metric.WithDescription("things"),
			},
			keys:     nil,
			desc:     "things",
			unit:     "",
			resource: resource.Resource{},
		},
		{
			name: "unit",
			opts: []metric.Option{
				metric.WithUnit("s"),
			},
			keys:     nil,
			desc:     "",
			unit:     "s",
			resource: resource.Resource{},
		},
		{
			name: "unit override",
			opts: []metric.Option{
				metric.WithUnit("s"),
				metric.WithUnit("h"),
			},
			keys:     nil,
			desc:     "",
			unit:     "h",
			resource: resource.Resource{},
		},
		{
			name: "resource override",
			opts: []metric.Option{
				metric.WithResource(*resource.New(key.New("name").String("test-name"))),
			},
			keys:     nil,
			desc:     "",
			unit:     "",
			resource: *resource.New(key.New("name").String("test-name")),
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing counter case %s (%d)", tt.name, idx)
		if diff := cmp.Diff(metric.Configure(tt.opts), metric.Config{
			Description: tt.desc,
			Unit:        tt.unit,
			Keys:        tt.keys,
			Resource:    tt.resource,
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
		labels := meter.Labels()
		c.Add(ctx, 42, labels)
		boundInstrument := c.Bind(labels)
		boundInstrument.Add(ctx, 42)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		t.Log("Testing float counter")
		checkBatches(t, ctx, labels, mockSDK, core.Float64NumberKind, c.SyncImpl())
	}
	{
		mockSDK, meter := mockTest.NewMeter()
		c := Must(meter).NewInt64Counter("test.counter.int")
		ctx := context.Background()
		labels := meter.Labels()
		c.Add(ctx, 42, labels)
		boundInstrument := c.Bind(labels)
		boundInstrument.Add(ctx, 42)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		t.Log("Testing int counter")
		checkBatches(t, ctx, labels, mockSDK, core.Int64NumberKind, c.SyncImpl())
	}
}

func TestMeasure(t *testing.T) {
	{
		mockSDK, meter := mockTest.NewMeter()
		m := Must(meter).NewFloat64Measure("test.measure.float")
		ctx := context.Background()
		labels := meter.Labels()
		m.Record(ctx, 42, labels)
		boundInstrument := m.Bind(labels)
		boundInstrument.Record(ctx, 42)
		meter.RecordBatch(ctx, labels, m.Measurement(42))
		t.Log("Testing float measure")
		checkBatches(t, ctx, labels, mockSDK, core.Float64NumberKind, m.SyncImpl())
	}
	{
		mockSDK, meter := mockTest.NewMeter()
		m := Must(meter).NewInt64Measure("test.measure.int")
		ctx := context.Background()
		labels := meter.Labels()
		m.Record(ctx, 42, labels)
		boundInstrument := m.Bind(labels)
		boundInstrument.Record(ctx, 42)
		meter.RecordBatch(ctx, labels, m.Measurement(42))
		t.Log("Testing int measure")
		checkBatches(t, ctx, labels, mockSDK, core.Int64NumberKind, m.SyncImpl())
	}
}

func TestObserver(t *testing.T) {
	{
		mockSDK, meter := mockTest.NewMeter()
		labels := meter.Labels()
		o := Must(meter).RegisterFloat64Observer("test.observer.float", func(result metric.Float64ObserverResult) {
			result.Observe(42, labels)
		})
		t.Log("Testing float observer")

		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, core.Float64NumberKind, o.AsyncImpl())
	}
	{
		mockSDK, meter := mockTest.NewMeter()
		labels := meter.Labels()
		o := Must(meter).RegisterInt64Observer("test.observer.int", func(result metric.Int64ObserverResult) {
			result.Observe(42, labels)
		})
		t.Log("Testing int observer")
		mockSDK.RunAsyncInstruments()
		checkObserverBatch(t, labels, mockSDK, core.Int64NumberKind, o.AsyncImpl())
	}
}

func checkBatches(t *testing.T, ctx context.Context, labels metric.LabelSet, mock *mockTest.Meter, kind core.NumberKind, instrument metric.InstrumentImpl) {
	t.Helper()
	if len(mock.MeasurementBatches) != 3 {
		t.Errorf("Expected 3 recorded measurement batches, got %d", len(mock.MeasurementBatches))
	}
	ourInstrument := instrument.Implementation().(*mockTest.Sync)
	ourLabelSet := labels.(*mockTest.LabelSet)
	minLen := 3
	if minLen > len(mock.MeasurementBatches) {
		minLen = len(mock.MeasurementBatches)
	}
	for i := 0; i < minLen; i++ {
		got := mock.MeasurementBatches[i]
		if got.Ctx != ctx {
			d := func(c context.Context) string {
				return fmt.Sprintf("(ptr: %p, ctx %#v)", c, c)
			}
			t.Errorf("Wrong recorded context in batch %d, expected %s, got %s", i, d(ctx), d(got.Ctx))
		}
		if got.LabelSet != ourLabelSet {
			d := func(l *mockTest.LabelSet) string {
				return fmt.Sprintf("(ptr: %p, labels %#v)", l, l.Labels)
			}
			t.Errorf("Wrong recorded label set in batch %d, expected %s, got %s", i, d(ourLabelSet), d(got.LabelSet))
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

func checkObserverBatch(t *testing.T, labels metric.LabelSet, mock *mockTest.Meter, kind core.NumberKind, observer metric.AsyncImpl) {
	t.Helper()
	assert.Len(t, mock.MeasurementBatches, 1)
	if len(mock.MeasurementBatches) < 1 {
		return
	}
	o := observer.Implementation().(*mockTest.Async)
	if !assert.NotNil(t, o) {
		return
	}
	ourLabelSet := labels.(*mockTest.LabelSet)
	got := mock.MeasurementBatches[0]
	assert.Equal(t, ourLabelSet, got.LabelSet)
	assert.Len(t, got.Measurements, 1)
	if len(got.Measurements) < 1 {
		return
	}
	measurement := got.Measurements[0]
	assert.Equal(t, o, measurement.Instrument.Implementation().(*mockTest.Async))
	ft := fortyTwo(t, kind)
	assert.Equal(t, 0, measurement.Number.CompareNumber(kind, ft))
}

func fortyTwo(t *testing.T, kind core.NumberKind) core.Number {
	t.Helper()
	switch kind {
	case core.Int64NumberKind:
		return core.NewInt64Number(42)
	case core.Float64NumberKind:
		return core.NewFloat64Number(42)
	}
	t.Errorf("Invalid value kind %q", kind)
	return core.NewInt64Number(0)
}

type testWrappedMeter struct {
}

var _ metric.MeterImpl = testWrappedMeter{}

func (testWrappedMeter) Labels(...core.KeyValue) metric.LabelSet {
	return nil
}

func (testWrappedMeter) RecordBatch(context.Context, metric.LabelSet, ...metric.Measurement) {
}

func (testWrappedMeter) NewSyncInstrument(_ metric.Descriptor) (metric.SyncImpl, error) {
	return nil, nil
}

func (testWrappedMeter) NewAsyncInstrument(_ metric.Descriptor, _ func(func(core.Number, metric.LabelSet))) (metric.AsyncImpl, error) {
	return nil, errors.New("Test wrap error")
}

func TestWrappedInstrumentError(t *testing.T) {
	impl := &testWrappedMeter{}
	meter := metric.WrapMeterImpl(impl)

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
