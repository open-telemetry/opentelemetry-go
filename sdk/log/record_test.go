// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

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
	v := log.BoolValue(true)
	r := new(Record)
	r.SetBody(v)
	assert.True(t, v.Equal(r.Body()))
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
	assert.True(t, res.Equal(&got))
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
	for i := 1; i < attributesInlineCount*5; i++ {
		r := new(Record)
		r.attributeCountLimit = 1

		attrs := make([]log.KeyValue, i)
		attrs[0] = log.Bool("only key different then the rest", true)
		r.AddAttributes(attrs...)
		assert.Equalf(t, i-1, r.DroppedAttributes(), "%d: AddAttributes", i)

		r.SetAttributes(attrs...)
		assert.Equalf(t, i-1, r.DroppedAttributes(), "%d: SetAttributes", i)
	}
}

func TestRecordCompactAttr(t *testing.T) {
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
			r := new(Record)
			r.setAttributes(tc.attrs)
			r.compactAttr()

			var i int
			r.WalkAttributes(func(kv log.KeyValue) bool {
				if assert.Lessf(t, i, len(tc.want), "additional: %v", kv) {
					want := tc.want[i]
					assert.Truef(t, kv.Equal(want), "%d: want %v, got %v", i, want, kv)
				}
				i++
				return true
			})
		})
	}
}
