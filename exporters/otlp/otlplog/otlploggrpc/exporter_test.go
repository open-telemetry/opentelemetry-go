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

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	collogpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	logpb "go.opentelemetry.io/proto/otlp/logs/v1"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
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
	errClient := errors.New("client")

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
			err:     errClient,
			wantErr: errClient,
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

			err := e.Export(t.Context(), tc.logs)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.logs, got)
			assert.Equal(t, 1, mockCli.uploads)
		})
	}
}

func TestExporterShutdown(t *testing.T) {
	ctx := t.Context()
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
	ctx := t.Context()
	e, err := New(ctx)
	require.NoError(t, err, "New")

	assert.NoError(t, e.ForceFlush(ctx), "ForceFlush")
}

func TestExporterConcurrentSafe(t *testing.T) {
	e := newExporter(&mockClient{})

	const goroutines = 10

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(t.Context())
	var runs atomic.Uint64
	for range goroutines {
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
					runs.Add(1)
				}
			}
		}()
	}

	for runs.Load() == 0 {
		runtime.Gosched()
	}

	_ = e.Shutdown(ctx)
	cancel()
	wg.Wait()
}

// TestExporter runs integration test against the real OTLP collector.
func TestExporter(t *testing.T) {
	t.Run("ExporterHonorsContextErrors", func(t *testing.T) {
		t.Run("Export", testCtxErrs(func() func(context.Context) error {
			c, _ := clientFactory(t, nil)
			e := newExporter(c)
			return func(ctx context.Context) error {
				return e.Export(ctx, []sdklog.Record{{}})
			}
		}))

		t.Run("Shutdown", testCtxErrs(func() func(context.Context) error {
			c, _ := clientFactory(t, nil)
			e := newExporter(c)
			return e.Shutdown
		}))
	})

	t.Run("Export", func(t *testing.T) {
		ctx := t.Context()
		c, coll := clientFactory(t, nil)
		e := newExporter(c)

		require.NoError(t, e.Export(ctx, records))
		require.NoError(t, e.Shutdown(ctx))
		got := coll.Collect().Dump()
		require.Len(t, got, 1, "upload of one ResourceLogs")
		require.Len(t, got[0].ScopeLogs, 1, "upload of one ScopeLogs")
		require.Len(t, got[0].ScopeLogs[0].LogRecords, 2, "upload of two ScopeLogs")

		// Check body
		assert.Equal(t, "A", got[0].ScopeLogs[0].LogRecords[0].Body.GetStringValue())
		assert.Equal(t, "B", got[0].ScopeLogs[0].LogRecords[1].Body.GetStringValue())
	})

	t.Run("PartialSuccess", func(t *testing.T) {
		const n, msg = 2, "bad data"
		rCh := make(chan exportResult, 3)
		rCh <- exportResult{
			Response: &collogpb.ExportLogsServiceResponse{
				PartialSuccess: &collogpb.ExportLogsPartialSuccess{
					RejectedLogRecords: n,
					ErrorMessage:       msg,
				},
			},
		}
		rCh <- exportResult{
			Response: &collogpb.ExportLogsServiceResponse{
				PartialSuccess: &collogpb.ExportLogsPartialSuccess{
					// Should not be logged.
					RejectedLogRecords: 0,
					ErrorMessage:       "",
				},
			},
		}
		rCh <- exportResult{
			Response: &collogpb.ExportLogsServiceResponse{},
		}

		ctx := t.Context()
		c, _ := clientFactory(t, rCh)
		e := newExporter(c)

		assert.ErrorIs(t, e.Export(ctx, records), internal.PartialSuccess{})
		assert.NoError(t, e.Export(ctx, records))
		assert.NoError(t, e.Export(ctx, records))
	})
}
