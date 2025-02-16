// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"fmt"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func BenchmarkSpanProcessorOnEndWithMetricsSDK(b *testing.B) {
	b.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")
	// TODO: set global meterprovider
	for _, bb := range []struct {
		batchSize  int
		spansCount int
	}{
		{batchSize: 10, spansCount: 10},
		{batchSize: 10, spansCount: 100},
		{batchSize: 100, spansCount: 10},
		{batchSize: 100, spansCount: 100},
	} {
		b.Run(fmt.Sprintf("batch: %d, spans: %d", bb.batchSize, bb.spansCount), func(b *testing.B) {
			bsp := sdktrace.NewBatchSpanProcessor(
				tracetest.NewNoopExporter(),
				sdktrace.WithMaxExportBatchSize(bb.batchSize),
			)
			snap := tracetest.SpanStub{}.Snapshot()

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				// Ensure the export happens for every run
				for j := 0; j < bb.spansCount; j++ {
					bsp.OnEnd(snap)
				}
			}
		})
	}
}
