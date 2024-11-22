// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bytes"
	"context"
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

func TestTruncateAttr(t *testing.T) {
	const key = "key"

	strAttr := attribute.String(key, "value")
	strSliceAttr := attribute.StringSlice(key, []string{"value-0", "value-1"})

	tests := []struct {
		limit      int
		attr, want attribute.KeyValue
	}{
		{
			limit: -1,
			attr:  strAttr,
			want:  strAttr,
		},
		{
			limit: -1,
			attr:  strSliceAttr,
			want:  strSliceAttr,
		},
		{
			limit: 0,
			attr:  attribute.Bool(key, true),
			want:  attribute.Bool(key, true),
		},
		{
			limit: 0,
			attr:  attribute.BoolSlice(key, []bool{true, false}),
			want:  attribute.BoolSlice(key, []bool{true, false}),
		},
		{
			limit: 0,
			attr:  attribute.Int(key, 42),
			want:  attribute.Int(key, 42),
		},
		{
			limit: 0,
			attr:  attribute.IntSlice(key, []int{42, -1}),
			want:  attribute.IntSlice(key, []int{42, -1}),
		},
		{
			limit: 0,
			attr:  attribute.Int64(key, 42),
			want:  attribute.Int64(key, 42),
		},
		{
			limit: 0,
			attr:  attribute.Int64Slice(key, []int64{42, -1}),
			want:  attribute.Int64Slice(key, []int64{42, -1}),
		},
		{
			limit: 0,
			attr:  attribute.Float64(key, 42),
			want:  attribute.Float64(key, 42),
		},
		{
			limit: 0,
			attr:  attribute.Float64Slice(key, []float64{42, -1}),
			want:  attribute.Float64Slice(key, []float64{42, -1}),
		},
		{
			limit: 0,
			attr:  strAttr,
			want:  attribute.String(key, ""),
		},
		{
			limit: 0,
			attr:  strSliceAttr,
			want:  attribute.StringSlice(key, []string{"", ""}),
		},
		{
			limit: 0,
			attr:  attribute.Stringer(key, bytes.NewBufferString("value")),
			want:  attribute.String(key, ""),
		},
		{
			limit: 1,
			attr:  strAttr,
			want:  attribute.String(key, "v"),
		},
		{
			limit: 1,
			attr:  strSliceAttr,
			want:  attribute.StringSlice(key, []string{"v", "v"}),
		},
		{
			limit: 5,
			attr:  strAttr,
			want:  strAttr,
		},
		{
			limit: 7,
			attr:  strSliceAttr,
			want:  strSliceAttr,
		},
		{
			limit: 6,
			attr:  attribute.StringSlice(key, []string{"value", "value-1"}),
			want:  attribute.StringSlice(key, []string{"value", "value-"}),
		},
		{
			limit: 128,
			attr:  strAttr,
			want:  strAttr,
		},
		{
			limit: 128,
			attr:  strSliceAttr,
			want:  strSliceAttr,
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%s->%s(limit:%d)", test.attr.Key, test.attr.Value.Emit(), test.limit)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, truncateAttr(test.limit, test.attr))
		})
	}
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
	for i := 0; i < 100; i++ {
		attr := attribute.String(fmt.Sprintf("hello.attrib%d", i), fmt.Sprintf("goodbye.attrib%d", i))
		attrs = append(attrs, attr)
	}

	ctx := context.Background()
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
	tracer := NewTracerProvider().Tracer("")
	ctx := trace.ContextWithSpanContext(context.Background(), trace.SpanContext{})

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
}
