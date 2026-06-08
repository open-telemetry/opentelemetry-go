// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracehttp

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryAfterUsesSeconds(t *testing.T) {
	err := newResponseError(http.Header{"Retry-After": {"10"}}, nil)
	_, throttle := evaluate(err)
	assert.Equal(t, 10*time.Second, throttle)
}
