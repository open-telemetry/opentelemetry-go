// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrdedup_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/internal/attrdedup"
)

func BenchmarkValue(b *testing.B) {
	uniqueMap := attribute.MapValue(
		attribute.String("one", "1"),
		attribute.String("two", "2"),
		attribute.String("three", "3"),
	)
	duplicateMap := attribute.MapValue(
		attribute.String("one", "1"),
		attribute.String("one", "2"),
		attribute.String("two", "3"),
	)
	nestedMapInSlice := attribute.SliceValue(
		attribute.MapValue(
			attribute.String("one", "1"),
			attribute.String("one", "2"),
			attribute.String("two", "3"),
		),
	)

	b.Run("FastPath", func(b *testing.B) {
		for b.Loop() {
			attrdedup.Value(uniqueMap, false)
		}
	})
	b.Run("DuplicateMap", func(b *testing.B) {
		for b.Loop() {
			attrdedup.Value(duplicateMap, false)
		}
	})
	b.Run("NestedMapInSlice", func(b *testing.B) {
		for b.Loop() {
			attrdedup.Value(nestedMapInSlice, false)
		}
	})
	b.Run("AllowKeyDuplication", func(b *testing.B) {
		for b.Loop() {
			attrdedup.Value(duplicateMap, true)
		}
	})
}
