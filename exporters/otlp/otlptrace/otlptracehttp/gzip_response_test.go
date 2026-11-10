// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracehttp_test

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
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func TestGzipResponses(t *testing.T) {
	partial, err := proto.Marshal(&coltracepb.ExportTraceServiceResponse{
		PartialSuccess: &coltracepb.ExportTracePartialSuccess{RejectedSpans: 2, ErrorMessage: "two rejected"},
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
				opts := []otlptracehttp.Option{
					otlptracehttp.WithEndpointURL(srv.URL),
					otlptracehttp.WithMaxResponseSize(tc.limit),
					otlptracehttp.WithRetry(
						otlptracehttp.RetryConfig{
							Enabled:         true,
							InitialInterval: time.Millisecond,
							MaxInterval:     time.Millisecond,
							MaxElapsedTime:  time.Second,
						},
					),
				}
				if explicit {
					opts = append(opts, otlptracehttp.WithHeaders(map[string]string{"Accept-Encoding": "gzip"}))
				}
				client := otlptracehttp.NewClient(opts...)
				require.NoError(t, client.Start(t.Context()))
				t.Cleanup(func() { _ = client.Stop(t.Context()) })
				exportErr := client.UploadTraces(t.Context(), []*tracepb.ResourceSpans{{}})
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
