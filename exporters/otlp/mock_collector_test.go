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

package otlp_test

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"

	colmetricpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/collector/metrics/v1"
	coltracepb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/collector/trace/v1"
	commonpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/common/v1"
	metricpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/metrics/v1"
	resourcepb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/resource/v1"
	tracepb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/trace/v1"
)

func makeMockCollector(t *testing.T) *mockCol {
	return &mockCol{
		t: t,
		traceSvc: &mockTraceService{
			rsm: map[string]*tracepb.ResourceSpans{},
		},
		metricSvc: &mockMetricService{},
	}
}

type mockTraceService struct {
	mu      sync.RWMutex
	rsm     map[string]*tracepb.ResourceSpans
	headers metadata.MD
}

func (mts *mockTraceService) getHeaders() metadata.MD {
	return mts.headers
}

func (mts *mockTraceService) getSpans() []*tracepb.Span {
	mts.mu.RLock()
	defer mts.mu.RUnlock()
	spans := []*tracepb.Span{}
	for _, rs := range mts.rsm {
		spans = append(spans, rs.InstrumentationLibrarySpans[0].Spans...)
	}
	return spans
}

func (mts *mockTraceService) getResourceSpans() []*tracepb.ResourceSpans {
	mts.mu.RLock()
	defer mts.mu.RUnlock()
	rss := make([]*tracepb.ResourceSpans, 0, len(mts.rsm))
	for _, rs := range mts.rsm {
		rss = append(rss, rs)
	}
	return rss
}

func (mts *mockTraceService) Export(ctx context.Context, exp *coltracepb.ExportTraceServiceRequest) (*coltracepb.ExportTraceServiceResponse, error) {
	mts.mu.Lock()
	mts.headers, _ = metadata.FromIncomingContext(ctx)
	defer mts.mu.Unlock()
	rss := exp.GetResourceSpans()
	for _, rs := range rss {
		rstr := resourceString(rs.Resource)
		existingRs, ok := mts.rsm[rstr]
		if !ok {
			mts.rsm[rstr] = rs
			// TODO (rghetia): Add support for library Info.
			if len(rs.InstrumentationLibrarySpans) == 0 {
				rs.InstrumentationLibrarySpans = []*tracepb.InstrumentationLibrarySpans{
					{
						Spans: []*tracepb.Span{},
					},
				}
			}
		} else {
			if len(rs.InstrumentationLibrarySpans) > 0 {
				existingRs.InstrumentationLibrarySpans[0].Spans =
					append(existingRs.InstrumentationLibrarySpans[0].Spans,
						rs.InstrumentationLibrarySpans[0].GetSpans()...)
			}
		}
	}
	return &coltracepb.ExportTraceServiceResponse{}, nil
}

func resourceString(res *resourcepb.Resource) string {
	sAttrs := sortedAttributes(res.GetAttributes())
	rstr := ""
	for _, attr := range sAttrs {
		rstr = rstr + attr.String()

	}
	return rstr
}

func sortedAttributes(attrs []*commonpb.KeyValue) []*commonpb.KeyValue {
	sort.Slice(attrs[:], func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})
	return attrs
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
		// TODO (rghetia) handle multiple resource and library info.
		if len(rm.InstrumentationLibraryMetrics) > 0 {
			mms.metrics = append(mms.metrics, rm.InstrumentationLibraryMetrics[0].Metrics...)
		}
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

func (mc *mockCol) getResourceSpans() []*tracepb.ResourceSpans {
	return mc.traceSvc.getResourceSpans()
}

func (mc *mockCol) getHeaders() metadata.MD {
	return mc.traceSvc.getHeaders()
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
