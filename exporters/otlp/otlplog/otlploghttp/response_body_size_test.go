// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxResponseBodySizeOption(t *testing.T) {
	cfg := newConfig([]Option{WithMaxResponseBodySize(2)})
	assert.Equal(t, int64(2), cfg.maxResponseBodySize.Value)

	for _, size := range []int64{0, -1} {
		cfg := newConfig([]Option{WithMaxResponseBodySize(size)})
		assert.Equal(t, defaultMaxResponseBodySize, cfg.maxResponseBodySize.Value)
	}
}
