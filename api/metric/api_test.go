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
				WithNonMonotonic(true),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  true,
		},
		{
			name: "nonmonotonic, but not really",
			opts: []CounterOptionApplier{
				WithNonMonotonic(true),
				WithNonMonotonic(false),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
	}
	checkCounterDescriptor := func(tt testcase, vk ValueKind, d *Descriptor) {
		e := descriptor{
			name: tt.name,
			keys: tt.keys,
			desc: tt.desc,
			unit: tt.unit,
			alt:  tt.alt,
			kind: CounterKind,
			vk:   vk,
		}
		checkDescriptor(t, e, d)
	}
	for idx, tt := range testcases {
		t.Logf("Testing counter case %s (%d)", tt.name, idx)
		f := NewFloat64Counter(tt.name, tt.opts...)
		checkCounterDescriptor(tt, Float64ValueKind, f.Descriptor())
		i := NewInt64Counter(tt.name, tt.opts...)
		checkCounterDescriptor(tt, Int64ValueKind, i.Descriptor())
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
	checkGaugeDescriptor := func(tt testcase, vk ValueKind, d *Descriptor) {
		e := descriptor{
			name: tt.name,
			keys: tt.keys,
			desc: tt.desc,
			unit: tt.unit,
			alt:  tt.alt,
			kind: GaugeKind,
			vk:   vk,
		}
		checkDescriptor(t, e, d)
	}
	for idx, tt := range testcases {
		t.Logf("Testing gauge case %s (%d)", tt.name, idx)
		f := NewFloat64Gauge(tt.name, tt.opts...)
		checkGaugeDescriptor(tt, Float64ValueKind, f.Descriptor())
		i := NewInt64Gauge(tt.name, tt.opts...)
		checkGaugeDescriptor(tt, Int64ValueKind, i.Descriptor())
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
			name: "signed",
			opts: []MeasureOptionApplier{
				WithSigned(true),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  true,
		},
		{
			name: "signed, but not really",
			opts: []MeasureOptionApplier{
				WithSigned(true),
				WithSigned(false),
			},
			keys: nil,
			desc: "",
			unit: "",
			alt:  false,
		},
	}
	checkMeasureDescriptor := func(tt testcase, vk ValueKind, d *Descriptor) {
		e := descriptor{
			name: tt.name,
			keys: tt.keys,
			desc: tt.desc,
			unit: tt.unit,
			alt:  tt.alt,
			kind: MeasureKind,
			vk:   vk,
		}
		checkDescriptor(t, e, d)
	}
	for idx, tt := range testcases {
		t.Logf("Testing measure case %s (%d)", tt.name, idx)
		f := NewFloat64Measure(tt.name, tt.opts...)
		checkMeasureDescriptor(tt, Float64ValueKind, f.Descriptor())
		i := NewInt64Measure(tt.name, tt.opts...)
		checkMeasureDescriptor(tt, Int64ValueKind, i.Descriptor())
	}
}

type descriptor struct {
	name string
	keys []core.Key
	desc string
	unit unit.Unit
	alt  bool
	kind Kind
	vk   ValueKind
}

func checkDescriptor(t *testing.T, e descriptor, d *Descriptor) {
	if e.name != d.Name() {
		t.Errorf("Expected name %q, got %q", e.name, d.Name())
	}
	if len(e.keys) != len(d.Keys()) {
		t.Errorf("Expected %d key(s), got %d", len(e.keys), len(d.Keys()))
	}
	minLen := len(e.keys)
	if minLen > len(d.Keys()) {
		minLen = len(d.Keys())
	}
	for i := 0; i < minLen; i++ {
		if e.keys[i] != d.Keys()[i] {
			t.Errorf("Expected key %q, got %q", e.keys[i], d.Keys()[i])
		}
	}
	if e.desc != d.Description() {
		t.Errorf("Expected description %q, got %q", e.desc, d.Description())
	}
	if e.unit != d.Unit() {
		t.Errorf("Expected unit %q, got %q", e.unit, d.Unit())
	}
	if e.alt != d.Alternate() {
		t.Errorf("Expected alternate %v, got %v", e.alt, d.Alternate())
	}
	if e.vk != d.ValueKind() {
		t.Errorf("Expected value kind %q, got %q", e.vk, d.ValueKind())
	}
	if e.kind != d.Kind() {
		t.Errorf("Expected kind %q, got %q", e.kind, d.Kind())
	}
}

func TestCounter(t *testing.T) {
	{
		c := NewFloat64Counter("ajwaj")
		meter := newMockMeter()
		ctx := context.Background()
		labels := meter.DefineLabels(ctx)
		c.Add(ctx, 42, labels)
		handle := c.GetHandle(labels)
		handle.Add(ctx, 42)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		t.Log("Testing float counter")
		checkBatches(t, ctx, labels, meter, c.Descriptor())
	}
	{
		c := NewInt64Counter("ajwaj")
		meter := newMockMeter()
		ctx := context.Background()
		labels := meter.DefineLabels(ctx)
		c.Add(ctx, 42, labels)
		handle := c.GetHandle(labels)
		handle.Add(ctx, 42)
		meter.RecordBatch(ctx, labels, c.Measurement(42))
		t.Log("Testing int counter")
		checkBatches(t, ctx, labels, meter, c.Descriptor())
	}
}

func TestGauge(t *testing.T) {
	{
		g := NewFloat64Gauge("ajwaj")
		meter := newMockMeter()
		ctx := context.Background()
		labels := meter.DefineLabels(ctx)
		g.Set(ctx, 42, labels)
		handle := g.GetHandle(labels)
		handle.Set(ctx, 42)
		meter.RecordBatch(ctx, labels, g.Measurement(42))
		t.Log("Testing float gauge")
		checkBatches(t, ctx, labels, meter, g.Descriptor())
	}
	{
		g := NewInt64Gauge("ajwaj")
		meter := newMockMeter()
		ctx := context.Background()
		labels := meter.DefineLabels(ctx)
		g.Set(ctx, 42, labels)
		handle := g.GetHandle(labels)
		handle.Set(ctx, 42)
		meter.RecordBatch(ctx, labels, g.Measurement(42))
		t.Log("Testing int gauge")
		checkBatches(t, ctx, labels, meter, g.Descriptor())
	}
}

func TestMeasure(t *testing.T) {
	{
		m := NewFloat64Measure("ajwaj")
		meter := newMockMeter()
		ctx := context.Background()
		labels := meter.DefineLabels(ctx)
		m.Record(ctx, 42, labels)
		handle := m.GetHandle(labels)
		handle.Record(ctx, 42)
		meter.RecordBatch(ctx, labels, m.Measurement(42))
		t.Log("Testing float measure")
		checkBatches(t, ctx, labels, meter, m.Descriptor())
	}
	{
		m := NewInt64Measure("ajwaj")
		meter := newMockMeter()
		ctx := context.Background()
		labels := meter.DefineLabels(ctx)
		m.Record(ctx, 42, labels)
		handle := m.GetHandle(labels)
		handle.Record(ctx, 42)
		meter.RecordBatch(ctx, labels, m.Measurement(42))
		t.Log("Testing int measure")
		checkBatches(t, ctx, labels, meter, m.Descriptor())
	}
}

func checkBatches(t *testing.T, ctx context.Context, labels LabelSet, meter *mockMeter, descriptor *Descriptor) {
	t.Helper()
	if len(meter.measurementBatches) != 3 {
		t.Errorf("Expected 3 recorded measurement batches, got %d", len(meter.measurementBatches))
	}
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
			if measurement.Descriptor != descriptor {
				d := func(d *Descriptor) string {
					return fmt.Sprintf("(ptr: %p, descriptor %#v)", d, d)
				}
				t.Errorf("Wrong recorded descriptor in measurement %d in batch %d, expected %s, got %s", j, i, d(descriptor), d(measurement.Descriptor))
			}
			ft := fortyTwo(t, descriptor.ValueKind())
			if measurement.Value.RawCompare(ft.AsRaw(), descriptor.ValueKind()) != 0 {
				t.Errorf("Wrong recorded value in measurement %d in batch %d, expected %s, got %s", j, i, ft.Emit(descriptor.ValueKind()), measurement.Value.Emit(descriptor.ValueKind()))
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
