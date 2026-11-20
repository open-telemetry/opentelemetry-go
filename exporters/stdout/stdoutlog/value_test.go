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

func BenchmarkExporterBody(b *testing.B) {
	for _, tc := range []struct {
		name  string
		value attribute.Value
	}{
		{"short string", attribute.StringValue("test")},
		{"escaped string", attribute.StringValue("quote\"\\\n\t\x00<>&\u2028\u2029\xff")},
		{"byte slice", attribute.ByteSliceValue([]byte{0, 1, 255})},
		{"NaN", attribute.Float64Value(math.NaN())},
		{"finite float", attribute.Float64Value(1e-7)},
		{"empty", attribute.Value{}},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.Setenv("OTEL_GO_X_OBSERVABILITY", "false")
			var record sdklog.Record
			record.SetBody(tc.value)
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
