// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracehttp

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
)

type responseBodyErrorReader struct {
	err   error
	first bool
}

func (r *responseBodyErrorReader) Read(p []byte) (int, error) {
	if r.first {
		return 0, r.err
	}
	r.first = true
	p[0] = 'x'
	return 1, nil
}

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
	assert.Equal(t, time.Duration(1<<63-1), throttle)
}

func TestCopyResponseBodyReadErrors(t *testing.T) {
	readErr := errors.New("read error")

	var dst bytes.Buffer
	err := copyResponseBody(&dst, io.NopCloser(&responseBodyErrorReader{err: readErr, first: true}), 0)
	assert.ErrorIs(t, err, readErr)

	dst.Reset()
	err = copyResponseBody(&dst, io.NopCloser(&responseBodyErrorReader{err: readErr, first: true}), 1)
	assert.ErrorIs(t, err, readErr)

	dst.Reset()
	err = copyResponseBody(&dst, io.NopCloser(&responseBodyErrorReader{err: readErr}), 1)
	assert.ErrorIs(t, err, readErr)
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
