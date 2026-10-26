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
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	collogpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	logpb "go.opentelemetry.io/proto/otlp/logs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel/attribute"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

var records []sdklog.Record

func init() {
	var r sdklog.Record
	r.SetTimestamp(ts)
	r.SetBody(attribute.StringValue("A"))
	records = append(records, r)

	r.SetBody(attribute.StringValue("B"))
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

	r := make([]sdklog.Record, 1)
	assert.ErrorIs(t, e.Export(ctx, r), sdklog.ErrExporterShutdown, "Export on Shutdown Exporter")
	assert.NoError(t, e.ForceFlush(ctx), "ForceFlush on Shutdown Exporter")
	assert.NoError(t, e.Shutdown(ctx), "Shutdown on Shutdown Exporter")
}

func TestExporterExportConcurrentShutdown(t *testing.T) {
	orig := transformResourceLogs
	started := make(chan struct{})
	resume := make(chan struct{})
	transformResourceLogs = func([]sdklog.Record) []*logpb.ResourceLogs {
		close(started)
		<-resume
		return make([]*logpb.ResourceLogs, 1)
	}
	t.Cleanup(func() { transformResourceLogs = orig })

	c := &mockClient{}
	e := newExporter(c)
	ctx := t.Context()
	errCh := make(chan error, 1)
	go func() {
		errCh <- e.Export(ctx, make([]sdklog.Record, 1))
	}()

	<-started
	assert.NoError(t, e.Shutdown(ctx), "Shutdown Exporter")
	close(resume)

	assert.ErrorIs(t, <-errCh, sdklog.ErrExporterShutdown)
	assert.Zero(t, c.uploads, "client UploadLogs calls")
}

type shutdownTrackingClient struct {
	logClient
	shutdowns atomic.Int32
	done      chan struct{}
}

func (c *shutdownTrackingClient) Shutdown(ctx context.Context) error {
	c.shutdowns.Add(1)
	defer close(c.done)
	return c.logClient.Shutdown(ctx)
}

func TestExporterShutdownDuringExport(t *testing.T) {
	for _, canceled := range []bool{false, true} {
		name := "DeadlineExceeded"
		if canceled {
			name = "Canceled"
		}
		t.Run(name, func(t *testing.T) {
			for _, owned := range []bool{false, true} {
				name := "InjectedConnection"
				if owned {
					name = "OwnedConnection"
				}
				t.Run(name, func(t *testing.T) {
					started := make(chan struct{})
					release := make(chan struct{})
					var startedOnce, releaseOnce sync.Once
					finishExport := func() { releaseOnce.Do(func() { close(release) }) }
					coll, err := newGRPCCollector(t.Context(), "", nil, grpc.UnaryInterceptor(
						func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
							startedOnce.Do(func() { close(started) })
							select {
							case <-release:
								return handler(ctx, req)
							case <-ctx.Done():
								return nil, ctx.Err()
							}
						},
					))
					require.NoError(t, err)
					t.Cleanup(coll.srv.Stop)
					t.Cleanup(finishExport)

					endpoint := coll.listener.Addr().String()
					opts := []Option{WithEndpoint(endpoint), WithInsecure(), WithTimeout(5 * time.Second)}
					if !owned {
						conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
						require.NoError(t, err)
						t.Cleanup(func() { assert.NoError(t, conn.Close()) })
						opts = append(opts, WithGRPCConn(conn))
					}
					c, err := newClient(newConfig(opts))
					require.NoError(t, err)
					conn := c.conn
					tracked := &shutdownTrackingClient{logClient: c, done: make(chan struct{})}
					e := newExporter(tracked)
					t.Cleanup(func() { _ = e.Shutdown(t.Context()) })

					exportDone := make(chan error, 1)
					go func() { exportDone <- e.Export(t.Context(), records) }()
					<-started

					var ctx context.Context
					var cancel context.CancelFunc
					if canceled {
						ctx, cancel = context.WithCancel(t.Context())
					} else {
						ctx, cancel = context.WithTimeout(t.Context(), 10*time.Millisecond)
					}
					t.Cleanup(cancel)
					shutdownDone := make(chan error, 1)
					go func() { shutdownDone <- e.Shutdown(ctx) }()
					require.Eventually(t, e.stopped.Load, time.Second, time.Millisecond)
					if canceled {
						cancel()
					}
					<-ctx.Done()

					var shutdownErr error
					returned := false
					select {
					case shutdownErr = <-shutdownDone:
						returned = true
					case <-time.After(time.Second):
						t.Error("Shutdown did not return after its context expired while Export was in flight")
					}

					assert.Zero(t, tracked.shutdowns.Load(), "cleanup must wait for the active export")
					assert.NotEqual(t, connectivity.Shutdown, conn.GetState())
					assert.ErrorIs(t, e.Export(t.Context(), records), sdklog.ErrExporterShutdown)
					assert.NoError(t, e.Shutdown(t.Context()), "repeated shutdown must not wait for cleanup")

					finishExport()
					require.NoError(t, <-exportDone)
					if !returned {
						shutdownErr = <-shutdownDone
					}
					assert.ErrorIs(t, shutdownErr, ctx.Err())
					<-tracked.done
					assert.EqualValues(t, 1, tracked.shutdowns.Load())
					if owned {
						assert.Equal(t, connectivity.Shutdown, conn.GetState())
					} else {
						_, err := collogpb.NewLogsServiceClient(conn).
							Export(t.Context(), &collogpb.ExportLogsServiceRequest{})
						assert.NoError(t, err, "injected connection remains usable")
					}
				})
			}
		})
	}
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
		wg.Go(func() {
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
		})
	}

	for runs.Load() == 0 {
		runtime.Gosched()
	}

	_ = e.Shutdown(ctx)
	cancel()
	wg.Wait()
}

// TestExporter runs an integration test against the real OTLP collector.
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
