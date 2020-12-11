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

package otlp // import "go.opentelemetry.io/otel/exporters/otlp"

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"

	colmetricpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/metrics/v1"
	coltracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/trace/v1"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/export/trace"
)

type grpcDriver struct {
	connection *grpcConnection

	lock          sync.Mutex
	metricsClient colmetricpb.MetricsServiceClient
	tracesClient  coltracepb.TraceServiceClient
}

func (d *grpcDriver) getMetricsClient() grpcMetricsClient {
	d.lock.Lock()
	client := d.metricsClient
	d.lock.Unlock()
	return grpcMetricsClient{
		grpcClientBase: grpcClientBase{
			clientLock: &d.lock,
			connection: d.connection,
		},
		client: client,
	}
}

func (d *grpcDriver) getTracesClient() grpcTracesClient {
	d.lock.Lock()
	client := d.tracesClient
	d.lock.Unlock()
	return grpcTracesClient{
		grpcClientBase: grpcClientBase{
			clientLock: &d.lock,
			connection: d.connection,
		},
		client: client,
	}
}

func (d *grpcDriver) handleNewConnection(cc *grpc.ClientConn) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if cc != nil {
		d.metricsClient = colmetricpb.NewMetricsServiceClient(cc)
		d.tracesClient = coltracepb.NewTraceServiceClient(cc)
	} else {
		d.metricsClient = nil
		d.tracesClient = nil
	}
}

func NewGRPCDriver(opts ...GRPCConnectionOption) ProtocolDriver {
	cfg := grpcConnectionConfig{
		collectorAddr:     fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorPort),
		grpcServiceConfig: DefaultGRPCServiceConfig,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	d := &grpcDriver{}
	d.connection = newGRPCConnection(cfg, d.handleNewConnection)
	return d
}

func (d *grpcDriver) Start(ctx context.Context) error {
	d.connection.startConnection(ctx)
	return nil
}

func (d *grpcDriver) Stop(ctx context.Context) error {
	return d.connection.shutdown(ctx)
}

func (d *grpcDriver) ExportMetrics(ctx context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector) error {
	return uploadMetrics(ctx, cps, selector, d.getMetricsClient())
}

func (d *grpcDriver) ExportTraces(ctx context.Context, ss []*tracesdk.SpanSnapshot) error {
	return uploadTraces(ctx, ss, d.getTracesClient())
}
