// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	colprofilepb "go.opentelemetry.io/proto/otlp/collector/profiles/v1development"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type testProfilesServer struct {
	colprofilepb.UnimplementedProfilesServiceServer

	mu       sync.Mutex
	requests []*colprofilepb.ExportProfilesServiceRequest
	mds      []metadata.MD
	handler  func(context.Context, *colprofilepb.ExportProfilesServiceRequest) (*colprofilepb.ExportProfilesServiceResponse, error)
}

func (s *testProfilesServer) Export(
	ctx context.Context,
	req *colprofilepb.ExportProfilesServiceRequest,
) (*colprofilepb.ExportProfilesServiceResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	s.mu.Lock()
	s.requests = append(s.requests, req)
	s.mds = append(s.mds, md.Copy())
	h := s.handler
	s.mu.Unlock()

	if h != nil {
		return h(ctx, req)
	}
	return &colprofilepb.ExportProfilesServiceResponse{}, nil
}

func (s *testProfilesServer) requestCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.requests)
}

func startTestGRPCServer(t *testing.T, srv *testProfilesServer, opts ...grpc.ServerOption) string {
	t.Helper()
	var lc net.ListenConfig
	lis, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	gs := grpc.NewServer(opts...)
	colprofilepb.RegisterProfilesServiceServer(gs, srv)
	go func() { _ = gs.Serve(lis) }()
	t.Cleanup(func() { gs.Stop() })
	return lis.Addr().String()
}

func TestClientContextCancellationAndTimeout(t *testing.T) {
	t.Run("pre_canceled_and_expired_upload_skip_server", func(t *testing.T) {
		srv := &testProfilesServer{}
		addr := startTestGRPCServer(t, srv)

		exp, err := New(t.Context(), WithEndpoint(addr), WithInsecure())
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		canceledCtx, cancel := context.WithCancel(t.Context())
		cancel()
		assert.ErrorIs(t, exp.Export(canceledCtx, nonEmptyResourceProfiles()), context.Canceled)

		expiredCtx, expire := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
		defer expire()
		assert.ErrorIs(t, exp.Export(expiredCtx, nonEmptyResourceProfiles()), context.DeadlineExceeded)

		assert.Zero(t, srv.requestCount(), "canceled or expired contexts must not send RPCs")
	})

	t.Run("shutdown_honors_context_error_and_closes_conn", func(t *testing.T) {
		srv := &testProfilesServer{}
		addr := startTestGRPCServer(t, srv)

		c, err := newClient(newConfig([]Option{WithEndpoint(addr), WithInsecure()}))
		require.NoError(t, err)
		conn := c.conn

		canceledCtx, cancel := context.WithCancel(t.Context())
		cancel()
		assert.ErrorIs(t, c.Shutdown(canceledCtx), context.Canceled)
		assert.Nil(t, c.conn)
		// Closing conn a second time confirms Shutdown closed it despite canceled ctx.
		assert.Error(t, conn.Close())
	})

	t.Run("shutdown_propagates_conn_close_error", func(t *testing.T) {
		srv := &testProfilesServer{}
		addr := startTestGRPCServer(t, srv)

		c, err := newClient(newConfig([]Option{WithEndpoint(addr), WithInsecure()}))
		require.NoError(t, err)
		// Pre-close owned connection so c.Shutdown observes a close error.
		require.NoError(t, c.conn.Close())

		err = c.Shutdown(t.Context())
		assert.Error(t, err)
	})

	t.Run("with_timeout_enforces_deadline_on_blocking_server", func(t *testing.T) {
		release := make(chan struct{})
		t.Cleanup(func() { close(release) })

		var attempts atomic.Int32
		srv := &testProfilesServer{
			handler: func(ctx context.Context, _ *colprofilepb.ExportProfilesServiceRequest) (*colprofilepb.ExportProfilesServiceResponse, error) {
				attempts.Add(1)
				select {
				case <-ctx.Done():
					return nil, status.FromContextError(ctx.Err()).Err()
				case <-release:
					return &colprofilepb.ExportProfilesServiceResponse{}, nil
				}
			},
		}
		addr := startTestGRPCServer(t, srv)

		exp, err := New(
			t.Context(),
			WithEndpoint(addr),
			WithInsecure(),
			WithTimeout(50*time.Millisecond),
			WithRetry(RetryConfig{
				Enabled:         true,
				InitialInterval: 5 * time.Millisecond,
				MaxInterval:     10 * time.Millisecond,
				MaxElapsedTime:  time.Second,
			}),
		)
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		err = exp.Export(t.Context(), nonEmptyResourceProfiles())
		require.Error(t, err)
		assert.ErrorContains(t, err, "exporter export timeout")
		assert.Equal(t, codes.DeadlineExceeded, status.Code(err))
		assert.GreaterOrEqual(t, attempts.Load(), int32(1))
	})

	t.Run("export_context_deadline_behavior", func(t *testing.T) {
		cPos := &client{exportTimeout: time.Second}
		ctxPos, cancelPos := cPos.exportContext(t.Context())
		defer cancelPos()
		_, hasDeadline := ctxPos.Deadline()
		assert.True(t, hasDeadline)

		for _, d := range []time.Duration{0, -time.Second} {
			cNoDeadline := &client{exportTimeout: d}
			ctxNoDeadline, cancelNoDeadline := cNoDeadline.exportContext(t.Context())
			_, hasDeadline = ctxNoDeadline.Deadline()
			cancelNoDeadline()
			assert.False(t, hasDeadline, "timeout=%v should not add a deadline", d)
		}
	})
}
