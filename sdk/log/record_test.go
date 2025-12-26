// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"
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
		body            log.Value
		want            log.Value
	}{
		{
			name: "boolean value",
			body: log.BoolValue(true),
			want: log.BoolValue(true),
		},
		{
			name: "slice",
			body: log.SliceValue(log.BoolValue(true), log.BoolValue(false)),
			want: log.SliceValue(log.BoolValue(true), log.BoolValue(false)),
		},
		{
			name: "map",
			body: log.MapValue(
				log.Bool("0", true),
				log.Int64("1", 2), // This should be removed
				log.Float64("2", 3.0),
				log.String("3", "forth"),
				log.Slice("4", log.Int64Value(1)),
				log.Map("5", log.Int("key", 2)),
				log.Bytes("6", []byte("six")),
				log.Int64("1", 3),
			),
			want: log.MapValue(
				log.Bool("0", true),
				log.Float64("2", 3.0),
				log.String("3", "forth"),
				log.Slice("4", log.Int64Value(1)),
				log.Map("5", log.Int("key", 2)),
				log.Bytes("6", []byte("six")),
				log.Int64("1", 3),
			),
		},
		{
			name: "nested map",
			body: log.MapValue(
				log.Map("key",
					log.Int64("key", 1),
					log.Int64("key", 2),
				),
			),
			want: log.MapValue(
				log.Map("key",
					log.Int64("key", 2),
				),
			),
		},
		{
			name:            "map - allow duplicates",
			allowDuplicates: true,
			body: log.MapValue(
				log.Int64("1", 2),
				log.Int64("1", 3),
			),
			want: log.MapValue(
				log.Int64("1", 2),
				log.Int64("1", 3),
			),
		},
		{
			name: "slice with nested deduplication",
			body: log.SliceValue(
				log.MapValue(log.String("key", "value1"), log.String("key", "value2")),
				log.StringValue("normal"),
				log.SliceValue(
					log.MapValue(log.String("nested", "val1"), log.String("nested", "val2")),
				),
			),
			want: log.SliceValue(
				log.MapValue(log.String("key", "value2")),
				log.StringValue("normal"),
				log.SliceValue(
					log.MapValue(log.String("nested", "val2")),
				),
			),
		},
		{
			name: "empty slice",
			body: log.SliceValue(),
			want: log.SliceValue(),
		},
		{
			name: "empty map",
			body: log.MapValue(),
			want: log.MapValue(),
		},
		{
			name: "single key map",
			body: log.MapValue(log.String("single", "value")),
			want: log.MapValue(log.String("single", "value")),
		},
		{
			name: "slice with no deduplication needed",
			body: log.SliceValue(
				log.StringValue("value1"),
				log.StringValue("value2"),
				log.MapValue(log.String("unique1", "val1")),
				log.MapValue(log.String("unique2", "val2")),
			),
			want: log.SliceValue(
				log.StringValue("value1"),
				log.StringValue("value2"),
				log.MapValue(log.String("unique1", "val1")),
				log.MapValue(log.String("unique2", "val2")),
			),
		},
		{
			name: "deeply nested slice and map structure",
			body: log.SliceValue(
				log.MapValue(
					log.String("outer", "value"),
					log.Slice("inner_slice",
						log.MapValue(log.String("deep", "value1"), log.String("deep", "value2")),
					),
				),
			),
			want: log.SliceValue(
				log.MapValue(
					log.String("outer", "value"),
					log.Slice("inner_slice",
						log.MapValue(log.String("deep", "value2")),
					),
				),
			),
		},
		{
			name:            "slice with duplicates allowed",
			allowDuplicates: true,
			body: log.SliceValue(
				log.MapValue(log.String("key", "value1"), log.String("key", "value2")),
			),
			want: log.SliceValue(
				log.MapValue(log.String("key", "value1"), log.String("key", "value2")),
			),
		},
		{
			name: "string value",
			body: log.StringValue("test"),
			want: log.StringValue("test"),
		},
		{
			name: "boolean value without deduplication",
			body: log.BoolValue(true),
			want: log.BoolValue(true),
		},
		{
			name: "integer value",
			body: log.Int64Value(42),
			want: log.Int64Value(42),
		},
		{
			name: "float value",
			body: log.Float64Value(3.14),
			want: log.Float64Value(3.14),
		},
		{
			name: "bytes value",
			body: log.BytesValue([]byte("test")),
			want: log.BytesValue([]byte("test")),
		},
		{
			name: "empty slice",
			body: log.SliceValue(),
			want: log.SliceValue(),
		},
		{
			name: "slice without nested deduplication",
			body: log.SliceValue(log.StringValue("test"), log.BoolValue(true)),
			want: log.SliceValue(log.StringValue("test"), log.BoolValue(true)),
		},
		{
			name: "slice with nested deduplication needed",
			body: log.SliceValue(log.MapValue(log.String("key", "value1"), log.String("key", "value2"))),
			want: log.SliceValue(log.MapValue(log.String("key", "value2"))),
		},
		{
			name: "empty map",
			body: log.MapValue(),
			want: log.MapValue(),
		},
		{
			name: "single key map",
			body: log.MapValue(log.String("key", "value")),
			want: log.MapValue(log.String("key", "value")),
		},
		{
			name: "map with duplicate keys",
			body: log.MapValue(log.String("key", "value1"), log.String("key", "value2")),
			want: log.MapValue(log.String("key", "value2")),
		},
		{
			name: "map without duplicates",
			body: log.MapValue(log.String("key1", "value1"), log.String("key2", "value2")),
			want: log.MapValue(log.String("key1", "value1"), log.String("key2", "value2")),
		},
		{
			name: "map with nested slice deduplication",
			body: log.MapValue(
				log.Slice("slice", log.MapValue(log.String("nested", "val1"), log.String("nested", "val2"))),
			),
			want: log.MapValue(
				log.Slice("slice", log.MapValue(log.String("nested", "val2"))),
			),
		},
		{
			name: "deeply nested structure with deduplication",
			body: log.SliceValue(
				log.MapValue(
					log.Map("nested",
						log.String("key", "value1"),
						log.String("key", "value2"),
					),
				),
			),
			want: log.SliceValue(
				log.MapValue(
					log.Map("nested",
						log.String("key", "value2"),
					),
				),
			),
		},
		{
			name: "deeply nested structure without deduplication",
			body: log.SliceValue(
				log.MapValue(
					log.Map("nested",
						log.String("key1", "value1"),
						log.String("key2", "value2"),
					),
				),
			),
			want: log.SliceValue(
				log.MapValue(
					log.Map("nested",
						log.String("key1", "value1"),
						log.String("key2", "value2"),
					),
				),
			),
		},
		{
			name: "string value for collection deduplication",
			body: log.StringValue("test"),
			want: log.StringValue("test"),
		},
		{
			name: "boolean value for collection deduplication",
			body: log.BoolValue(true),
			want: log.BoolValue(true),
		},
		{
			name: "empty slice for collection deduplication",
			body: log.SliceValue(),
			want: log.SliceValue(),
		},
		{
			name: "slice without nested deduplication for collection testing",
			body: log.SliceValue(log.StringValue("test"), log.BoolValue(true)),
			want: log.SliceValue(log.StringValue("test"), log.BoolValue(true)),
		},
		{
			name: "slice with nested map requiring deduplication",
			body: log.SliceValue(
				log.MapValue(log.String("key", "value1"), log.String("key", "value2")),
				log.StringValue("normal"),
			),
			want: log.SliceValue(
				log.MapValue(log.String("key", "value2")),
				log.StringValue("normal"),
			),
		},
		{
			name: "deeply nested slice with map deduplication",
			body: log.SliceValue(
				log.SliceValue(
					log.MapValue(log.String("deep", "value1"), log.String("deep", "value2")),
				),
			),
			want: log.SliceValue(
				log.SliceValue(
					log.MapValue(log.String("deep", "value2")),
				),
			),
		},
		{
			name: "empty map for collection deduplication",
			body: log.MapValue(),
			want: log.MapValue(),
		},
		{
			name: "map with nested slice containing duplicates",
			body: log.MapValue(
				log.String("outer", "value"),
				log.Slice("nested_slice",
					log.MapValue(log.String("inner", "val1"), log.String("inner", "val2")),
				),
			),
			want: log.MapValue(
				log.String("outer", "value"),
				log.Slice("nested_slice",
					log.MapValue(log.String("inner", "val2")),
				),
			),
		},
		{
			name: "map with key duplication and nested value deduplication",
			body: log.MapValue(
				log.String("key1", "value1"),
				log.String("key1", "value2"), // key dedup
				log.Slice("slice",
					log.MapValue(log.String("nested", "val1"), log.String("nested", "val2")), // nested value dedup
				),
			),
			want: log.MapValue(
				log.String("key1", "value2"),
				log.Slice("slice",
					log.MapValue(log.String("nested", "val2")),
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
			if !got.Equal(tc.want) {
				t.Errorf("r.Body() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRecordAttributes(t *testing.T) {
	attrs := []log.KeyValue{
		log.Bool("0", true),
		log.Int64("1", 2),
		log.Float64("2", 3.0),
		log.String("3", "forth"),
		log.Slice("4", log.Int64Value(1)),
		log.Map("5", log.Int("key", 2)),
		log.Bytes("6", []byte("six")),
	}
	r := new(Record)
	r.attributeValueLengthLimit = -1
	r.SetAttributes(attrs...)
	r.SetAttributes(attrs[:2]...) // Overwrite existing.
	r.AddAttributes(attrs[2:]...)

	assert.Equal(t, len(attrs), r.AttributesLen(), "attribute length")

	for n := range attrs {
		var i int
		r.WalkAttributes(func(log.KeyValue) bool {
			i++
			return i <= n
		})
		assert.Equalf(t, n+1, i, "WalkAttributes did not stop at %d", n+1)
	}

	var i int
	r.WalkAttributes(func(kv log.KeyValue) bool {
		assert.Truef(t, kv.Equal(attrs[i]), "%d: %v != %v", i, kv, attrs[i])
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
	val0 := log.BoolValue(true)
	attr0 := log.Bool("0", true)
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
	val1 := log.IntValue(1)
	attr1 := log.Int64("1", 2)
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
	assert.True(t, val0.Equal(r0.Body()))
	assert.Equal(t, traceID0, r0.TraceID())
	assert.Equal(t, spanID0, r0.SpanID())
	assert.Equal(t, flag0, r0.TraceFlags())
	r0.WalkAttributes(func(kv log.KeyValue) bool {
		return assert.Truef(t, kv.Equal(attr0), "%v != %v", kv, attr0)
	})

	assert.Equal(t, now1, r1.Timestamp())
	assert.Equal(t, now1, r1.ObservedTimestamp())
	assert.Equal(t, sev1, r1.Severity())
	assert.Equal(t, text1, r1.SeverityText())
	assert.True(t, val1.Equal(r1.Body()))
	assert.Equal(t, traceID1, r1.TraceID())
	assert.Equal(t, spanID1, r1.SpanID())
	assert.Equal(t, flag1, r1.TraceFlags())
	r1.WalkAttributes(func(kv log.KeyValue) bool {
		return assert.Truef(t, kv.Equal(attr1), "%v != %v", kv, attr1)
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

		attrs := make([]log.KeyValue, i)
		attrs[0] = log.Bool("only key different then the rest", true)
		assert.False(t, called, "non-dropped attributed logged")

		r.AddAttributes(attrs...)
		wantDropped := 0
		if i > 1 {
			wantDropped = 1
		}
		assert.Equalf(t, wantDropped, r.DroppedAttributes(), "%d: AddAttributes", i)
		if i > 1 {
			assert.True(t, called, "dropped attributes not logged")
		}

		called = false
		logAttrDropped = func() { called = true }

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
		attrs []log.KeyValue
		want  []log.KeyValue
	}{
		{
			name:  "EmptyKey",
			attrs: make([]log.KeyValue, 10),
			want:  make([]log.KeyValue, 10),
		},
		{
			name: "MapKey",
			attrs: []log.KeyValue{
				log.Map("key", log.Int("key", 5), log.Int("key", 10)),
			},
			want: []log.KeyValue{
				log.Map("key", log.Int("key", 5), log.Int("key", 10)),
			},
		},
		{
			name: "NonEmptyKey",
			attrs: []log.KeyValue{
				log.Bool("key", true),
				log.Int64("key", 1),
				log.Bool("key", false),
				log.Float64("key", 2.),
				log.String("key", "3"),
				log.Slice("key", log.Int64Value(4)),
				log.Map("key", log.Int("key", 5)),
				log.Bytes("key", []byte("six")),
				log.Bool("key", false),
			},
			want: []log.KeyValue{
				log.Bool("key", true),
				log.Int64("key", 1),
				log.Bool("key", false),
				log.Float64("key", 2.),
				log.String("key", "3"),
				log.Slice("key", log.Int64Value(4)),
				log.Map("key", log.Int("key", 5)),
				log.Bytes("key", []byte("six")),
				log.Bool("key", false),
			},
		},
		{
			name: "Multiple",
			attrs: []log.KeyValue{
				log.Bool("a", true),
				log.Int64("b", 1),
				log.Bool("a", false),
				log.Float64("c", 2.),
				log.String("b", "3"),
				log.Slice("d", log.Int64Value(4)),
				log.Map("a", log.Int("key", 5)),
				log.Bytes("d", []byte("six")),
				log.Bool("e", true),
				log.Int("f", 1),
				log.Int("f", 2),
				log.Int("f", 3),
				log.Float64("b", 0.0),
				log.Float64("b", 0.0),
				log.String("g", "G"),
				log.String("h", "H"),
				log.String("g", "GG"),
				log.Bool("a", false),
			},
			want: []log.KeyValue{
				// Order is important here.
				log.Bool("a", true),
				log.Int64("b", 1),
				log.Bool("a", false),
				log.Float64("c", 2.),
				log.String("b", "3"),
				log.Slice("d", log.Int64Value(4)),
				log.Map("a", log.Int("key", 5)),
				log.Bytes("d", []byte("six")),
				log.Bool("e", true),
				log.Int("f", 1),
				log.Int("f", 2),
				log.Int("f", 3),
				log.Float64("b", 0.0),
				log.Float64("b", 0.0),
				log.String("g", "G"),
				log.String("h", "H"),
				log.String("g", "GG"),
				log.Bool("a", false),
			},
		},
		{
			name: "NoDuplicate",
			attrs: func() []log.KeyValue {
				out := make([]log.KeyValue, attributesInlineCount*2)
				for i := range out {
					out[i] = log.Bool(strconv.Itoa(i), true)
				}
				return out
			}(),
			want: func() []log.KeyValue {
				out := make([]log.KeyValue, attributesInlineCount*2)
				for i := range out {
					out[i] = log.Bool(strconv.Itoa(i), true)
				}
				return out
			}(),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			validate := func(t *testing.T, r *Record, want []log.KeyValue) {
				t.Helper()

				var i int
				r.WalkAttributes(func(kv log.KeyValue) bool {
					if assert.Lessf(t, i, len(want), "additional: %v", kv) {
						want := want[i]
						assert.Truef(t, kv.Equal(want), "%d: want %v, got %v", i, want, kv)
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
		attrs []log.KeyValue
		want  []log.KeyValue
	}{
		{
			name:  "EmptyKey",
			attrs: make([]log.KeyValue, 10),
			want:  make([]log.KeyValue, 1),
		},
		{
			name: "NonEmptyKey",
			attrs: []log.KeyValue{
				log.Bool("key", true),
				log.Int64("key", 1),
				log.Bool("key", false),
				log.Float64("key", 2.),
				log.String("key", "3"),
				log.Slice("key", log.Int64Value(4)),
				log.Map("key", log.Int("key", 5)),
				log.Bytes("key", []byte("six")),
				log.Bool("key", false),
			},
			want: []log.KeyValue{
				log.Bool("key", false),
			},
		},
		{
			name: "Multiple",
			attrs: []log.KeyValue{
				log.Bool("a", true),
				log.Int64("b", 1),
				log.Bool("a", false),
				log.Float64("c", 2.),
				log.String("b", "3"),
				log.Slice("d", log.Int64Value(4)),
				log.Map("a", log.Int("key", 5)),
				log.Bytes("d", []byte("six")),
				log.Bool("e", true),
				log.Int("f", 1),
				log.Int("f", 2),
				log.Int("f", 3),
				log.Float64("b", 0.0),
				log.Float64("b", 0.0),
				log.String("g", "G"),
				log.String("h", "H"),
				log.String("g", "GG"),
				log.Bool("a", false),
			},
			want: []log.KeyValue{
				// Order is important here.
				log.Bool("a", false),
				log.Float64("b", 0.0),
				log.Float64("c", 2.),
				log.Bytes("d", []byte("six")),
				log.Bool("e", true),
				log.Int("f", 3),
				log.String("g", "GG"),
				log.String("h", "H"),
			},
		},
		{
			name: "NoDuplicate",
			attrs: func() []log.KeyValue {
				out := make([]log.KeyValue, attributesInlineCount*2)
				for i := range out {
					out[i] = log.Bool(strconv.Itoa(i), true)
				}
				return out
			}(),
			want: func() []log.KeyValue {
				out := make([]log.KeyValue, attributesInlineCount*2)
				for i := range out {
					out[i] = log.Bool(strconv.Itoa(i), true)
				}
				return out
			}(),
		},
		{
			name: "AttributeWithDuplicateKeys",
			attrs: []log.KeyValue{
				log.String("duplicate", "first"),
				log.String("unique", "value"),
				log.String("duplicate", "second"),
			},
			want: []log.KeyValue{
				log.String("duplicate", "second"),
				log.String("unique", "value"),
			},
		},
		{
			name: "ManyDuplicateKeys",
			attrs: []log.KeyValue{
				log.String("key", "value1"),
				log.String("key", "value2"),
				log.String("key", "value3"),
				log.String("key", "value4"),
				log.String("key", "value5"),
			},
			want: []log.KeyValue{
				log.String("key", "value5"),
			},
		},
		{
			name: "InterleavedDuplicates",
			attrs: []log.KeyValue{
				log.String("a", "a1"),
				log.String("b", "b1"),
				log.String("a", "a2"),
				log.String("c", "c1"),
				log.String("b", "b2"),
			},
			want: []log.KeyValue{
				log.String("a", "a2"),
				log.String("b", "b2"),
				log.String("c", "c1"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			validate := func(t *testing.T, r *Record) {
				t.Helper()

				var i int
				r.WalkAttributes(func(kv log.KeyValue) bool {
					if assert.Lessf(t, i, len(tc.want), "additional: %v", kv) {
						want := tc.want[i]
						assert.Truef(t, kv.Equal(want), "%d: want %v, got %v", i, want, kv)
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
		input, want      log.Value
		wantDroppedAttrs int
	}{
		{
			// No de-duplication
			name: "Slice",
			input: log.SliceValue(
				log.BoolValue(true),
				log.BoolValue(true),
				log.Float64Value(1.3),
				log.Float64Value(1.3),
				log.Int64Value(43),
				log.Int64Value(43),
				log.BytesValue([]byte("hello")),
				log.BytesValue([]byte("hello")),
				log.StringValue("foo"),
				log.StringValue("foo"),
				log.SliceValue(log.StringValue("baz")),
				log.SliceValue(log.StringValue("baz")),
				log.MapValue(log.String("a", "qux")),
				log.MapValue(log.String("a", "qux")),
			),
			want: log.SliceValue(
				log.BoolValue(true),
				log.BoolValue(true),
				log.Float64Value(1.3),
				log.Float64Value(1.3),
				log.Int64Value(43),
				log.Int64Value(43),
				log.BytesValue([]byte("hello")),
				log.BytesValue([]byte("hello")),
				log.StringValue("foo"),
				log.StringValue("foo"),
				log.SliceValue(log.StringValue("baz")),
				log.SliceValue(log.StringValue("baz")),
				log.MapValue(log.String("a", "qux")),
				log.MapValue(log.String("a", "qux")),
			),
		},
		{
			name: "Map",
			input: log.MapValue(
				log.Bool("a", true),
				log.Int64("b", 1),
				log.Bool("a", false),
				log.Float64("c", 2.),
				log.String("b", "3"),
				log.Slice("d", log.Int64Value(4)),
				log.Map("a", log.Int("key", 5)),
				log.Bytes("d", []byte("six")),
				log.Bool("e", true),
				log.Int("f", 1),
				log.Int("f", 2),
				log.Int("f", 3),
				log.Float64("b", 0.0),
				log.Float64("b", 0.0),
				log.String("g", "G"),
				log.String("h", "H"),
				log.String("g", "GG"),
				log.Bool("a", false),
			),
			want: log.MapValue(
				// Order is important here.
				log.Bool("a", false),
				log.Float64("b", 0.0),
				log.Float64("c", 2.),
				log.Bytes("d", []byte("six")),
				log.Bool("e", true),
				log.Int("f", 3),
				log.String("g", "GG"),
				log.String("h", "H"),
			),
			wantDroppedAttrs: 0, // Deduplication doesn't count as dropped
		},
		{
			name:             "EmptyMap",
			input:            log.MapValue(),
			want:             log.MapValue(),
			wantDroppedAttrs: 0,
		},
		{
			name:             "SingleKeyMap",
			input:            log.MapValue(log.String("key1", "value1")),
			want:             log.MapValue(log.String("key1", "value1")),
			wantDroppedAttrs: 0,
		},
		{
			name:             "EmptySlice",
			input:            log.SliceValue(),
			want:             log.SliceValue(),
			wantDroppedAttrs: 0,
		},
		{
			name: "SliceWithNestedDedup",
			input: log.SliceValue(
				log.MapValue(log.String("key", "value1"), log.String("key", "value2")),
				log.StringValue("normal"),
			),
			want: log.SliceValue(
				log.MapValue(log.String("key", "value2")),
				log.StringValue("normal"),
			),
			wantDroppedAttrs: 0, // Nested deduplication doesn't count as dropped
		},
		{
			name: "NestedSliceInMap",
			input: log.MapValue(
				log.Slice("slice_key",
					log.MapValue(log.String("nested", "value1"), log.String("nested", "value2")),
				),
			),
			want: log.MapValue(
				log.Slice("slice_key",
					log.MapValue(log.String("nested", "value2")),
				),
			),
			wantDroppedAttrs: 0, // Nested deduplication doesn't count as dropped
		},
		{
			name: "DeeplyNestedStructure",
			input: log.MapValue(
				log.Map("level1",
					log.Map("level2",
						log.Slice("level3",
							log.MapValue(log.String("deep", "value1"), log.String("deep", "value2")),
						),
					),
				),
			),
			want: log.MapValue(
				log.Map("level1",
					log.Map("level2",
						log.Slice("level3",
							log.MapValue(log.String("deep", "value2")),
						),
					),
				),
			),
			wantDroppedAttrs: 0, // Deeply nested deduplication doesn't count as dropped
		},
		{
			name: "NestedMapWithoutDuplicateKeys",
			input: log.SliceValue((log.MapValue(
				log.String("key1", "value1"),
				log.String("key2", "value2"),
			))),
			want: log.SliceValue(log.MapValue(
				log.String("key1", "value1"),
				log.String("key2", "value2"),
			)),
			wantDroppedAttrs: 0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			const key = "key"
			kv := log.KeyValue{Key: key, Value: tc.input}
			r := Record{attributeValueLengthLimit: -1}

			t.Run("AddAttributes", func(t *testing.T) {
				r.AddAttributes(kv)
				assertKV(t, r, log.KeyValue{Key: key, Value: tc.want})
				assert.Equal(t, tc.wantDroppedAttrs, r.DroppedAttributes())
			})

			t.Run("SetAttributes", func(t *testing.T) {
				r.SetAttributes(kv)
				assertKV(t, r, log.KeyValue{Key: key, Value: tc.want})
				assert.Equal(t, tc.wantDroppedAttrs, r.DroppedAttributes())
			})
		})
	}
}

func TestLogKeyValuePairDroppedOnDeduplication(t *testing.T) {
	origKeyValueDropped := logKeyValuePairDropped
	origAttrDropped := logAttrDropped
	t.Cleanup(func() {
		logKeyValuePairDropped = origKeyValueDropped
		logAttrDropped = origAttrDropped
	})

	testCases := []struct {
		name                string
		setupRecord         func() *Record
		operation           func(*Record)
		wantKeyValueDropped bool
		wantAttrDropped     bool
		wantDroppedCount    int
		wantAttributeCount  int
		description         string
	}{
		{
			name: "SetAttributes with duplicate keys",
			setupRecord: func() *Record {
				r := &Record{
					attributeValueLengthLimit: -1,
					allowDupKeys:              false,
				}
				return r
			},
			operation: func(r *Record) {
				r.SetAttributes(
					log.String("key", "value1"),
					log.String("key", "value2"),
				)
			},
			wantKeyValueDropped: true,
			wantAttrDropped:     false,
			wantDroppedCount:    0,
			wantAttributeCount:  1,
			description:         "SetAttributes with duplicate keys should call logKeyValuePairDropped",
		},
		{
			name: "AddAttributes with duplicate keys in new attrs",
			setupRecord: func() *Record {
				r := &Record{
					attributeValueLengthLimit: -1,
					allowDupKeys:              false,
				}
				return r
			},
			operation: func(r *Record) {
				r.AddAttributes(
					log.String("key", "value1"),
					log.String("key", "value2"),
				)
			},
			wantKeyValueDropped: true,
			wantAttrDropped:     false,
			wantDroppedCount:    0,
			wantAttributeCount:  1,
			description:         "AddAttributes with duplicate keys should call logKeyValuePairDropped",
		},
		{
			name: "AddAttributes with duplicate between existing and new",
			setupRecord: func() *Record {
				r := &Record{
					attributeValueLengthLimit: -1,
					allowDupKeys:              false,
				}
				r.SetAttributes(log.String("key1", "value1"))
				return r
			},
			operation: func(r *Record) {
				r.AddAttributes(log.String("key1", "value2"))
			},
			wantKeyValueDropped: true,
			wantAttrDropped:     false,
			wantDroppedCount:    0,
			wantAttributeCount:  1,
			description:         "AddAttributes with duplicate between existing and new should call logKeyValuePairDropped",
		},
		{
			name: "AddAttributes with nested map duplicates",
			setupRecord: func() *Record {
				r := &Record{
					attributeValueLengthLimit: -1,
					allowDupKeys:              false,
				}
				return r
			},
			operation: func(r *Record) {
				r.AddAttributes(
					log.Map("outer",
						log.String("nested", "value1"),
						log.String("nested", "value2"),
					),
				)
			},
			wantKeyValueDropped: true,
			wantAttrDropped:     false,
			wantDroppedCount:    0,
			wantAttributeCount:  1,
			description:         "Nested map duplicates should call logKeyValuePairDropped",
		},
		{
			name: "SetAttributes with limit reached (no duplicates)",
			setupRecord: func() *Record {
				r := &Record{
					attributeValueLengthLimit: -1,
					attributeCountLimit:       2,
					allowDupKeys:              false,
				}
				return r
			},
			operation: func(r *Record) {
				r.SetAttributes(
					log.String("key1", "value1"),
					log.String("key2", "value2"),
					log.String("key3", "value3"),
				)
			},
			wantKeyValueDropped: false,
			wantAttrDropped:     true,
			wantDroppedCount:    1,
			wantAttributeCount:  2,
			description:         "Limit reached without duplicates should call logAttrDropped",
		},
		{
			name: "SetAttributes with both duplicates and limit",
			setupRecord: func() *Record {
				r := &Record{
					attributeValueLengthLimit: -1,
					attributeCountLimit:       2,
					allowDupKeys:              false,
				}
				return r
			},
			operation: func(r *Record) {
				r.SetAttributes(
					log.String("key1", "value1"),
					log.String("key1", "value2"),
					log.String("key2", "value3"),
					log.String("key3", "value4"),
					log.String("key3", "value5"),
				)
			},
			wantKeyValueDropped: true,
			wantAttrDropped:     true,
			wantDroppedCount:    1,
			wantAttributeCount:  2,
			description:         "Both duplicates and limit should call both log functions",
		},
		{
			name: "AddAttributes no duplicates no limit",
			setupRecord: func() *Record {
				r := &Record{
					attributeValueLengthLimit: -1,
					attributeCountLimit:       10,
					allowDupKeys:              false,
				}
				return r
			},
			operation: func(r *Record) {
				r.AddAttributes(
					log.String("key1", "value1"),
					log.String("key2", "value2"),
				)
			},
			wantKeyValueDropped: false,
			wantAttrDropped:     false,
			wantDroppedCount:    0,
			wantAttributeCount:  2,
			description:         "No duplicates and no limit should not call any log function",
		},
		{
			name: "SetAttributes with allowDupKeys enabled",
			setupRecord: func() *Record {
				r := &Record{
					attributeValueLengthLimit: -1,
					allowDupKeys:              true,
				}
				return r
			},
			operation: func(r *Record) {
				r.SetAttributes(
					log.String("key", "value1"),
					log.String("key", "value2"),
				)
			},
			wantKeyValueDropped: false,
			wantAttrDropped:     false,
			wantDroppedCount:    0,
			wantAttributeCount:  2,
			description:         "With allowDupKeys=true, should not call logKeyValuePairDropped",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keyValueDroppedCalled := false
			attrDroppedCalled := false

			logKeyValuePairDropped = sync.OnceFunc(func() {
				keyValueDroppedCalled = true
			})
			logAttrDropped = sync.OnceFunc(func() {
				attrDroppedCalled = true
			})

			r := tc.setupRecord()
			tc.operation(r)

			assert.Equal(t, tc.wantKeyValueDropped, keyValueDroppedCalled,
				"logKeyValuePairDropped call mismatch: %s", tc.description)
			assert.Equal(t, tc.wantAttrDropped, attrDroppedCalled,
				"logAttrDropped call mismatch: %s", tc.description)
			assert.Equal(t, tc.wantDroppedCount, r.DroppedAttributes(),
				"DroppedAttributes count mismatch: %s", tc.description)
			assert.Equal(t, tc.wantAttributeCount, r.AttributesLen(),
				"AttributesLen mismatch: %s", tc.description)
		})
	}
}

func TestDroppedCountExcludesDeduplication(t *testing.T) {
	testCases := []struct {
		name                string
		attrs               []log.KeyValue
		attributeCountLimit int
		wantDropped         int
		wantAttrCount       int
		description         string
	}{
		{
			name: "Only deduplication, no limit",
			attrs: []log.KeyValue{
				log.String("key", "value1"),
				log.String("key", "value2"),
				log.String("key", "value3"),
			},
			attributeCountLimit: 0,
			wantDropped:         0,
			wantAttrCount:       1,
			description:         "Multiple duplicates should not increase dropped count",
		},
		{
			name: "Deduplication and limit",
			attrs: []log.KeyValue{
				log.String("a", "value1"),
				log.String("a", "value2"),
				log.String("b", "value3"),
				log.String("c", "value4"),
				log.String("c", "value5"),
			},
			attributeCountLimit: 2,
			wantDropped:         1, // Only limit drops count (3 unique -> 2 kept)
			wantAttrCount:       2,
			description:         "Deduplication then limit: only limit drops should count",
		},
		{
			name: "No deduplication, only limit",
			attrs: []log.KeyValue{
				log.String("a", "value1"),
				log.String("b", "value2"),
				log.String("c", "value3"),
				log.String("d", "value4"),
			},
			attributeCountLimit: 2,
			wantDropped:         2, // 2 attributes dropped due to limit
			wantAttrCount:       2,
			description:         "Only limit without deduplication",
		},
		{
			name: "Complex: multiple duplicates and limit",
			attrs: []log.KeyValue{
				log.String("a", "a1"),
				log.String("b", "b1"),
				log.String("a", "a2"),
				log.String("c", "c1"),
				log.String("b", "b2"),
				log.String("d", "d1"),
				log.String("a", "a3"),
				log.String("e", "e1"),
			},
			attributeCountLimit: 3,
			wantDropped:         2, // 5 unique keys, limit 3, so 2 dropped
			wantAttrCount:       3,
			description:         "Complex scenario with multiple duplicates and limit",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &Record{
				attributeValueLengthLimit: -1,
				attributeCountLimit:       tc.attributeCountLimit,
				allowDupKeys:              false,
			}

			r.SetAttributes(tc.attrs...)

			assert.Equal(t, tc.wantDropped, r.DroppedAttributes(),
				"%s: dropped count mismatch", tc.description)
			assert.Equal(t, tc.wantAttrCount, r.AttributesLen(),
				"%s: attribute count mismatch", tc.description)
		})
	}
}

func TestSetDroppedAndAddDroppedWithZero(t *testing.T) {
	origAttrDropped := logAttrDropped
	t.Cleanup(func() {
		logAttrDropped = origAttrDropped
	})

	t.Run("setDropped(0) does not call log", func(t *testing.T) {
		called := false
		logAttrDropped = sync.OnceFunc(func() {
			called = true
		})

		r := &Record{}
		r.setDropped(0)

		assert.False(t, called, "logAttrDropped should not be called for setDropped(0)")
		assert.Equal(t, 0, r.DroppedAttributes())
	})

	t.Run("setDropped(n>0) calls log", func(t *testing.T) {
		called := false
		logAttrDropped = sync.OnceFunc(func() {
			called = true
		})

		r := &Record{}
		r.setDropped(5)

		assert.True(t, called, "logAttrDropped should be called for setDropped(n>0)")
		assert.Equal(t, 5, r.DroppedAttributes())
	})

	t.Run("addDropped(0) does not call log", func(t *testing.T) {
		called := false
		logAttrDropped = sync.OnceFunc(func() {
			called = true
		})

		r := &Record{}
		r.addDropped(0)

		assert.False(t, called, "logAttrDropped should not be called for addDropped(0)")
		assert.Equal(t, 0, r.DroppedAttributes())
	})

	t.Run("addDropped(n>0) calls log", func(t *testing.T) {
		called := false
		logAttrDropped = sync.OnceFunc(func() {
			called = true
		})

		r := &Record{}
		r.addDropped(3)

		assert.True(t, called, "logAttrDropped should be called for addDropped(n>0)")
		assert.Equal(t, 3, r.DroppedAttributes())
	})

	t.Run("multiple addDropped accumulates", func(t *testing.T) {
		called := false
		logAttrDropped = sync.OnceFunc(func() {
			called = true
		})

		r := &Record{}
		r.addDropped(2)
		assert.True(t, called, "first addDropped should call log")

		// Reset for second call
		called = false
		logAttrDropped = sync.OnceFunc(func() {
			called = true
		})

		r.addDropped(3)
		assert.True(t, called, "second addDropped should call log")
		assert.Equal(t, 5, r.DroppedAttributes(), "dropped count should accumulate")
	})
}

func TestApplyAttrLimitsTruncation(t *testing.T) {
	testcases := []struct {
		name        string
		limit       int
		input, want log.Value
	}{
		{
			name:  "Empty",
			limit: 0,
			input: log.Value{},
			want:  log.Value{},
		},
		{
			name:  "Bool",
			limit: 0,
			input: log.BoolValue(true),
			want:  log.BoolValue(true),
		},
		{
			name:  "Float64",
			limit: 0,
			input: log.Float64Value(1.3),
			want:  log.Float64Value(1.3),
		},
		{
			name:  "Int64",
			limit: 0,
			input: log.Int64Value(43),
			want:  log.Int64Value(43),
		},
		{
			name:  "Bytes",
			limit: 0,
			input: log.BytesValue([]byte("foo")),
			want:  log.BytesValue([]byte("foo")),
		},
		{
			name:  "String",
			limit: 0,
			input: log.StringValue("foo"),
			want:  log.StringValue(""),
		},
		{
			name:  "Slice",
			limit: 0,
			input: log.SliceValue(
				log.BoolValue(true),
				log.Float64Value(1.3),
				log.Int64Value(43),
				log.BytesValue([]byte("hello")),
				log.StringValue("foo"),
				log.StringValue("bar"),
				log.SliceValue(log.StringValue("baz")),
				log.MapValue(log.String("a", "qux")),
			),
			want: log.SliceValue(
				log.BoolValue(true),
				log.Float64Value(1.3),
				log.Int64Value(43),
				log.BytesValue([]byte("hello")),
				log.StringValue(""),
				log.StringValue(""),
				log.SliceValue(log.StringValue("")),
				log.MapValue(log.String("a", "")),
			),
		},
		{
			name:  "Map",
			limit: 0,
			input: log.MapValue(
				log.Bool("0", true),
				log.Float64("1", 1.3),
				log.Int64("2", 43),
				log.Bytes("3", []byte("hello")),
				log.String("4", "foo"),
				log.String("5", "bar"),
				log.Slice("6", log.StringValue("baz")),
				log.Map("7", log.String("a", "qux")),
			),
			want: log.MapValue(
				log.Bool("0", true),
				log.Float64("1", 1.3),
				log.Int64("2", 43),
				log.Bytes("3", []byte("hello")),
				log.String("4", ""),
				log.String("5", ""),
				log.Slice("6", log.StringValue("")),
				log.Map("7", log.String("a", "")),
			),
		},
		{
			name:  "LongStringTruncated",
			limit: 5,
			input: log.StringValue("This is a very long string that should be truncated"),
			want:  log.StringValue("This "),
		},
		{
			name:  "LongBytesNotTruncated",
			limit: 5,
			input: log.BytesValue([]byte("This is a very long byte array")),
			want:  log.BytesValue([]byte("This is a very long byte array")),
		},
		{
			name:  "TruncationInNestedMap",
			limit: 3,
			input: log.MapValue(
				log.String("short", "ok"),
				log.String("long", "toolong"),
			),
			want: log.MapValue(
				log.String("short", "ok"),
				log.String("long", "too"),
			),
		},
		{
			name:  "TruncationInNestedSlice",
			limit: 4,
			input: log.SliceValue(
				log.StringValue("good"),
				log.StringValue("toolong"),
			),
			want: log.SliceValue(
				log.StringValue("good"),
				log.StringValue("tool"),
			),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			const key = "key"
			kv := log.KeyValue{Key: key, Value: tc.input}
			r := Record{attributeValueLengthLimit: tc.limit}

			t.Run("AddAttributes", func(t *testing.T) {
				r.AddAttributes(kv)
				assertKV(t, r, log.KeyValue{Key: key, Value: tc.want})
			})

			t.Run("SetAttributes", func(t *testing.T) {
				r.SetAttributes(kv)
				assertKV(t, r, log.KeyValue{Key: key, Value: tc.want})
			})
		})
	}
}

func assertKV(t *testing.T, r Record, kv log.KeyValue) {
	t.Helper()

	var kvs []log.KeyValue
	r.WalkAttributes(func(kv log.KeyValue) bool {
		kvs = append(kvs, kv)
		return true
	})

	require.Len(t, kvs, 1)
	assert.Truef(t, kv.Equal(kvs[0]), "%s != %s", kv, kvs[0])
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
	attrs := []log.KeyValue{
		log.String("attr1", "very long value that will be truncated"),
		log.String("attr2", "another very long value that will be truncated"),
		log.String("attr3", "yet another very long value that will be truncated"),
		log.String("attr4", "more very long value that will be truncated"),
		log.String("attr5", "extra very long value that will be truncated"),
		log.String("attr6", "additional very long value that will be truncated"),
		log.String("attr7", "more additional very long value that will be truncated"),
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
	var gotAttrs []log.KeyValue
	r.WalkAttributes(func(kv log.KeyValue) bool {
		gotAttrs = append(gotAttrs, kv)
		return true
	})
	wantAttr := []log.KeyValue{
		log.String("attr1", "very long value that"),
		log.String("attr2", "another very long va"),
		log.String("attr3", "yet another very lon"),
		log.String("attr4", "more very long value"),
		log.String("attr5", "extra very long valu"),
		log.String("attr6", "additional very long"),
		log.String("attr7", "more additional very"),
	}
	if !slices.EqualFunc(gotAttrs, wantAttr, func(a, b log.KeyValue) bool { return a.Equal(b) }) {
		t.Errorf("Attributes do not match.\ngot:\n%v\nwant:\n%v", printKVs(gotAttrs), printKVs(wantAttr))
	}
}

func TestRecordMethodsInputConcurrentSafe(t *testing.T) {
	nestedSlice := log.Slice("nested_slice",
		log.SliceValue(log.StringValue("nested_inner1"), log.StringValue("nested_inner2")),
		log.StringValue("nested_outer"),
	)

	nestedMap := log.Map("nested_map",
		log.String("nested_key1", "nested_value1"),
		log.Map("nested_map", log.String("nested_inner_key", "nested_inner_value")),
		log.String("nested_key1", "duplicate"), // This will trigger dedup.
	)

	dedupAttributes := []log.KeyValue{
		log.String("dedup_key1", "dedup_value1"),
		log.String("dedup_key2", "dedup_value2"),
		log.String("dedup_key1", "duplicate"),    // This will trigger the dedup.
		log.String("dedup_key3", "dedup_value3"), // This will trigger attr count limit.
	}

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := &Record{
				attributeValueLengthLimit: 10,
				attributeCountLimit:       4,
				allowDupKeys:              false,
			}

			r.SetAttributes(nestedSlice)
			r.AddAttributes(nestedMap)
			r.AddAttributes(dedupAttributes...)
			r.SetBody(nestedMap.Value)

			var gotAttrs []log.KeyValue
			r.WalkAttributes(func(kv log.KeyValue) bool {
				gotAttrs = append(gotAttrs, kv)
				return true
			})
			wantAttr := []log.KeyValue{
				log.Slice("nested_slice",
					log.SliceValue(log.StringValue("nested_inn"), log.StringValue("nested_inn")),
					log.StringValue("nested_out"),
				),
				log.Map("nested_map",
					log.String("nested_key1", "duplicate"),
					log.Map("nested_map", log.String("nested_inner_key", "nested_inn")),
				),
				log.String("dedup_key1", "duplicate"),
				log.String("dedup_key2", "dedup_valu"),
			}
			if !slices.EqualFunc(gotAttrs, wantAttr, func(a, b log.KeyValue) bool { return a.Equal(b) }) {
				t.Errorf("Attributes do not match.\ngot:\n%v\nwant:\n%v", printKVs(gotAttrs), printKVs(wantAttr))
			}

			gotBody := r.Body()
			wantBody := log.MapValue(
				log.String("nested_key1", "duplicate"),
				log.Map("nested_map", log.String("nested_inner_key", "nested_inner_value")),
			)
			if !gotBody.Equal(wantBody) {
				t.Errorf("Body does not match.\ngot:\n%v\nwant:\n%v", gotBody, wantBody)
			}
		}()
	}

	wg.Wait()
}

func printKVs(kvs []log.KeyValue) string {
	var sb strings.Builder
	for _, kv := range kvs {
		_, _ = sb.WriteString(fmt.Sprintf("%s: %s\n", kv.Key, kv.Value))
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
					log.String(fmt.Sprintf("key-%d", tt.attrCount), "value"),
				)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				record.WalkAttributes(func(log.KeyValue) bool {
					return true
				})
			}
		})
	}
}

func BenchmarkAddAttributes(b *testing.B) {
	// Simple attribute (no deduplication or limits).
	singleKV := log.String("key", "value")

	// Attributes with no duplicates.
	uniqueAttrs := []log.KeyValue{
		log.String("key1", "value1"),
		log.String("key2", "value2"),
		log.String("key3", "value3"),
		log.String("key4", "value4"),
		log.String("key5", "value5"),
	}

	// Attributes with duplicates that trigger deduplication.
	dupAttrs := []log.KeyValue{
		log.String("key1", "value1"),
		log.String("key2", "value2"),
		log.String("key1", "duplicate1"), // duplicate key
		log.String("key3", "value3"),
		log.String("key2", "duplicate2"), // duplicate key
	}

	// Large number of attributes to trigger count limits.
	manyAttrs := make([]log.KeyValue, 20)
	for i := range manyAttrs {
		manyAttrs[i] = log.String(fmt.Sprintf("key%d", i), "value")
	}

	// Attributes with long values to trigger value length limits.
	longValueAttrs := []log.KeyValue{
		log.String("short", "short"),
		log.String("long1", strings.Repeat("a", 50)),
		log.String("long2", strings.Repeat("b", 100)),
	}

	// Attributes with nested maps that have duplicates (triggers recursive deduplication).
	nestedDupAttrs := []log.KeyValue{
		log.String("simple", "value"),
		log.Map("map1",
			log.String("inner1", "value1"),
			log.String("inner2", "value2"),
			log.String("inner1", "duplicate"), // duplicate in nested map
		),
		log.Map("map2",
			log.String("key", "original"),
			log.Map("deeply_nested",
				log.String("deep1", "value1"),
				log.String("deep2", "value2"),
				log.String("deep1", "duplicate_deep"), // duplicate in deeply nested map
			),
			log.String("key", "overwrite"), // duplicate key at this level
		),
		log.Slice("slice_with_maps",
			log.MapValue(
				log.String("slice_key", "value1"),
				log.String("slice_key", "duplicate"), // duplicate in slice element
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
	singleKV := log.String("key", "value")

	// Attributes with no duplicates.
	uniqueAttrs := []log.KeyValue{
		log.String("key1", "value1"),
		log.String("key2", "value2"),
		log.String("key3", "value3"),
		log.String("key4", "value4"),
		log.String("key5", "value5"),
	}

	// Attributes with duplicates that trigger deduplication.
	dupAttrs := []log.KeyValue{
		log.String("key1", "value1"),
		log.String("key2", "value2"),
		log.String("key1", "duplicate1"), // duplicate key
		log.String("key3", "value3"),
		log.String("key2", "duplicate2"), // duplicate key
	}

	// Large number of attributes to trigger count limits.
	manyAttrs := make([]log.KeyValue, 20)
	for i := range manyAttrs {
		manyAttrs[i] = log.String(fmt.Sprintf("key%d", i), "value")
	}

	// Attributes with long values to trigger value length limits.
	longValueAttrs := []log.KeyValue{
		log.String("short", "short"),
		log.String("long1", strings.Repeat("a", 50)),
		log.String("long2", strings.Repeat("b", 100)),
	}

	// Attributes with nested maps that have duplicates (triggers recursive deduplication).
	nestedDupAttrs := []log.KeyValue{
		log.String("simple", "value"),
		log.Map("map1",
			log.String("inner1", "value1"),
			log.String("inner2", "value2"),
			log.String("inner1", "duplicate"), // duplicate in nested map
		),
		log.Map("map2",
			log.String("key", "original"),
			log.Map("deeply_nested",
				log.String("deep1", "value1"),
				log.String("deep2", "value2"),
				log.String("deep1", "duplicate_deep"), // duplicate in deeply nested map
			),
			log.String("key", "overwrite"), // duplicate key at this level
		),
		log.Slice("slice_with_maps",
			log.MapValue(
				log.String("slice_key", "value1"),
				log.String("slice_key", "duplicate"), // duplicate in slice element
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
				log.String("existing1", "value1"),
				log.String("existing2", "value2"),
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
	simpleValue := log.StringValue("simple string value")

	// Map with unique keys (no deduplication needed).
	uniqueMapValue := log.MapValue(
		log.Bool("bool_key", true),
		log.Float64("float_key", 3.14),
		log.String("string_key", "value"),
		log.Slice("slice_key", log.Int64Value(1), log.Int64Value(2)),
		log.Map("nested_key", log.Int("inner", 42)),
		log.Bytes("bytes_key", []byte("data")),
	)

	// Map with duplicate keys (triggers deduplication).
	dupMapValue := log.MapValue(
		log.String("key1", "value1"),
		log.String("key2", "value2"),
		log.String("key1", "duplicate1"), // duplicate key
		log.String("key3", "value3"),
		log.String("key2", "duplicate2"), // duplicate key
	)

	// Nested map with duplicates.
	nestedDupMapValue := log.MapValue(
		log.String("outer1", "value1"),
		log.Map("nested",
			log.String("inner1", "value1"),
			log.String("inner2", "value2"),
			log.String("inner1", "duplicate"), // duplicate in nested map
		),
		log.Slice("slice_with_maps",
			log.MapValue(
				log.String("slice_key", "value1"),
				log.String("slice_key", "duplicate"), // duplicate in slice element
			),
		),
	)

	// Map with long string values (triggers value length limits).
	longValueMapValue := log.MapValue(
		log.String("short", "short"),
		log.String("long1", strings.Repeat("a", 50)),
		log.String("long2", strings.Repeat("b", 100)),
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
