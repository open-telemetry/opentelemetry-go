// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	logpb "go.opentelemetry.io/proto/otlp/logs/v1"
)

var records []sdklog.Record

func init() {
	var r sdklog.Record
	r.SetTimestamp(ts)
	r.SetBody(log.StringValue("A"))
	records = append(records, r)

	r.SetBody(log.StringValue("B"))
	records = append(records, r)
}

type mockClient struct {
	err error

	uploads int
}

func (m *mockClient) UploadLogs(context.Context, []*logpb.ResourceLogs) error {
	m.uploads++
	return m.err
}

func (m *mockClient) Shutdown(context.Context) error {
	return m.err
}

func TestExporterExport(t *testing.T) {
	testCases := []struct {
		name string
		logs []sdklog.Record
		err  error

		wantLogs []sdklog.Record
		wantErr  error
	}{
		{
			name:     "NoError",
			logs:     make([]sdklog.Record, 2),
			wantLogs: make([]sdklog.Record, 2),
		},
		{
			name:    "Error",
			logs:    make([]sdklog.Record, 2),
			err:     errors.New("test"),
			wantErr: errors.New("test"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orig := transformResourceLogs
			var got []sdklog.Record
			transformResourceLogs = func(r []sdklog.Record) []*logpb.ResourceLogs {
				got = r
				return make([]*logpb.ResourceLogs, len(r))
			}
			t.Cleanup(func() { transformResourceLogs = orig })

			mockCli := mockClient{err: tc.err}

			e := newExporter(&mockCli)

			err := e.Export(context.Background(), tc.logs)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.logs, got)
			assert.Equal(t, 1, mockCli.uploads)
		})
	}
}

func TestExporterShutdown(t *testing.T) {
	ctx := context.Background()
	e, err := New(ctx)
	require.NoError(t, err, "New")
	assert.NoError(t, e.Shutdown(ctx), "Shutdown Exporter")

	// After Shutdown is called, calls to Export, Shutdown, or ForceFlush
	// should perform no operation and return nil error.
	r := make([]sdklog.Record, 1)
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
	e := newExporter(&mockClient{})

	const goroutines = 10

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	runs := new(uint64)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := make([]sdklog.Record, 1)
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
