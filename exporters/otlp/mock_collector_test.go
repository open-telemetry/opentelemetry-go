// Copyright 2020, OpenTelemetry Authors
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

package otlp_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"

	colmetricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/metrics/v1"
	coltracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/trace/v1"
	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
	tracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/trace/v1"
)

func makeMockCollector(t *testing.T) *mockCol {
	return &mockCol{
		t:         t,
		traceSvc:  &mockTraceService{},
		metricSvc: &mockMetricService{},
	}
}

type mockTraceService struct {
	mu    sync.RWMutex
	spans []*tracepb.Span
}

func (mts *mockTraceService) getSpans() []*tracepb.Span {
	mts.mu.RLock()
	spans := append([]*tracepb.Span{}, mts.spans...)
	mts.mu.RUnlock()

	return spans
}

func (mts *mockTraceService) Export(ctx context.Context, exp *coltracepb.ExportTraceServiceRequest) (*coltracepb.ExportTraceServiceResponse, error) {
	resourceSpans := exp.GetResourceSpans()
	// TODO (rghetia): handle Resources
	mts.mu.Lock()
	for _, rs := range resourceSpans {
		mts.spans = append(mts.spans, rs.Spans...)
	}
	mts.mu.Unlock()
	return &coltracepb.ExportTraceServiceResponse{}, nil
}

type mockMetricService struct {
	mu      sync.RWMutex
	metrics []*metricpb.Metric
}

func (mms *mockMetricService) getMetrics() []*metricpb.Metric {
	// copy in order to not change.
	m := make([]*metricpb.Metric, 0, len(mms.metrics))
	mms.mu.RLock()
	defer mms.mu.RUnlock()
	return append(m, mms.metrics...)
}

func (mms *mockMetricService) Export(ctx context.Context, exp *colmetricpb.ExportMetricsServiceRequest) (*colmetricpb.ExportMetricsServiceResponse, error) {
	mms.mu.Lock()
	for _, rm := range exp.GetResourceMetrics() {
		mms.metrics = append(mms.metrics, rm.Metrics...)
	}
	mms.mu.Unlock()
	return &colmetricpb.ExportMetricsServiceResponse{}, nil
}

type mockCol struct {
	t *testing.T

	traceSvc  *mockTraceService
	metricSvc *mockMetricService

	address  string
	stopFunc func() error
	stopOnce sync.Once
}

var _ coltracepb.TraceServiceServer = (*mockTraceService)(nil)
var _ colmetricpb.MetricsServiceServer = (*mockMetricService)(nil)

var errAlreadyStopped = fmt.Errorf("already stopped")

func (mc *mockCol) stop() error {
	var err = errAlreadyStopped
	mc.stopOnce.Do(func() {
		if mc.stopFunc != nil {
			err = mc.stopFunc()
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

func (mc *mockCol) getSpans() []*tracepb.Span {
	return mc.traceSvc.getSpans()
}

func (mc *mockCol) getMetrics() []*metricpb.Metric {
	return mc.metricSvc.getMetrics()
}

// runMockCol is a helper function to create a mockCol
func runMockCol(t *testing.T) *mockCol {
	return runMockColAtAddr(t, "localhost:0")
}

func runMockColAtAddr(t *testing.T, addr string) *mockCol {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to get an address: %v", err)
	}

	srv := grpc.NewServer()
	mc := makeMockCollector(t)
	coltracepb.RegisterTraceServiceServer(srv, mc.traceSvc)
	colmetricpb.RegisterMetricsServiceServer(srv, mc.metricSvc)
	go func() {
		_ = srv.Serve(ln)
	}()

	deferFunc := func() error {
		srv.Stop()
		return ln.Close()
	}

	_, collectorPortStr, _ := net.SplitHostPort(ln.Addr().String())

	mc.address = "localhost:" + collectorPortStr
	mc.stopFunc = deferFunc

	return mc
}
