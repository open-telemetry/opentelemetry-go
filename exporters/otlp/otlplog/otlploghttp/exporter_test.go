// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/log"
	logpb "go.opentelemetry.io/proto/otlp/logs/v1"
)

func TestExporterExportErrors(t *testing.T) {
	errUpload := errors.New("upload")
	c := &client{
		uploadLogs: func(context.Context, []*logpb.ResourceLogs) error {
			return errUpload
		},
	}

	e, err := newExporter(c, config{})
	require.NoError(t, err, "New")

	err = e.Export(context.Background(), make([]log.Record, 1))
	assert.ErrorIs(t, err, errUpload)
}

func TestExporterExport(t *testing.T) {
	var uploads int
	c := &client{
		uploadLogs: func(context.Context, []*logpb.ResourceLogs) error {
			uploads++
			return nil
		},
	}

	orig := transformResourceLogs
	var got []log.Record
	transformResourceLogs = func(r []log.Record) []*logpb.ResourceLogs {
		got = r
		return make([]*logpb.ResourceLogs, 1)
	}
	t.Cleanup(func() { transformResourceLogs = orig })

	e, err := newExporter(c, config{})
	require.NoError(t, err, "New")

	ctx := context.Background()
	want := make([]log.Record, 1)
	assert.NoError(t, e.Export(ctx, want))

	assert.Equal(t, 1, uploads, "client UploadLogs calls")
	assert.Equal(t, want, got, "transformed log records")
}

func TestExporterShutdown(t *testing.T) {
	ctx := context.Background()
	e, err := New(ctx)
	require.NoError(t, err, "New")
	assert.NoError(t, e.Shutdown(ctx), "Shutdown Exporter")

	// After Shutdown is called, calls to Export, Shutdown, or ForceFlush
	// should perform no operation and return nil error.
	r := make([]log.Record, 1)
	assert.NoError(t, e.Export(ctx, r), "Export on Shutdown Exporter")
	assert.NoError(t, e.ForceFlush(ctx), "ForceFlush on Shutdown Exporter")
	assert.NoError(t, e.Shutdown(ctx), "Shutdown on Shutdown Exporter")
}

func TestExporterForceFlush(t *testing.T) {
	ctx := context.Background()
	e, err := New(ctx)
	require.NoError(t, err, "New")

	assert.NoError(t, e.ForceFlush(ctx), "ForceFlush")
}

func TestExporterConcurrentSafe(t *testing.T) {
	ctx := context.Background()
	e, err := New(ctx)
	require.NoError(t, err, "newExporter")

	const goroutines = 10

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	runs := new(uint64)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := make([]log.Record, 1)
			for {
				select {
				case <-ctx.Done():
					return
				default:
					_ = e.Export(ctx, r)
					_ = e.ForceFlush(ctx)
					atomic.AddUint64(runs, 1)
				}
			}
		}()
	}

	for atomic.LoadUint64(runs) == 0 {
		runtime.Gosched()
	}

	_ = e.Shutdown(ctx)
	cancel()
	wg.Wait()
}
