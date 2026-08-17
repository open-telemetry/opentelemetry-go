// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleProcessorConcurrentSafe(t *testing.T) {
	var active atomic.Int32
	call := func() {
		assert.Equal(t, int32(1), active.Add(1), "concurrent exporter calls")
		time.Sleep(time.Microsecond)
		active.Add(-1)
	}
	e := &testExporter{
		ExportFunc: func(context.Context, []Record) error {
			call()
			return nil
		},
		ForceFlushFunc: func(context.Context) error {
			call()
			return nil
		},
		ShutdownFunc: func(context.Context) error {
			call()
			return nil
		},
	}
	s := NewSimpleProcessor(e)

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			_ = s.OnEmit(t.Context(), new(Record))
			_ = s.ForceFlush(t.Context())
			_ = s.Shutdown(t.Context())
		})
	}
	wg.Wait()

	assert.Zero(t, active.Load(), "active exporter calls")
}
