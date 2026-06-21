// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrnorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringNeedsTruncation(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		input string
		want  bool
	}{
		// Negative limit: truncation is never needed.
		{name: "no_limit", limit: -1, input: "hello", want: false},
		// ASCII: string is well under the limit.
		{name: "ascii/under_limit", limit: 10, input: "hello", want: false},
		// ASCII: string length equals the limit exactly.
		{name: "ascii/exact_limit", limit: 5, input: "hello", want: false},
		// ASCII: string length exceeds the limit.
		{name: "ascii/exceeds", limit: 3, input: "hello", want: true},
		// Multibyte: byte length exceeds limit but rune count does not.
		{name: "multibyte/under_limit", limit: 3, input: "日本語", want: false},
		// Multibyte: rune count exceeds limit.
		{name: "multibyte/exceeds", limit: 2, input: "日本語", want: true},
		// Invalid UTF-8: byte length exceeds limit, invalid byte would be stripped on truncation.
		{name: "invalid_utf8/exceeds", limit: 2, input: "ab\xff", want: true},
		// Invalid UTF-8: byte length within limit, Truncate returns string unchanged.
		{name: "invalid_utf8/under_limit", limit: 5, input: "ab\xff", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, StringNeedsTruncation(tt.limit, tt.input))
		})
	}
}

func BenchmarkTruncate(b *testing.B) {
	run := func(limit int, input string) func(b *testing.B) {
		return func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				var out string
				for pb.Next() {
					out = Truncate(limit, input)
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

		// Truncation with mixed valid and invalid UTF-8 characters
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

				got := Truncate(g.limit, g.input)
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
