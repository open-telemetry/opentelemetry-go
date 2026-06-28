// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

func valueEqual(a, b attribute.Value) bool {
	return reflect.DeepEqual(a, b)
}

func keyValueEqual(a, b attribute.KeyValue) bool {
	return a.Key == b.Key && valueEqual(a.Value, b.Value)
}

func TestRecordEventName(t *testing.T) {
	const text = "testing text"

	r := new(Record)
	r.SetEventName(text)
	assert.Equal(t, text, r.EventName())
}

func TestRecordTimestamp(t *testing.T) {
	now := time.Now()
	r := new(Record)
	r.SetTimestamp(now)
	assert.Equal(t, now, r.Timestamp())
}

func TestRecordObservedTimestamp(t *testing.T) {
	now := time.Now()
	r := new(Record)
	r.SetObservedTimestamp(now)
	assert.Equal(t, now, r.ObservedTimestamp())
}

func TestRecordSeverity(t *testing.T) {
	s := log.SeverityInfo
	r := new(Record)
	r.SetSeverity(s)
	assert.Equal(t, s, r.Severity())
}

func TestRecordSeverityText(t *testing.T) {
	text := "text"
	r := new(Record)
	r.SetSeverityText(text)
	assert.Equal(t, text, r.SeverityText())
}

func TestRecordBody(t *testing.T) {
	testcases := []struct {
		name            string
		allowDuplicates bool
		body            attribute.Value
		want            attribute.Value
	}{
		{
			name: "boolean value",
			body: attribute.BoolValue(true),
			want: attribute.BoolValue(true),
		},
		{
			name: "slice",
			body: attribute.SliceValue(attribute.BoolValue(true), attribute.BoolValue(false)),
			want: attribute.SliceValue(attribute.BoolValue(true), attribute.BoolValue(false)),
		},
		{
			name: "map",
			body: attribute.MapValue(
				attribute.Bool("0", true),
				attribute.Int64("1", 2),
				attribute.Float64("2", 3.0),
				attribute.String("3", "forth"),
				attribute.Slice("4", attribute.Int64Value(1)),
				attribute.Map("5", attribute.Int("key", 2)),
				attribute.ByteSlice("6", []byte("six")),
				attribute.Int64("1", 3),
			),
			want: attribute.MapValue(
				attribute.Bool("0", true),
				attribute.Float64("2", 3.0),
				attribute.String("3", "forth"),
				attribute.Slice("4", attribute.Int64Value(1)),
				attribute.Map("5", attribute.Int("key", 2)),
				attribute.ByteSlice("6", []byte("six")),
				attribute.Int64("1", 3),
			),
		},
		{
			name: "nested map",
			body: attribute.MapValue(
				attribute.Map(
					"key",
					attribute.Int64("key", 1),
					attribute.Int64("key", 2),
				),
			),

			want: attribute.MapValue(
				attribute.Map(
					"key",
					attribute.Int64("key", 2),
				),
			),
		},
		{
			name:            "map - allow duplicates",
			allowDuplicates: true,
			body: attribute.MapValue(
				attribute.Int64("1", 2),
				attribute.Int64("1", 3),
			),
			want: attribute.MapValue(
				attribute.Int64("1", 2),
				attribute.Int64("1", 3),
			),
		},
		{
			name: "slice with nested deduplication",
			body: attribute.SliceValue(
				attribute.MapValue(attribute.String("key", "value1"), attribute.String("key", "value2")),
				attribute.StringValue("normal"),
				attribute.SliceValue(
					attribute.MapValue(attribute.String("nested", "val1"), attribute.String("nested", "val2")),
				),
			),
			want: attribute.SliceValue(
				attribute.MapValue(attribute.String("key", "value2")),
				attribute.StringValue("normal"),
				attribute.SliceValue(
					attribute.MapValue(attribute.String("nested", "val2")),
				),
			),
		},
		{
			name: "empty slice",
			body: attribute.SliceValue(),
			want: attribute.SliceValue(),
		},
		{
			name: "empty map",
			body: attribute.MapValue(),
			want: attribute.MapValue(),
		},
		{
			name: "single key map",
			body: attribute.MapValue(attribute.String("single", "value")),
			want: attribute.MapValue(attribute.String("single", "value")),
		},
		{
			name: "slice with no deduplication needed",
			body: attribute.SliceValue(
				attribute.StringValue("value1"),
				attribute.StringValue("value2"),
				attribute.MapValue(attribute.String("unique1", "val1")),
				attribute.MapValue(attribute.String("unique2", "val2")),
			),
			want: attribute.SliceValue(
				attribute.StringValue("value1"),
				attribute.StringValue("value2"),
				attribute.MapValue(attribute.String("unique1", "val1")),
				attribute.MapValue(attribute.String("unique2", "val2")),
			),
		},
		{
			name: "deeply nested slice and map structure",
			body: attribute.SliceValue(
				attribute.MapValue(
					attribute.String("outer", "value"),
					attribute.Slice(
						"inner_slice",
						attribute.MapValue(attribute.String("deep", "value1"), attribute.String("deep", "value2")),
					),
				),
			),

			want: attribute.SliceValue(
				attribute.MapValue(
					attribute.String("outer", "value"),
					attribute.Slice(
						"inner_slice",
						attribute.MapValue(attribute.String("deep", "value2")),
					),
				),
			),
		},
		{
			name:            "slice with duplicates allowed",
			allowDuplicates: true,
			body: attribute.SliceValue(
				attribute.MapValue(attribute.String("key", "value1"), attribute.String("key", "value2")),
			),

			want: attribute.SliceValue(
				attribute.MapValue(attribute.String("key", "value1"), attribute.String("key", "value2")),
			),
		},
		{
			name: "string value",
			body: attribute.StringValue("test"),
			want: attribute.StringValue("test"),
		},
		{
			name: "boolean value without deduplication",
			body: attribute.BoolValue(true),
			want: attribute.BoolValue(true),
		},
		{
			name: "integer value",
			body: attribute.Int64Value(42),
			want: attribute.Int64Value(42),
		},
		{
			name: "float value",
			body: attribute.Float64Value(3.14),
			want: attribute.Float64Value(3.14),
		},
		{
			name: "bytes value",
			body: attribute.ByteSliceValue([]byte("test")),
			want: attribute.ByteSliceValue([]byte("test")),
		},
		{
			name: "empty slice",
			body: attribute.SliceValue(),
			want: attribute.SliceValue(),
		},
		{
			name: "slice without nested deduplication",
			body: attribute.SliceValue(attribute.StringValue("test"), attribute.BoolValue(true)),
			want: attribute.SliceValue(attribute.StringValue("test"), attribute.BoolValue(true)),
		},
		{
			name: "slice with nested deduplication needed",
			body: attribute.SliceValue(
				attribute.MapValue(attribute.String("key", "value1"), attribute.String("key", "value2")),
			),
			want: attribute.SliceValue(attribute.MapValue(attribute.String("key", "value2"))),
		},
		{
			name: "empty map",
			body: attribute.MapValue(),
			want: attribute.MapValue(),
		},
		{
			name: "single key map",
			body: attribute.MapValue(attribute.String("key", "value")),
			want: attribute.MapValue(attribute.String("key", "value")),
		},
		{
			name: "map with duplicate keys",
			body: attribute.MapValue(attribute.String("key", "value1"), attribute.String("key", "value2")),
			want: attribute.MapValue(attribute.String("key", "value2")),
		},
		{
			name: "map without duplicates",
			body: attribute.MapValue(attribute.String("key1", "value1"), attribute.String("key2", "value2")),
			want: attribute.MapValue(attribute.String("key1", "value1"), attribute.String("key2", "value2")),
		},
		{
			name: "map with nested slice deduplication",
			body: attribute.MapValue(
				attribute.Slice(
					"slice",
					attribute.MapValue(attribute.String("nested", "val1"), attribute.String("nested", "val2")),
				),
			),

			want: attribute.MapValue(
				attribute.Slice("slice", attribute.MapValue(attribute.String("nested", "val2"))),
			),
		},
		{
			name: "deeply nested structure with deduplication",
			body: attribute.SliceValue(
				attribute.MapValue(
					attribute.Map(
						"nested",
						attribute.String("key", "value1"),
						attribute.String("key", "value2"),
					),
				),
			),

			want: attribute.SliceValue(
				attribute.MapValue(
					attribute.Map(
						"nested",
						attribute.String("key", "value2"),
					),
				),
			),
		},
		{
			name: "deeply nested structure without deduplication",
			body: attribute.SliceValue(
				attribute.MapValue(
					attribute.Map(
						"nested",
						attribute.String("key1", "value1"),
						attribute.String("key2", "value2"),
					),
				),
			),

			want: attribute.SliceValue(
				attribute.MapValue(
					attribute.Map(
						"nested",
						attribute.String("key1", "value1"),
						attribute.String("key2", "value2"),
					),
				),
			),
		},
		{
			name: "string value for collection deduplication",
			body: attribute.StringValue("test"),
			want: attribute.StringValue("test"),
		},
		{
			name: "boolean value for collection deduplication",
			body: attribute.BoolValue(true),
			want: attribute.BoolValue(true),
		},
		{
			name: "empty slice for collection deduplication",
			body: attribute.SliceValue(),
			want: attribute.SliceValue(),
		},
		{
			name: "slice without nested deduplication for collection testing",
			body: attribute.SliceValue(attribute.StringValue("test"), attribute.BoolValue(true)),
			want: attribute.SliceValue(attribute.StringValue("test"), attribute.BoolValue(true)),
		},
		{
			name: "slice with nested map requiring deduplication",
			body: attribute.SliceValue(
				attribute.MapValue(attribute.String("key", "value1"), attribute.String("key", "value2")),
				attribute.StringValue("normal"),
			),
			want: attribute.SliceValue(
				attribute.MapValue(attribute.String("key", "value2")),
				attribute.StringValue("normal"),
			),
		},
		{
			name: "deeply nested slice with map deduplication",
			body: attribute.SliceValue(
				attribute.SliceValue(
					attribute.MapValue(attribute.String("deep", "value1"), attribute.String("deep", "value2")),
				),
			),

			want: attribute.SliceValue(
				attribute.SliceValue(
					attribute.MapValue(attribute.String("deep", "value2")),
				),
			),
		},
		{
			name: "empty map for collection deduplication",
			body: attribute.MapValue(),
			want: attribute.MapValue(),
		},
		{
			name: "map with nested slice containing duplicates",
			body: attribute.MapValue(
				attribute.String("outer", "value"),
				attribute.Slice(
					"nested_slice",
					attribute.MapValue(attribute.String("inner", "val1"), attribute.String("inner", "val2")),
				),
			),
			want: attribute.MapValue(
				attribute.String("outer", "value"),
				attribute.Slice(
					"nested_slice",
					attribute.MapValue(attribute.String("inner", "val2")),
				),
			),
		},
		{
			name: "map with key duplication and nested value deduplication",
			body: attribute.MapValue(
				attribute.String("key1", "value1"),
				attribute.String("key1", "value2"),
				attribute.Slice(
					"slice",
					attribute.MapValue(attribute.String("nested", "val1"), attribute.String("nested", "val2")),
				),
			),
			want: attribute.MapValue(
				attribute.String("key1", "value2"),
				attribute.Slice(
					"slice",
					attribute.MapValue(attribute.String("nested", "val2")),
				),
			),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := new(Record)
			r.allowDupKeys = tc.allowDuplicates
			r.SetBody(tc.body)
			got := r.Body()
			if !valueEqual(got, tc.want) {
				t.Errorf("r.Body() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRecordAttributes(t *testing.T) {
	attrs := []attribute.KeyValue{
		attribute.Bool("0", true),
		attribute.Int64("1", 2),
		attribute.Float64("2", 3.0),
		attribute.String("3", "forth"),
		attribute.Slice("4", attribute.Int64Value(1)),
		attribute.Map("5", attribute.Int("key", 2)),
		attribute.ByteSlice("6", []byte("six")),
	}
	r := new(Record)
	r.attributeValueLengthLimit = -1
	r.SetAttributes(attrs...)
	r.SetAttributes(attrs[:2]...) // Overwrite existing.
	r.AddAttributes(attrs[2:]...)

	assert.Equal(t, len(attrs), r.AttributesLen(), "attribute length")

	for n := range attrs {
		var i int
		r.WalkAttributes(func(attribute.KeyValue) bool {
			i++
			return i <= n
		})
		assert.Equalf(t, n+1, i, "WalkAttributes did not stop at %d", n+1)
	}

	var i int
	r.WalkAttributes(func(kv attribute.KeyValue) bool {
		assert.Truef(t, keyValueEqual(kv, attrs[i]), "%d: %v != %v", i, kv, attrs[i])
		i++
		return true
	})
}

func TestRecordTraceID(t *testing.T) {
	id := trace.TraceID([16]byte{1})
	r := new(Record)
	r.SetTraceID(id)
	assert.Equal(t, id, r.TraceID())
}

func TestRecordSpanID(t *testing.T) {
	id := trace.SpanID([8]byte{1})
	r := new(Record)
	r.SetSpanID(id)
	assert.Equal(t, id, r.SpanID())
}

func TestRecordTraceFlags(t *testing.T) {
	flag := trace.FlagsSampled
	r := new(Record)
	r.SetTraceFlags(flag)
	assert.Equal(t, flag, r.TraceFlags())
}

func TestRecordResource(t *testing.T) {
	r := new(Record)
	assert.NotPanics(t, func() { r.Resource() })

	res := resource.NewSchemaless(attribute.Bool("key", true))
	r.resource = res
	got := r.Resource()
	assert.Equal(t, res, got)
}

func TestRecordInstrumentationScope(t *testing.T) {
	r := new(Record)
	assert.NotPanics(t, func() { r.InstrumentationScope() })

	scope := instrumentation.Scope{Name: "testing"}
	r.scope = &scope
	assert.Equal(t, scope, r.InstrumentationScope())
}

func TestRecordClone(t *testing.T) {
	now0 := time.Now()
	sev0 := log.SeverityInfo
	text0 := "text"
	val0 := attribute.BoolValue(true)
	attr0 := attribute.Bool("0", true)
	traceID0 := trace.TraceID([16]byte{1})
	spanID0 := trace.SpanID([8]byte{1})
	flag0 := trace.FlagsSampled

	r0 := new(Record)
	r0.SetTimestamp(now0)
	r0.SetObservedTimestamp(now0)
	r0.SetSeverity(sev0)
	r0.SetSeverityText(text0)
	r0.SetBody(val0)
	r0.SetAttributes(attr0)
	r0.SetTraceID(traceID0)
	r0.SetSpanID(spanID0)
	r0.SetTraceFlags(flag0)

	now1 := now0.Add(time.Second)
	sev1 := log.SeverityDebug
	text1 := "string"
	val1 := attribute.IntValue(1)
	attr1 := attribute.Int64("1", 2)
	traceID1 := trace.TraceID([16]byte{2})
	spanID1 := trace.SpanID([8]byte{2})
	flag1 := trace.TraceFlags(2)

	r1 := r0.Clone()
	r1.SetTimestamp(now1)
	r1.SetObservedTimestamp(now1)
	r1.SetSeverity(sev1)
	r1.SetSeverityText(text1)
	r1.SetBody(val1)
	r1.SetAttributes(attr1)
	r1.SetTraceID(traceID1)
	r1.SetSpanID(spanID1)
	r1.SetTraceFlags(flag1)

	assert.Equal(t, now0, r0.Timestamp())
	assert.Equal(t, now0, r0.ObservedTimestamp())
	assert.Equal(t, sev0, r0.Severity())
	assert.Equal(t, text0, r0.SeverityText())
	assert.Equal(t, val0, r0.Body())
	assert.Equal(t, traceID0, r0.TraceID())
	assert.Equal(t, spanID0, r0.SpanID())
	assert.Equal(t, flag0, r0.TraceFlags())
	r0.WalkAttributes(func(kv attribute.KeyValue) bool {
		return assert.Truef(t, keyValueEqual(kv, attr0), "%v != %v", kv, attr0)
	})

	assert.Equal(t, now1, r1.Timestamp())
	assert.Equal(t, now1, r1.ObservedTimestamp())
	assert.Equal(t, sev1, r1.Severity())
	assert.Equal(t, text1, r1.SeverityText())
	assert.Equal(t, val1, r1.Body())
	assert.Equal(t, traceID1, r1.TraceID())
	assert.Equal(t, spanID1, r1.SpanID())
	assert.Equal(t, flag1, r1.TraceFlags())
	r1.WalkAttributes(func(kv attribute.KeyValue) bool {
		return assert.Truef(t, keyValueEqual(kv, attr1), "%v != %v", kv, attr1)
	})
}

func TestRecordDroppedAttributes(t *testing.T) {
	orig := logAttrDropped
	t.Cleanup(func() { logAttrDropped = orig })

	for i := 1; i < attributesInlineCount*5; i++ {
		var called bool
		logAttrDropped = func() { called = true }

		r := new(Record)
		r.attributeCountLimit = 1

		attrs := make([]attribute.KeyValue, i)
		attrs[0] = attribute.Bool("only key different then the rest", true)

		r.AddAttributes(attrs...)
		// Deduplication doesn't count as dropped.
		wantDropped := 0
		if i > 1 {
			wantDropped = 1
		}
		assert.Equalf(t, wantDropped, r.DroppedAttributes(), "%d: AddAttributes", i)
		if i <= 1 {
			assert.False(t, called, "%d: dropped attributes logged", i)
		} else {
			assert.True(t, called, "%d: dropped attributes not logged", i)
		}

		r.AddAttributes(attrs...)
		wantDropped = 0
		if i > 1 {
			wantDropped = 2
		}
		assert.Equalf(t, wantDropped, r.DroppedAttributes(), "%d: second AddAttributes", i)

		r.SetAttributes(attrs...)
		wantDropped = 0
		if i > 1 {
			wantDropped = 1
		}
		assert.Equalf(t, wantDropped, r.DroppedAttributes(), "%d: SetAttributes", i)
	}
}

func TestRecordAttrAllowDuplicateAttributes(t *testing.T) {
	testcases := []struct {
		name  string
		attrs []attribute.KeyValue
		want  []attribute.KeyValue
	}{
		{
			name:  "EmptyKey",
			attrs: make([]attribute.KeyValue, 10),
			want:  make([]attribute.KeyValue, 10),
		},
		{
			name: "MapKey",
			attrs: []attribute.KeyValue{
				attribute.Map("key", attribute.Int("key", 5), attribute.Int("key", 10)),
			},
			want: []attribute.KeyValue{
				attribute.Map("key", attribute.Int("key", 5), attribute.Int("key", 10)),
			},
		},
		{
			name: "NonEmptyKey",
			attrs: []attribute.KeyValue{
				attribute.Bool("key", true),
				attribute.Int64("key", 1),
				attribute.Bool("key", false),
				attribute.Float64("key", 2.),
				attribute.String("key", "3"),
				attribute.Slice("key", attribute.Int64Value(4)),
				attribute.Map("key", attribute.Int("key", 5)),
				attribute.ByteSlice("key", []byte("six")),
				attribute.Bool("key", false),
			},
			want: []attribute.KeyValue{
				attribute.Bool("key", true),
				attribute.Int64("key", 1),
				attribute.Bool("key", false),
				attribute.Float64("key", 2.),
				attribute.String("key", "3"),
				attribute.Slice("key", attribute.Int64Value(4)),
				attribute.Map("key", attribute.Int("key", 5)),
				attribute.ByteSlice("key", []byte("six")),
				attribute.Bool("key", false),
			},
		},
		{
			name: "Multiple",
			attrs: []attribute.KeyValue{
				attribute.Bool("a", true),
				attribute.Int64("b", 1),
				attribute.Bool("a", false),
				attribute.Float64("c", 2.),
				attribute.String("b", "3"),
				attribute.Slice("d", attribute.Int64Value(4)),
				attribute.Map("a", attribute.Int("key", 5)),
				attribute.ByteSlice("d", []byte("six")),
				attribute.Bool("e", true),
				attribute.Int("f", 1),
				attribute.Int("f", 2),
				attribute.Int("f", 3),
				attribute.Float64("b", 0.0),
				attribute.Float64("b", 0.0),
				attribute.String("g", "G"),
				attribute.String("h", "H"),
				attribute.String("g", "GG"),
				attribute.Bool("a", false),
			},
			want: []attribute.KeyValue{
				attribute.Bool("a", true),
				attribute.Int64("b", 1),
				attribute.Bool("a", false),
				attribute.Float64("c", 2.),
				attribute.String("b", "3"),
				attribute.Slice("d", attribute.Int64Value(4)),
				attribute.Map("a", attribute.Int("key", 5)),
				attribute.ByteSlice("d", []byte("six")),
				attribute.Bool("e", true),
				attribute.Int("f", 1),
				attribute.Int("f", 2),
				attribute.Int("f", 3),
				attribute.Float64("b", 0.0),
				attribute.Float64("b", 0.0),
				attribute.String("g", "G"),
				attribute.String("h", "H"),
				attribute.String("g", "GG"),
				attribute.Bool("a", false),
			},
		},
		{
			name: "NoDuplicate",
			attrs: func() []attribute.KeyValue {
				out := make([]attribute.KeyValue, attributesInlineCount*2)
				for i := range out {
					out[i] = attribute.Bool(strconv.Itoa(i), true)
				}
				return out
			}(),
			want: func() []attribute.KeyValue {
				out := make([]attribute.KeyValue, attributesInlineCount*2)
				for i := range out {
					out[i] = attribute.Bool(strconv.Itoa(i), true)
				}
				return out
			}(),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			validate := func(t *testing.T, r *Record, want []attribute.KeyValue) {
				t.Helper()

				var i int
				r.WalkAttributes(func(kv attribute.KeyValue) bool {
					if assert.Lessf(t, i, len(want), "additional: %v", kv) {
						want := want[i]
						assert.Truef(t, keyValueEqual(kv, want), "%d: want %v, got %v", i, want, kv)
					}
					i++
					return true
				})
			}

			t.Run("SetAttributes", func(t *testing.T) {
				r := new(Record)
				r.allowDupKeys = true
				r.attributeValueLengthLimit = -1
				r.SetAttributes(tc.attrs...)
				validate(t, r, tc.want)
			})

			t.Run("AddAttributes/Empty", func(t *testing.T) {
				r := new(Record)
				r.allowDupKeys = true
				r.attributeValueLengthLimit = -1
				r.AddAttributes(tc.attrs...)
				validate(t, r, tc.want)
			})

			t.Run("AddAttributes/Twice", func(t *testing.T) {
				r := new(Record)
				r.allowDupKeys = true
				r.attributeValueLengthLimit = -1
				r.AddAttributes(tc.attrs...)
				r.AddAttributes(tc.attrs...)
				want := append(tc.want, tc.want...)
				validate(t, r, want)
			})
		})
	}
}

func TestRecordAttrDeduplication(t *testing.T) {
	testcases := []struct {
		name  string
		attrs []attribute.KeyValue
		want  []attribute.KeyValue
	}{
		{
			name:  "EmptyKey",
			attrs: make([]attribute.KeyValue, 10),
			want:  make([]attribute.KeyValue, 1),
		},
		{
			name: "NonEmptyKey",
			attrs: []attribute.KeyValue{
				attribute.Bool("key", true),
				attribute.Int64("key", 1),
				attribute.Bool("key", false),
				attribute.Float64("key", 2.),
				attribute.String("key", "3"),
				attribute.Slice("key", attribute.Int64Value(4)),
				attribute.Map("key", attribute.Int("key", 5)),
				attribute.ByteSlice("key", []byte("six")),
				attribute.Bool("key", false),
			},
			want: []attribute.KeyValue{
				attribute.Bool("key", false),
			},
		},
		{
			name: "Multiple",
			attrs: []attribute.KeyValue{
				attribute.Bool("a", true),
				attribute.Int64("b", 1),
				attribute.Bool("a", false),
				attribute.Float64("c", 2.),
				attribute.String("b", "3"),
				attribute.Slice("d", attribute.Int64Value(4)),
				attribute.Map("a", attribute.Int("key", 5)),
				attribute.ByteSlice("d", []byte("six")),
				attribute.Bool("e", true),
				attribute.Int("f", 1),
				attribute.Int("f", 2),
				attribute.Int("f", 3),
				attribute.Float64("b", 0.0),
				attribute.Float64("b", 0.0),
				attribute.String("g", "G"),
				attribute.String("h", "H"),
				attribute.String("g", "GG"),
				attribute.Bool("a", false),
			},
			want: []attribute.KeyValue{
				attribute.Bool("a", false),
				attribute.Float64("b", 0.0),
				attribute.Float64("c", 2.),
				attribute.ByteSlice("d", []byte("six")),
				attribute.Bool("e", true),
				attribute.Int("f", 3),
				attribute.String("g", "GG"),
				attribute.String("h", "H"),
			},
		},
		{
			name: "NoDuplicate",
			attrs: func() []attribute.KeyValue {
				out := make([]attribute.KeyValue, attributesInlineCount*2)
				for i := range out {
					out[i] = attribute.Bool(strconv.Itoa(i), true)
				}
				return out
			}(),
			want: func() []attribute.KeyValue {
				out := make([]attribute.KeyValue, attributesInlineCount*2)
				for i := range out {
					out[i] = attribute.Bool(strconv.Itoa(i), true)
				}
				return out
			}(),
		},
		{
			name: "AttributeWithDuplicateKeys",
			attrs: []attribute.KeyValue{
				attribute.String("duplicate", "first"),
				attribute.String("unique", "value"),
				attribute.String("duplicate", "second"),
			},
			want: []attribute.KeyValue{
				attribute.String("duplicate", "second"),
				attribute.String("unique", "value"),
			},
		},
		{
			name: "ManyDuplicateKeys",
			attrs: []attribute.KeyValue{
				attribute.String("key", "value1"),
				attribute.String("key", "value2"),
				attribute.String("key", "value3"),
				attribute.String("key", "value4"),
				attribute.String("key", "value5"),
			},
			want: []attribute.KeyValue{
				attribute.String("key", "value5"),
			},
		},
		{
			name: "InterleavedDuplicates",
			attrs: []attribute.KeyValue{
				attribute.String("a", "a1"),
				attribute.String("b", "b1"),
				attribute.String("a", "a2"),
				attribute.String("c", "c1"),
				attribute.String("b", "b2"),
			},
			want: []attribute.KeyValue{
				attribute.String("a", "a2"),
				attribute.String("b", "b2"),
				attribute.String("c", "c1"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			validate := func(t *testing.T, r *Record) {
				t.Helper()

				var i int
				r.WalkAttributes(func(kv attribute.KeyValue) bool {
					if assert.Lessf(t, i, len(tc.want), "additional: %v", kv) {
						want := tc.want[i]
						assert.Truef(t, keyValueEqual(kv, want), "%d: want %v, got %v", i, want, kv)
					}
					i++
					return true
				})
			}

			t.Run("SetAttributes", func(t *testing.T) {
				r := new(Record)
				r.attributeValueLengthLimit = -1
				r.SetAttributes(tc.attrs...)
				validate(t, r)
			})

			t.Run("AddAttributes/Empty", func(t *testing.T) {
				r := new(Record)
				r.attributeValueLengthLimit = -1
				r.AddAttributes(tc.attrs...)
				validate(t, r)
			})

			t.Run("AddAttributes/Duplicates", func(t *testing.T) {
				r := new(Record)
				r.attributeValueLengthLimit = -1
				r.AddAttributes(tc.attrs...)
				r.AddAttributes(tc.attrs...)
				validate(t, r)
			})
		})
	}
}

func TestApplyAttrLimitsDeduplication(t *testing.T) {
	testcases := []struct {
		name             string
		limit            int
		input, want      attribute.Value
		wantDroppedAttrs int
	}{
		{
			// No de-duplication
			name: "Slice",
			input: attribute.SliceValue(
				attribute.BoolValue(true),
				attribute.BoolValue(true),
				attribute.Float64Value(1.3),
				attribute.Float64Value(1.3),
				attribute.Int64Value(43),
				attribute.Int64Value(43),
				attribute.ByteSliceValue([]byte("hello")),
				attribute.ByteSliceValue([]byte("hello")),
				attribute.StringValue("foo"),
				attribute.StringValue("foo"),
				attribute.SliceValue(attribute.StringValue("baz")),
				attribute.SliceValue(attribute.StringValue("baz")),
				attribute.MapValue(attribute.String("a", "qux")),
				attribute.MapValue(attribute.String("a", "qux")),
			),
			want: attribute.SliceValue(
				attribute.BoolValue(true),
				attribute.BoolValue(true),
				attribute.Float64Value(1.3),
				attribute.Float64Value(1.3),
				attribute.Int64Value(43),
				attribute.Int64Value(43),
				attribute.ByteSliceValue([]byte("hello")),
				attribute.ByteSliceValue([]byte("hello")),
				attribute.StringValue("foo"),
				attribute.StringValue("foo"),
				attribute.SliceValue(attribute.StringValue("baz")),
				attribute.SliceValue(attribute.StringValue("baz")),
				attribute.MapValue(attribute.String("a", "qux")),
				attribute.MapValue(attribute.String("a", "qux")),
			),
		},
		{
			name: "Map",
			input: attribute.MapValue(
				attribute.Bool("a", true),
				attribute.Int64("b", 1),
				attribute.Bool("a", false),
				attribute.Float64("c", 2.),
				attribute.String("b", "3"),
				attribute.Slice("d", attribute.Int64Value(4)),
				attribute.Map("a", attribute.Int("key", 5)),
				attribute.ByteSlice("d", []byte("six")),
				attribute.Bool("e", true),
				attribute.Int("f", 1),
				attribute.Int("f", 2),
				attribute.Int("f", 3),
				attribute.Float64("b", 0.0),
				attribute.Float64("b", 0.0),
				attribute.String("g", "G"),
				attribute.String("h", "H"),
				attribute.String("g", "GG"),
				attribute.Bool("a", false),
			),
			want: attribute.MapValue(

				attribute.Bool("a", false),
				attribute.Float64("b", 0.0),
				attribute.Float64("c", 2.),
				attribute.ByteSlice("d", []byte("six")),
				attribute.Bool("e", true),
				attribute.Int("f", 3),
				attribute.String("g", "GG"),
				attribute.String("h", "H"),
			),
			wantDroppedAttrs: 0, // Deduplication doesn't count as dropped
		},
		{
			name:             "EmptyMap",
			input:            attribute.MapValue(),
			want:             attribute.MapValue(),
			wantDroppedAttrs: 0,
		},
		{
			name:             "SingleKeyMap",
			input:            attribute.MapValue(attribute.String("key1", "value1")),
			want:             attribute.MapValue(attribute.String("key1", "value1")),
			wantDroppedAttrs: 0,
		},
		{
			name:             "EmptySlice",
			input:            attribute.SliceValue(),
			want:             attribute.SliceValue(),
			wantDroppedAttrs: 0,
		},
		{
			name: "SliceWithNestedDedup",
			input: attribute.SliceValue(
				attribute.MapValue(attribute.String("key", "value1"), attribute.String("key", "value2")),
				attribute.StringValue("normal"),
			),
			want: attribute.SliceValue(
				attribute.MapValue(attribute.String("key", "value2")),
				attribute.StringValue("normal"),
			),
			wantDroppedAttrs: 0, // Nested deduplication doesn't count as dropped
		},
		{
			name: "NestedSliceInMap",
			input: attribute.MapValue(
				attribute.Slice(
					"slice_key",
					attribute.MapValue(attribute.String("nested", "value1"), attribute.String("nested", "value2")),
				),
			),

			want: attribute.MapValue(
				attribute.Slice(
					"slice_key",
					attribute.MapValue(attribute.String("nested", "value2")),
				),
			),

			wantDroppedAttrs: 0, // Nested deduplication doesn't count as dropped
		},
		{
			name: "DeeplyNestedStructure",
			input: attribute.MapValue(
				attribute.Map(
					"level1",
					attribute.Map(
						"level2",
						attribute.Slice(
							"level3",
							attribute.MapValue(attribute.String("deep", "value1"), attribute.String("deep", "value2")),
						),
					),
				),
			),

			want: attribute.MapValue(
				attribute.Map(
					"level1",
					attribute.Map(
						"level2",
						attribute.Slice(
							"level3",
							attribute.MapValue(attribute.String("deep", "value2")),
						),
					),
				),
			),

			wantDroppedAttrs: 0, // Deeply nested deduplication doesn't count as dropped
		},
		{
			name: "NestedMapWithoutDuplicateKeys",
			input: attribute.SliceValue(attribute.MapValue(
				attribute.String("key1", "value1"),
				attribute.String("key2", "value2"),
			)),

			want: attribute.SliceValue(attribute.MapValue(
				attribute.String("key1", "value1"),
				attribute.String("key2", "value2"),
			)),

			wantDroppedAttrs: 0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			const key = "key"
			kv := attribute.KeyValue{Key: key, Value: tc.input}
			r := Record{attributeValueLengthLimit: -1}

			t.Run("AddAttributes", func(t *testing.T) {
				r.AddAttributes(kv)
				assertKV(t, r, attribute.KeyValue{Key: key, Value: tc.want})
				assert.Equal(t, tc.wantDroppedAttrs, r.DroppedAttributes())
			})

			t.Run("SetAttributes", func(t *testing.T) {
				r.SetAttributes(kv)
				assertKV(t, r, attribute.KeyValue{Key: key, Value: tc.want})
				assert.Equal(t, tc.wantDroppedAttrs, r.DroppedAttributes())
			})
		})
	}
}

func TestDeduplicationBehavior(t *testing.T) {
	origKeyValueDropped := logKeyValuePairDropped
	origAttrDropped := logAttrDropped
	t.Cleanup(func() {
		logKeyValuePairDropped = origKeyValueDropped
		logAttrDropped = origAttrDropped
	})

	testCases := []struct {
		name                string
		attributeCountLimit int
		allowDupKeys        bool
		attrs               []attribute.KeyValue
		wantKeyValueDropped bool
		wantAttrDropped     bool
		wantDroppedCount    int
		wantAttributeCount  int
	}{
		{
			name:                "Duplicate keys only",
			attrs:               []attribute.KeyValue{attribute.String("key", "v1"), attribute.String("key", "v2")},
			wantKeyValueDropped: true,
			wantDroppedCount:    0, // Deduplication doesn't count
			wantAttributeCount:  1,
		},
		{
			name:                "Limit exceeded only",
			attributeCountLimit: 2,
			attrs: []attribute.KeyValue{
				attribute.String("a", "v1"),
				attribute.String("b", "v2"),
				attribute.String("c", "v3"),
			},
			wantAttrDropped:    true,
			wantDroppedCount:   1,
			wantAttributeCount: 2,
		},
		{
			name:                "Both duplicates and limit",
			attributeCountLimit: 2,
			attrs: []attribute.KeyValue{
				attribute.String("a", "v1"),
				attribute.String("a", "v2"),
				attribute.String("b", "v3"),
				attribute.String("c", "v4"),
			},
			wantKeyValueDropped: true,
			wantAttrDropped:     true,
			wantDroppedCount:    1, // Only limit drops count
			wantAttributeCount:  2,
		},
		{
			name:                "allowDupKeys=true",
			allowDupKeys:        true,
			attrs:               []attribute.KeyValue{attribute.String("key", "v1"), attribute.String("key", "v2")},
			wantKeyValueDropped: false,
			wantDroppedCount:    0,
			wantAttributeCount:  2,
		},
		{
			name: "Nested map duplicates",
			attrs: []attribute.KeyValue{
				attribute.Map("outer", attribute.String("nested", "v1"), attribute.String("nested", "v2")),
			},
			wantKeyValueDropped: true,
			wantDroppedCount:    0,
			wantAttributeCount:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keyValueDroppedCalled := false
			attrDroppedCalled := false

			logKeyValuePairDropped = sync.OnceFunc(func() { keyValueDroppedCalled = true })
			logAttrDropped = sync.OnceFunc(func() { attrDroppedCalled = true })

			r := &Record{
				attributeValueLengthLimit: -1,
				attributeCountLimit:       tc.attributeCountLimit,
				allowDupKeys:              tc.allowDupKeys,
			}

			r.SetAttributes(tc.attrs...)

			assert.Equal(t, tc.wantKeyValueDropped, keyValueDroppedCalled)
			assert.Equal(t, tc.wantAttrDropped, attrDroppedCalled)
			assert.Equal(t, tc.wantDroppedCount, r.DroppedAttributes())
			assert.Equal(t, tc.wantAttributeCount, r.AttributesLen())
		})
	}
}

func TestApplyAttrLimitsTruncation(t *testing.T) {
	testcases := []struct {
		name        string
		limit       int
		input, want attribute.Value
	}{
		{
			name:  "Empty",
			limit: 0,
			input: attribute.Value{},
			want:  attribute.Value{},
		},
		{
			name:  "Bool",
			limit: 0,
			input: attribute.BoolValue(true),
			want:  attribute.BoolValue(true),
		},
		{
			name:  "Float64",
			limit: 0,
			input: attribute.Float64Value(1.3),
			want:  attribute.Float64Value(1.3),
		},
		{
			name:  "Int64",
			limit: 0,
			input: attribute.Int64Value(43),
			want:  attribute.Int64Value(43),
		},
		{
			name:  "Bytes",
			limit: 0,
			input: attribute.ByteSliceValue([]byte("foo")),
			want:  attribute.ByteSliceValue([]byte("")),
		},
		{
			name:  "String",
			limit: 0,
			input: attribute.StringValue("foo"),
			want:  attribute.StringValue(""),
		},
		{
			name:  "StringSlice",
			limit: 3,
			input: attribute.StringSliceValue([]string{"ok", "toolong"}),
			want:  attribute.StringSliceValue([]string{"ok", "too"}),
		},
		{
			name:  "NestedSliceNoTruncation",
			limit: 10,
			input: attribute.SliceValue(
				attribute.StringValue("short"),
				attribute.StringSliceValue([]string{"ok", "fine"}),
			),
			want: attribute.SliceValue(
				attribute.StringValue("short"),
				attribute.StringSliceValue([]string{"ok", "fine"}),
			),
		},
		{
			name:  "NestedMapNoTruncation",
			limit: 10,
			input: attribute.MapValue(
				attribute.String("short", "ok"),
				attribute.StringSlice("strings", []string{"ok", "fine"}),
			),
			want: attribute.MapValue(
				attribute.String("short", "ok"),
				attribute.StringSlice("strings", []string{"ok", "fine"}),
			),
		},
		{
			name:  "StringSliceInNestedSlice",
			limit: 3,
			input: attribute.SliceValue(
				attribute.StringSliceValue([]string{"ok", "toolong"}),
			),
			want: attribute.SliceValue(
				attribute.StringSliceValue([]string{"ok", "too"}),
			),
		},
		{
			name:  "Slice",
			limit: 0,
			input: attribute.SliceValue(
				attribute.BoolValue(true),
				attribute.Float64Value(1.3),
				attribute.Int64Value(43),
				attribute.ByteSliceValue([]byte("hello")),
				attribute.StringValue("foo"),
				attribute.StringValue("bar"),
				attribute.SliceValue(attribute.StringValue("baz")),
				attribute.MapValue(attribute.String("a", "qux")),
			),
			want: attribute.SliceValue(
				attribute.BoolValue(true),
				attribute.Float64Value(1.3),
				attribute.Int64Value(43),
				attribute.ByteSliceValue([]byte("")),
				attribute.StringValue(""),
				attribute.StringValue(""),
				attribute.SliceValue(attribute.StringValue("")),
				attribute.MapValue(attribute.String("a", "")),
			),
		},
		{
			name:  "Map",
			limit: 0,
			input: attribute.MapValue(
				attribute.Bool("0", true),
				attribute.Float64("1", 1.3),
				attribute.Int64("2", 43),
				attribute.ByteSlice("3", []byte("hello")),
				attribute.String("4", "foo"),
				attribute.String("5", "bar"),
				attribute.Slice("6", attribute.StringValue("baz")),
				attribute.Map("7", attribute.String("a", "qux")),
			),
			want: attribute.MapValue(
				attribute.Bool("0", true),
				attribute.Float64("1", 1.3),
				attribute.Int64("2", 43),
				attribute.ByteSlice("3", []byte("")),
				attribute.String("4", ""),
				attribute.String("5", ""),
				attribute.Slice("6", attribute.StringValue("")),
				attribute.Map("7", attribute.String("a", "")),
			),
		},
		{
			name:  "LongStringTruncated",
			limit: 5,
			input: attribute.StringValue("This is a very long string that should be truncated"),
			want:  attribute.StringValue("This "),
		},
		{
			name:  "LongBytesTruncated",
			limit: 5,
			input: attribute.ByteSliceValue([]byte("This is a very long byte array")),
			want:  attribute.ByteSliceValue([]byte("This ")),
		},
		{
			name:  "TruncationInNestedMap",
			limit: 3,
			input: attribute.MapValue(
				attribute.String("short", "ok"),
				attribute.String("long", "toolong"),
			),
			want: attribute.MapValue(
				attribute.String("short", "ok"),
				attribute.String("long", "too"),
			),
		},
		{
			name:  "TruncationInNestedSlice",
			limit: 4,
			input: attribute.SliceValue(
				attribute.StringValue("good"),
				attribute.StringValue("toolong"),
			),
			want: attribute.SliceValue(
				attribute.StringValue("good"),
				attribute.StringValue("tool"),
			),
		},
		{
			name:  "TruncationInNestedSliceOfBytes",
			limit: 4,
			input: attribute.SliceValue(
				attribute.ByteSliceValue([]byte("good")),
				attribute.ByteSliceValue([]byte("toolong")),
			),
			want: attribute.SliceValue(
				attribute.ByteSliceValue([]byte("good")),
				attribute.ByteSliceValue([]byte("tool")),
			),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			const key = "key"
			kv := attribute.KeyValue{Key: key, Value: tc.input}
			r := Record{attributeValueLengthLimit: tc.limit}

			t.Run("AddAttributes", func(t *testing.T) {
				r.AddAttributes(kv)
				assertKV(t, r, attribute.KeyValue{Key: key, Value: tc.want})
			})

			t.Run("SetAttributes", func(t *testing.T) {
				r.SetAttributes(kv)
				assertKV(t, r, attribute.KeyValue{Key: key, Value: tc.want})
			})
		})
	}
}

func assertKV(t *testing.T, r Record, kv attribute.KeyValue) {
	t.Helper()

	var kvs []attribute.KeyValue
	r.WalkAttributes(func(kv attribute.KeyValue) bool {
		kvs = append(kvs, kv)
		return true
	})

	require.Len(t, kvs, 1)
	assert.Truef(t, keyValueEqual(kv, kvs[0]), "%s != %s", kv, kvs[0])
}

func TestTruncate(t *testing.T) {
	type group struct {
		limit    int
		input    string
		expected string
	}

	tests := []struct {
		name   string
		groups []group
	}{
		// Edge case: limit is negative, no truncation should occur
		{
			name: "NoTruncation",
			groups: []group{
				{-1, "No truncation!", "No truncation!"},
			},
		},

		// Edge case: string is already shorter than the limit, no truncation
		// should occur
		{
			name: "ShortText",
			groups: []group{
				{10, "Short text", "Short text"},
				{15, "Short text", "Short text"},
				{100, "Short text", "Short text"},
			},
		},

		// Edge case: truncation happens with ASCII characters only
		{
			name: "ASCIIOnly",
			groups: []group{
				{1, "Hello World!", "H"},
				{5, "Hello World!", "Hello"},
				{12, "Hello World!", "Hello World!"},
			},
		},

		// Truncation including multi-byte characters (UTF-8)
		{
			name: "ValidUTF-8",
			groups: []group{
				{7, "Hello, 世界", "Hello, "},
				{8, "Hello, 世界", "Hello, 世"},
				{2, "こんにちは", "こん"},
				{3, "こんにちは", "こんに"},
				{5, "こんにちは", "こんにちは"},
				{12, "こんにちは", "こんにちは"},
			},
		},

		// Truncation with invalid UTF-8 characters
		{
			name: "InvalidUTF-8",
			groups: []group{
				{11, "Invalid\x80text", "Invalidtext"},
				// Do not modify invalid text if equal to limit.
				{11, "Valid text\x80", "Valid text\x80"},
				// Do not modify invalid text if under limit.
				{15, "Valid text\x80", "Valid text\x80"},
				{5, "Hello\x80World", "Hello"},
				{11, "Hello\x80World\x80!", "HelloWorld!"},
				{15, "Hello\x80World\x80Test", "HelloWorldTest"},
				{15, "Hello\x80\x80\x80World\x80Test", "HelloWorldTest"},
				{15, "\x80\x80\x80Hello\x80\x80\x80World\x80Test\x80\x80", "HelloWorldTest"},
			},
		},

		// Truncation with mixed validn and invalid UTF-8 characters
		{
			name: "MixedUTF-8",
			groups: []group{
				{6, "€"[0:2] + "hello€€", "hello€"},
				{6, "€" + "€"[0:2] + "hello", "€hello"},
				{11, "Valid text\x80📜", "Valid text📜"},
				{11, "Valid text📜\x80", "Valid text📜"},
				{14, "😊 Hello\x80World🌍🚀", "😊 HelloWorld🌍🚀"},
				{14, "😊\x80 Hello\x80World🌍🚀", "😊 HelloWorld🌍🚀"},
				{14, "😊\x80 Hello\x80World🌍\x80🚀", "😊 HelloWorld🌍🚀"},
				{14, "😊\x80 Hello\x80World🌍\x80🚀\x80", "😊 HelloWorld🌍🚀"},
				{14, "\x80😊\x80 Hello\x80World🌍\x80🚀\x80", "😊 HelloWorld🌍🚀"},
			},
		},

		// Edge case: empty string, should return empty string
		{
			name: "Empty",
			groups: []group{
				{5, "", ""},
			},
		},

		// Edge case: limit is 0, should return an empty string
		{
			name: "Zero",
			groups: []group{
				{0, "Some text", ""},
				{0, "", ""},
			},
		},
	}

	for _, tt := range tests {
		for _, g := range tt.groups {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := truncate(g.limit, g.input)
				assert.Equalf(
					t, g.expected, got,
					"input: %q([]rune%v))\ngot: %q([]rune%v)\nwant %q([]rune%v)",
					g.input, []rune(g.input),
					got, []rune(got),
					g.expected, []rune(g.expected),
				)
			})
		}
	}
}

func TestRecordAddAttributesDoesNotMutateInput(t *testing.T) {
	attrs := []attribute.KeyValue{
		attribute.String("attr1", "very long value that will be truncated"),
		attribute.String("attr2", "another very long value that will be truncated"),
		attribute.String("attr3", "yet another very long value that will be truncated"),
		attribute.String("attr4", "more very long value that will be truncated"),
		attribute.String("attr5", "extra very long value that will be truncated"),
		attribute.String("attr6", "additional very long value that will be truncated"),
		attribute.String("attr7", "more additional very long value that will be truncated"),
	}

	originalValues := make([]string, len(attrs))
	for i, kv := range attrs {
		originalValues[i] = kv.Value.AsString()
	}

	r := &Record{
		attributeValueLengthLimit: 20, // Short limit to trigger truncation.
		attributeCountLimit:       -1, // No count limit.
		allowDupKeys:              false,
	}

	r.AddAttributes(attrs...)

	// Verify that the original shared slice was not mutated
	for i, kv := range attrs {
		if kv.Value.AsString() != originalValues[i] {
			t.Errorf("Input slice was mutated! Attribute %d: original=%q, current=%q",
				i, originalValues[i], kv.Value.AsString())
		}
	}

	// Verify that the record has the truncated values
	var gotAttrs []attribute.KeyValue
	r.WalkAttributes(func(kv attribute.KeyValue) bool {
		gotAttrs = append(gotAttrs, kv)
		return true
	})
	wantAttr := []attribute.KeyValue{
		attribute.String("attr1", "very long value that"),
		attribute.String("attr2", "another very long va"),
		attribute.String("attr3", "yet another very lon"),
		attribute.String("attr4", "more very long value"),
		attribute.String("attr5", "extra very long valu"),
		attribute.String("attr6", "additional very long"),
		attribute.String("attr7", "more additional very"),
	}
	if !slices.EqualFunc(gotAttrs, wantAttr, keyValueEqual) {
		t.Errorf("Attributes do not match.\ngot:\n%v\nwant:\n%v", printKVs(gotAttrs), printKVs(wantAttr))
	}
}

func TestRecordMethodsInputConcurrentSafe(t *testing.T) {
	nestedSlice := attribute.Slice(
		"nested_slice",
		attribute.SliceValue(attribute.StringValue("nested_inner1"), attribute.StringValue("nested_inner2")),
		attribute.StringValue("nested_outer"),
	)

	nestedMap := attribute.Map(
		"nested_map",
		attribute.String("nested_key1", "nested_value1"),
		attribute.Map("nested_map", attribute.String("nested_inner_key", "nested_inner_value")),
		attribute.String("nested_key1", "duplicate"),
	)

	dedupAttributes := []attribute.KeyValue{
		attribute.String("dedup_key1", "dedup_value1"),
		attribute.String("dedup_key2", "dedup_value2"),
		attribute.String("dedup_key1", "duplicate"),
		attribute.String("dedup_key3", "dedup_value3"),
	}

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			r := &Record{
				attributeValueLengthLimit: 10,
				attributeCountLimit:       4,
				allowDupKeys:              false,
			}

			r.SetAttributes(nestedSlice)
			r.AddAttributes(nestedMap)
			r.AddAttributes(dedupAttributes...)
			r.SetBody(nestedMap.Value)

			var gotAttrs []attribute.KeyValue
			r.WalkAttributes(func(kv attribute.KeyValue) bool {
				gotAttrs = append(gotAttrs, kv)
				return true
			})
			wantAttr := []attribute.KeyValue{
				attribute.Slice(
					"nested_slice",
					attribute.SliceValue(attribute.StringValue("nested_inn"), attribute.StringValue("nested_inn")),
					attribute.StringValue("nested_out"),
				),
				attribute.Map(
					"nested_map",
					attribute.String("nested_key1", "duplicate"),
					attribute.Map("nested_map", attribute.String("nested_inner_key", "nested_inn")),
				),
				attribute.String("dedup_key1", "duplicate"),
				attribute.String("dedup_key2", "dedup_valu"),
			}
			if !slices.EqualFunc(gotAttrs, wantAttr, keyValueEqual) {
				t.Errorf("Attributes do not match.\ngot:\n%v\nwant:\n%v", printKVs(gotAttrs), printKVs(wantAttr))
			}

			gotBody := r.Body()
			wantBody := attribute.MapValue(
				attribute.String("nested_key1", "duplicate"),
				attribute.Map("nested_map", attribute.String("nested_inner_key", "nested_inner_value")),
			)
			if !valueEqual(gotBody, wantBody) {
				t.Errorf("Body does not match.\ngot:\n%v\nwant:\n%v", gotBody, wantBody)
			}
		})
	}

	wg.Wait()
}

func printKVs(kvs []attribute.KeyValue) string {
	var sb strings.Builder
	for _, kv := range kvs {
		_, _ = fmt.Fprintf(&sb, "%s: %s\n", kv.Key, kv.Value)
	}
	return sb.String()
}

func BenchmarkTruncate(b *testing.B) {
	run := func(limit int, input string) func(b *testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				var out string
				for pb.Next() {
					out = truncate(limit, input)
				}
				_ = out
			})
		}
	}
	b.Run("Unlimited", run(-1, "hello 😊 world 🌍🚀"))
	b.Run("Zero", run(0, "Some text"))
	b.Run("Short", run(10, "Short Text"))
	b.Run("ASCII", run(5, "Hello, World!"))
	b.Run("ValidUTF-8", run(10, "hello 😊 world 🌍🚀"))
	b.Run("InvalidUTF-8", run(6, "€"[0:2]+"hello€€"))
	b.Run("MixedUTF-8", run(14, "\x80😊\x80 Hello\x80World🌍\x80🚀\x80"))
}

func BenchmarkWalkAttributes(b *testing.B) {
	for _, tt := range []struct {
		attrCount int
	}{
		{attrCount: 1},
		{attrCount: 10},
		{attrCount: 100},
		{attrCount: 1000},
	} {
		b.Run(fmt.Sprintf("%d attributes", tt.attrCount), func(b *testing.B) {
			record := &Record{}
			for i := 0; i < tt.attrCount; i++ {
				record.SetAttributes(
					attribute.String(fmt.Sprintf("key-%d", tt.attrCount), "value"),
				)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				record.WalkAttributes(func(attribute.KeyValue) bool {
					return true
				})
			}
		})
	}
}

func BenchmarkAddAttributes(b *testing.B) {
	// Simple attribute (no deduplication or limits).
	singleKV := attribute.String("key", "value")

	// Attributes with no duplicates.
	uniqueAttrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.String("key2", "value2"),
		attribute.String("key3", "value3"),
		attribute.String("key4", "value4"),
		attribute.String("key5", "value5"),
	}

	// Attributes with duplicates that trigger deduplication.
	dupAttrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.String("key2", "value2"),
		attribute.String("key1", "duplicate1"),
		attribute.String("key3", "value3"),
		attribute.String("key2", "duplicate2"),
	}

	// Large number of attributes to trigger count limits.
	manyAttrs := make([]attribute.KeyValue, 20)
	for i := range manyAttrs {
		manyAttrs[i] = attribute.String(fmt.Sprintf("key%d", i), "value")
	}

	// Attributes with long values to trigger value length limits.
	longValueAttrs := []attribute.KeyValue{
		attribute.String("short", "short"),
		attribute.String("long1", strings.Repeat("a", 50)),
		attribute.String("long2", strings.Repeat("b", 100)),
	}

	// Attributes with nested maps that have duplicates (triggers recursive deduplication).
	nestedDupAttrs := []attribute.KeyValue{
		attribute.String("simple", "value"),
		attribute.Map(
			"map1",
			attribute.String("inner1", "value1"),
			attribute.String("inner2", "value2"),
			attribute.String("inner1", "duplicate"),
		),
		attribute.Map(
			"map2",
			attribute.String("key", "original"),
			attribute.Map(
				"deeply_nested",
				attribute.String("deep1", "value1"),
				attribute.String("deep2", "value2"),
				attribute.String("deep1", "duplicate_deep"),
			),
			attribute.String("key", "overwrite"),
		),
		attribute.Slice(
			"slice_with_maps",
			attribute.MapValue(
				attribute.String("slice_key", "value1"),
				attribute.String("slice_key", "duplicate"),
			),
		),
	}

	// Adding a single attribute with no limits applied.
	b.Run("Single/NoLimits", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(singleKV)
		}
	})

	// Adding a single attribute with duplicate keys allowed (faster path).
	b.Run("Single/AllowDuplicates", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(singleKV)
		}
	})

	// Adding multiple unique attributes with no limits applied.
	b.Run("Unique/NoLimits", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(uniqueAttrs...)
		}
	})

	// Adding multiple unique attributes with duplicate keys allowed (faster path).
	b.Run("Unique/AllowDuplicates", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(uniqueAttrs...)
		}
	})

	// Adding attributes with duplicates that trigger deduplication logic.
	b.Run("Deduplication/Enabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(dupAttrs...)
		}
	})

	// Adding nested maps with duplicates that trigger recursive deduplication.
	b.Run("NestedDeduplication/Enabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(nestedDupAttrs...)
		}
	})

	// Adding nested maps with duplicates with deduplication disabled.
	b.Run("NestedDeduplication/Disabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(nestedDupAttrs...)
		}
	})

	// Adding attributes with duplicates with deduplication disabled.
	b.Run("Deduplication/Disabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(dupAttrs...)
		}
	})

	// Adding more attributes than the count limit allows (triggers dropping).
	b.Run("CountLimit/Hit", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = 10 // Less than manyAttrs length.
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(manyAttrs...)
		}
	})

	// Adding attributes within the count limit (no dropping).
	b.Run("CountLimit/NotHit", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = 100 // More than manyAttrs length.
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(manyAttrs...)
		}
	})

	// Adding attributes with long values that trigger truncation.
	b.Run("ValueLimit/Hit", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = 20 // Less than long values.
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(longValueAttrs...)
		}
	})

	// Adding attributes with values within the length limit (no truncation).
	b.Run("ValueLimit/NotHit", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = 200 // More than long values.
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].AddAttributes(longValueAttrs...)
		}
	})
}

func BenchmarkSetAttributes(b *testing.B) {
	// Simple attribute (no deduplication or limits).
	singleKV := attribute.String("key", "value")

	// Attributes with no duplicates.
	uniqueAttrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.String("key2", "value2"),
		attribute.String("key3", "value3"),
		attribute.String("key4", "value4"),
		attribute.String("key5", "value5"),
	}

	// Attributes with duplicates that trigger deduplication.
	dupAttrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.String("key2", "value2"),
		attribute.String("key1", "duplicate1"),
		attribute.String("key3", "value3"),
		attribute.String("key2", "duplicate2"),
	}

	// Large number of attributes to trigger count limits.
	manyAttrs := make([]attribute.KeyValue, 20)
	for i := range manyAttrs {
		manyAttrs[i] = attribute.String(fmt.Sprintf("key%d", i), "value")
	}

	// Attributes with long values to trigger value length limits.
	longValueAttrs := []attribute.KeyValue{
		attribute.String("short", "short"),
		attribute.String("long1", strings.Repeat("a", 50)),
		attribute.String("long2", strings.Repeat("b", 100)),
	}

	// Attributes with nested maps that have duplicates (triggers recursive deduplication).
	nestedDupAttrs := []attribute.KeyValue{
		attribute.String("simple", "value"),
		attribute.Map(
			"map1",
			attribute.String("inner1", "value1"),
			attribute.String("inner2", "value2"),
			attribute.String("inner1", "duplicate"),
		),
		attribute.Map(
			"map2",
			attribute.String("key", "original"),
			attribute.Map(
				"deeply_nested",
				attribute.String("deep1", "value1"),
				attribute.String("deep2", "value2"),
				attribute.String("deep1", "duplicate_deep"),
			),
			attribute.String("key", "overwrite"),
		),
		attribute.Slice(
			"slice_with_maps",
			attribute.MapValue(
				attribute.String("slice_key", "value1"),
				attribute.String("slice_key", "duplicate"),
			),
		),
	}

	// Setting a single attribute with no limits applied.
	b.Run("Single/NoLimits", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(singleKV)
		}
	})

	// Setting a single attribute with duplicate keys allowed (faster path).
	b.Run("Single/AllowDuplicates", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(singleKV)
		}
	})

	// Setting multiple unique attributes with no limits applied.
	b.Run("Unique/NoLimits", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(uniqueAttrs...)
		}
	})

	// Setting multiple unique attributes with duplicate keys allowed (faster path).
	b.Run("Unique/AllowDuplicates", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(uniqueAttrs...)
		}
	})

	// Setting attributes with duplicates that trigger deduplication logic.
	b.Run("Deduplication/Enabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(dupAttrs...)
		}
	})

	// Setting attributes with duplicates with deduplication disabled.
	b.Run("Deduplication/Disabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(dupAttrs...)
		}
	})

	// Setting nested maps with duplicates that trigger recursive deduplication.
	b.Run("NestedDeduplication/Enabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(nestedDupAttrs...)
		}
	})

	// Setting nested maps with duplicates with deduplication disabled.
	b.Run("NestedDeduplication/Disabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(nestedDupAttrs...)
		}
	})

	// Setting more attributes than the count limit allows (triggers dropping).
	b.Run("CountLimit/Hit", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = 10 // Less than manyAttrs length.
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(manyAttrs...)
		}
	})

	// Setting attributes within the count limit (no dropping).
	b.Run("CountLimit/NotHit", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = 100 // More than manyAttrs length.
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(manyAttrs...)
		}
	})

	// Setting attributes with long values that trigger truncation.
	b.Run("ValueLimit/Hit", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = 20 // Less than long values.
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(longValueAttrs...)
		}
	})

	// Setting attributes with values within the length limit (no truncation).
	b.Run("ValueLimit/NotHit", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = 200 // More than long values.
			records[i].attributeCountLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(longValueAttrs...)
		}
	})

	// Setting attributes on a record that already has existing attributes (tests overwrite behavior).
	b.Run("Overwrite/Existing", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
			records[i].attributeCountLimit = -1
			// Pre-populate with existing attributes
			records[i].AddAttributes(
				attribute.String("existing1", "value1"),
				attribute.String("existing2", "value2"),
			)
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetAttributes(uniqueAttrs...)
		}
	})
}

func BenchmarkSetBody(b *testing.B) {
	// Simple value (no deduplication or limits).
	simpleValue := attribute.StringValue("simple string value")

	// Map with unique keys (no deduplication needed).
	uniqueMapValue := attribute.MapValue(
		attribute.Bool("bool_key", true),
		attribute.Float64("float_key", 3.14),
		attribute.String("string_key", "value"),
		attribute.Slice("slice_key", attribute.Int64Value(1), attribute.Int64Value(2)),
		attribute.Map("nested_key", attribute.Int("inner", 42)),
		attribute.ByteSlice("bytes_key", []byte("data")),
	)

	// Map with duplicate keys (triggers deduplication).
	dupMapValue := attribute.MapValue(
		attribute.String("key1", "value1"),
		attribute.String("key2", "value2"),
		attribute.String("key1", "duplicate1"),
		attribute.String("key3", "value3"),
		attribute.String("key2", "duplicate2"),
	)

	// Nested map with duplicates.
	nestedDupMapValue := attribute.MapValue(
		attribute.String("outer1", "value1"),
		attribute.Map(
			"nested",
			attribute.String("inner1", "value1"),
			attribute.String("inner2", "value2"),
			attribute.String("inner1", "duplicate"),
		),
		attribute.Slice(
			"slice_with_maps",
			attribute.MapValue(
				attribute.String("slice_key", "value1"),
				attribute.String("slice_key", "duplicate"),
			),
		),
	)

	// Map with long string values (triggers value length limits).
	longValueMapValue := attribute.MapValue(
		attribute.String("short", "short"),
		attribute.String("long1", strings.Repeat("a", 50)),
		attribute.String("long2", strings.Repeat("b", 100)),
	)

	// Setting a simple string value with no limits applied.
	b.Run("Simple/NoLimits", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetBody(simpleValue)
		}
	})

	// Setting a simple string value with limits applied.
	b.Run("Simple/WithLimits", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = 20
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetBody(simpleValue)
		}
	})

	// Setting a map value with unique keys and no limits applied.
	b.Run("UniqueMap/NoLimits", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetBody(uniqueMapValue)
		}
	})

	// Setting a map value with duplicate keys allowed (faster path).
	b.Run("UniqueMap/AllowDuplicates", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetBody(uniqueMapValue)
		}
	})

	// Setting a map with duplicate keys that triggers deduplication logic.
	b.Run("Deduplication/Enabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetBody(dupMapValue)
		}
	})

	// Setting a map with duplicate keys with deduplication disabled.
	b.Run("Deduplication/Disabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetBody(dupMapValue)
		}
	})

	// Setting nested maps with duplicates (tests recursive deduplication).
	b.Run("NestedDeduplication/Enabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetBody(nestedDupMapValue)
		}
	})

	// Setting nested maps with duplicates with deduplication disabled.
	b.Run("NestedDeduplication/Disabled", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
			records[i].attributeValueLengthLimit = -1
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetBody(nestedDupMapValue)
		}
	})

	// Setting map with long string values that trigger truncation.
	b.Run("ValueLimit/Hit", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = 30 // Less than long values.
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetBody(longValueMapValue)
		}
	})

	// Setting map with values within the length limit (no truncation).
	b.Run("ValueLimit/NoHit", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].attributeValueLengthLimit = 200 // More than long values.
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := range b.N {
			records[i].SetBody(longValueMapValue)
		}
	})
}
