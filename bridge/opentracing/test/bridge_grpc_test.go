// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"net"
	"testing"
	"time"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	testpb "github.com/opentracing-contrib/go-grpc/test/otgrpc_testing"
	ot "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ototel "go.opentelemetry.io/otel/bridge/opentracing"
	"go.opentelemetry.io/otel/bridge/opentracing/internal"
	"go.opentelemetry.io/otel/propagation"
)

type testGRPCServer struct{}

func (*testGRPCServer) UnaryCall(ctx context.Context, r *testpb.SimpleRequest) (*testpb.SimpleResponse, error) {
	return &testpb.SimpleResponse{Payload: r.Payload * 2}, nil
}

func (*testGRPCServer) StreamingOutputCall(*testpb.SimpleRequest, testpb.TestService_StreamingOutputCallServer) error {
	return nil
}

func (*testGRPCServer) StreamingInputCall(testpb.TestService_StreamingInputCallServer) error {
	return nil
}

func (*testGRPCServer) StreamingBidirectionalCall(testpb.TestService_StreamingBidirectionalCallServer) error {
	return nil
}

func startTestGRPCServer(t *testing.T, tracer ot.Tracer) (*grpc.Server, net.Addr) {
	lis, _ := net.Listen("tcp", ":0")
	server := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)),
	)
	testpb.RegisterTestServiceServer(server, &testGRPCServer{})

	go func() {
		err := server.Serve(lis)
		require.NoError(t, err)
	}()

	return server, lis.Addr()
}

func TestBridgeTracer_ExtractAndInject_gRPC(t *testing.T) {
	tracer := internal.NewMockTracer()

	bridge := ototel.NewBridgeTracer()
	bridge.SetOpenTelemetryTracer(tracer)
	bridge.SetTextMapPropagator(propagation.TraceContext{})

	srv, addr := startTestGRPCServer(t, bridge)
	defer srv.Stop()

	conn, err := grpc.NewClient(
		addr.String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(bridge)),
	)
	require.NoError(t, err)
	cli := testpb.NewTestServiceClient(conn)

	ctx, cx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cx()
	res, err := cli.UnaryCall(ctx, &testpb.SimpleRequest{Payload: 42})
	require.NoError(t, err)
	assert.EqualValues(t, 84, res.Payload)

	checkSpans := func() bool {
		return len(tracer.FinishedSpans) == 2
	}
	require.Eventuallyf(t, checkSpans, 5*time.Second, 5*time.Millisecond, "expecting two spans")
	assert.Equal(t,
		tracer.FinishedSpans[0].SpanContext().TraceID(),
		tracer.FinishedSpans[1].SpanContext().TraceID(),
		"expecting same trace ID",
	)
}
