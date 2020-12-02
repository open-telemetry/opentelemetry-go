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
	"sync"

	"google.golang.org/grpc"

	colmetricpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/metrics/v1"
	coltracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/trace/v1"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/export/trace"
)

type grpcSingleConnectionDriver struct {
	connection *grpcConnection

	lock          sync.Mutex
	metricsClient colmetricpb.MetricsServiceClient
	tracesClient  coltracepb.TraceServiceClient
}

func (d *grpcSingleConnectionDriver) getMetricsClient() grpcMetricsClient {
	return grpcMetricsClient{
		grpcClientBase: grpcClientBase{
			clientLock: &d.lock,
			connection: d.connection,
		},
		client: d.metricsClient,
	}
}

func (d *grpcSingleConnectionDriver) getTracesClient() grpcTracesClient {
	return grpcTracesClient{
		grpcClientBase: grpcClientBase{
			clientLock: &d.lock,
			connection: d.connection,
		},
		client: d.tracesClient,
	}
}

func (d *grpcSingleConnectionDriver) handleNewConnection(cc *grpc.ClientConn) {
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

func NewGRPCSingleConnectionDriver(cfg GRPCConnectionConfig) ProtocolDriver {
	d := &grpcSingleConnectionDriver{}
	d.connection = newGRPCConnection(cfg, d.handleNewConnection)
	return d
}

func (d *grpcSingleConnectionDriver) Start(ctx context.Context) error {
	d.connection.startConnection(ctx)
	return nil
}

func (d *grpcSingleConnectionDriver) Stop(ctx context.Context) error {
	return d.connection.shutdown(ctx)
}

func (d *grpcSingleConnectionDriver) ExportMetrics(ctx context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector) error {
	return uploadMetrics(ctx, cps, selector, d.getMetricsClient())
}

func (d *grpcSingleConnectionDriver) ExportTraces(ctx context.Context, ss []*tracesdk.SpanSnapshot) error {
	return uploadTraces(ctx, ss, d.getTracesClient())
}
