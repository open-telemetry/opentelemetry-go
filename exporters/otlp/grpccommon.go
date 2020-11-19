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

	colmetricpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/metrics/v1"
	coltracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/trace/v1"
	metricpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/metrics/v1"
	tracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/trace/v1"
	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/export/trace"
)

type grpcClientBase struct {
	clientLock *sync.Mutex
	connection *grpcConnection
}

func (c grpcClientBase) check() error {
	if !c.connection.connected() {
		return errDisconnected
	}
	return nil
}

func (c grpcClientBase) unifyContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return c.connection.contextWithStop(ctx)
}

type grpcMetricsClient struct {
	grpcClientBase
	client colmetricpb.MetricsServiceClient
}

func (c grpcMetricsClient) uploadMetrics(ctx context.Context, protoMetrics []*metricpb.ResourceMetrics) error {
	ctx = c.connection.contextWithMetadata(ctx)
	err := func() error {
		c.clientLock.Lock()
		defer c.clientLock.Unlock()
		if c.client == nil {
			return errNoClient
		}
		_, err := c.client.Export(ctx, &colmetricpb.ExportMetricsServiceRequest{
			ResourceMetrics: protoMetrics,
		})
		return err
	}()
	if err != nil {
		c.connection.setStateDisconnected(err)
	}
	return err
}

type grpcTracesClient struct {
	grpcClientBase
	client coltracepb.TraceServiceClient
}

func (c grpcTracesClient) uploadTraces(ctx context.Context, protoSpans []*tracepb.ResourceSpans) error {
	ctx = c.connection.contextWithMetadata(ctx)
	err := func() error {
		c.clientLock.Lock()
		defer c.clientLock.Unlock()
		if c.client == nil {
			return errNoClient
		}
		_, err := c.client.Export(ctx, &coltracepb.ExportTraceServiceRequest{
			ResourceSpans: protoSpans,
		})
		return err
	}()
	if err != nil {
		c.connection.setStateDisconnected(err)
	}
	return err
}

func uploadMetrics(ctx context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector, client grpcMetricsClient) error {
	if err := client.check(); err != nil {
		return err
	}
	ctx, cancel := client.unifyContext(ctx)
	defer cancel()

	rms, err := transform.CheckpointSet(ctx, selector, cps, 1)
	if err != nil {
		return err
	}
	if len(rms) == 0 {
		return nil
	}

	return client.uploadMetrics(ctx, rms)
}

func uploadTraces(ctx context.Context, ss []*tracesdk.SpanSnapshot, client grpcTracesClient) error {
	if err := client.check(); err != nil {
		return err
	}
	ctx, cancel := client.unifyContext(ctx)
	defer cancel()

	protoSpans := transform.SpanData(ss)
	if len(protoSpans) == 0 {
		return nil
	}

	return client.uploadTraces(ctx, protoSpans)
}
