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

package repro

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	cmpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	rpb "go.opentelemetry.io/proto/otlp/resource/v1"
)

func TestStdoutExporter(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	const manualMetricName = "test instrument"

	// Run a fake stdout
	path := filepath.Join(t.TempDir(), "stdout.txt")
	file, err := os.Create(path)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, file.Close()) })
	prev := os.Stdout
	t.Cleanup(func() { os.Stdout = prev })
	os.Stdout = file

	// ACT
	// Setup metrics provider
	exp, err := stdoutmetric.New()
	require.NoError(t, err)
	provider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exp)))
	global.SetMeterProvider(provider)

	// Add runtime metrics instrumentation.
	err = runtime.Start()
	require.NoError(t, err)

	// Add manual metric
	cnt, err := global.MeterProvider().Meter(t.Name()).Int64Counter(manualMetricName)
	require.NoError(t, err)
	cnt.Add(ctx, 123)

	// Shutdown to flush all spans from SDK.
	require.NoError(t, provider.Shutdown(context.Background()))

	// ASSERT
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	got := string(b)
	assert.Contains(t, got, manualMetricName)
	assert.Contains(t, got, "runtime.uptime")
}

func TestOTLPExporter(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	const manualMetricName = "test instrument"

	// Run a fake gRPC Collector
	coll := &collector{}
	coll.Start(t)

	// ACT
	// Setup metrics provider
	exp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(coll.Endpoint),
		otlpmetricgrpc.WithInsecure())
	require.NoError(t, err)
	provider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exp)))
	global.SetMeterProvider(provider)

	// Add runtime metrics instrumentation.
	err = runtime.Start()
	require.NoError(t, err)

	// Add manual metric
	cnt, err := global.MeterProvider().Meter(t.Name()).Int64Counter(manualMetricName)
	require.NoError(t, err)
	cnt.Add(ctx, 123)

	// Shutdown to flush all spans from SDK.
	require.NoError(t, provider.Shutdown(context.Background()))

	// ASSERT
	got := coll.ExportedMetrics()
	assert.NotNil(t, got)
	assertHasMetric(t, got, manualMetricName)
	assertHasMetric(t, got, "runtime.uptime")
}

func assertHasMetric(t *testing.T, got *metricsExportRequest, name string) {
	t.Helper()
	for _, m := range got.Metrics {
		if m.Name == name {
			return
		}
	}

	// Not found. Generate assertion failure.
	var gotMetrics []string
	for _, m := range got.Metrics {
		gotMetrics = append(gotMetrics, m.Name)
	}
	assert.Failf(t, "should contain metric", "want: %v, got: %v", name, gotMetrics)
}

type (
	collector struct {
		Endpoint string

		metricsService *collectorMetricsServiceServer
		grpcSrv        *grpc.Server
	}

	collectorMetricsServiceServer struct {
		cmpb.UnimplementedMetricsServiceServer

		mtx  sync.Mutex
		data *metricsExportRequest
	}

	metricsExportRequest struct {
		Header   metadata.MD
		Resource *rpb.Resource
		Metrics  []*mpb.Metric
	}
)

func (coll *collector) Start(t *testing.T) {
	if coll.Endpoint == "" {
		coll.Endpoint = "localhost:0"
	}
	ln, err := net.Listen("tcp", coll.Endpoint)
	require.NoError(t, err)
	coll.Endpoint = ln.Addr().String() // set actual endpoint

	coll.metricsService = &collectorMetricsServiceServer{}

	coll.grpcSrv = grpc.NewServer()
	cmpb.RegisterMetricsServiceServer(coll.grpcSrv, coll.metricsService)
	errCh := make(chan error, 1)

	// Serve and then stop during cleanup.
	t.Cleanup(func() {
		coll.grpcSrv.GracefulStop()
		assert.NoError(t, <-errCh)
	})
	go func() { errCh <- coll.grpcSrv.Serve(ln) }()

	// Wait until gRPC server is up.
	dialOpts := []grpc.DialOption{grpc.WithBlock()}
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(coll.Endpoint, dialOpts...)
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func (coll *collector) ExportedMetrics() *metricsExportRequest {
	// stop to make sure all requests are processed and synchronized
	coll.grpcSrv.GracefulStop()

	defer coll.metricsService.mtx.Unlock()
	coll.metricsService.mtx.Lock()
	return coll.metricsService.data
}

func (cmss *collectorMetricsServiceServer) Export(ctx context.Context, exp *cmpb.ExportMetricsServiceRequest) (*cmpb.ExportMetricsServiceResponse, error) {
	rs := exp.ResourceMetrics[0]
	scopeMetrics := rs.ScopeMetrics[0]
	headers, _ := metadata.FromIncomingContext(ctx)

	cmss.mtx.Lock()
	if cmss.data == nil {
		// headers and resource should be the same. set them once
		cmss.data = &metricsExportRequest{
			Header:   headers,
			Resource: rs.GetResource(),
		}
	}
	// concat all metrics
	cmss.data.Metrics = append(cmss.data.Metrics, scopeMetrics.GetMetrics()...)
	cmss.mtx.Unlock()

	return &cmpb.ExportMetricsServiceResponse{}, nil
}
