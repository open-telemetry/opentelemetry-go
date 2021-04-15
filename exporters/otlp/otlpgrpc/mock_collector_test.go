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

package otlpgrpc_test

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"

	"google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel/exporters/otlp/internal/otlptest"
	collectormetricpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	collectortracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

func makeMockCollector(t *testing.T) *mockCollector {
	return &mockCollector{
		t: t,
		traceSvc: &mockTraceService{
			storage: otlptest.NewSpansStorage(),
		},
		metricSvc: &mockMetricService{
			storage: otlptest.NewMetricsStorage(),
		},
	}
}

type mockTraceService struct {
	collectortracepb.UnimplementedTraceServiceServer

	mu      sync.RWMutex
	storage otlptest.SpansStorage
	headers metadata.MD
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
	reply := &collectortracepb.ExportTraceServiceResponse{}
	mts.mu.Lock()
	defer mts.mu.Unlock()
	mts.headers, _ = metadata.FromIncomingContext(ctx)
	mts.storage.AddSpans(exp)
	return reply, nil
}

type mockMetricService struct {
	collectormetricpb.UnimplementedMetricsServiceServer

	mu      sync.RWMutex
	storage otlptest.MetricsStorage
}

func (mms *mockMetricService) getMetrics() []*metricpb.Metric {
	mms.mu.RLock()
	defer mms.mu.RUnlock()
	return mms.storage.GetMetrics()
}

func (mms *mockMetricService) Export(ctx context.Context, exp *collectormetricpb.ExportMetricsServiceRequest) (*collectormetricpb.ExportMetricsServiceResponse, error) {
	reply := &collectormetricpb.ExportMetricsServiceResponse{}
	mms.mu.Lock()
	defer mms.mu.Unlock()
	mms.storage.AddMetrics(exp)
	return reply, nil
}

type mockCollector struct {
	t *testing.T

	traceSvc  *mockTraceService
	metricSvc *mockMetricService

	endpoint string
	ln       *listener
	stopFunc func()
	stopOnce sync.Once
}

var _ collectortracepb.TraceServiceServer = (*mockTraceService)(nil)
var _ collectormetricpb.MetricsServiceServer = (*mockMetricService)(nil)

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

	// Wait for services to finish reading/writing.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// Getting the lock ensures the traceSvc is done flushing.
		mc.traceSvc.mu.Lock()
		defer mc.traceSvc.mu.Unlock()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		// Getting the lock ensures the metricSvc is done flushing.
		mc.metricSvc.mu.Lock()
		defer mc.metricSvc.mu.Unlock()
		wg.Done()
	}()
	wg.Wait()
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

func (mc *mockCollector) getMetrics() []*metricpb.Metric {
	return mc.metricSvc.getMetrics()
}

func (mc *mockCollector) GetMetrics() []*metricpb.Metric {
	return mc.getMetrics()
}

// WaitForConn will wait indefintely for a connection to be estabilished
// with the mockCollector before returning.
func (mc *mockCollector) WaitForConn() {
	for {
		select {
		case <-mc.ln.C:
			return
		default:
			runtime.Gosched()
		}
	}
}

// runMockCollector is a helper function to create a mock Collector
func runMockCollector(t *testing.T) *mockCollector {
	return runMockCollectorAtEndpoint(t, "127.0.0.1:0")
}

func runMockCollectorAtEndpoint(t *testing.T, endpoint string) *mockCollector {
	ln, err := net.Listen("tcp", endpoint)
	if err != nil {
		t.Fatalf("Failed to get an endpoint: %v", err)
	}

	srv := grpc.NewServer()
	mc := makeMockCollector(t)
	collectortracepb.RegisterTraceServiceServer(srv, mc.traceSvc)
	collectormetricpb.RegisterMetricsServiceServer(srv, mc.metricSvc)
	mc.ln = newListener(ln)
	go func() {
		_ = srv.Serve((net.Listener)(mc.ln))
	}()

	mc.endpoint = ln.Addr().String()
	mc.stopFunc = srv.Stop

	return mc
}

type listener struct {
	wrapped net.Listener

	C      chan struct{}
	closed chan struct{}
}

func newListener(wrapped net.Listener) *listener {
	return &listener{
		wrapped: wrapped,
		C:       make(chan struct{}, 1),
		closed:  make(chan struct{}),
	}
}

func (l *listener) Accept() (net.Conn, error) {
	select {
	case <-l.closed:
		close(l.C)
		return nil, syscall.EINVAL
	default:
	}

	conn, err := l.wrapped.Accept()
	select {
	case l.C <- struct{}{}:
	default:
		// If C is full move on.
	}
	return conn, err
}

func (l *listener) Close() error {
	close(l.closed)
	return l.wrapped.Close()
}

func (l *listener) Addr() net.Addr { return l.wrapped.Addr() }
