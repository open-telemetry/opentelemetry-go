// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracehttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"testing"
	"time"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
)

func TestRetryAfterUsesSeconds(t *testing.T) {
	err := newResponseError(http.Header{"Retry-After": {"10"}}, nil)
	_, throttle := evaluate(err)
	assert.Equal(t, 10*time.Second, throttle)
}

func TestRetryAfterUsesHTTPDate(t *testing.T) {
	date := time.Now().UTC().Add(time.Hour).Format(http.TimeFormat)
	err := newResponseError(http.Header{"Retry-After": {date}}, nil)
	_, throttle := evaluate(err)
	assert.Greater(t, throttle, 59*time.Minute)
	assert.LessOrEqual(t, throttle, time.Hour)
}

func TestRetryAfterSecondsOverflow(t *testing.T) {
	err := newResponseError(http.Header{"Retry-After": {"9223372036854775807"}}, nil)
	_, throttle := evaluate(err)
	assert.Equal(t, time.Duration(math.MaxInt64), throttle)
}

func TestClientMarshalLogDoesNotIncludeEndpointConfig(t *testing.T) {
	const sensitiveEndpoint = "user:pass@collector.internal:4318"

	var buf bytes.Buffer
	logger := funcr.New(func(_, args string) {
		_, _ = buf.WriteString(args)
	}, funcr.Options{})

	client := NewClient(WithEndpoint(sensitiveEndpoint), WithInsecure())
	logger.Info("client", "config", client)

	logged := buf.String()
	assert.Contains(t, logged, "otlptracehttp")
	assert.NotContains(t, logged, sensitiveEndpoint)
	assert.NotContains(t, logged, "Insecure")
}

func TestRetryableNetworkError(t *testing.T) {
	live := context.Background()
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	expiredCtx, cancelExpired := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancelExpired()

	tests := []struct {
		name string
		ctx  context.Context
		err  error
		want bool
	}{
		{name: "nil", ctx: live},

		// Cancellation: nothing would be delivered.
		{name: "explicitly canceled error", ctx: live, err: context.Canceled},
		{name: "canceled export context", ctx: canceledCtx, err: io.EOF},
		{name: "export context deadline passed", ctx: expiredCtx, err: &timeoutError{}},

		// A client-level timeout is transient while the export is still running:
		// http.Client.Timeout wraps context.DeadlineExceeded just like the export
		// context's own deadline does, so the context has to decide.
		{
			name: "client timeout while the export is live",
			ctx:  live,
			err:  &url.Error{Op: "Post", URL: "https://example.com", Err: &timeoutError{}},
			want: true,
		},

		// Transient disconnects.
		{name: "EOF", ctx: live, err: io.EOF, want: true},
		{name: "unexpected EOF", ctx: live, err: io.ErrUnexpectedEOF, want: true},
		{
			name: "url wrapped EOF",
			ctx:  live,
			err:  &url.Error{Op: "Post", URL: "https://example.com", Err: io.EOF},
			want: true,
		},
		{
			name: "connection refused",
			ctx:  live,
			err: &url.Error{
				Op:  "Post",
				URL: "https://example.com",
				Err: &net.OpError{Op: "dial", Net: "tcp", Err: syscall.ECONNREFUSED},
			},
			want: true,
		},
		{
			name: "connection reset",
			ctx:  live,
			err: &url.Error{
				Op:  "Post",
				URL: "https://example.com",
				Err: &net.OpError{Op: "read", Net: "tcp", Err: syscall.ECONNRESET},
			},
			want: true,
		},

		// Permanent configuration and trust failures must fail fast.
		{
			name: "invalid port",
			ctx:  live,
			err: &url.Error{
				Op:  "Post",
				URL: "https://example.com",
				Err: &net.OpError{Op: "dial", Net: "tcp",
					Err: &net.AddrError{Err: "invalid port", Addr: "127.0.0.1:99999"}},
			},
		},
		{
			name: "tls handshake alert",
			ctx:  live,
			err: &url.Error{
				Op:  "Post",
				URL: "https://example.com",
				Err: &net.OpError{Op: "remote error", Net: "tcp", Err: tls.AlertError(42)},
			},
		},
		{
			name: "certificate verification",
			ctx:  live,
			err: &url.Error{
				Op:  "Post",
				URL: "https://example.com",
				Err: &tls.CertificateVerificationError{Err: x509.UnknownAuthorityError{}},
			},
		},
		{name: "non-network error stays permanent", ctx: live, err: errors.New("invalid configuration")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, retryableNetworkError(tt.ctx, tt.err))
		})
	}
}

// timeoutError implements net.Error with Timeout() true.
type timeoutError struct{}

func (*timeoutError) Error() string   { return "i/o timeout" }
func (*timeoutError) Timeout() bool   { return true }
func (*timeoutError) Temporary() bool { return true }
