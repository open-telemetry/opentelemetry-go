// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bytes"
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
	bytesAttr := attribute.ByteSlice(key, []byte("value"))
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
			limit: -1,
			attr:  bytesAttr,
			want:  bytesAttr,
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
			limit: 0,
			attr:  bytesAttr,
			want:  attribute.ByteSlice(key, []byte{}),
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
			limit: 1,
			attr:  bytesAttr,
			want:  attribute.ByteSlice(key, []byte("v")),
		},
		{
			limit: 5,
			attr:  strAttr,
			want:  strAttr,
		},
		{
			limit: 5,
			attr:  bytesAttr,
			want:  bytesAttr,
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
		{
			limit: 128,
			attr:  bytesAttr,
			want:  bytesAttr,
		},
		{
			// Multi-byte string: byte length (9) exceeds limit (5) but rune count (3) does not.
			// Must not be truncated.
			limit: 5,
			attr:  attribute.String(key, "日本語"),
			want:  attribute.String(key, "日本語"),
		},
		{
			// Multi-byte string: both byte length and rune count exceed limit.
			// Must be truncated to limit runes.
			limit: 2,
			attr:  attribute.String(key, "日本語"),
			want:  attribute.String(key, "日本"),
		},
		{
			// STRINGSLICE with multi-byte elements: byte lengths exceed limit but rune counts do not.
			// Must not be truncated.
			limit: 1,
			attr:  attribute.StringSlice(key, []string{"日", "本"}),
			want:  attribute.StringSlice(key, []string{"日", "本"}),
		},
		// SLICE cases
		{
			limit: -1,
			attr:  attribute.Slice(key, attribute.StringValue("value")),
			want:  attribute.Slice(key, attribute.StringValue("value")),
		},
		{
			limit: 0,
			attr:  attribute.Slice(key, attribute.BoolValue(true), attribute.StringValue("value")),
			want:  attribute.Slice(key, attribute.BoolValue(true), attribute.StringValue("")),
		},
		{
			limit: 5,
			attr:  attribute.Slice(key, attribute.StringValue("value"), attribute.StringValue("toolong")),
			want:  attribute.Slice(key, attribute.StringValue("value"), attribute.StringValue("toolo")),
		},
		{
			// Nested SLICE: recursive truncation.
			limit: 1,
			attr:  attribute.Slice(key, attribute.SliceValue(attribute.StringValue("value"))),
			want:  attribute.Slice(key, attribute.SliceValue(attribute.StringValue("v"))),
		},
		{
			// STRINGSLICE within SLICE: each string element is truncated.
			limit: 2,
			attr:  attribute.Slice(key, attribute.StringSliceValue([]string{"abc", "de"})),
			want:  attribute.Slice(key, attribute.StringSliceValue([]string{"ab", "de"})),
		},
		{
			// STRINGSLICE within SLICE where all strings fit: no change.
			// Exercises needsTruncation(STRINGSLICE) exhausting the loop without
			// finding an over-limit string, returning false.
			limit: 7,
			attr:  attribute.Slice(key, attribute.StringSliceValue([]string{"value-0", "value-1"})),
			want:  attribute.Slice(key, attribute.StringSliceValue([]string{"value-0", "value-1"})),
		},
		{
			// Mixed SLICE: STRINGSLICE (all strings fit) + STRING (too long).
			// Exercises recursive truncation over mixed slice elements: the
			// STRINGSLICE element remains unchanged because each string fits
			// within the limit, while the sibling STRING element is truncated.
			limit: 3,
			attr: attribute.Slice(
				key,
				attribute.StringSliceValue([]string{"ab", "cd"}),
				attribute.StringValue("toolong"),
			),
			want: attribute.Slice(
				key,
				attribute.StringSliceValue([]string{"ab", "cd"}),
				attribute.StringValue("too"),
			),
		},
		{
			// Nested SLICE (no truncation needed) alongside STRING (needs truncation).
			// Exercises the truncateValue SLICE branch early-return path: truncateValue
			// is called recursively on the nested SLICE but returns it unchanged because
			// none of its elements require truncation.
			limit: 3,
			attr: attribute.Slice(
				key,
				attribute.SliceValue(attribute.BoolValue(true)),
				attribute.StringValue("toolong"),
			),
			want: attribute.Slice(
				key,
				attribute.SliceValue(attribute.BoolValue(true)),
				attribute.StringValue("too"),
			),
		},
		{
			// Multi-byte string whose byte length exceeds the limit but rune count
			// does not: must not be truncated (guards use rune count, not byte length).
			limit: 3,
			attr:  attribute.Slice(key, attribute.StringValue("日本語")), // 3 runes, 9 bytes
			want:  attribute.Slice(key, attribute.StringValue("日本語")),
		},
		{
			// SLICE with invalid UTF-8 where rune count equals the limit:
			// invalid byte is dropped.
			limit: 2,
			attr:  attribute.Slice(key, attribute.StringValue("日\x80")), // 2 runes (日 + invalid byte), 4 bytes
			want:  attribute.Slice(key, attribute.StringValue("日")),
		},
		{
			// BYTESLICE within SLICE: each byte slice is truncated.
			limit: 2,
			attr:  attribute.Slice(key, attribute.ByteSliceValue([]byte{1, 2, 3})),
			want:  attribute.Slice(key, attribute.ByteSliceValue([]byte{1, 2})),
		},
		{
			// BYTESLICE within SLICE: no truncation needed.
			limit: 5,
			attr:  attribute.Slice(key, attribute.ByteSliceValue([]byte{1, 2})),
			want:  attribute.Slice(key, attribute.ByteSliceValue([]byte{1, 2})),
		},
		{
			// Mixed SLICE: BYTESLICE + STRING (both need truncation).
			limit: 2,
			attr: attribute.Slice(
				key,
				attribute.ByteSliceValue([]byte{1, 2, 3}),
				attribute.StringValue("abc"),
			),
			want: attribute.Slice(
				key,
				attribute.ByteSliceValue([]byte{1, 2}),
				attribute.StringValue("ab"),
			),
		},
		// MAP cases
		{
			limit: -1,
			attr:  attribute.Map(key, attribute.String("value", "value")),
			want:  attribute.Map(key, attribute.String("value", "value")),
		},
		{
			limit: 0,
			attr:  attribute.Map(key, attribute.Bool("ok", true), attribute.String("value", "value")),
			want:  attribute.Map(key, attribute.Bool("ok", true), attribute.String("value", "")),
		},
		{
			limit: 5,
			attr: attribute.Map(
				key,
				attribute.String("short", "value"),
				attribute.String("long", "toolong"),
			),
			want: attribute.Map(
				key,
				attribute.String("short", "value"),
				attribute.String("long", "toolo"),
			),
		},
		{
			// STRINGSLICE within MAP: each string element is truncated.
			limit: 2,
			attr:  attribute.Map(key, attribute.StringSlice("strings", []string{"abc", "de"})),
			want:  attribute.Map(key, attribute.StringSlice("strings", []string{"ab", "de"})),
		},
		{
			// BYTESLICE within MAP: each byte slice is truncated.
			limit: 2,
			attr:  attribute.Map(key, attribute.ByteSlice("bytes", []byte{1, 2, 3})),
			want:  attribute.Map(key, attribute.ByteSlice("bytes", []byte{1, 2})),
		},
		{
			// Nested MAP: recursive truncation.
			limit: 1,
			attr:  attribute.Map(key, attribute.Map("map", attribute.String("nested", "value"))),
			want:  attribute.Map(key, attribute.Map("map", attribute.String("nested", "v"))),
		},
		{
			// SLICE within MAP: recursive truncation.
			limit: 2,
			attr: attribute.Map(
				key,
				attribute.Slice(
					"slice",
					attribute.StringValue("abc"),
					attribute.MapValue(attribute.String("nested", "abc")),
				),
			),
			want: attribute.Map(
				key,
				attribute.Slice(
					"slice",
					attribute.StringValue("ab"),
					attribute.MapValue(attribute.String("nested", "ab")),
				),
			),
		},
		{
			// MAP within SLICE: recursive truncation.
			limit: 2,
			attr:  attribute.Slice(key, attribute.MapValue(attribute.String("nested", "value"))),
			want:  attribute.Slice(key, attribute.MapValue(attribute.String("nested", "va"))),
		},
		{
			// Multi-byte string whose byte length exceeds the limit but rune count
			// does not: must not be truncated.
			limit: 3,
			attr:  attribute.Map(key, attribute.String("string", "日本語")), // 3 runes, 9 bytes
			want:  attribute.Map(key, attribute.String("string", "日本語")),
		},
		{
			// MAP with invalid UTF-8 where rune count equals the limit:
			// invalid byte is dropped.
			limit: 2,
			attr:  attribute.Map(key, attribute.String("string", "日\x80")), // 2 runes, 4 bytes
			want:  attribute.Map(key, attribute.String("string", "日")),
		},
		{
			// Duplicate MAP entries are truncated but not dropped.
			limit: 2,
			attr: attribute.Map(
				key,
				attribute.String("dup", "abc"),
				attribute.String("dup", "de"),
			),
			want: attribute.Map(
				key,
				attribute.String("dup", "ab"),
				attribute.String("dup", "de"),
			),
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%s->%s(limit:%d)", test.attr.Key, test.attr.Value.String(), test.limit)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.want, truncateAttr(test.limit, test.attr))
		})
	}
}

func TestStringSliceNeedsTruncation(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		value attribute.Value
		want  bool
	}{
		// Empty slice, no elements to exceed the limit.
		{
			name:  "empty",
			limit: 0,
			value: attribute.StringSliceValue([]string{}),
			want:  false,
		},
		// Single element within limit.
		{
			name:  "one/fits",
			limit: 5,
			value: attribute.StringSliceValue([]string{"hello"}),
			want:  false,
		},
		// Single element exceeds limit.
		{
			name:  "one/exceeds",
			limit: 3,
			value: attribute.StringSliceValue([]string{"hello"}),
			want:  true,
		},
		// Two elements, both fit.
		{
			name:  "two/fits",
			limit: 5,
			value: attribute.StringSliceValue([]string{"hello", "world"}),
			want:  false,
		},
		// Two elements, second exceeds limit.
		{
			name:  "two/second-exceeds",
			limit: 3,
			value: attribute.StringSliceValue([]string{"hi", "hello"}),
			want:  true,
		},
		// Three elements, all fit.
		{
			name:  "three/fits",
			limit: 5,
			value: attribute.StringSliceValue([]string{"one", "two", "three"}),
			want:  false,
		},
		// Three elements, middle exceeds limit.
		{
			name:  "three/middle-exceeds",
			limit: 3,
			value: attribute.StringSliceValue([]string{"ab", "abcd", "cd"}),
			want:  true,
		},
		// Four elements (reflect path), all fit.
		{
			name:  "four/fits",
			limit: 5,
			value: attribute.StringSliceValue([]string{"a", "b", "c", "d"}),
			want:  false,
		},
		// Four elements (reflect path), last exceeds limit.
		{
			name:  "four/last-exceeds",
			limit: 3,
			value: attribute.StringSliceValue([]string{"a", "b", "c", "abcd"}),
			want:  true,
		},
		// Multi-byte runes: byte length exceeds limit but rune count does not.
		{
			name:  "multibyte/rune-count-fits",
			limit: 3,
			value: attribute.StringSliceValue([]string{"日本語"}), // 3 runes, 9 bytes
			want:  false,
		},
		// Multi-byte runes: rune count exceeds limit.
		{
			name:  "multibyte/rune-count-exceeds",
			limit: 2,
			value: attribute.StringSliceValue([]string{"日本語"}), // 3 runes
			want:  true,
		},
		// Invalid UTF-8: byte length exceeds limit, invalid byte would be stripped.
		{
			name:  "invalid-utf8/exceeds",
			limit: 2,
			value: attribute.StringSliceValue([]string{"ab\xff"}), // 3 bytes, invalid
			want:  true,
		},
		// Invalid UTF-8: byte length within limit, returned unchanged by truncate.
		{
			name:  "invalid-utf8/fits",
			limit: 5,
			value: attribute.StringSliceValue([]string{"ab\xff"}), // 3 bytes <= 5
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, stringSliceNeedsTruncation(tt.limit, tt.value))
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

func BenchmarkTruncateAttr(b *testing.B) {
	const key = "key"

	strAttr := attribute.String(key, "value")
	bytesAttr := attribute.ByteSlice(key, []byte("value"))
	strSliceAttr := attribute.StringSlice(key, []string{"value-0", "value-1"})

	run := func(limit int, attr attribute.KeyValue) func(b *testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				var out attribute.KeyValue
				for pb.Next() {
					out = truncateAttr(limit, attr)
				}
				_ = out
			})
		}
	}

	b.Run("String", run(3, strAttr))
	b.Run("StringSlice", run(3, strSliceAttr))
	b.Run("ByteSlice", run(3, bytesAttr))
	b.Run("String/Limit0", run(0, strAttr))
	b.Run("StringSlice/Limit0", run(0, strSliceAttr))
	b.Run("ByteSlice/Limit0", run(0, bytesAttr))
	b.Run("String/Unlimited", run(-1, strAttr))
	b.Run("StringSlice/Unlimited", run(-1, strSliceAttr))
	b.Run("ByteSlice/Unlimited", run(-1, bytesAttr))
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
