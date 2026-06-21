// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrnorm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestNeedsTruncation(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		value attribute.Value
		want  bool
	}{
		// STRING: delegates to StringNeedsTruncation.
		{name: "string/under_limit", limit: 10, value: attribute.StringValue("hello"), want: false},
		{name: "string/exceeds", limit: 3, value: attribute.StringValue("hello"), want: true},

		// BYTESLICE: truncated by byte length, not rune count.
		{name: "byteslice/under_limit", limit: 5, value: attribute.ByteSliceValue([]byte("hello")), want: false},
		{name: "byteslice/exceeds", limit: 3, value: attribute.ByteSliceValue([]byte("hello")), want: true},

		// STRINGSLICE: delegates to StringSliceNeedsTruncation.
		{name: "stringslice/under_limit", limit: 5, value: attribute.StringSliceValue([]string{"hi"}), want: false},
		{name: "stringslice/exceeds", limit: 3, value: attribute.StringSliceValue([]string{"hello"}), want: true},

		// SLICE: recurses into elements.
		{name: "slice/under_limit", limit: 5, value: attribute.SliceValue(attribute.StringValue("hi")), want: false},
		{name: "slice/exceeds", limit: 3, value: attribute.SliceValue(attribute.StringValue("hello")), want: true},

		// MAP: recurses into values.
		{name: "map/under_limit", limit: 5, value: attribute.MapValue(attribute.String("k", "hi")), want: false},
		{name: "map/exceeds", limit: 3, value: attribute.MapValue(attribute.String("k", "hello")), want: true},

		// Non-string types: never need truncation.
		{name: "int64", limit: 1, value: attribute.Int64Value(42), want: false},
		{name: "float64", limit: 1, value: attribute.Float64Value(3.14), want: false},
		{name: "bool", limit: 1, value: attribute.BoolValue(true), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NeedsTruncation(tt.limit, tt.value))
		})
	}
}
