// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetrichttp

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

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

func TestCopyResponseBodyMaxInt64StaysBounded(t *testing.T) {
	var dst bytes.Buffer
	err := copyResponseBody(&dst, io.NopCloser(bytes.NewReader([]byte("ok"))), 1<<63-1)
	assert.NoError(t, err)
	assert.Equal(t, "ok", dst.String())
}
