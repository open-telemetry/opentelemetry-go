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
	"go.opentelemetry.io/otel/metric/noop"
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

func BenchmarkCounterIncrement(b *testing.B) {
	ctx := b.Context()
	for _, mp := range []struct {
		name     string
		provider func() metric.MeterProvider
	}{
		{
			name:     "NoOpMeterProvider",
			provider: func() metric.MeterProvider { return noop.NewMeterProvider() },
		},
		{
			name: "DefaultMeterProvider",
			provider: func() metric.MeterProvider {
				return sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewManualReader()))
			},
		},
		{
			name: "FilteredMeterProvider",
			provider: func() metric.MeterProvider {
				view := sdkmetric.NewView(
					sdkmetric.Instrument{
						Name: "test.counter",
					},
					// Filter out one attribute from each call, but don't change cardinality.
					sdkmetric.Stream{AttributeFilter: attribute.NewDenyKeysFilter("b")},
				)
				return sdkmetric.NewMeterProvider(sdkmetric.WithView(view), sdkmetric.WithReader(sdkmetric.NewManualReader()))
			},
		},
	} {
		b.Run(mp.name, func(b *testing.B) {
			for _, attrsLen := range []int{2, 5, 10} {
				attrPool := sync.Pool{
					New: func() any {
						// Pre-allocate common capacity
						s := make([]attribute.KeyValue, 0, attrsLen)
						// Return a pointer to avoid extra allocation on Put().
						return &s
					},
				}
				b.Run(fmt.Sprintf("Attributes/%d", attrsLen), func(b *testing.B) {
					for _, cardinality := range []int{1, 10, 100} {
						b.Run(fmt.Sprintf("Cardinality/%d", cardinality), func(b *testing.B) {
							b.Run("PrecomputedWithAttributeSet", func(b *testing.B) {
								counter := testCounter(b, mp.provider())
								opts := make([][]metric.AddOption, cardinality)
								for i := range cardinality {
									opts[i] = []metric.AddOption{metric.WithAttributeSet(attribute.NewSet(getAttributes(attrsLen, cardinality, i)...))}
								}
								b.ReportAllocs()
								b.RunParallel(func(pb *testing.PB) {
									i := 0
									for pb.Next() {
										counter.Add(ctx, 1, opts[i%cardinality]...)
										i++
									}
								})
							})
							b.Run("PrecomputedWithAttributes", func(b *testing.B) {
								counter := testCounter(b, mp.provider())
								opts := make([][]metric.AddOption, cardinality)
								for i := range cardinality {
									opts[i] = []metric.AddOption{metric.WithAttributes(getAttributes(attrsLen, cardinality, i)...)}
								}
								b.ReportAllocs()
								b.RunParallel(func(pb *testing.PB) {
									i := 0
									for pb.Next() {
										counter.Add(ctx, 1, opts[i%cardinality]...)
										i++
									}
								})
							})
							// Based on https://github.com/open-telemetry/opentelemetry-go/blob/main/CONTRIBUTING.md#attribute-and-option-allocation-management
							b.Run("DynamicWithAttributeSet", func(b *testing.B) {
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
											appendAttributes(attrsLen, cardinality, i, attrsSlice)
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
							b.Run("DynamicWithAttributes", func(b *testing.B) {
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
											appendAttributes(attrsLen, cardinality, i, attrsSlice)
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
		})
	}
}

func getAttributes(number, cardinality, index int) []attribute.KeyValue {
	kvs := make([]attribute.KeyValue, 0, number)
	appendAttributes(number, cardinality, index, &kvs)
	return kvs
}

func appendAttributes(number, cardinality, index int, kvs *[]attribute.KeyValue) {
	switch number {
	case 2:
		*kvs = append(*kvs,
			attribute.Int("a", index%cardinality),
			attribute.String("b", "b"),
		)
	case 5:
		*kvs = append(*kvs,
			attribute.Int("a", index%cardinality),
			attribute.String("b", "b"),
			attribute.String("c", "c"),
			attribute.String("d", "d"),
			attribute.String("e", "e"),
		)
	case 10:
		*kvs = append(*kvs,
			attribute.Int("a", index%cardinality),
			attribute.String("b", "b"),
			attribute.String("c", "c"),
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
