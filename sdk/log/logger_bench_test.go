// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/log"
)

func BenchmarkLoggerEmit(b *testing.B) {
	logger := newTestLogger(b)

	r := log.Record{}
	r.SetTimestamp(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	r.SetObservedTimestamp(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	r.SetBody(log.StringValue("testing body value"))
	r.SetSeverity(log.SeverityInfo)
	r.SetSeverityText("testing text")

	r.AddAttributes(
		log.String("k1", "str"),
		log.Float64("k2", 1.0),
		log.Int("k3", 2),
		log.Bool("k4", true),
		log.Bytes("k5", []byte{1}),
	)

	r10 := r
	r10.AddAttributes(
		log.String("k6", "str"),
		log.Float64("k7", 1.0),
		log.Int("k8", 2),
		log.Bool("k9", true),
		log.Bytes("k10", []byte{1}),
	)

	require.Equal(b, 5, r.AttributesLen())
	require.Equal(b, 10, r10.AttributesLen())

	b.Run("5 attributes", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Emit(context.Background(), r)
			}
		})
	})

	b.Run("10 attributes", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Emit(context.Background(), r10)
			}
		})
	})
}

func BenchmarkLoggerEnabled(b *testing.B) {
	logger := newTestLogger(b)
	ctx := context.Background()
	param := log.EnabledParameters{Severity: log.SeverityDebug}
	var enabled bool

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		enabled = logger.Enabled(ctx, param)
	}

	_ = enabled
}

func newTestLogger(t testing.TB) log.Logger {
	provider := NewLoggerProvider(
		WithProcessor(newFltrProcessor("0", false)),
		WithProcessor(newFltrProcessor("1", true)),
	)
	return provider.Logger(t.Name())
}

func BenchmarkLoggerRetrieval(b *testing.B) {
	prepopulatedValues := []int{0, 100, 500}
	for _, prepop := range prepopulatedValues {
		b.Run(fmt.Sprintf("Retrieve different logger each time single-threaded with %d prepopulated values", prepop), func(b *testing.B) {
			benchmarkLoggerRetrieval(b, prepop)
		})
	}

	for _, prepop := range prepopulatedValues {
		b.Run(fmt.Sprintf("Retrieve different logger each time parallel with %d prepopulated values", prepop), func(b *testing.B) {
			benchmarkLoggerRetrievalParallel(b, prepop)
		})
	}

	for _, prepop := range prepopulatedValues {
		b.Run(fmt.Sprintf("Retrieve same logger each time single-threaded with %d prepopulated values", prepop), func(b *testing.B) {
			benchmarkLoggerRetrievalWithSameValues(b, prepop)
		})
	}

	for _, prepop := range prepopulatedValues {
		b.Run(fmt.Sprintf("Retrieve same logger each time parallel with %d prepopulated values", prepop), func(b *testing.B) {
			benchmarkLoggerRetrievalWithSameValuesParallel(b, prepop)
		})
	}
}

func benchmarkLoggerRetrieval(b *testing.B, prepopulate int) {
	provider := NewLoggerProvider()

	// Prepopulate the provider before measuring
	for i := 0; i < prepopulate; i++ {
		provider.Logger(strconv.Itoa(i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			provider.Logger(strconv.Itoa(rand.Int()))
		}
	})
}

func benchmarkLoggerRetrievalParallel(b *testing.B, prepopulate int) {
	provider := NewLoggerProvider()

	// Prepopulate the provider before measuring
	for i := 0; i < prepopulate; i++ {
		provider.Logger(strconv.Itoa(i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			provider.Logger(strconv.Itoa(rand.Int()))
		}
	})
}

func benchmarkLoggerRetrievalWithSameValues(b *testing.B, prepopulate int) {
	provider := NewLoggerProvider()

	// Prepopulate the provider before measuring
	for i := 0; i < prepopulate; i++ {
		provider.Logger(strconv.Itoa(i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		provider.Logger("100") // Lookup a fixed logger
	}
}

func benchmarkLoggerRetrievalWithSameValuesParallel(b *testing.B, prepopulate int) {
	provider := NewLoggerProvider()

	// Prepopulate the provider before starting parallel execution
	for i := 0; i < prepopulate; i++ {
		provider.Logger(strconv.Itoa(i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			provider.Logger("100") // Lookup a fixed logger
		}
	})
}
