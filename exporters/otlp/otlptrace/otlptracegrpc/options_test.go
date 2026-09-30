// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracegrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/otlpconfig"
)

func TestCompressorToCompression(t *testing.T) {
	tests := []struct {
		compressor string
		want       otlpconfig.Compression
		wantErr    bool
	}{
		{"gzip", otlpconfig.GzipCompression, false},
		{"none", otlpconfig.NoCompression, false},
		{"", otlpconfig.NoCompression, false},
		{"unknown", otlpconfig.NoCompression, true},
	}
	for _, tc := range tests {
		t.Run(tc.compressor, func(t *testing.T) {
			var errs []error
			orig := otel.GetErrorHandler()
			otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) { errs = append(errs, err) }))
			t.Cleanup(func() { otel.SetErrorHandler(orig) })

			assert.Equal(t, tc.want, compressorToCompression(tc.compressor))
			if tc.wantErr {
				assert.Len(t, errs, 1)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
