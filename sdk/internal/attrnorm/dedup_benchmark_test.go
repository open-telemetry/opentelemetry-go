// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrnorm_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/internal/attrnorm"
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
				attrnorm.ValueDedup(value.value)
			}
		})
	}
}

func BenchmarkValueWithDepthLimit(b *testing.B) {
	values := []struct {
		name       string
		value      attribute.Value
		depthLimit int
	}{
		{
			name:       "ScalarNoop",
			value:      attribute.StringValue("value"),
			depthLimit: 2,
		},
		{
			name: "NestedNoop",
			value: attribute.MapValue(
				attribute.Map(
					"level1",
					attribute.String("leaf", "value"),
				),
			),
			depthLimit: 2,
		},
		{
			name: "LimitHit",
			value: attribute.MapValue(
				attribute.Map(
					"level1",
					attribute.Map(
						"level2",
						attribute.String("leaf", "value"),
					),
				),
			),
			depthLimit: 2,
		},
	}

	for _, value := range values {
		b.Run(value.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				attrnorm.ValueWithDepthLimit(value.value, value.depthLimit)
			}
		})
	}
}

func BenchmarkKeyValuesWithDepthLimit(b *testing.B) {
	values := []struct {
		name       string
		values     []attribute.KeyValue
		depthLimit int
	}{
		{
			name: "ScalarNoop",
			values: []attribute.KeyValue{
				attribute.String("one", "1"),
				attribute.Int("two", 2),
				attribute.Bool("three", true),
				attribute.Float64("four", 4.0),
			},
			depthLimit: 2,
		},
		{
			name: "NestedNoop",
			values: []attribute.KeyValue{
				attribute.String("top", "value"),
				attribute.Map(
					"nested",
					attribute.String("leaf", "value"),
				),
				attribute.String("tail", "value"),
			},
			depthLimit: 2,
		},
		{
			name: "LimitHit",
			values: []attribute.KeyValue{
				attribute.String("top", "value"),
				attribute.Map(
					"nested",
					attribute.Map(
						"level1",
						attribute.Map(
							"level2",
							attribute.String("leaf", "value"),
						),
					),
				),
				attribute.String("tail", "value"),
			},
			depthLimit: 2,
		},
	}

	for _, value := range values {
		b.Run(value.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				attrnorm.KeyValuesWithDepthLimit(value.values, value.depthLimit)
			}
		})
	}
}

func BenchmarkSetWithDepthLimit(b *testing.B) {
	values := []struct {
		name       string
		set        attribute.Set
		depthLimit int
	}{
		{
			name: "ScalarNoop",
			set: attribute.NewSet(
				attribute.String("one", "1"),
				attribute.Int("two", 2),
				attribute.Bool("three", true),
				attribute.Float64("four", 4.0),
			),
			depthLimit: 2,
		},
		{
			name: "NestedNoop",
			set: attribute.NewSet(
				attribute.String("top", "value"),
				attribute.Map(
					"nested",
					attribute.String("leaf", "value"),
				),
				attribute.String("tail", "value"),
			),
			depthLimit: 2,
		},
		{
			name: "LimitHit",
			set: attribute.NewSet(
				attribute.String("top", "value"),
				attribute.Map(
					"nested",
					attribute.Map(
						"level1",
						attribute.Map(
							"level2",
							attribute.String("leaf", "value"),
						),
					),
				),
				attribute.String("tail", "value"),
			),
			depthLimit: 2,
		},
	}

	for _, value := range values {
		b.Run(value.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				attrnorm.SetWithDepthLimit(value.set, value.depthLimit)
			}
		})
	}
}
