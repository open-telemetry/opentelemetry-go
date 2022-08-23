// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlptracegrpc_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/otlptracetest"
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
}

type mockConfig struct {
	errors   []error
	endpoint string
	partial  *collectortracepb.ExportTracePartialSuccess
}

var _ collectortracepb.TraceServiceServer = (*mockTraceService)(nil)

var errAlreadyStopped = fmt.Errorf("already stopped")

func (mc *mockCollector) stop() error {
	var err = errAlreadyStopped
	mc.stopOnce.Do(func() {
		err = nil
		if mc.stopFunc != nil {
			mc.stopFunc()
		}
	})
	// Give it sometime to shutdown.
	<-time.After(160 * time.Millisecond)

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
	return runMockCollectorAtEndpoint(t, "localhost:0")
}

func runMockCollectorAtEndpoint(t *testing.T, endpoint string) *mockCollector {
	return runMockCollectorWithConfig(t, &mockConfig{endpoint: endpoint})
}

func runMockCollectorWithConfig(t *testing.T, mockConfig *mockConfig) *mockCollector {
	ln, err := net.Listen("tcp", mockConfig.endpoint)
	if err != nil {
		t.Fatalf("Failed to get an endpoint: %v", err)
	}

	srv := grpc.NewServer()
	mc := makeMockCollector(t, mockConfig)
	collectortracepb.RegisterTraceServiceServer(srv, mc.traceSvc)
	go func() {
		_ = srv.Serve(ln)
	}()

	mc.endpoint = ln.Addr().String()
	mc.stopFunc = srv.Stop

	return mc
}
