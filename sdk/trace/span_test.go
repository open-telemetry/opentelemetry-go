// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

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
