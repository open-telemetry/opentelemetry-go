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

func testCounter(b *testing.B, mp metric.MeterProvider) metric.Float64Counter {
	meter := mp.Meter(scopeName)
	counter, err := meter.Float64Counter("test.counter")
	assert.NoError(b, err)
	return counter
}

var (
	addOptPool = &sync.Pool{
		New: func() any {
			const n = 1 // WithAttributeSet
			o := make([]metric.AddOption, 0, n)
			// Return a pointer to avoid extra allocation on Put().
			return &o
		},
	}
)

func BenchmarkCounterAdd(b *testing.B) {
	ctx := b.Context()
	for _, mp := range []struct {
		name     string
		provider func() metric.MeterProvider
	}{
		{
			name: "NoFilter",
			provider: func() metric.MeterProvider {
				return sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewManualReader()))
			},
		},
		{
			name: "Filtered",
			provider: func() metric.MeterProvider {
				view := sdkmetric.NewView(
					sdkmetric.Instrument{
						Name: "test.counter",
					},
					// Filter out one attribute from each call.
					sdkmetric.Stream{AttributeFilter: attribute.NewDenyKeysFilter("a")},
				)
				return sdkmetric.NewMeterProvider(sdkmetric.WithView(view), sdkmetric.WithReader(sdkmetric.NewManualReader()))
			},
		},
	} {
		b.Run(mp.name, func(b *testing.B) {
			for _, attrsLen := range []int{1, 5, 10} {
				attrPool := sync.Pool{
					New: func() any {
						// Pre-allocate common capacity
						s := make([]attribute.KeyValue, 0, attrsLen)
						// Return a pointer to avoid extra allocation on Put().
						return &s
					},
				}
				b.Run(fmt.Sprintf("Attributes/%d", attrsLen), func(b *testing.B) {
					b.Run("Precomputed/WithAttributeSet", func(b *testing.B) {
						counter := testCounter(b, mp.provider())
						precomputedOpts := []metric.AddOption{metric.WithAttributeSet(attribute.NewSet(getAttributes(attrsLen)...))}
						b.ReportAllocs()
						b.RunParallel(func(pb *testing.PB) {
							i := 0
							for pb.Next() {
								counter.Add(ctx, 1, precomputedOpts...)
								i++
							}
						})
					})
					b.Run("Precomputed/WithAttributes", func(b *testing.B) {
						counter := testCounter(b, mp.provider())
						precomputedOpts := []metric.AddOption{metric.WithAttributes(getAttributes(attrsLen)...)}
						b.ReportAllocs()
						b.RunParallel(func(pb *testing.PB) {
							i := 0
							for pb.Next() {
								counter.Add(ctx, 1, precomputedOpts...)
								i++
							}
						})
					})
					// Based on https://github.com/open-telemetry/opentelemetry-go/blob/main/CONTRIBUTING.md#attribute-and-option-allocation-management
					b.Run("Dynamic/WithAttributeSet", func(b *testing.B) {
						counter := testCounter(b, mp.provider())
						b.ReportAllocs()
						b.RunParallel(func(pb *testing.PB) {
							i := 0
							for pb.Next() {
								// Wrap in a function so we can use defer.
								func() {
									attrsSlice := attrPool.Get().(*[]attribute.KeyValue)
									defer func() {
										*attrsSlice = (*attrsSlice)[:0] // Reset.
										attrPool.Put(attrsSlice)
									}()
									appendAttributes(attrsLen, attrsSlice)
									addOpt := addOptPool.Get().(*[]metric.AddOption)
									defer func() {
										*addOpt = (*addOpt)[:0]
										addOptPool.Put(addOpt)
									}()
									set := attribute.NewSet(*attrsSlice...)
									*addOpt = append(*addOpt, metric.WithAttributeSet(set))
									counter.Add(ctx, 1, *addOpt...)
								}()
								i++
							}
						})
					})
					b.Run("Dynamic/WithAttributes", func(b *testing.B) {
						counter := testCounter(b, mp.provider())
						b.ReportAllocs()
						b.RunParallel(func(pb *testing.PB) {
							i := 0
							for pb.Next() {
								// Wrap in a function so we can use defer.
								func() {
									attrsSlice := attrPool.Get().(*[]attribute.KeyValue)
									defer func() {
										*attrsSlice = (*attrsSlice)[:0] // Reset.
										attrPool.Put(attrsSlice)
									}()
									appendAttributes(attrsLen, attrsSlice)
									addOpt := addOptPool.Get().(*[]metric.AddOption)
									defer func() {
										*addOpt = (*addOpt)[:0]
										addOptPool.Put(addOpt)
									}()
									counter.Add(ctx, 1, metric.WithAttributes(*attrsSlice...))
								}()
								i++
							}
						})
					})
				})
			}
		})
	}
}

func getAttributes(number int) []attribute.KeyValue {
	kvs := make([]attribute.KeyValue, 0, number)
	appendAttributes(number, &kvs)
	return kvs
}

func appendAttributes(number int, kvs *[]attribute.KeyValue) {
	switch number {
	case 1:
		*kvs = append(*kvs,
			attribute.String("a", "a"),
		)
	case 5:
		*kvs = append(*kvs,
			attribute.String("a", "a"),
			attribute.String("b", "b"),
			attribute.String("c", "c"),
			attribute.String("d", "d"),
			attribute.String("e", "e"),
		)
	case 10:
		*kvs = append(*kvs,
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
		)
	default:
		panic("unknown number of attributes")
	}
}
