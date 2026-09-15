// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
