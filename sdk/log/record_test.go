// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"
	"strconv"
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
			name: "Bool",
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
			name: "nestedMap",
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
		assert.Equalf(t, i-1, r.DroppedAttributes(), "%d: AddAttributes", i)
		assert.True(t, called, "dropped attributes not logged")

		r.AddAttributes(attrs...)
		assert.Equalf(t, 2*i-1, r.DroppedAttributes(), "%d: second AddAttributes", i)

		r.SetAttributes(attrs...)
		assert.Equalf(t, i-1, r.DroppedAttributes(), "%d: SetAttributes", i)
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
			wantDroppedAttrs: 10,
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

func BenchmarkSetAddAttributes(b *testing.B) {
	kv := log.String("key", "value")

	b.Run("SetAttributes", func(b *testing.B) {
		records := make([]Record, b.N)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			records[i].SetAttributes(kv)
		}
	})

	b.Run("SetAttributes/AllowDuplicates", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			records[i].SetAttributes(kv)
		}
	})

	b.Run("AddAttributes", func(b *testing.B) {
		records := make([]Record, b.N)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			records[i].AddAttributes(kv)
		}
	})

	b.Run("AddAttributes/AllowDuplicates", func(b *testing.B) {
		records := make([]Record, b.N)
		for i := range records {
			records[i].allowDupKeys = true
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			records[i].AddAttributes(kv)
		}
	})
}

func BenchmarkSetBody(b *testing.B) {
	b.Run("SetBody", func(b *testing.B) {
		records := make([]Record, b.N)

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			records[i].SetBody(log.MapValue(
				log.Bool("0", true),
				log.Float64("2", 3.0),
				log.String("3", "forth"),
				log.Slice("4", log.Int64Value(1)),
				log.Map("5", log.Int("key", 2)),
				log.Bytes("6", []byte("six")),
				log.Int64("1", 3),
			))
		}
	})
}
