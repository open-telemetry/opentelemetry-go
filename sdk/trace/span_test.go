// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.45.0"
	"go.opentelemetry.io/otel/trace"
)

func newRecordingSpanWithAttributes(tb testing.TB, count int) *recordingSpan {
	tb.Helper()

	limits := NewSpanLimits()
	limits.AttributeCountLimit = count
	tp := NewTracerProvider(WithSpanLimits(limits))
	_, span := tp.Tracer(tb.Name()).Start(tb.Context(), "span")
	tb.Cleanup(func() { span.End() })

	attrs := make([]attribute.KeyValue, count)
	for i := range attrs {
		attrs[i] = attribute.Int(fmt.Sprintf("key-%d", i), i)
	}
	span.SetAttributes(attrs...)
	return span.(*recordingSpan)
}

func runConcurrently(fns ...func()) {
	var ready, done sync.WaitGroup
	ready.Add(len(fns))
	done.Add(len(fns))
	start := make(chan struct{})
	for _, fn := range fns {
		go func() {
			defer done.Done()
			ready.Done()
			<-start
			fn()
		}()
	}
	ready.Wait()
	close(start)
	done.Wait()
}

func TestRecordingSpanConcurrentSafeAttributes(t *testing.T) {
	const attrCount = 8_192
	s := newRecordingSpanWithAttributes(t, attrCount)
	attrs := s.Attributes()

	var total int
	runConcurrently(
		func() {
			for range 100 {
				for _, attr := range attrs {
					total += len(attr.Key)
				}
			}
		},
		func() {
			for range 100 {
				_ = s.Attributes()
			}
		},
	)
	assert.Positive(t, total)
}

func TestRecordingSpanConcurrentSafeSetAttributes(t *testing.T) {
	const attrCount = 8_192
	s := newRecordingSpanWithAttributes(t, attrCount)
	attrs := s.Attributes()
	updates := make([]attribute.KeyValue, attrCount)
	for i := range updates {
		updates[i] = attribute.Int(fmt.Sprintf("key-%d", i), i+1)
	}

	var total int
	runConcurrently(
		func() {
			for range 100 {
				for _, attr := range attrs {
					total += len(attr.Key)
				}
			}
		},
		func() {
			for range 100 {
				s.SetAttributes(updates...)
			}
		},
	)
	assert.Positive(t, total)
}

func TestRecordingSpanAttributeViewStable(t *testing.T) {
	s := newRecordingSpanWithAttributes(t, 1)
	before := s.Attributes()

	s.SetAttributes(attribute.Int("key-0", 1))
	after := s.Attributes()

	assert.Equal(t, int64(0), before[0].Value.AsInt64())
	assert.Equal(t, int64(1), after[0].Value.AsInt64())
}

func TestRecordingSpanAttributeViewStableWithSpareCapacity(t *testing.T) {
	limits := NewSpanLimits()
	limits.AttributeCountLimit = 2
	tp := NewTracerProvider(WithSpanLimits(limits))
	_, span := tp.Tracer(t.Name()).Start(t.Context(), "span")
	t.Cleanup(func() { span.End() })
	s := span.(*recordingSpan)
	s.SetAttributes(attribute.Int("key", 0))
	before := s.Attributes()

	s.SetAttributes(attribute.Int("key", 1))
	after := s.Attributes()
	afterSnapshot := s.snapshot().Attributes()

	require.Len(t, before, 1)
	assert.Equal(t, int64(0), before[0].Value.AsInt64())
	require.Len(t, after, 1)
	assert.Equal(t, int64(1), after[0].Value.AsInt64())
	require.Len(t, afterSnapshot, 1)
	assert.Equal(t, int64(1), afterSnapshot[0].Value.AsInt64())
}

func TestRecordingSpanEmptyAttributeViewStable(t *testing.T) {
	limits := NewSpanLimits()
	limits.AttributeCountLimit = 1
	tp := NewTracerProvider(WithSpanLimits(limits))
	_, span := tp.Tracer(t.Name()).Start(t.Context(), "span")
	t.Cleanup(func() { span.End() })
	s := span.(*recordingSpan)
	s.SetAttributes(attribute.KeyValue{})
	before := s.Attributes()
	require.Positive(t, cap(before))
	before = before[:1]

	s.SetAttributes(attribute.Int("key", 1))

	assert.False(t, before[0].Valid())
}

func TestRecordingSpanSnapshotAttributesStable(t *testing.T) {
	s := newRecordingSpanWithAttributes(t, 1)
	before := s.snapshot()

	s.SetAttributes(attribute.Int("key-0", 1))
	after := s.snapshot()

	assert.Equal(t, int64(0), before.Attributes()[0].Value.AsInt64())
	assert.Equal(t, int64(1), after.Attributes()[0].Value.AsInt64())
}

func TestRecordingSpanAndEndedSnapshotConcurrentSafeAttributes(t *testing.T) {
	const attrCount = 8_192
	exporter := NewTestExporter()
	limits := NewSpanLimits()
	limits.AttributeCountLimit = attrCount
	tp := NewTracerProvider(WithSyncer(exporter), WithSpanLimits(limits))
	_, span := tp.Tracer(t.Name()).Start(t.Context(), "span")
	s := span.(*recordingSpan)
	attrs := make([]attribute.KeyValue, attrCount)
	for i := range attrs {
		attrs[i] = attribute.Int(fmt.Sprintf("key-%d", i), i)
	}
	s.SetAttributes(attrs...)
	s.End()

	ended, ok := exporter.GetSpan("span")
	if !ok {
		t.Fatal("ended span not exported")
	}
	endedAttrs := ended.Attributes()
	var total int
	runConcurrently(
		func() {
			for range 100 {
				for _, attr := range endedAttrs {
					total += len(attr.Key)
				}
			}
		},
		func() {
			for range 100 {
				_ = s.Attributes()
			}
		},
	)
	assert.Positive(t, total)
}

func TestSetStatus(t *testing.T) {
	tests := []struct {
		name        string
		span        recordingSpan
		code        codes.Code
		description string
		expected    Status
	}{
		{
			"Error and description should overwrite Unset",
			recordingSpan{},
			codes.Error,
			"description",
			Status{Code: codes.Error, Description: "description"},
		},
		{
			"Ok should overwrite Unset and ignore description",
			recordingSpan{},
			codes.Ok,
			"description",
			Status{Code: codes.Ok},
		},
		{
			"Error and description should return error and overwrite description",
			recordingSpan{status: Status{Code: codes.Error, Description: "d1"}},
			codes.Error,
			"d2",
			Status{Code: codes.Error, Description: "d2"},
		},
		{
			"Ok should overwrite error and remove description",
			recordingSpan{status: Status{Code: codes.Error, Description: "d1"}},
			codes.Ok,
			"d2",
			Status{Code: codes.Ok},
		},
		{
			"Error and description should be ignored when already Ok",
			recordingSpan{status: Status{Code: codes.Ok}},
			codes.Error,
			"d2",
			Status{Code: codes.Ok},
		},
		{
			"Ok should be noop when already Ok",
			recordingSpan{status: Status{Code: codes.Ok}},
			codes.Ok,
			"d2",
			Status{Code: codes.Ok},
		},
		{
			"Unset should be noop when already Ok",
			recordingSpan{status: Status{Code: codes.Ok}},
			codes.Unset,
			"d2",
			Status{Code: codes.Ok},
		},
		{
			"Unset should be noop when already Error",
			recordingSpan{status: Status{Code: codes.Error, Description: "d1"}},
			codes.Unset,
			"d2",
			Status{Code: codes.Error, Description: "d1"},
		},
	}

	for i := range tests {
		tc := &tests[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.span.SetStatus(tc.code, tc.description)
			assert.Equal(t, tc.expected, tc.span.status)
		})
	}
}

func TestLogDropAttrs(t *testing.T) {
	orig := logDropAttrs
	t.Cleanup(func() { logDropAttrs = orig })

	var called bool
	logDropAttrs = func() { called = true }

	s := &recordingSpan{}
	s.addDroppedAttr(1)
	assert.True(t, called, "logDropAttrs not called")

	called = false
	s.addDroppedAttr(1)
	assert.False(t, called, "logDropAttrs called multiple times for same Span")
}

func BenchmarkRecordingSpanSetAttributes(b *testing.B) {
	var attrs []attribute.KeyValue
	for i := range 100 {
		attr := attribute.String(fmt.Sprintf("hello.attrib%d", i), fmt.Sprintf("goodbye.attrib%d", i))
		attrs = append(attrs, attr)
	}

	ctx := b.Context()
	for _, limit := range []bool{false, true} {
		b.Run(fmt.Sprintf("WithLimit/%t", limit), func(b *testing.B) {
			b.ReportAllocs()
			sl := NewSpanLimits()
			if limit {
				sl.AttributeCountLimit = 50
			}
			tp := NewTracerProvider(WithSampler(AlwaysSample()), WithSpanLimits(sl))
			tracer := tp.Tracer("tracer")

			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_, span := tracer.Start(ctx, "span")
				span.SetAttributes(attrs...)
				span.End()
			}
		})
	}
}

var benchmarkRecordingSpanAttributes []attribute.KeyValue

func BenchmarkRecordingSpanAttributes(b *testing.B) {
	for _, count := range []int{1, 8, 32, 128} {
		b.Run(fmt.Sprintf("Repeated/%d", count), func(b *testing.B) {
			s := newRecordingSpanWithAttributes(b, count)
			_ = s.Attributes()

			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				benchmarkRecordingSpanAttributes = s.Attributes()
			}
		})

		b.Run(fmt.Sprintf("ReadAfterSet/%d", count), func(b *testing.B) {
			s := newRecordingSpanWithAttributes(b, count)
			updates := make([]attribute.KeyValue, count)
			for i := range updates {
				updates[i] = attribute.Int(fmt.Sprintf("key-%d", i), i+1)
			}
			_ = s.Attributes()

			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				s.SetAttributes(updates...)
				benchmarkRecordingSpanAttributes = s.Attributes()
			}
		})
	}
}

func BenchmarkSpanEnd(b *testing.B) {
	cases := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "Default",
		},
		{
			name: "ObservabilityEnabled",
			env: map[string]string{
				"OTEL_GO_X_OBSERVABILITY": "True",
			},
		},
	}

	ctx := trace.ContextWithSpanContext(b.Context(), trace.SpanContext{})

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for k, v := range c.env {
				b.Setenv(k, v)
			}

			tracer := NewTracerProvider().Tracer("")

			spans := make([]trace.Span, b.N)
			for i := 0; i < b.N; i++ {
				_, span := tracer.Start(ctx, "")
				spans[i] = span
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				spans[i].End()
			}
		})
	}
}

func BenchmarkRecordingSpanRecordError(b *testing.B) {
	for _, tc := range []struct {
		name    string
		workers int
		ended   bool
	}{
		{name: "Active"},
		{name: "Ended", ended: true},
		{name: "Contended", workers: 8},
	} {
		b.Run(tc.name, func(b *testing.B) {
			span := startSpan(NewTracerProvider(), b.Name()).(*recordingSpan)
			if tc.ended {
				span.End()
			}
			err := errors.New("benchmark error")

			b.ReportAllocs()
			b.ResetTimer()
			if tc.workers == 0 {
				for b.Loop() {
					span.RecordError(err)
				}
				return
			}

			var wg sync.WaitGroup
			wg.Add(tc.workers)
			for i := range tc.workers {
				start := i * b.N / tc.workers
				end := (i + 1) * b.N / tc.workers
				go func(start, end int) {
					defer wg.Done()
					for i := start; i < end; i++ {
						span.RecordError(err)
					}
				}(start, end)
			}
			wg.Wait()
		})
	}
}

func nanFloat64SliceAttr() attribute.KeyValue {
	return attribute.Float64Slice("nan_slice", []float64{math.NaN(), 1.5})
}

func assertNaNFloat64Slice(t *testing.T, kv attribute.KeyValue, rest ...float64) {
	t.Helper()
	assertNaNFloat64SliceValue(t, kv.Value, rest...)
}

func assertNaNFloat64SliceValue(t *testing.T, v attribute.Value, rest ...float64) {
	t.Helper()
	got := v.AsFloat64Slice()
	require.Len(t, got, 1+len(rest))
	assert.Truef(t, math.IsNaN(got[0]), "got %v", got)
	for i, want := range rest {
		assert.Equal(t, want, got[i+1])
	}
}

func assertNormAttrsUnchanged(t *testing.T, kvs []attribute.KeyValue, depthLimit, lengthLimit int) {
	t.Helper()
	require.NotEmpty(t, kvs)

	var out []attribute.KeyValue
	allocs := testing.AllocsPerRun(1000, func() {
		out = normAttrs(kvs, depthLimit, lengthLimit)
	})
	assert.Zero(t, allocs)
	require.Len(t, out, len(kvs))
	assert.Same(t, &kvs[0], &out[0], "normAttrs copied an unchanged attribute list")
}

func TestNormAttrsNaNFloat64SliceDoesNotAllocate(t *testing.T) {
	const lengthLimit = 16
	depthLimit := DefaultAttributeValueDepthLimit

	t.Run("TopLevel", func(t *testing.T) {
		kvs := []attribute.KeyValue{nanFloat64SliceAttr()}
		_, changed := normAttr(kvs[0], depthLimit, lengthLimit)
		assert.False(t, changed)
		assertNormAttrsUnchanged(t, kvs, depthLimit, lengthLimit)
		assertNaNFloat64Slice(t, kvs[0], 1.5)
	})

	t.Run("WideSlice", func(t *testing.T) {
		kvs := []attribute.KeyValue{
			attribute.Float64Slice("nan_slice", []float64{math.NaN(), 1, 2, 3, 4, 5}),
		}
		assertNormAttrsUnchanged(t, kvs, depthLimit, lengthLimit)
		assertNaNFloat64Slice(t, kvs[0], 1, 2, 3, 4, 5)
	})

	t.Run("WithPrimitiveArrays", func(t *testing.T) {
		kvs := []attribute.KeyValue{
			attribute.Bool("ok", true),
			nanFloat64SliceAttr(),
			attribute.Int64Slice("ints", []int64{1, 2, 3, 4, 5, 6}),
			attribute.BoolSlice("flags", []bool{true, false, true, false}),
		}
		assertNormAttrsUnchanged(t, kvs, depthLimit, lengthLimit)
		assertNaNFloat64Slice(t, kvs[1], 1.5)
	})
}

func TestNormAttrsNestedNaNFloat64SliceDoesNotAllocate(t *testing.T) {
	const lengthLimit = 16
	depthLimit := DefaultAttributeValueDepthLimit
	nan := attribute.Float64SliceValue([]float64{math.NaN()})

	tests := []struct {
		name string
		kv   attribute.KeyValue
	}{
		{
			name: "Slice",
			kv:   attribute.Slice("nested", nan),
		},
		{
			name: "Map",
			kv:   attribute.Map("nested", attribute.Float64Slice("nan", []float64{math.NaN()})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kvs := []attribute.KeyValue{tt.kv}
			_, changed := normAttr(tt.kv, depthLimit, lengthLimit)
			assert.False(t, changed)
			assertNormAttrsUnchanged(t, kvs, depthLimit, lengthLimit)

			switch tt.kv.Value.Type() {
			case attribute.SLICE:
				elems := kvs[0].Value.AsSlice()
				require.Len(t, elems, 1)
				assertNaNFloat64SliceValue(t, elems[0])
			case attribute.MAP:
				entries := kvs[0].Value.AsMap()
				require.Len(t, entries, 1)
				assertNaNFloat64SliceValue(t, entries[0].Value)
			}
		})
	}
}

func TestNormAttrsNaNFloat64SliceCopyOnWrite(t *testing.T) {
	const lengthLimit = 16
	long := strings.Repeat("a", lengthLimit+4)

	t.Run("OverLimitString", func(t *testing.T) {
		kvs := []attribute.KeyValue{
			nanFloat64SliceAttr(),
			attribute.String("long", long),
		}
		_, nanChanged := normAttr(kvs[0], DefaultAttributeValueDepthLimit, lengthLimit)
		assert.False(t, nanChanged)
		_, changed := normAttr(kvs[1], DefaultAttributeValueDepthLimit, lengthLimit)
		assert.True(t, changed)

		caller := append([]attribute.KeyValue(nil), kvs...)
		out := normAttrs(kvs, DefaultAttributeValueDepthLimit, lengthLimit)

		require.Len(t, out, len(kvs))
		assert.NotSame(t, &kvs[0], &out[0])
		assertNaNFloat64Slice(t, out[0], 1.5)
		assert.Equal(t, strings.Repeat("a", lengthLimit), out[1].Value.AsString())
		assertNaNFloat64Slice(t, kvs[0], 1.5)
		assert.Equal(t, long, kvs[1].Value.AsString())
		assert.Equal(t, caller[1], kvs[1])
	})

	t.Run("DuplicateMapKeys", func(t *testing.T) {
		kvs := []attribute.KeyValue{
			nanFloat64SliceAttr(),
			attribute.Map("dups",
				attribute.Float64Slice("nan", []float64{math.NaN()}),
				attribute.String("k", "abcdef"),
				attribute.String("k", "zzzzzz"),
			),
		}
		_, nanChanged := normAttr(kvs[0], DefaultAttributeValueDepthLimit, lengthLimit)
		assert.False(t, nanChanged)
		_, changed := normAttr(kvs[1], DefaultAttributeValueDepthLimit, lengthLimit)
		assert.True(t, changed)

		out := normAttrs(kvs, DefaultAttributeValueDepthLimit, lengthLimit)

		require.Len(t, out, len(kvs))
		assert.NotSame(t, &kvs[0], &out[0])
		assertNaNFloat64Slice(t, out[0], 1.5)
		gotMap := out[1].Value.AsMap()
		require.Len(t, gotMap, 2)
		assert.Equal(t, attribute.String("k", "zzzzzz"), gotMap[0])
		assertNaNFloat64SliceValue(t, gotMap[1].Value)

		callerMap := kvs[1].Value.AsMap()
		require.Len(t, callerMap, 3)
		assert.Equal(t, "abcdef", callerMap[0].Value.AsString())
		assert.Equal(t, "zzzzzz", callerMap[1].Value.AsString())
		assertNaNFloat64SliceValue(t, callerMap[2].Value)
	})

	t.Run("DepthLimit", func(t *testing.T) {
		kvs := []attribute.KeyValue{
			nanFloat64SliceAttr(),
			attribute.Map("nested",
				attribute.String("keep", "short"),
				attribute.Map("child", attribute.String("leaf", long)),
			),
		}
		_, nanChanged := normAttr(kvs[0], 1, lengthLimit)
		assert.False(t, nanChanged)
		_, changed := normAttr(kvs[1], 1, lengthLimit)
		assert.True(t, changed)

		out := normAttrs(kvs, 1, lengthLimit)

		require.Len(t, out, len(kvs))
		assert.NotSame(t, &kvs[0], &out[0])
		assertNaNFloat64Slice(t, out[0], 1.5)
		assert.Equal(t, []attribute.KeyValue{
			{Key: "child"},
			attribute.String("keep", "short"),
		}, out[1].Value.AsMap())

		callerMap := kvs[1].Value.AsMap()
		require.Len(t, callerMap, 2)
		assert.Equal(t, attribute.Key("child"), callerMap[0].Key)
		child := callerMap[0].Value.AsMap()
		require.Len(t, child, 1)
		assert.Equal(t, long, child[0].Value.AsString())
		assert.Equal(t, "short", callerMap[1].Value.AsString())
	})
}

func TestNaNFloat64SliceEventAndLinkAttributes(t *testing.T) {
	limits := NewSpanLimits()
	limits.AttributeValueLengthLimit = 16

	rec := new(recorder)
	tp := NewTracerProvider(WithRawSpanLimits(limits), WithSpanProcessor(rec))
	tracer := tp.Tracer(t.Name())
	ctx := t.Context()
	sc := testLinkSpanContext()

	linkAttrs := []attribute.KeyValue{nanFloat64SliceAttr()}
	_, span := tracer.Start(ctx, "span", trace.WithLinks(trace.Link{
		SpanContext: sc,
		Attributes:  linkAttrs,
	}))
	rs := span.(*recordingSpan)
	require.Len(t, rs.links.queue, 1)
	assert.Same(t, &linkAttrs[0], &rs.links.queue[0].Attributes[0])
	assertNaNFloat64Slice(t, rs.links.queue[0].Attributes[0], 1.5)

	eventAttrs := []attribute.KeyValue{nanFloat64SliceAttr()}
	span.AddEvent("event", trace.WithAttributes(eventAttrs...))
	require.Len(t, rs.events.queue, 1)
	assertNaNFloat64Slice(t, rs.events.queue[0].Attributes[0], 1.5)
	assertNaNFloat64Slice(t, eventAttrs[0], 1.5)

	addLinkAttrs := []attribute.KeyValue{nanFloat64SliceAttr()}
	span.AddLink(trace.Link{SpanContext: sc, Attributes: addLinkAttrs})
	require.Len(t, rs.links.queue, 2)
	assert.Same(t, &addLinkAttrs[0], &rs.links.queue[1].Attributes[0])
	assertNaNFloat64Slice(t, addLinkAttrs[0], 1.5)

	recordAttrs := []attribute.KeyValue{nanFloat64SliceAttr()}
	span.RecordError(errors.New("nan-msg"), trace.WithAttributes(recordAttrs...))
	require.Len(t, rs.events.queue, 2)
	recorded := rs.events.queue[1].Attributes
	require.GreaterOrEqual(t, len(recorded), 3)
	assertNaNFloat64Slice(t, recorded[0], 1.5)
	assert.Equal(t, semconv.ExceptionTypeKey, recorded[1].Key)
	assert.LessOrEqual(t, len(recorded[1].Value.AsString()), limits.AttributeValueLengthLimit)
	assert.Equal(t, semconv.ExceptionMessage("nan-msg"), recorded[2])
	assertNaNFloat64Slice(t, recordAttrs[0], 1.5)

	span.End()
	require.NoError(t, tp.Shutdown(ctx))
	require.Len(t, *rec, 1)
	got := (*rec)[0]
	require.Len(t, got.Links(), 2)
	assertNaNFloat64Slice(t, got.Links()[0].Attributes[0], 1.5)
	assertNaNFloat64Slice(t, got.Links()[1].Attributes[0], 1.5)
	require.Len(t, got.Events(), 2)
	assertNaNFloat64Slice(t, got.Events()[0].Attributes[0], 1.5)
	assertNaNFloat64Slice(t, got.Events()[1].Attributes[0], 1.5)
}
