// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global

import (
	"testing"
)

func BenchmarkStartEndSpanNoSDK(b *testing.B) {
	// Compare with BenchmarkStartEndSpan() in
	// ../../sdk/trace/benchmark_test.go.
	ResetForTest(b)
	t := TracerProvider().Tracer("Benchmark StartEndSpan")
	ctx := b.Context()

	for b.Loop() {
		_, span := t.Start(ctx, "/foo")
		span.End()
	}
}
