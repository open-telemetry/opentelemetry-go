// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package benchmark // import "go.opentelemetry.io/otel/internal/benchmark"

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const scopeName = "go.opentelemetry.op/otel/internal/benchmark"

func testCounter(b *testing.B) metric.Float64Counter {
	rdr := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(rdr))
	meter := provider.Meter(scopeName)
	counter, err := meter.Float64Counter("test.counter")
	assert.NoError(b, err)
	return counter
}

var (
	attributes = []attribute.KeyValue{
		attribute.String("a", "a"),
		attribute.String("b", "b"),
		attribute.String("c", "c"),
		attribute.String("d", "d"),
		attribute.String("e", "e"),
		attribute.String("f", "f"),
		attribute.String("g", "g"),
		attribute.String("h", "h"),
		attribute.String("i", "i"),
		attribute.String("j", "j"),
	}
	addOptPool = &sync.Pool{
		New: func() any {
			const n = 1 // WithAttributeSet
			o := make([]metric.AddOption, 0, n)
			// Return a pointer to avoid extra allocation on Put().
			return &o
		},
	}
)

func BenchmarkCounterIncrement(b *testing.B) {
	ctx := b.Context()
	for _, attrsLen := range []int{1, 3, 10} {
		attrs := attributes[:attrsLen]
		attrPool := sync.Pool{
			New: func() any {
				// Pre-allocate common capacity
				s := make([]attribute.KeyValue, 0, attrsLen)
				// Return a pointer to avoid extra allocation on Put().
				return &s
			},
		}
		b.Run(fmt.Sprintf("Attributes/%d", attrsLen), func(b *testing.B) {
			b.Run("PrecomputedAttributeSet", func(b *testing.B) {
				counter := testCounter(b)
				precomputedOpts := []metric.AddOption{metric.WithAttributeSet(attribute.NewSet(attrs...))}
				b.ReportAllocs()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						counter.Add(ctx, 1, precomputedOpts...)
					}
				})
			})
			// Based on https://github.com/open-telemetry/opentelemetry-go/blob/main/CONTRIBUTING.md#attribute-and-option-allocation-management
			b.Run("OptimizedDynamicAttributeSet", func(b *testing.B) {
				counter := testCounter(b)
				b.ReportAllocs()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						// Wrap in a function so we can use defer.
						func() {
							attrsSlice := attrPool.Get().(*[]attribute.KeyValue)
							defer func() {
								*attrsSlice = (*attrsSlice)[:0] // Reset.
								attrPool.Put(attrsSlice)
							}()
							*attrsSlice = append(*attrsSlice, attrs...)
							addOpt := addOptPool.Get().(*[]metric.AddOption)
							defer func() {
								*addOpt = (*addOpt)[:0]
								addOptPool.Put(addOpt)
							}()
							set := attribute.NewSet(*attrsSlice...)
							*addOpt = append(*addOpt, metric.WithAttributeSet(set))
							counter.Add(ctx, 1, *addOpt...)
						}()
					}
				})
			})
			b.Run("NaiveDynamicAttributeSet", func(b *testing.B) {
				counter := testCounter(b)
				b.ReportAllocs()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						counter.Add(ctx, 1, metric.WithAttributeSet(attribute.NewSet(attrs...)))
					}
				})
			})
			b.Run("NewDynamicWithAttributes", func(b *testing.B) {
				counter := testCounter(b)
				b.ReportAllocs()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						counter.AddWithAttributes(ctx, 1, attrs)
					}
				})
			})
		})
	}
}
