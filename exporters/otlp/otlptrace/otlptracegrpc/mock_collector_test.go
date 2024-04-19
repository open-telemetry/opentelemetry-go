// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracegrpc_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/otlptracetest"
	collectortracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

func makeMockCollector(t *testing.T, mockConfig *mockConfig) *mockCollector {
	return &mockCollector{
		t: t,
		traceSvc: &mockTraceService{
			storage: otlptracetest.NewSpansStorage(),
			errors:  mockConfig.errors,
			partial: mockConfig.partial,
		},
		stopped: make(chan struct{}),
	}
}

type mockTraceService struct {
	collectortracepb.UnimplementedTraceServiceServer

	errors      []error
	partial     *collectortracepb.ExportTracePartialSuccess
	requests    int
	mu          sync.RWMutex
	storage     otlptracetest.SpansStorage
	headers     metadata.MD
	exportBlock chan struct{}
}

func (mts *mockTraceService) getHeaders() metadata.MD {
	mts.mu.RLock()
	defer mts.mu.RUnlock()
	return mts.headers
}

func (mts *mockTraceService) getSpans() []*tracepb.Span {
	mts.mu.RLock()
	defer mts.mu.RUnlock()
	return mts.storage.GetSpans()
}

func (mts *mockTraceService) getResourceSpans() []*tracepb.ResourceSpans {
	mts.mu.RLock()
	defer mts.mu.RUnlock()
	return mts.storage.GetResourceSpans()
}

func (mts *mockTraceService) Export(ctx context.Context, exp *collectortracepb.ExportTraceServiceRequest) (*collectortracepb.ExportTraceServiceResponse, error) {
	mts.mu.Lock()
	defer func() {
		mts.requests++
		mts.mu.Unlock()
	}()

	if mts.exportBlock != nil {
		// Do this with the lock held so the mockCollector.Stop does not
		// abandon cleaning up resources.
		<-mts.exportBlock
	}

	reply := &collectortracepb.ExportTraceServiceResponse{
		PartialSuccess: mts.partial,
	}
	if mts.requests < len(mts.errors) {
		idx := mts.requests
		return reply, mts.errors[idx]
	}

	mts.headers, _ = metadata.FromIncomingContext(ctx)
	mts.storage.AddSpans(exp)
	return reply, nil
}

type mockCollector struct {
	t *testing.T

	traceSvc *mockTraceService

	endpoint string
	stopFunc func()
	stopOnce sync.Once
	stopped  chan struct{}
}

type mockConfig struct {
	errors   []error
	endpoint string
	partial  *collectortracepb.ExportTracePartialSuccess
}

var _ collectortracepb.TraceServiceServer = (*mockTraceService)(nil)

var errAlreadyStopped = fmt.Errorf("already stopped")

func (mc *mockCollector) stop() error {
	err := errAlreadyStopped
	mc.stopOnce.Do(func() {
		err = nil
		if mc.stopFunc != nil {
			mc.stopFunc()
		}
	})
	// Wait until gRPC server is down.
	<-mc.stopped

	// Getting the lock ensures the traceSvc is done flushing.
	mc.traceSvc.mu.Lock()
	defer mc.traceSvc.mu.Unlock()

	return err
}

func (mc *mockCollector) Stop() error {
	return mc.stop()
}

func (mc *mockCollector) getSpans() []*tracepb.Span {
	return mc.traceSvc.getSpans()
}

func (mc *mockCollector) getResourceSpans() []*tracepb.ResourceSpans {
	return mc.traceSvc.getResourceSpans()
}

func (mc *mockCollector) GetResourceSpans() []*tracepb.ResourceSpans {
	return mc.getResourceSpans()
}

func (mc *mockCollector) getHeaders() metadata.MD {
	return mc.traceSvc.getHeaders()
}

// runMockCollector is a helper function to create a mock Collector.
func runMockCollector(t *testing.T) *mockCollector {
	t.Helper()
	return runMockCollectorAtEndpoint(t, "localhost:0")
}

func runMockCollectorAtEndpoint(t *testing.T, endpoint string) *mockCollector {
	t.Helper()
	return runMockCollectorWithConfig(t, &mockConfig{endpoint: endpoint})
}

func runMockCollectorWithConfig(t *testing.T, mockConfig *mockConfig) *mockCollector {
	t.Helper()
	ln, err := net.Listen("tcp", mockConfig.endpoint)
	require.NoError(t, err, "net.Listen")

	srv := grpc.NewServer()
	mc := makeMockCollector(t, mockConfig)
	collectortracepb.RegisterTraceServiceServer(srv, mc.traceSvc)
	go func() {
		_ = srv.Serve(ln)
		close(mc.stopped)
	}()

	mc.endpoint = ln.Addr().String()
	mc.stopFunc = srv.Stop
	return mc
}
