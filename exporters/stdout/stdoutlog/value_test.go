// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func BenchmarkExporterFloat64Slice(b *testing.B) {
	for _, name := range []string{"finite", "nonfinite"} {
		for _, size := range []int{8, 1024} {
			b.Run(fmt.Sprintf("%s/%d", name, size), func(b *testing.B) {
				b.Setenv("OTEL_GO_X_OBSERVABILITY", "false")
				values := make([]float64, size)
				for i := range values {
					values[i] = float64(i) + 0.5
				}
				if name == "nonfinite" {
					values[0], values[size/2], values[size-1] = math.NaN(), math.Inf(1), math.Inf(-1)
				}
				var record sdklog.Record
				record.SetBody(attribute.Float64SliceValue(values))
				records := []sdklog.Record{record}
				exporter, err := New(WithWriter(io.Discard), WithoutTimestamps())
				require.NoError(b, err)
				ctx := b.Context()

				b.ReportAllocs()
				b.ResetTimer()
				for b.Loop() {
					if err := exporter.Export(ctx, records); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
