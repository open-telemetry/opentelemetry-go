// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/logtest"
)

func BenchmarkExporterScopeAttributes(b *testing.B) {
	for _, tc := range []struct {
		name  string
		attrs []attribute.KeyValue
	}{
		{name: "Empty"},
		{
			name: "Primitive",
			attrs: []attribute.KeyValue{
				attribute.String("name", "test"),
				attribute.Bool("enabled", true),
				attribute.Int64("count", 10),
			},
		},
		{
			name: "Nested",
			attrs: []attribute.KeyValue{
				attribute.Map("nested",
					attribute.ByteSlice("bytes", []byte("bin")),
					attribute.Slice("array", attribute.StringValue("test"), attribute.Float64Value(1.5)),
				),
			},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.Setenv("OTEL_GO_X_OBSERVABILITY", "false")
			exporter, err := New(WithWriter(io.Discard), WithoutTimestamps())
			require.NoError(b, err)
			factory := logtest.RecordFactory{InstrumentationScope: &instrumentation.Scope{
				Name:       "test",
				Attributes: attribute.NewSet(tc.attrs...),
			}}
			records := []sdklog.Record{factory.NewRecord()}
			b.ReportAllocs()
			for b.Loop() {
				if err := exporter.Export(b.Context(), records); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
