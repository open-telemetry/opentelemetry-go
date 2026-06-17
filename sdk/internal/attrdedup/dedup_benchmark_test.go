// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrdedup_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/internal/attrdedup"
)

func BenchmarkValue(b *testing.B) {
	values := []struct {
		name  string
		value attribute.Value
	}{
		{
			name: "FastPath",
			value: attribute.MapValue(
				attribute.String("one", "1"),
				attribute.String("two", "2"),
				attribute.String("three", "3"),
			),
		},
		{
			name: "DuplicateMap",
			value: attribute.MapValue(
				attribute.String("one", "1"),
				attribute.String("one", "2"),
				attribute.String("two", "3"),
			),
		},
		{
			name: "NestedMapInSlice",
			value: attribute.SliceValue(
				attribute.MapValue(
					attribute.String("one", "1"),
					attribute.String("one", "2"),
					attribute.String("two", "3"),
				),
			),
		},
	}

	for _, value := range values {
		b.Run(value.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				attrdedup.Value(value.value)
			}
		})
	}
}
