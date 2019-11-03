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
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/unit"
	mock "go.opentelemetry.io/otel/internal/metric"

	"github.com/google/go-cmp/cmp"
)

func TestCounterOptions(t *testing.T) {
	type testcase struct {
		name string
		opts []metric.CounterOptionApplier
		keys []core.Key
		desc string
		unit unit.Unit
		alt  bool
	}
	testcases := []testcase{
		{
			name: "no opts",
			opts: nil,
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
		{
			name: "keys keys keys",
			opts: []metric.CounterOptionApplier{
				metric.WithKeys(key.New("foo"), key.New("foo2")),
				metric.WithKeys(key.New("bar"), key.New("bar2")),
				metric.WithKeys(key.New("baz"), key.New("baz2")),
			},
			keys: []core.Key{
				key.New("foo"), key.New("foo2"),
				key.New("bar"), key.New("bar2"),
				key.New("baz"), key.New("baz2"),
			},
			desc: "",
			unit: "",
			alt:  false,
		},
		{
			name: "description",
			opts: []metric.CounterOptionApplier{
				metric.WithDescription("stuff"),
			},
			keys: nil,
			desc: "stuff",
			unit: "",
			alt:  false,
		},
		{
			name: "description override",
			opts: []metric.CounterOptionApplier{
				metric.WithDescription("stuff"),
				metric.WithDescription("things"),
			},
			keys: nil,
			desc: "things",
			unit: "",
			alt:  false,
		},
		{
			name: "unit",
			opts: []metric.CounterOptionApplier{
				metric.WithUnit("s"),
			},
			keys: nil,
			desc: "",
			unit: "s",
			alt:  false,
		},
		{
			name: "unit override",
			opts: []metric.CounterOptionApplier{
				metric.WithUnit("s"),
				metric.WithUnit("h"),
			},
			keys: nil,
			desc: "",
			unit: "h",
			alt:  false,
		},
		{
			name: "nonmonotonic",
			opts: []metric.CounterOptionApplier{
				metric.WithMonotonic(false),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  true,
		},
		{
			name: "nonmonotonic, but not really",
			opts: []metric.CounterOptionApplier{
				metric.WithMonotonic(false),
				metric.WithMonotonic(true),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing counter case %s (%d)", tt.name, idx)
		opts := &metric.Options{}
		metric.ApplyCounterOptions(opts, tt.opts...)
		checkOptions(t, opts, &metric.Options{
			Description: tt.desc,
			Unit:        tt.unit,
			Keys:        tt.keys,
			Alternate:   tt.alt,
		})
	}
}

func TestGaugeOptions(t *testing.T) {
	type testcase struct {
		name string
		opts []metric.GaugeOptionApplier
		keys []core.Key
		desc string
		unit unit.Unit
		alt  bool
	}
	testcases := []testcase{
		{
			name: "no opts",
			opts: nil,
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
		{
			name: "keys keys keys",
			opts: []metric.GaugeOptionApplier{
				metric.WithKeys(key.New("foo"), key.New("foo2")),
				metric.WithKeys(key.New("bar"), key.New("bar2")),
				metric.WithKeys(key.New("baz"), key.New("baz2")),
			},
			keys: []core.Key{
				key.New("foo"), key.New("foo2"),
				key.New("bar"), key.New("bar2"),
				key.New("baz"), key.New("baz2"),
			},
			desc: "",
			unit: "",
			alt:  false,
		},
		{
			name: "description",
			opts: []metric.GaugeOptionApplier{
				metric.WithDescription("stuff"),
			},
			keys: nil,
			desc: "stuff",
			unit: "",
			alt:  false,
		},
		{
			name: "description override",
			opts: []metric.GaugeOptionApplier{
				metric.WithDescription("stuff"),
				metric.WithDescription("things"),
			},
			keys: nil,
			desc: "things",
			unit: "",
			alt:  false,
		},
		{
			name: "unit",
			opts: []metric.GaugeOptionApplier{
				metric.WithUnit("s"),
			},
			keys: nil,
			desc: "",
			unit: "s",
			alt:  false,
		},
		{
			name: "unit override",
			opts: []metric.GaugeOptionApplier{
				metric.WithUnit("s"),
				metric.WithUnit("h"),
			},
			keys: nil,
			desc: "",
			unit: "h",
			alt:  false,
		},
		{
			name: "monotonic",
			opts: []metric.GaugeOptionApplier{
				metric.WithMonotonic(true),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  true,
		},
		{
			name: "monotonic, but not really",
			opts: []metric.GaugeOptionApplier{
				metric.WithMonotonic(true),
				metric.WithMonotonic(false),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing gauge case %s (%d)", tt.name, idx)
		opts := &metric.Options{}
		metric.ApplyGaugeOptions(opts, tt.opts...)
		checkOptions(t, opts, &metric.Options{
			Description: tt.desc,
			Unit:        tt.unit,
			Keys:        tt.keys,
			Alternate:   tt.alt,
		})
	}
}

func TestMeasureOptions(t *testing.T) {
	type testcase struct {
		name string
		opts []metric.MeasureOptionApplier
		keys []core.Key
		desc string
		unit unit.Unit
		alt  bool
	}
	testcases := []testcase{
		{
			name: "no opts",
			opts: nil,
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
		{
			name: "keys keys keys",
			opts: []metric.MeasureOptionApplier{
				metric.WithKeys(key.New("foo"), key.New("foo2")),
				metric.WithKeys(key.New("bar"), key.New("bar2")),
				metric.WithKeys(key.New("baz"), key.New("baz2")),
			},
			keys: []core.Key{
				key.New("foo"), key.New("foo2"),
				key.New("bar"), key.New("bar2"),
				key.New("baz"), key.New("baz2"),
			},
			desc: "",
			unit: "",
			alt:  false,
		},
		{
			name: "description",
			opts: []metric.MeasureOptionApplier{
				metric.WithDescription("stuff"),
			},
			keys: nil,
			desc: "stuff",
			unit: "",
			alt:  false,
		},
		{
			name: "description override",
			opts: []metric.MeasureOptionApplier{
				metric.WithDescription("stuff"),
				metric.WithDescription("things"),
			},
			keys: nil,
			desc: "things",
			unit: "",
			alt:  false,
		},
		{
			name: "unit",
			opts: []metric.MeasureOptionApplier{
				metric.WithUnit("s"),
			},
			keys: nil,
			desc: "",
			unit: "s",
			alt:  false,
		},
		{
			name: "unit override",
			opts: []metric.MeasureOptionApplier{
				metric.WithUnit("s"),
				metric.WithUnit("h"),
			},
			keys: nil,
			desc: "",
			unit: "h",
			alt:  false,
		},
		{
			name: "not absolute",
			opts: []metric.MeasureOptionApplier{
				metric.WithAbsolute(false),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  true,
		},
		{
			name: "not absolute, but not really",
			opts: []metric.MeasureOptionApplier{
				metric.WithAbsolute(false),
				metric.WithAbsolute(true),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing measure case %s (%d)", tt.name, idx)
		opts := &metric.Options{}
		metric.ApplyMeasureOptions(opts, tt.opts...)
		checkOptions(t, opts, &metric.Options{
			Description: tt.desc,
			Unit:        tt.unit,
			Keys:        tt.keys,
			Alternate:   tt.alt,
		})
	}
}

func checkOptions(t *testing.T, got *metric.Options, expected *metric.Options) {
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Errorf("Compare options: -got +want %s", diff)
	}
}

func TestCounter(t *testing.T) {
	{
		meter := mock.NewMeter()
		c := meter.NewFloat64Counter("ajwaj")
		ctx := context.Background()
		labels := meter.Labels()
		c.Add(ctx, 42, labels)
		handle := c.AcquireHandle(labels)
		handle.Add(ctx, 42)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		t.Log("Testing float counter")
		checkBatches(t, ctx, labels, meter, core.Float64NumberKind, c.Impl())
	}
	{
		meter := mock.NewMeter()
		c := meter.NewInt64Counter("ajwaj")
		ctx := context.Background()
		labels := meter.Labels()
		c.Add(ctx, 42, labels)
		handle := c.AcquireHandle(labels)
		handle.Add(ctx, 42)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		t.Log("Testing int counter")
		checkBatches(t, ctx, labels, meter, core.Int64NumberKind, c.Impl())
	}
}

func TestGauge(t *testing.T) {
	{
		meter := mock.NewMeter()
		g := meter.NewFloat64Gauge("ajwaj")
		ctx := context.Background()
		labels := meter.Labels()
		g.Set(ctx, 42, labels)
		handle := g.AcquireHandle(labels)
		handle.Set(ctx, 42)
		meter.RecordBatch(ctx, labels, g.Measurement(42))
		t.Log("Testing float gauge")
		checkBatches(t, ctx, labels, meter, core.Float64NumberKind, g.Impl())
	}
	{
		meter := mock.NewMeter()
		g := meter.NewInt64Gauge("ajwaj")
		ctx := context.Background()
		labels := meter.Labels()
		g.Set(ctx, 42, labels)
		handle := g.AcquireHandle(labels)
		handle.Set(ctx, 42)
		meter.RecordBatch(ctx, labels, g.Measurement(42))
		t.Log("Testing int gauge")
		checkBatches(t, ctx, labels, meter, core.Int64NumberKind, g.Impl())
	}
}

func TestMeasure(t *testing.T) {
	{
		meter := mock.NewMeter()
		m := meter.NewFloat64Measure("ajwaj")
		ctx := context.Background()
		labels := meter.Labels()
		m.Record(ctx, 42, labels)
		handle := m.AcquireHandle(labels)
		handle.Record(ctx, 42)
		meter.RecordBatch(ctx, labels, m.Measurement(42))
		t.Log("Testing float measure")
		checkBatches(t, ctx, labels, meter, core.Float64NumberKind, m.Impl())
	}
	{
		meter := mock.NewMeter()
		m := meter.NewInt64Measure("ajwaj")
		ctx := context.Background()
		labels := meter.Labels()
		m.Record(ctx, 42, labels)
		handle := m.AcquireHandle(labels)
		handle.Record(ctx, 42)
		meter.RecordBatch(ctx, labels, m.Measurement(42))
		t.Log("Testing int measure")
		checkBatches(t, ctx, labels, meter, core.Int64NumberKind, m.Impl())
	}
}

func checkBatches(t *testing.T, ctx context.Context, labels metric.LabelSet, meter *mock.Meter, kind core.NumberKind, instrument metric.InstrumentImpl) {
	t.Helper()
	if len(meter.MeasurementBatches) != 3 {
		t.Errorf("Expected 3 recorded measurement batches, got %d", len(meter.MeasurementBatches))
	}
	ourInstrument := instrument.(*mock.Instrument)
	ourLabelSet := labels.(*mock.LabelSet)
	minLen := 3
	if minLen > len(meter.MeasurementBatches) {
		minLen = len(meter.MeasurementBatches)
	}
	for i := 0; i < minLen; i++ {
		got := meter.MeasurementBatches[i]
		if got.Ctx != ctx {
			d := func(c context.Context) string {
				return fmt.Sprintf("(ptr: %p, ctx %#v)", c, c)
			}
			t.Errorf("Wrong recorded context in batch %d, expected %s, got %s", i, d(ctx), d(got.Ctx))
		}
		if got.LabelSet != ourLabelSet {
			d := func(l *mock.LabelSet) string {
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
			if measurement.Instrument != ourInstrument {
				d := func(i *mock.Instrument) string {
					return fmt.Sprintf("(ptr: %p, instrument %#v)", i, i)
				}
				t.Errorf("Wrong recorded instrument in measurement %d in batch %d, expected %s, got %s", j, i, d(ourInstrument), d(measurement.Instrument))
			}
			ft := fortyTwo(t, kind)
			if measurement.Number.CompareNumber(kind, ft) != 0 {
				t.Errorf("Wrong recorded value in measurement %d in batch %d, expected %s, got %s", j, i, ft.Emit(kind), measurement.Number.Emit(kind))
			}
		}
	}
}

func fortyTwo(t *testing.T, kind core.NumberKind) core.Number {
	switch kind {
	case core.Int64NumberKind:
		return core.NewInt64Number(42)
	case core.Float64NumberKind:
		return core.NewFloat64Number(42)
	}
	t.Errorf("Invalid value kind %q", kind)
	return core.NewInt64Number(0)
}
