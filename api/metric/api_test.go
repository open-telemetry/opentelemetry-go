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

package metric

import (
	"context"
	"fmt"
	"testing"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/api/unit"

	"github.com/google/go-cmp/cmp"
)

func TestCounterOptions(t *testing.T) {
	type testcase struct {
		name string
		opts []CounterOptionApplier
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
			opts: []CounterOptionApplier{
				WithKeys(key.New("foo"), key.New("foo2")),
				WithKeys(key.New("bar"), key.New("bar2")),
				WithKeys(key.New("baz"), key.New("baz2")),
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
			opts: []CounterOptionApplier{
				WithDescription("stuff"),
			},
			keys: nil,
			desc: "stuff",
			unit: "",
			alt:  false,
		},
		{
			name: "description override",
			opts: []CounterOptionApplier{
				WithDescription("stuff"),
				WithDescription("things"),
			},
			keys: nil,
			desc: "things",
			unit: "",
			alt:  false,
		},
		{
			name: "unit",
			opts: []CounterOptionApplier{
				WithUnit("s"),
			},
			keys: nil,
			desc: "",
			unit: "s",
			alt:  false,
		},
		{
			name: "unit override",
			opts: []CounterOptionApplier{
				WithUnit("s"),
				WithUnit("h"),
			},
			keys: nil,
			desc: "",
			unit: "h",
			alt:  false,
		},
		{
			name: "nonmonotonic",
			opts: []CounterOptionApplier{
				WithMonotonic(false),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  true,
		},
		{
			name: "nonmonotonic, but not really",
			opts: []CounterOptionApplier{
				WithMonotonic(false),
				WithMonotonic(true),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing counter case %s (%d)", tt.name, idx)
		opts := &Options{}
		ApplyCounterOptions(opts, tt.opts...)
		checkOptions(t, opts, &Options{
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
		opts []GaugeOptionApplier
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
			opts: []GaugeOptionApplier{
				WithKeys(key.New("foo"), key.New("foo2")),
				WithKeys(key.New("bar"), key.New("bar2")),
				WithKeys(key.New("baz"), key.New("baz2")),
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
			opts: []GaugeOptionApplier{
				WithDescription("stuff"),
			},
			keys: nil,
			desc: "stuff",
			unit: "",
			alt:  false,
		},
		{
			name: "description override",
			opts: []GaugeOptionApplier{
				WithDescription("stuff"),
				WithDescription("things"),
			},
			keys: nil,
			desc: "things",
			unit: "",
			alt:  false,
		},
		{
			name: "unit",
			opts: []GaugeOptionApplier{
				WithUnit("s"),
			},
			keys: nil,
			desc: "",
			unit: "s",
			alt:  false,
		},
		{
			name: "unit override",
			opts: []GaugeOptionApplier{
				WithUnit("s"),
				WithUnit("h"),
			},
			keys: nil,
			desc: "",
			unit: "h",
			alt:  false,
		},
		{
			name: "monotonic",
			opts: []GaugeOptionApplier{
				WithMonotonic(true),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  true,
		},
		{
			name: "monotonic, but not really",
			opts: []GaugeOptionApplier{
				WithMonotonic(true),
				WithMonotonic(false),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing gauge case %s (%d)", tt.name, idx)
		opts := &Options{}
		ApplyGaugeOptions(opts, tt.opts...)
		checkOptions(t, opts, &Options{
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
		opts []MeasureOptionApplier
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
			opts: []MeasureOptionApplier{
				WithKeys(key.New("foo"), key.New("foo2")),
				WithKeys(key.New("bar"), key.New("bar2")),
				WithKeys(key.New("baz"), key.New("baz2")),
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
			opts: []MeasureOptionApplier{
				WithDescription("stuff"),
			},
			keys: nil,
			desc: "stuff",
			unit: "",
			alt:  false,
		},
		{
			name: "description override",
			opts: []MeasureOptionApplier{
				WithDescription("stuff"),
				WithDescription("things"),
			},
			keys: nil,
			desc: "things",
			unit: "",
			alt:  false,
		},
		{
			name: "unit",
			opts: []MeasureOptionApplier{
				WithUnit("s"),
			},
			keys: nil,
			desc: "",
			unit: "s",
			alt:  false,
		},
		{
			name: "unit override",
			opts: []MeasureOptionApplier{
				WithUnit("s"),
				WithUnit("h"),
			},
			keys: nil,
			desc: "",
			unit: "h",
			alt:  false,
		},
		{
			name: "not absolute",
			opts: []MeasureOptionApplier{
				WithAbsolute(false),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  true,
		},
		{
			name: "not absolute, but not really",
			opts: []MeasureOptionApplier{
				WithAbsolute(false),
				WithAbsolute(true),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
	}
	for idx, tt := range testcases {
		t.Logf("Testing measure case %s (%d)", tt.name, idx)
		opts := &Options{}
		ApplyMeasureOptions(opts, tt.opts...)
		checkOptions(t, opts, &Options{
			Description: tt.desc,
			Unit:        tt.unit,
			Keys:        tt.keys,
			Alternate:   tt.alt,
		})
	}
}

func checkOptions(t *testing.T, got *Options, expected *Options) {
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Errorf("Compare options: -got +want %s", diff)
	}
}

func TestCounter(t *testing.T) {
	{
		meter := newMockMeter()
		c := meter.NewFloat64Counter("ajwaj")
		ctx := context.Background()
		labels := meter.Labels(ctx)
		c.Add(ctx, 42, labels)
		handle := c.AcquireHandle(labels)
		handle.Add(ctx, 42)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		t.Log("Testing float counter")
		checkBatches(t, ctx, labels, meter, Float64ValueKind, c.instrument)
	}
	{
		meter := newMockMeter()
		c := meter.NewInt64Counter("ajwaj")
		ctx := context.Background()
		labels := meter.Labels(ctx)
		c.Add(ctx, 42, labels)
		handle := c.AcquireHandle(labels)
		handle.Add(ctx, 42)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		t.Log("Testing int counter")
		checkBatches(t, ctx, labels, meter, Int64ValueKind, c.instrument)
	}
}

func TestGauge(t *testing.T) {
	{
		meter := newMockMeter()
		g := meter.NewFloat64Gauge("ajwaj")
		ctx := context.Background()
		labels := meter.Labels(ctx)
		g.Set(ctx, 42, labels)
		handle := g.AcquireHandle(labels)
		handle.Set(ctx, 42)
		meter.RecordBatch(ctx, labels, g.Measurement(42))
		t.Log("Testing float gauge")
		checkBatches(t, ctx, labels, meter, Float64ValueKind, g.instrument)
	}
	{
		meter := newMockMeter()
		g := meter.NewInt64Gauge("ajwaj")
		ctx := context.Background()
		labels := meter.Labels(ctx)
		g.Set(ctx, 42, labels)
		handle := g.AcquireHandle(labels)
		handle.Set(ctx, 42)
		meter.RecordBatch(ctx, labels, g.Measurement(42))
		t.Log("Testing int gauge")
		checkBatches(t, ctx, labels, meter, Int64ValueKind, g.instrument)
	}
}

func TestMeasure(t *testing.T) {
	{
		meter := newMockMeter()
		m := meter.NewFloat64Measure("ajwaj")
		ctx := context.Background()
		labels := meter.Labels(ctx)
		m.Record(ctx, 42, labels)
		handle := m.AcquireHandle(labels)
		handle.Record(ctx, 42)
		meter.RecordBatch(ctx, labels, m.Measurement(42))
		t.Log("Testing float measure")
		checkBatches(t, ctx, labels, meter, Float64ValueKind, m.instrument)
	}
	{
		meter := newMockMeter()
		m := meter.NewInt64Measure("ajwaj")
		ctx := context.Background()
		labels := meter.Labels(ctx)
		m.Record(ctx, 42, labels)
		handle := m.AcquireHandle(labels)
		handle.Record(ctx, 42)
		meter.RecordBatch(ctx, labels, m.Measurement(42))
		t.Log("Testing int measure")
		checkBatches(t, ctx, labels, meter, Int64ValueKind, m.instrument)
	}
}

func checkBatches(t *testing.T, ctx context.Context, labels LabelSet, meter *mockMeter, kind ValueKind, instrument Instrument) {
	t.Helper()
	if len(meter.measurementBatches) != 3 {
		t.Errorf("Expected 3 recorded measurement batches, got %d", len(meter.measurementBatches))
	}
	ourInstrument := instrument.(*mockInstrument)
	ourLabelSet := labels.(*mockLabelSet)
	minLen := 3
	if minLen > len(meter.measurementBatches) {
		minLen = len(meter.measurementBatches)
	}
	for i := 0; i < minLen; i++ {
		got := meter.measurementBatches[i]
		if got.ctx != ctx {
			d := func(c context.Context) string {
				return fmt.Sprintf("(ptr: %p, ctx %#v)", c, c)
			}
			t.Errorf("Wrong recorded context in batch %d, expected %s, got %s", i, d(ctx), d(got.ctx))
		}
		if got.labelSet != ourLabelSet {
			d := func(l *mockLabelSet) string {
				return fmt.Sprintf("(ptr: %p, labels %#v)", l, l.labels)
			}
			t.Errorf("Wrong recorded label set in batch %d, expected %s, got %s", i, d(ourLabelSet), d(got.labelSet))
		}
		if len(got.measurements) != 1 {
			t.Errorf("Expected 1 measurement in batch %d, got %d", i, len(got.measurements))
		}
		minMLen := 1
		if minMLen > len(got.measurements) {
			minMLen = len(got.measurements)
		}
		for j := 0; j < minMLen; j++ {
			measurement := got.measurements[j]
			if measurement.instrument != ourInstrument {
				d := func(i *mockInstrument) string {
					return fmt.Sprintf("(ptr: %p, instrument %#v)", i, i)
				}
				t.Errorf("Wrong recorded instrument in measurement %d in batch %d, expected %s, got %s", j, i, d(ourInstrument), d(measurement.instrument))
			}
			ft := fortyTwo(t, kind)
			if measurement.value.RawCompare(ft.AsRaw(), kind) != 0 {
				t.Errorf("Wrong recorded value in measurement %d in batch %d, expected %s, got %s", j, i, ft.Emit(kind), measurement.value.Emit(kind))
			}
		}
	}
}

func fortyTwo(t *testing.T, kind ValueKind) MeasurementValue {
	switch kind {
	case Int64ValueKind:
		return NewInt64MeasurementValue(42)
	case Float64ValueKind:
		return NewFloat64MeasurementValue(42)
	}
	t.Errorf("Invalid value kind %q", kind)
	return NewInt64MeasurementValue(0)
}
