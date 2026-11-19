// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile_test

import (
	"testing"

	"go.opentelemetry.io/otel/sdk/profile"
)

// Recommended CLI invocation:
// go test ./sdk/profile -run '^$' -bench '^BenchmarkRuntimeProducerProduce$' -benchtime 1s -count 6 -cpu 2 -timeout 999m -benchmem | tee sdk_profile_bench.txt.
func BenchmarkRuntimeProducerProduce(b *testing.B) {
	rp := profile.NewRuntimeProducer()
	b.Cleanup(func() { _ = rp.Shutdown(b.Context()) })

	ctx := b.Context()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if _, err := rp.Produce(ctx); err != nil {
			b.Fatal(err)
		}
	}
}
