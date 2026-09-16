// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package finish

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryShutdown(t *testing.T) {
	var calls atomic.Int64
	registry := &Registry{}
	registry.Register(func() { calls.Add(1) })

	registry.Shutdown()
	registry.Shutdown()
	registry.Register(func() { calls.Add(1) })

	assert.Equal(t, int64(2), calls.Load())
}
