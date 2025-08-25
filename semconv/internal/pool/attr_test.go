// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pool_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv/internal/pool"
)

func TestGetAttrs(t *testing.T) {
	var a *[]attribute.KeyValue
	for i := 0; i <= pool.MaxCapAttrsSmall*2; i++ {
		assert.NotPanics(t, func() { a = pool.GetAttrSlice(i) })

		require.NotNil(t, a)
		assert.Empty(t, *a)
		assert.GreaterOrEqual(t, cap(*a), i)

		pool.PutAttrSlice(a)
	}
}

func BenchmarkAttrsPoolSmall(b *testing.B) {
	kv := attribute.String("k", "v")
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a := pool.GetAttrSlice(1)
			*a = append(*a, kv)
			_ = a
			pool.PutAttrSlice(a)
		}
	})
}

func BenchmarkAttrsPoolMedium(b *testing.B) {
	kvs := make([]attribute.KeyValue, 2*pool.MaxCapAttrsSmall)
	for i := range kvs {
		kvs[i] = attribute.Int("k"+strconv.Itoa(i), i)
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a := pool.GetAttrSlice(len(kvs))
			*a = append(*a, kvs...)
			_ = a
			pool.PutAttrSlice(a)
		}
	})
}
