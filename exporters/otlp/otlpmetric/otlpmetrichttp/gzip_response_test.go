// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetrichttp

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	colmetricpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/oconf"
)

func TestGzipResponses(t *testing.T) {
	partial, err := proto.Marshal(&colmetricpb.ExportMetricsServiceResponse{
		PartialSuccess: &colmetricpb.ExportMetricsPartialSuccess{RejectedDataPoints: 2, ErrorMessage: "two rejected"},
	})
	require.NoError(t, err)
	tests := []struct {
		name      string
		status    int
		body      []byte
		malformed bool
		limit     int64
		wantErr   string
	}{
		{name: "success", status: http.StatusOK, limit: 1024},
		{name: "partial success", status: http.StatusOK, body: partial, limit: 1024, wantErr: "two rejected"},
		{
			name:    "error response",
			status:  http.StatusBadRequest,
			body:    []byte("rejected request"),
			limit:   1024,
			wantErr: "body: rejected request",
		},
		{
			name:    "decompressed size limit",
			status:  http.StatusOK,
			body:    bytes.Repeat([]byte("x"), 1024),
			limit:   64,
			wantErr: "response body too large: exceeded 64 bytes",
		},
		{
			name:      "malformed success",
			status:    http.StatusOK,
			body:      []byte("not a gzip response"),
			malformed: true,
			limit:     1024,
			wantErr:   "gzip: invalid header",
		},
		{
			name:      "malformed retryable response",
			status:    http.StatusServiceUnavailable,
			body:      []byte("not a gzip response"),
			malformed: true,
			limit:     1024,
			wantErr:   "gzip: invalid header",
		},
	}
	for _, tc := range tests {
		for _, explicit := range []bool{false, true} {
			name := tc.name + "/automatic"
			if explicit {
				name = tc.name + "/explicit"
			}
			t.Run(name, func(t *testing.T) {
				var calls atomic.Int32
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					calls.Add(1)
					w.Header().Set("Content-Encoding", "gzip")
					w.Header().Set("Content-Type", "application/x-protobuf")
					w.WriteHeader(tc.status)
					if tc.malformed {
						_, _ = w.Write(tc.body)
						return
					}
					zw := gzip.NewWriter(w)
					_, _ = zw.Write(tc.body)
					_ = zw.Close()
				}))
				t.Cleanup(srv.Close)
				opts := []Option{
					WithEndpointURL(srv.URL),
					WithMaxResponseSize(tc.limit),
					WithRetry(
						RetryConfig{
							Enabled:         true,
							InitialInterval: time.Millisecond,
							MaxInterval:     time.Millisecond,
							MaxElapsedTime:  time.Second,
						},
					),
				}
				if explicit {
					opts = append(opts, WithHeaders(map[string]string{"Accept-Encoding": "gzip"}))
				}
				cfg := oconf.NewHTTPConfig(asHTTPOptions(opts)...)
				client, err := newClient(cfg)
				require.NoError(t, err)
				t.Cleanup(func() { _ = client.Shutdown(t.Context()) })
				exportErr := client.UploadMetrics(t.Context(), &mpb.ResourceMetrics{})
				if tc.wantErr == "" {
					assert.NoError(t, exportErr)
				} else {
					assert.ErrorContains(t, exportErr, tc.wantErr)
				}
				assert.Equal(t, int32(1), calls.Load(), "response must not be retried")
			})
		}
	}
}
