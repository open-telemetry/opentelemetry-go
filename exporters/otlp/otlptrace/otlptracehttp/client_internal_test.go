// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracehttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"syscall"
	"testing"
	"time"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	canceled, cancel := context.WithCancel(t.Context())
	cancel()

	tests := []struct {
		name string
		ctx  context.Context
		err  error
		want bool
	}{
		{name: "nil"},
		{name: "canceled sentinel", err: context.Canceled},
		{name: "canceled retry context", ctx: canceled, err: io.EOF},
		{
			// http.Client.Timeout errors match context.DeadlineExceeded.
			name: "deadline exceeded with live retry context",
			err:  context.DeadlineExceeded,
			want: true,
		},
		{name: "EOF", err: io.EOF, want: true},
		{name: "unexpected EOF", err: io.ErrUnexpectedEOF, want: true},
		{
			name: "url wrapped unexpected EOF",
			err:  &url.Error{Op: "Post", URL: "https://example.com", Err: io.ErrUnexpectedEOF},
			want: true,
		},
		{
			name: "connection refused",
			err: &url.Error{
				Op:  "Post",
				URL: "http://127.0.0.1:1",
				Err: &net.OpError{Op: "dial", Net: "tcp", Err: syscall.ECONNREFUSED},
			},
			want: true,
		},
		{
			name: "connection reset",
			err: &url.Error{
				Op:  "Post",
				URL: "http://127.0.0.1:1",
				Err: &net.OpError{Op: "read", Net: "tcp", Err: syscall.ECONNRESET},
			},
			want: true,
		},
		{
			name: "timeout",
			err:  &url.Error{Op: "Post", URL: "https://example.com", Err: timeoutError{}},
			want: true,
		},
		{
			name: "tls alert",
			err: &url.Error{
				Op:  "Post",
				URL: "https://example.com",
				Err: &net.OpError{Op: "remote error", Err: errors.New("tls: protocol version not supported")},
			},
		},
		{
			name: "invalid address",
			err: &url.Error{
				Op:  "Post",
				URL: "http://127.0.0.1:65536",
				Err: &net.OpError{Op: "dial", Net: "tcp", Err: &net.AddrError{Err: "invalid port", Addr: "65536"}},
			},
		},
		{
			name: "dns not found",
			err: &url.Error{
				Op:  "Post",
				URL: "http://collector.invalid",
				Err: &net.OpError{Op: "dial", Net: "tcp", Err: &net.DNSError{Err: "no such host", IsNotFound: true}},
			},
		},
		{
			name: "tls verification error",
			err: &url.Error{
				Op:  "Post",
				URL: "https://example.com",
				Err: errors.New("tls: failed to verify certificate"),
			},
		},
		{name: "generic permanent error", err: errors.New("boom")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = t.Context()
			}
			assert.Equal(t, tt.want, retryableNetworkError(ctx, tt.err))
		})
	}
}

// TestRetryableNetworkErrorRealErrors classifies errors produced by real
// net/http requests, so the classifier is checked against the actual error
// chains rather than hand-built ones.
func TestRetryableNetworkErrorRealErrors(t *testing.T) {
	post := func(t *testing.T, c *http.Client, target string) error {
		t.Helper()
		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, target, http.NoBody)
		require.NoError(t, err)
		resp, err := c.Do(req)
		if resp != nil {
			_ = resp.Body.Close()
		}
		require.Error(t, err)
		return err
	}

	t.Run("client timeout", func(t *testing.T) {
		release := make(chan struct{})
		srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			<-release
		}))
		t.Cleanup(srv.Close)
		t.Cleanup(func() { close(release) })

		err := post(t, &http.Client{Timeout: 10 * time.Millisecond}, srv.URL)
		assert.True(t, retryableNetworkError(t.Context(), err), "err: %v", err)
	})

	t.Run("tls protocol version rejected", func(t *testing.T) {
		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		srv.TLS = &tls.Config{MinVersion: tls.VersionTLS13}
		srv.StartTLS()
		t.Cleanup(srv.Close)

		c := srv.Client()
		tr := c.Transport.(*http.Transport)
		tr.TLSClientConfig.MaxVersion = tls.VersionTLS12

		err := post(t, c, srv.URL)
		assert.False(t, retryableNetworkError(t.Context(), err), "err: %v", err)
	})

	t.Run("invalid port", func(t *testing.T) {
		err := post(t, &http.Client{}, "http://127.0.0.1:65536/")
		assert.False(t, retryableNetworkError(t.Context(), err), "err: %v", err)
	})
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }
