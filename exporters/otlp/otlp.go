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

// This code was based on
// contrib.go.opencensus.io/exporter/ocagent/connection.go

import (
	"context"
	"errors"
	"sync"

	"google.golang.org/grpc"

	colmetricpb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/metrics/v1"
	coltracepb "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/collector/trace/v1"
	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	tracesdk "go.opentelemetry.io/otel/sdk/export/trace"
)

// Exporter is an OpenTelemetry exporter. It exports both traces and metrics
// from OpenTelemetry instrumented to code using OpenTelemetry protocol
// buffers to a configurable receiver.
type Exporter struct {
	// mu protects the non-atomic and non-channel variables
	mu sync.RWMutex
	// senderMu protects the concurrent unsafe sends on the shared gRPC client connection.
	senderMu       sync.Mutex
	started        bool
	traceExporter  coltracepb.TraceServiceClient
	metricExporter colmetricpb.MetricsServiceClient
	cc             *grpcConnection

	startOnce sync.Once
	stopOnce  sync.Once

	exportKindSelector metricsdk.ExportKindSelector
}

var _ tracesdk.SpanExporter = (*Exporter)(nil)
var _ metricsdk.Exporter = (*Exporter)(nil)

// newConfig initializes a config struct with default values and applies
// any ExporterOptions provided.
func newConfig(opts ...ExporterOption) config {
	cfg := config{
		grpcServiceConfig: DefaultGRPCServiceConfig,

		// Note: the default ExportKindSelector is specified
		// as Cumulative:
		// https://github.com/open-telemetry/opentelemetry-specification/issues/731
		exportKindSelector: metricsdk.CumulativeExportKindSelector(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// NewExporter constructs a new Exporter and starts it.
func NewExporter(ctx context.Context, opts ...ExporterOption) (*Exporter, error) {
	exp := NewUnstartedExporter(opts...)
	if err := exp.Start(ctx); err != nil {
		return nil, err
	}
	return exp, nil
}

// NewUnstartedExporter constructs a new Exporter and does not start it.
func NewUnstartedExporter(opts ...ExporterOption) *Exporter {
	e := new(Exporter)
	cfg := newConfig(opts...)
	e.exportKindSelector = cfg.exportKindSelector
	e.cc = newGRPCConnection(cfg, e.handleNewConnection)
	return e
}

func (e *Exporter) handleNewConnection(cc *grpc.ClientConn) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if cc != nil {
		e.metricExporter = colmetricpb.NewMetricsServiceClient(cc)
		e.traceExporter = coltracepb.NewTraceServiceClient(cc)
	} else {
		e.metricExporter = nil
		e.traceExporter = nil
	}
	return nil
}

var (
	errNoClient       = errors.New("no client")
	errAlreadyStarted = errors.New("already started")
	errDisconnected   = errors.New("exporter disconnected")
)

// Start dials to the collector, establishing a connection to it. It also
// initiates the Config and Trace services by sending over the initial
// messages that consist of the node identifier. Start invokes a background
// connector that will reattempt connections to the collector periodically
// if the connection dies.
func (e *Exporter) Start(ctx context.Context) error {
	var err = errAlreadyStarted
	e.startOnce.Do(func() {
		e.mu.Lock()
		e.started = true
		e.mu.Unlock()

		err = nil
		e.cc.startConnection(ctx)
	})

	return err
}

// Shutdown closes all connections and releases resources currently being used
// by the exporter. If the exporter is not started this does nothing.
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.mu.RLock()
	cc := e.cc
	started := e.started
	e.mu.RUnlock()

	if !started {
		return nil
	}

	var err error

	e.stopOnce.Do(func() {
		// Clean things up before checking this error.
		err = cc.shutdown(ctx)

		// At this point we can change the state variable started
		e.mu.Lock()
		e.started = false
		e.mu.Unlock()
	})

	return err
}

// Export implements the "go.opentelemetry.io/otel/sdk/export/metric".Exporter
// interface. It transforms and batches metric Records into OTLP Metrics and
// transmits them to the configured collector.
func (e *Exporter) Export(parent context.Context, cps metricsdk.CheckpointSet) error {
	ctx, cancel := e.cc.contextWithStop(parent)
	defer cancel()

	// Hardcode the number of worker goroutines to 1. We later will
	// need to see if there's a way to adjust that number for longer
	// running operations.
	rms, err := transform.CheckpointSet(ctx, e, cps, 1)
	if err != nil {
		return err
	}

	if !e.cc.connected() {
		return errDisconnected
	}

	err = func() error {
		e.senderMu.Lock()
		defer e.senderMu.Unlock()
		if e.metricExporter == nil {
			return errNoClient
		}
		_, err := e.metricExporter.Export(e.cc.contextWithMetadata(ctx), &colmetricpb.ExportMetricsServiceRequest{
			ResourceMetrics: rms,
		})
		return err
	}()
	if err != nil {
		e.cc.setStateDisconnected(err)
	}
	return err
}

// ExportKindFor reports back to the OpenTelemetry SDK sending this Exporter
// metric telemetry that it needs to be provided in a cumulative format.
func (e *Exporter) ExportKindFor(desc *metric.Descriptor, kind aggregation.Kind) metricsdk.ExportKind {
	return e.exportKindSelector.ExportKindFor(desc, kind)
}

// ExportSpans exports a batch of SpanData.
func (e *Exporter) ExportSpans(ctx context.Context, sds []*tracesdk.SpanData) error {
	return e.uploadTraces(ctx, sds)
}

func (e *Exporter) uploadTraces(ctx context.Context, sdl []*tracesdk.SpanData) error {
	ctx, cancel := e.cc.contextWithStop(ctx)
	defer cancel()

	if !e.cc.connected() {
		return nil
	}

	protoSpans := transform.SpanData(sdl)
	if len(protoSpans) == 0 {
		return nil
	}

	err := func() error {
		e.senderMu.Lock()
		defer e.senderMu.Unlock()
		if e.traceExporter == nil {
			return errNoClient
		}
		_, err := e.traceExporter.Export(e.cc.contextWithMetadata(ctx), &coltracepb.ExportTraceServiceRequest{
			ResourceSpans: protoSpans,
		})
		return err
	}()
	if err != nil {
		e.cc.setStateDisconnected(err)
	}
	return err
}
