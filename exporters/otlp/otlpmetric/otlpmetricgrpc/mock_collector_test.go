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

package otlpmetricgrpc_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otlpmetrictest"
	collectormetricpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

func makeMockCollector(t *testing.T, mockConfig *mockConfig) *mockCollector {
	return &mockCollector{
		t: t,
		metricSvc: &mockMetricService{
			storage: otlpmetrictest.NewMetricsStorage(),
			errors:  mockConfig.errors,
		},
	}
}

type mockMetricService struct {
	collectormetricpb.UnimplementedMetricsServiceServer

	requests int
	errors   []error

	headers metadata.MD
	mu      sync.RWMutex
	storage otlpmetrictest.MetricsStorage
	delay   time.Duration
}

func (mms *mockMetricService) getHeaders() metadata.MD {
	mms.mu.RLock()
	defer mms.mu.RUnlock()
	return mms.headers
}

func (mms *mockMetricService) getMetrics() []*metricpb.Metric {
	mms.mu.RLock()
	defer mms.mu.RUnlock()
	return mms.storage.GetMetrics()
}

func (mms *mockMetricService) Export(ctx context.Context, exp *collectormetricpb.ExportMetricsServiceRequest) (*collectormetricpb.ExportMetricsServiceResponse, error) {
	if mms.delay > 0 {
		time.Sleep(mms.delay)
	}

	mms.mu.Lock()
	defer func() {
		mms.requests++
		mms.mu.Unlock()
	}()

	reply := &collectormetricpb.ExportMetricsServiceResponse{}
	if mms.requests < len(mms.errors) {
		idx := mms.requests
		return reply, mms.errors[idx]
	}

	mms.headers, _ = metadata.FromIncomingContext(ctx)
	mms.storage.AddMetrics(exp)
	return reply, nil
}

type mockCollector struct {
	t *testing.T

	metricSvc *mockMetricService

	endpoint string
	stopFunc func()
	stopOnce sync.Once
}

type mockConfig struct {
	errors   []error
	endpoint string
}

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
	// Getting the lock ensures the metricSvc is done flushing.
	mc.metricSvc.mu.Lock()
	defer mc.metricSvc.mu.Unlock()
	return err
}

func (mc *mockCollector) Stop() error {
	return mc.stop()
}

func (mc *mockCollector) getHeaders() metadata.MD {
	return mc.metricSvc.getHeaders()
}

func (mc *mockCollector) getMetrics() []*metricpb.Metric {
	return mc.metricSvc.getMetrics()
}

func (mc *mockCollector) GetMetrics() []*metricpb.Metric {
	return mc.getMetrics()
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
	collectormetricpb.RegisterMetricsServiceServer(srv, mc.metricSvc)
	go func() {
		_ = srv.Serve(ln)
	}()

	mc.endpoint = ln.Addr().String()
	// srv.Stop calls Close on mc.ln.
	mc.stopFunc = srv.Stop

	return mc
}
