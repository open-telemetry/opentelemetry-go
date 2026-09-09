// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracehttp

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
)

func TestCopyResponseBody(t *testing.T) {
	readErr := errors.New("read error")
	tests := []struct {
		name    string
		src     io.Reader
		limit   int64
		want    string
		wantErr string
	}{
		{name: "Empty", src: strings.NewReader(""), limit: 1},
		{name: "UnderLimit", src: strings.NewReader("x"), limit: 2, want: "x"},
		{name: "AtLimit", src: strings.NewReader("x"), limit: 1, want: "x"},
		{
			name:    "OverLimit",
			src:     strings.NewReader("xx"),
			limit:   1,
			want:    "x",
			wantErr: "response body too large: exceeded 1 bytes",
		},
		{name: "MaxInt64", src: strings.NewReader("ok"), limit: 1<<63 - 1, want: "ok"},
		{name: "CopyError", src: iotest.ErrReader(readErr), limit: 1, wantErr: readErr.Error()},
		{
			name:    "ProbeError",
			src:     io.MultiReader(strings.NewReader("x"), iotest.ErrReader(readErr)),
			limit:   1,
			want:    "x",
			wantErr: readErr.Error(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var dst bytes.Buffer
			err := copyResponseBody(&dst, test.src, test.limit)
			if test.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.wantErr)
			}
			assert.Equal(t, test.want, dst.String())
		})
	}
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
