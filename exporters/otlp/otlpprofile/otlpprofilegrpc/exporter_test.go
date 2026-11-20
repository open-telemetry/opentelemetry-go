// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc_test

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
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc/internal"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/profile"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type mockProfilesServer struct {
	colprofilepb.UnimplementedProfilesServiceServer

	mu       sync.Mutex
	requests []*colprofilepb.ExportProfilesServiceRequest
	mds      []metadata.MD
	handler  func(context.Context, *colprofilepb.ExportProfilesServiceRequest) (*colprofilepb.ExportProfilesServiceResponse, error)
}

func (m *mockProfilesServer) Export(
	ctx context.Context,
	req *colprofilepb.ExportProfilesServiceRequest,
) (*colprofilepb.ExportProfilesServiceResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	m.mu.Lock()
	m.requests = append(m.requests, req)
	m.mds = append(m.mds, md.Copy())
	h := m.handler
	m.mu.Unlock()

	if h != nil {
		return h(ctx, req)
	}
	return &colprofilepb.ExportProfilesServiceResponse{}, nil
}

type compressionStatsHandler struct {
	mu           sync.Mutex
	compressions []string
}

func (*compressionStatsHandler) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	return ctx
}

func (h *compressionStatsHandler) HandleRPC(_ context.Context, s stats.RPCStats) {
	if ih, ok := s.(*stats.InHeader); ok {
		h.mu.Lock()
		h.compressions = append(h.compressions, ih.Compression)
		h.mu.Unlock()
	}
}

func (*compressionStatsHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

func (*compressionStatsHandler) HandleConn(context.Context, stats.ConnStats) {}

func startMockServer(t *testing.T, srv *mockProfilesServer, opts ...grpc.ServerOption) string {
	t.Helper()
	var lc net.ListenConfig
	lis, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	gs := grpc.NewServer(opts...)
	colprofilepb.RegisterProfilesServiceServer(gs, srv)
	go func() { _ = gs.Serve(lis) }()
	t.Cleanup(func() { gs.GracefulStop() })
	return lis.Addr().String()
}

func sampleResourceProfiles() *profiledata.ResourceProfiles {
	loc := &profiledata.Location{
		Address: 0x401000,
		Lines: []profiledata.Line{
			{Function: profiledata.Function{Name: "main.main", Filename: "main.go"}, Line: 42},
		},
	}
	return &profiledata.ResourceProfiles{
		Resource: resource.NewSchemaless(attribute.String("service.name", "otlp-grpc-test")),
		ScopeProfiles: []profiledata.ScopeProfiles{
			{
				Scope: instrumentation.Scope{Name: "go.opentelemetry.io/otel/sdk/profile/runtime/heap"},
				Profiles: []profiledata.Profile{
					{
						SampleType: profiledata.ValueType{Type: "inuse_space", Unit: "bytes"},
						Samples:    []profiledata.Sample{{Stack: []*profiledata.Location{loc}, Values: []int64{4096}}},
						Time:       time.Unix(1_700_000_000, 0).UTC(),
					},
				},
			},
		},
	}
}

func TestExporterEndToEndWithRuntimeProducer(t *testing.T) {
	mockSrv := &mockProfilesServer{}
	statsHandler := &compressionStatsHandler{}
	addr := startMockServer(t, mockSrv, grpc.StatsHandler(statsHandler))

	exp, err := otlpprofilegrpc.New(
		t.Context(),
		otlpprofilegrpc.WithEndpoint(addr),
		otlpprofilegrpc.WithInsecure(),
		otlpprofilegrpc.WithCompressor("gzip"),
		otlpprofilegrpc.WithHeaders(map[string]string{"x-custom": "hello"}),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		//nolint:usetesting // required to avoid getting a canceled context at cleanup.
		assert.NoError(t, exp.Shutdown(context.Background()))
	})

	mr := profile.NewManualReader()
	pp := profile.NewProfilerProvider(
		profile.WithResource(resource.NewSchemaless(attribute.String("service.name", "grpc-e2e"))),
		profile.WithProducer(profile.NewRuntimeProducer()),
		profile.WithReader(mr),
	)
	t.Cleanup(func() {
		//nolint:usetesting // required to avoid getting a canceled context at cleanup.
		assert.NoError(t, pp.Shutdown(context.Background()))
	})

	var rp profiledata.ResourceProfiles
	require.NoError(t, mr.Collect(t.Context(), &rp))
	require.NotEmpty(t, rp.ScopeProfiles)

	require.NoError(t, exp.Export(t.Context(), &rp))

	mockSrv.mu.Lock()
	require.Len(t, mockSrv.requests, 1)
	req := mockSrv.requests[0]
	md := mockSrv.mds[0]
	mockSrv.mu.Unlock()

	assert.Equal(t, []string{"1"}, md.Get("otlp-profiles-development-version"))
	assert.Equal(t, []string{"hello"}, md.Get("x-custom"))
	require.NotEmpty(t, md.Get("user-agent"))
	assert.Contains(t, md.Get("user-agent")[0], "OTel Go OTLP over gRPC profiles exporter/"+otlpprofilegrpc.Version())
	statsHandler.mu.Lock()
	assert.Equal(t, []string{"gzip"}, statsHandler.compressions)
	statsHandler.mu.Unlock()
	require.Len(t, req.ResourceProfiles, 1)
	assert.NotEmpty(t, req.ResourceProfiles[0].ScopeProfiles)
	require.NotNil(t, req.Dictionary)
	assert.NotEmpty(t, req.Dictionary.StringTable)
	assert.Empty(t, req.Dictionary.StringTable[0])

	require.NoError(t, exp.ForceFlush(t.Context()))
	require.NoError(t, exp.Shutdown(t.Context()))
	require.NoError(t, exp.Shutdown(t.Context()))
	assert.ErrorIs(t, exp.Export(t.Context(), &rp), profile.ErrExporterShutdown)
}

func TestExporterPartialSuccessRetriesAndLimits(t *testing.T) {
	t.Run("partial_success", func(t *testing.T) {
		var step atomic.Int32
		mockSrv := &mockProfilesServer{
			handler: func(context.Context, *colprofilepb.ExportProfilesServiceRequest) (*colprofilepb.ExportProfilesServiceResponse, error) {
				switch step.Add(1) {
				case 1:
					return &colprofilepb.ExportProfilesServiceResponse{
						PartialSuccess: &colprofilepb.ExportProfilesPartialSuccess{
							RejectedProfiles: 1,
							ErrorMessage:     "rejected profile",
						},
					}, nil
				case 2:
					return &colprofilepb.ExportProfilesServiceResponse{
						PartialSuccess: &colprofilepb.ExportProfilesPartialSuccess{},
					}, nil
				case 3:
					return &colprofilepb.ExportProfilesServiceResponse{
						PartialSuccess: &colprofilepb.ExportProfilesPartialSuccess{
							RejectedProfiles: 2,
						},
					}, nil
				case 4:
					return &colprofilepb.ExportProfilesServiceResponse{
						PartialSuccess: &colprofilepb.ExportProfilesPartialSuccess{
							ErrorMessage: "warning only",
						},
					}, nil
				case 5:
					return nil, status.Error(codes.Unavailable, "retry once")
				default:
					return &colprofilepb.ExportProfilesServiceResponse{
						PartialSuccess: &colprofilepb.ExportProfilesPartialSuccess{
							RejectedProfiles: 3,
							ErrorMessage:     "after retry",
						},
					}, nil
				}
			},
		}
		addr := startMockServer(t, mockSrv)

		exp, err := otlpprofilegrpc.New(
			t.Context(),
			otlpprofilegrpc.WithEndpoint(addr),
			otlpprofilegrpc.WithInsecure(),
			otlpprofilegrpc.WithRetry(otlpprofilegrpc.RetryConfig{
				Enabled:         true,
				InitialInterval: time.Millisecond,
				MaxInterval:     5 * time.Millisecond,
				MaxElapsedTime:  time.Second,
			}),
		)
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.ErrorIs(t, err, internal.PartialSuccess{})
		assert.ErrorContains(t, err, "rejected profile")

		assert.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.ErrorIs(t, err, internal.PartialSuccess{})
		assert.ErrorContains(t, err, "empty message")

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.ErrorIs(t, err, internal.PartialSuccess{})
		assert.ErrorContains(t, err, "warning only")

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.ErrorIs(t, err, internal.PartialSuccess{})
		assert.ErrorContains(t, err, "after retry")
	})

	t.Run("retries_unavailable_and_resource_exhausted_with_retry_info", func(t *testing.T) {
		throttledSt, err := status.New(codes.ResourceExhausted, "throttled").WithDetails(&errdetails.RetryInfo{
			RetryDelay: durationpb.New(time.Millisecond),
		})
		require.NoError(t, err)

		var attempts atomic.Int32
		mockSrv := &mockProfilesServer{
			handler: func(context.Context, *colprofilepb.ExportProfilesServiceRequest) (*colprofilepb.ExportProfilesServiceResponse, error) {
				switch attempts.Add(1) {
				case 1:
					return nil, status.Error(codes.Unavailable, "transient")
				case 2:
					return nil, throttledSt.Err()
				default:
					return &colprofilepb.ExportProfilesServiceResponse{}, nil
				}
			},
		}
		addr := startMockServer(t, mockSrv)

		exp, err := otlpprofilegrpc.New(
			t.Context(),
			otlpprofilegrpc.WithEndpoint(addr),
			otlpprofilegrpc.WithInsecure(),
			otlpprofilegrpc.WithRetry(otlpprofilegrpc.RetryConfig{
				Enabled:         true,
				InitialInterval: time.Millisecond,
				MaxInterval:     5 * time.Millisecond,
				MaxElapsedTime:  time.Second,
			}),
		)
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		require.NoError(t, exp.Export(t.Context(), sampleResourceProfiles()))
		assert.Equal(t, int32(3), attempts.Load())
	})

	t.Run("invalid_argument_not_retried", func(t *testing.T) {
		var attempts atomic.Int32
		mockSrv := &mockProfilesServer{
			handler: func(context.Context, *colprofilepb.ExportProfilesServiceRequest) (*colprofilepb.ExportProfilesServiceResponse, error) {
				attempts.Add(1)
				return nil, status.Error(codes.InvalidArgument, "unsupported development version")
			},
		}
		addr := startMockServer(t, mockSrv)

		exp, err := otlpprofilegrpc.New(
			t.Context(),
			otlpprofilegrpc.WithEndpoint(addr),
			otlpprofilegrpc.WithInsecure(),
		)
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp.Shutdown(context.Background()))
		})

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.Error(t, err)
		assert.Equal(t, int32(1), attempts.Load())
	})

	t.Run("max_request_size_and_external_conn", func(t *testing.T) {
		mockSrv := &mockProfilesServer{}
		addr := startMockServer(t, mockSrv)

		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		t.Cleanup(func() { assert.NoError(t, conn.Close()) })

		exp, err := otlpprofilegrpc.New(
			t.Context(),
			otlpprofilegrpc.WithGRPCConn(conn),
			otlpprofilegrpc.WithMaxRequestSize(8),
		)
		require.NoError(t, err)

		err = exp.Export(t.Context(), sampleResourceProfiles())
		require.Error(t, err)
		assert.ErrorContains(t, err, "request message too large")

		require.NoError(t, exp.Export(t.Context(), nil))
		require.NoError(t, exp.Export(t.Context(), &profiledata.ResourceProfiles{}))

		mockSrv.mu.Lock()
		assert.Empty(t, mockSrv.requests, "oversized, nil, and empty exports must not reach server")
		mockSrv.mu.Unlock()

		require.NoError(t, exp.Shutdown(t.Context()))

		// Verify Shutdown did not close the caller-owned conn and WithMaxRequestSize(0) disables the limit.
		exp2, err := otlpprofilegrpc.New(
			t.Context(),
			otlpprofilegrpc.WithGRPCConn(conn),
			otlpprofilegrpc.WithMaxRequestSize(0),
		)
		require.NoError(t, err)
		t.Cleanup(func() {
			//nolint:usetesting // required to avoid getting a canceled context at cleanup.
			assert.NoError(t, exp2.Shutdown(context.Background()))
		})

		require.NoError(t, exp2.Export(t.Context(), sampleResourceProfiles()))
		mockSrv.mu.Lock()
		assert.Len(t, mockSrv.requests, 1)
		mockSrv.mu.Unlock()
	})
}
