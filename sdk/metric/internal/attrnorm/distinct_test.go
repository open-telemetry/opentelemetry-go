// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrnorm

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestNewDistinct(t *testing.T) {
	tests := []struct {
		name   string
		kvs    []attribute.KeyValue
		filter attribute.Filter
	}{
		{name: "empty"},
		{
			name: "empty key",
			kvs: []attribute.KeyValue{
				{},
				attribute.String("", "last"),
			},
		},
		{
			name: "ordered",
			kvs: []attribute.KeyValue{
				attribute.String("a", "one"),
				attribute.Int("b", 2),
				attribute.Bool("c", true),
			},
		},
		{
			name: "reverse ordered",
			kvs: []attribute.KeyValue{
				attribute.Bool("c", true),
				attribute.Int("b", 2),
				attribute.String("a", "one"),
			},
		},
		{
			name: "late disorder",
			kvs: []attribute.KeyValue{
				attribute.String("a", "one"),
				attribute.Int("b", 2),
				attribute.Bool("d", true),
				attribute.String("c", "three"),
			},
		},
		{
			name: "duplicate keys",
			kvs: []attribute.KeyValue{
				attribute.String("b", "first"),
				attribute.String("a", "one"),
				attribute.String("b", "last"),
				attribute.String("a", "last"),
			},
		},
		{
			name: "filter after duplicate resolution",
			kvs: []attribute.KeyValue{
				attribute.String("key", "first"),
				attribute.String("key", "last"),
			},
			filter: func(kv attribute.KeyValue) bool {
				return kv.Value.AsString() == "first"
			},
		},
		{
			name: "filtered",
			kvs: []attribute.KeyValue{
				attribute.String("c", "three"),
				attribute.String("a", "one"),
				attribute.String("b", "two"),
			},
			filter: func(kv attribute.KeyValue) bool {
				return kv.Key != "b"
			},
		},
		{
			name: "all filtered",
			kvs: []attribute.KeyValue{
				attribute.String("a", "one"),
				attribute.String("b", "two"),
			},
			filter: func(attribute.KeyValue) bool { return false },
		},
		{
			name: "nested duplicate keys",
			kvs: []attribute.KeyValue{
				attribute.Map(
					"map",
					attribute.String("nested", "first"),
					attribute.String("nested", "last"),
				),
			},
		},
		{
			name: "large unordered",
			kvs:  reverseAttributes(stackAttributes + 1),
		},
		{
			name:   "large unordered filtered",
			kvs:    reverseAttributes(stackAttributes + 1),
			filter: acceptAll,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewDistinct(test.kvs, test.filter)
			want := referenceDistinct(test.kvs, test.filter)
			if got != want {
				t.Fatalf("NewDistinct() = %v, want %v", got, want)
			}
		})
	}
}

func TestNewDistinctDoesNotModifyInput(t *testing.T) {
	tests := []struct {
		name   string
		kvs    []attribute.KeyValue
		filter attribute.Filter
	}{
		{
			name: "nested",
			kvs: []attribute.KeyValue{
				attribute.String("b", "first"),
				attribute.Map(
					"a",
					attribute.String("nested", "first"),
					attribute.String("nested", "last"),
				),
				attribute.String("b", "last"),
			},
		},
		{name: "stack", kvs: reverseAttributes(stackAttributes)},
		{name: "stack filtered", kvs: reverseAttributes(stackAttributes), filter: acceptAll},
		{name: "allocated", kvs: reverseAttributes(stackAttributes + 1)},
		{
			name:   "allocated filtered",
			kvs:    reverseAttributes(stackAttributes + 1),
			filter: acceptAll,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := slices.Clone(test.kvs)

			NewDistinct(test.kvs, test.filter)

			assert.Equal(t, want, test.kvs, "NewDistinct modified its input")
		})
	}
}

func TestNewDistinctFilterOrder(t *testing.T) {
	kvs := []attribute.KeyValue{
		attribute.String("a", "first"),
		attribute.String("b", "first"),
		attribute.String("c", "first"),
		attribute.String("a", "last"),
	}
	var got []attribute.KeyValue

	NewDistinct(kvs, func(kv attribute.KeyValue) bool {
		got = append(got, kv)
		return true
	})

	want := []attribute.KeyValue{
		attribute.String("a", "last"),
		attribute.String("b", "first"),
		attribute.String("c", "first"),
	}
	assert.Equal(t, want, got, "filter input mismatch")
}

func TestNewDistinctTypicalInputsDoNotAllocate(t *testing.T) {
	tests := []struct {
		name   string
		kvs    []attribute.KeyValue
		filter attribute.Filter
	}{
		{name: "ordered", kvs: orderedAttributes(stackAttributes)},
		{
			name:   "ordered filtered",
			kvs:    orderedAttributes(stackAttributes),
			filter: acceptAll,
		},
		{name: "reverse ordered", kvs: reverseAttributes(stackAttributes)},
		{
			name: "reverse ordered filtered",
			kvs:  reverseAttributes(stackAttributes),
			filter: func(kv attribute.KeyValue) bool {
				return kv.Value.AsInt64()%2 == 0
			},
		},
		{name: "late disorder", kvs: lateDisorderedAttributes(stackAttributes)},
		{
			name: "duplicate keys",
			kvs: []attribute.KeyValue{
				attribute.Int("b", 1),
				attribute.Int("a", 2),
				attribute.Int("b", 3),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			allocs := testing.AllocsPerRun(1, func() {
				distinctSink = NewDistinct(test.kvs, test.filter)
			})
			if allocs != 0 {
				t.Fatalf("NewDistinct() allocations = %v, want 0", allocs)
			}
		})
	}
}

func BenchmarkNewDistinct(b *testing.B) {
	for _, n := range []int{1, 5, 10, 20} {
		b.Run(fmt.Sprintf("attributes=%d", n), func(b *testing.B) {
			benchmarks := []struct {
				name string
				kvs  []attribute.KeyValue
			}{
				{name: "ordered", kvs: orderedAttributes(n)},
				{name: "reverse", kvs: reverseAttributes(n)},
				{name: "late disorder", kvs: lateDisorderedAttributes(n)},
			}
			for _, benchmark := range benchmarks {
				filters := []struct {
					name   string
					filter attribute.Filter
				}{
					{name: "unfiltered"},
					{
						name: "filtered",
						filter: func(kv attribute.KeyValue) bool {
							return kv.Value.AsInt64()%2 == 0
						},
					},
				}
				for _, fltr := range filters {
					b.Run(benchmark.name+"/"+fltr.name, func(b *testing.B) {
						b.ReportAllocs()
						for b.Loop() {
							distinctSink = NewDistinct(benchmark.kvs, fltr.filter)
						}
					})
				}
			}
		})
	}
}

var distinctSink attribute.Distinct

func referenceDistinct(kvs []attribute.KeyValue, filter attribute.Filter) attribute.Distinct {
	normalized, _ := KeyValuesDedup(slices.Clone(kvs))
	set := attribute.NewSet(normalized...)
	set, _ = set.Filter(filter)
	return set.Equivalent()
}

func orderedAttributes(n int) []attribute.KeyValue {
	kvs := make([]attribute.KeyValue, n)
	for i := range n {
		kvs[i] = attribute.Int(fmt.Sprintf("%03d", i), i)
	}
	return kvs
}

func reverseAttributes(n int) []attribute.KeyValue {
	kvs := orderedAttributes(n)
	slices.Reverse(kvs)
	return kvs
}

func lateDisorderedAttributes(n int) []attribute.KeyValue {
	kvs := orderedAttributes(n)
	if n > 1 {
		kvs[n-2], kvs[n-1] = kvs[n-1], kvs[n-2]
	}
	return kvs
}

func acceptAll(attribute.KeyValue) bool {
	return true
}
