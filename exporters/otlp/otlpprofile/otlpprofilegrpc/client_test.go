// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	colprofilepb "go.opentelemetry.io/proto/otlp/collector/profiles/v1development"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
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

func TestRetryableAndThrottleDelay(t *testing.T) {
	t.Run("all_grpc_codes", func(t *testing.T) {
		cases := []struct {
			code          codes.Code
			wantRetryable bool
		}{
			{codes.OK, false},
			{codes.Canceled, true},
			{codes.Unknown, false},
			{codes.InvalidArgument, false},
			{codes.DeadlineExceeded, true},
			{codes.NotFound, false},
			{codes.AlreadyExists, false},
			{codes.PermissionDenied, false},
			{codes.ResourceExhausted, false},
			{codes.FailedPrecondition, false},
			{codes.Aborted, true},
			{codes.OutOfRange, true},
			{codes.Unimplemented, false},
			{codes.Internal, false},
			{codes.Unavailable, true},
			{codes.DataLoss, true},
			{codes.Unauthenticated, false},
		}

		for _, tc := range cases {
			t.Run(tc.code.String(), func(t *testing.T) {
				ok, delay := retryable(status.Error(tc.code, "msg"))
				assert.Equal(t, tc.wantRetryable, ok)
				assert.Zero(t, delay)
			})
		}

		ok, delay := retryable(errors.New("non-grpc error"))
		assert.False(t, ok)
		assert.Zero(t, delay)
	})

	t.Run("throttle_delay_details", func(t *testing.T) {
		stNone := status.New(codes.Unavailable, "no details")
		ok, d := throttleDelay(stNone)
		assert.False(t, ok)
		assert.Zero(t, d)

		stErrInfo, err := status.New(codes.ResourceExhausted, "errinfo").WithDetails(&errdetails.ErrorInfo{Reason: "LIMIT"})
		require.NoError(t, err)
		ok, d = retryable(stErrInfo.Err())
		assert.False(t, ok)
		assert.Zero(t, d)

		stMixed, err := status.New(codes.ResourceExhausted, "mixed").WithDetails(
			&errdetails.ErrorInfo{Reason: "LIMIT"},
			&errdetails.RetryInfo{RetryDelay: durationpb.New(25 * time.Millisecond)},
		)
		require.NoError(t, err)
		ok, d = retryable(stMixed.Err())
		assert.True(t, ok)
		assert.Equal(t, 25*time.Millisecond, d)

		stMulti, err := status.New(codes.Unavailable, "multi").WithDetails(
			&errdetails.RetryInfo{RetryDelay: durationpb.New(10 * time.Millisecond)},
			&errdetails.RetryInfo{RetryDelay: durationpb.New(99 * time.Millisecond)},
		)
		require.NoError(t, err)
		ok, d = retryable(stMulti.Err())
		assert.True(t, ok)
		assert.Equal(t, 10*time.Millisecond, d)
	})

	t.Run("resource_exhausted_without_retry_info_not_retried", func(t *testing.T) {
		var attempts atomic.Int32
		srv := &testProfilesServer{
			handler: func(context.Context, *colprofilepb.ExportProfilesServiceRequest) (*colprofilepb.ExportProfilesServiceResponse, error) {
				attempts.Add(1)
				return nil, status.Error(codes.ResourceExhausted, "quota exceeded")
			},
		}
		addr := startTestGRPCServer(t, srv)

		exp, err := New(t.Context(), WithEndpoint(addr), WithInsecure(), WithRetry(RetryConfig{
			Enabled:         true,
			InitialInterval: time.Millisecond,
			MaxInterval:     5 * time.Millisecond,
			MaxElapsedTime:  time.Second,
		}))
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		err = exp.Export(t.Context(), nonEmptyResourceProfiles())
		require.Error(t, err)
		assert.Equal(t, codes.ResourceExhausted, status.Code(err))
		assert.Equal(t, int32(1), attempts.Load())
	})

	t.Run("retry_disabled_does_not_retry_unavailable", func(t *testing.T) {
		var attempts atomic.Int32
		srv := &testProfilesServer{
			handler: func(context.Context, *colprofilepb.ExportProfilesServiceRequest) (*colprofilepb.ExportProfilesServiceResponse, error) {
				attempts.Add(1)
				return nil, status.Error(codes.Unavailable, "down")
			},
		}
		addr := startTestGRPCServer(t, srv)

		exp, err := New(t.Context(), WithEndpoint(addr), WithInsecure(), WithRetry(RetryConfig{Enabled: false}))
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		err = exp.Export(t.Context(), nonEmptyResourceProfiles())
		require.Error(t, err)
		assert.Equal(t, codes.Unavailable, status.Code(err))
		assert.Equal(t, int32(1), attempts.Load())
	})

	t.Run("retry_exhausted_returns_wrapped_error", func(t *testing.T) {
		var attempts atomic.Int32
		srv := &testProfilesServer{
			handler: func(context.Context, *colprofilepb.ExportProfilesServiceRequest) (*colprofilepb.ExportProfilesServiceResponse, error) {
				attempts.Add(1)
				return nil, status.Error(codes.Unavailable, "still down")
			},
		}
		addr := startTestGRPCServer(t, srv)

		exp, err := New(t.Context(), WithEndpoint(addr), WithInsecure(), WithRetry(RetryConfig{
			Enabled:         true,
			InitialInterval: time.Millisecond,
			MaxInterval:     2 * time.Millisecond,
			MaxElapsedTime:  15 * time.Millisecond,
		}))
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		err = exp.Export(t.Context(), nonEmptyResourceProfiles())
		require.Error(t, err)
		assert.ErrorContains(t, err, "max retry time elapsed")
		assert.Equal(t, codes.Unavailable, status.Code(err))
		assert.Greater(t, attempts.Load(), int32(1))
	})
}

func TestClientTLSAndDialOptions(t *testing.T) {
	newMTLSServer := func(t *testing.T, certPath, keyPath string) (string, *testProfilesServer) {
		t.Helper()
		serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
		require.NoError(t, err)
		caPEM, err := os.ReadFile(certPath)
		require.NoError(t, err)
		clientPool := x509.NewCertPool()
		require.True(t, clientPool.AppendCertsFromPEM(caPEM))

		srv := &testProfilesServer{}
		addr := startTestGRPCServer(t, srv, grpc.Creds(credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    clientPool,
			MinVersion:   tls.VersionTLS12,
		})))
		return addr, srv
	}

	t.Run("mtls_via_env_vars", func(t *testing.T) {
		certPath, keyPath := writeTestCertAndKey(t)
		addr, srv := newMTLSServer(t, certPath, keyPath)

		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_CERTIFICATE", certPath)
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_CLIENT_CERTIFICATE", certPath)
		t.Setenv("OTEL_EXPORTER_OTLP_PROFILES_CLIENT_KEY", keyPath)

		exp, err := New(t.Context(), WithEndpoint(addr))
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		require.NoError(t, exp.Export(t.Context(), nonEmptyResourceProfiles()))
		assert.Equal(t, 1, srv.requestCount())
	})

	t.Run("with_tls_credentials_overrides_insecure", func(t *testing.T) {
		certPath, keyPath := writeTestCertAndKey(t)
		addr, srv := newMTLSServer(t, certPath, keyPath)

		caPool, err := loadCertPool(certPath)
		require.NoError(t, err)
		certs, err := loadCertificates(certPath, keyPath)
		require.NoError(t, err)

		exp, err := New(
			t.Context(),
			WithEndpoint(addr),
			WithInsecure(),
			WithTLSCredentials(credentials.NewTLS(&tls.Config{
				RootCAs:      caPool,
				Certificates: certs,
				MinVersion:   tls.VersionTLS12,
			})),
		)
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		require.NoError(t, exp.Export(t.Context(), nonEmptyResourceProfiles()))
		assert.Equal(t, 1, srv.requestCount())
	})

	t.Run("with_insecure_overrides_env_certificate", func(t *testing.T) {
		certPath, _ := writeTestCertAndKey(t)
		t.Setenv("OTEL_EXPORTER_OTLP_CERTIFICATE", certPath)

		srv := &testProfilesServer{}
		addr := startTestGRPCServer(t, srv)

		exp, err := New(t.Context(), WithEndpoint(addr), WithInsecure())
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		require.NoError(t, exp.Export(t.Context(), nonEmptyResourceProfiles()))
		assert.Equal(t, 1, srv.requestCount())
	})

	t.Run("default_tls_fails_against_plaintext_server", func(t *testing.T) {
		srv := &testProfilesServer{}
		addr := startTestGRPCServer(t, srv)

		exp, err := New(
			t.Context(),
			WithEndpoint(addr),
			WithTimeout(200*time.Millisecond),
			WithRetry(RetryConfig{Enabled: false}),
		)
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		assert.Error(t, exp.Export(t.Context(), nonEmptyResourceProfiles()))
		assert.Zero(t, srv.requestCount())
	})

	t.Run("with_dial_option_overrides_credentials_and_user_agent", func(t *testing.T) {
		srv := &testProfilesServer{}
		addr := startTestGRPCServer(t, srv)

		exp, err := New(
			t.Context(),
			WithEndpoint(addr),
			WithServiceConfig(`{}`),
			WithReconnectionPeriod(time.Second),
			WithDialOption(
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithUserAgent("custom-profiles-ua/1.0"),
			),
		)
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		require.NoError(t, exp.Export(t.Context(), nonEmptyResourceProfiles()))
		srv.mu.Lock()
		require.Len(t, srv.mds, 1)
		ua := srv.mds[0].Get("user-agent")
		srv.mu.Unlock()
		require.NotEmpty(t, ua)
		assert.Contains(t, ua[0], "custom-profiles-ua/1.0")
	})

	t.Run("new_propagates_grpc_client_error", func(t *testing.T) {
		_, err := New(
			t.Context(),
			WithEndpoint("127.0.0.1:4317"),
			WithInsecure(),
			WithServiceConfig("{invalid-json"),
		)
		assert.Error(t, err)
	})
}

func TestClientMetadataAndDevelopmentVersion(t *testing.T) {
	t.Run("enforces_version_without_headers_option", func(t *testing.T) {
		srv := &testProfilesServer{}
		addr := startTestGRPCServer(t, srv)

		exp, err := New(t.Context(), WithEndpoint(addr), WithInsecure())
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		require.NoError(t, exp.Export(t.Context(), nonEmptyResourceProfiles()))
		srv.mu.Lock()
		require.Len(t, srv.mds, 1)
		assert.Equal(t, []string{profilesDevelopmentVersion}, srv.mds[0].Get(profilesDevelopmentVersionKey))
		srv.mu.Unlock()
	})

	t.Run("overrides_conflicting_headers_and_outgoing_context_without_mutating_client", func(t *testing.T) {
		srv := &testProfilesServer{}
		addr := startTestGRPCServer(t, srv)

		exp, err := New(
			t.Context(),
			WithEndpoint(addr),
			WithInsecure(),
			WithHeaders(map[string]string{
				"OTLP-Profiles-Development-Version": "999",
				"x-shared":                          "from-opt",
			}),
		)
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		ctx1 := metadata.AppendToOutgoingContext(
			t.Context(),
			profilesDevelopmentVersionKey, "777",
			"x-caller", "first",
			"x-shared", "from-ctx",
		)
		require.NoError(t, exp.Export(ctx1, nonEmptyResourceProfiles()))
		require.NoError(t, exp.Export(t.Context(), nonEmptyResourceProfiles()))

		srv.mu.Lock()
		require.Len(t, srv.mds, 2)
		assert.Equal(t, []string{profilesDevelopmentVersion}, srv.mds[0].Get(profilesDevelopmentVersionKey))
		assert.Equal(t, []string{"first"}, srv.mds[0].Get("x-caller"))
		assert.ElementsMatch(t, []string{"from-ctx", "from-opt"}, srv.mds[0].Get("x-shared"))

		assert.Equal(t, []string{profilesDevelopmentVersion}, srv.mds[1].Get(profilesDevelopmentVersionKey))
		assert.Empty(t, srv.mds[1].Get("x-caller"), "subsequent export must not retain prior caller metadata")
		assert.Equal(t, []string{"from-opt"}, srv.mds[1].Get("x-shared"))
		srv.mu.Unlock()
	})
}
