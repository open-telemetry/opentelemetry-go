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

// code in this package is mostly copied from contrib.go.opencensus.io/exporter/ocagent/connection.go
package otlp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"unsafe"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	colmetricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/metrics/v1"
	coltracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/trace/v1"
	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"

	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	tracesdk "go.opentelemetry.io/otel/sdk/export/trace"
)

type Exporter struct {
	// mu protects the non-atomic and non-channel variables
	mu sync.RWMutex
	// senderMu protects the concurrent unsafe sends on the shared gRPC client connection.
	senderMu          sync.Mutex
	started           bool
	traceExporter     coltracepb.TraceServiceClient
	metricExporter    colmetricpb.MetricsServiceClient
	grpcClientConn    *grpc.ClientConn
	lastConnectErrPtr unsafe.Pointer

	startOnce      sync.Once
	stopCh         chan bool
	disconnectedCh chan bool

	backgroundConnectionDoneCh chan bool

	c Config
}

var _ tracesdk.SpanBatcher = (*Exporter)(nil)
var _ metricsdk.Exporter = (*Exporter)(nil)

func configureOptions(cfg *Config, opts ...ExporterOption) {
	for _, opt := range opts {
		opt(cfg)
	}
}

func NewExporter(opts ...ExporterOption) (*Exporter, error) {
	exp := NewUnstartedExporter(opts...)
	if err := exp.Start(); err != nil {
		return nil, err
	}
	return exp, nil
}

func NewUnstartedExporter(opts ...ExporterOption) *Exporter {
	e := new(Exporter)
	e.c = Config{numWorkers: DefaultNumWorkers}
	configureOptions(&e.c, opts...)

	// TODO (rghetia): add resources

	return e
}

var (
	errAlreadyStarted = errors.New("already started")
	errNotStarted     = errors.New("not started")
)

// Start dials to the collector, establishing a connection to it. It also
// initiates the Config and Trace services by sending over the initial
// messages that consist of the node identifier. Start invokes a background
// connector that will reattempt connections to the collector periodically
// if the connection dies.
func (e *Exporter) Start() error {
	var err = errAlreadyStarted
	e.startOnce.Do(func() {
		e.mu.Lock()
		e.started = true
		e.disconnectedCh = make(chan bool, 1)
		e.stopCh = make(chan bool)
		e.backgroundConnectionDoneCh = make(chan bool)
		e.mu.Unlock()

		// An optimistic first connection attempt to ensure that
		// applications under heavy load can immediately process
		// data. See https://github.com/census-ecosystem/opencensus-go-exporter-ocagent/pull/63
		if err := e.connect(); err == nil {
			e.setStateConnected()
		} else {
			e.setStateDisconnected(err)
		}
		go e.indefiniteBackgroundConnection()

		err = nil
	})

	return err
}

func (e *Exporter) prepareCollectorAddress() string {
	if e.c.collectorAddr != "" {
		return e.c.collectorAddr
	}
	return fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorPort)
}

func (e *Exporter) enableConnections(cc *grpc.ClientConn) error {
	e.mu.RLock()
	started := e.started
	e.mu.RUnlock()

	if !started {
		return errNotStarted
	}

	e.mu.Lock()
	// If previous clientConn is same as the current then just return.
	// This doesn't happen right now as this func is only called with new ClientConn.
	// It is more about future-proofing.
	if e.grpcClientConn == cc {
		e.mu.Unlock()
		return nil
	}
	// If the previous clientConn was non-nil, close it
	if e.grpcClientConn != nil {
		_ = e.grpcClientConn.Close()
	}
	e.grpcClientConn = cc
	e.traceExporter = coltracepb.NewTraceServiceClient(cc)
	e.metricExporter = colmetricpb.NewMetricsServiceClient(cc)
	e.mu.Unlock()

	return nil
}

func (e *Exporter) dialToCollector() (*grpc.ClientConn, error) {
	addr := e.prepareCollectorAddress()
	var dialOpts []grpc.DialOption
	if e.c.clientCredentials != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(e.c.clientCredentials))
	} else if e.c.canDialInsecure {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	}
	if e.c.compressor != "" {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor(e.c.compressor)))
	}
	if len(e.c.grpcDialOptions) != 0 {
		dialOpts = append(dialOpts, e.c.grpcDialOptions...)
	}

	ctx := context.Background()
	if len(e.c.headers) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, metadata.New(e.c.headers))
	}
	return grpc.DialContext(ctx, addr, dialOpts...)
}

// Stop shuts down all the connections and resources
// related to the exporter.
// If the exporter is not started then this func does nothing.
func (e *Exporter) Stop() error {
	e.mu.RLock()
	cc := e.grpcClientConn
	started := e.started
	e.mu.RUnlock()

	if !started {
		return nil
	}

	// Now close the underlying gRPC connection.
	var err error
	if cc != nil {
		err = cc.Close()
	}

	// At this point we can change the state variable started
	e.mu.Lock()
	e.started = false
	e.mu.Unlock()
	close(e.stopCh)

	// Ensure that the backgroundConnector returns
	<-e.backgroundConnectionDoneCh

	return err
}

// Export implements the "go.opentelemetry.io/otel/sdk/export/metric".Exporter
// interface. It transforms metric Records into OTLP Metrics and transmits them.
func (e *Exporter) Export(ctx context.Context, cps metricsdk.CheckpointSet) error {
	// Seed records into the work processing pool.
	records := make(chan metricsdk.Record)
	go func() {
		_ = cps.ForEach(func(record metricsdk.Record) (err error) {
			select {
			case <-e.stopCh:
			case <-ctx.Done():
			case records <- record:
			}
			return
		})
		close(records)
	}()

	// Allow all errors to be collected and returned singularly.
	errCh := make(chan error)
	var errStrings []string
	go func() {
		for err := range errCh {
			if err != nil {
				errStrings = append(errStrings, err.Error())
			}
		}
	}()

	// Start the work processing pool.
	processed := make(chan *metricpb.Metric)
	var wg sync.WaitGroup
	for i := uint(0); i < e.c.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.processMetrics(ctx, processed, errCh, records)
		}()
	}
	go func() {
		wg.Wait()
		close(processed)
	}()

	// Synchronosly collected the processed records and transmit.
	e.uploadMetrics(ctx, processed, errCh)

	// Now that all processing is done, handle any errors seen.
	close(errCh)
	if len(errStrings) > 0 {
		return fmt.Errorf("errors exporting:\n -%s", strings.Join(errStrings, "\n -"))
	}
	return nil
}

func (e *Exporter) processMetrics(ctx context.Context, out chan<- *metricpb.Metric, errCh chan<- error, in <-chan metricsdk.Record) {
	for r := range in {
		m, err := transform.Record(r)
		if err != nil {
			if err == aggregator.ErrNoData {
				// The Aggregator was checkpointed before the first value
				// was set, skipping.
				continue
			}
			select {
			case <-e.stopCh:
				return
			case <-ctx.Done():
				return
			case errCh <- err:
				continue
			}
		}

		select {
		case <-e.stopCh:
			return
		case <-ctx.Done():
			return
		case out <- m:
		}
	}
}

func (e *Exporter) uploadMetrics(ctx context.Context, in <-chan *metricpb.Metric, errCh chan<- error) {
	var protoMetrics []*metricpb.Metric
	for m := range in {
		protoMetrics = append(protoMetrics, m)
	}

	if len(protoMetrics) == 0 {
		return
	}
	if !e.connected() {
		return
	}

	rm := []*metricpb.ResourceMetrics{
		{
			Resource: nil,
			InstrumentationLibraryMetrics: []*metricpb.InstrumentationLibraryMetrics{
				{
					Metrics: protoMetrics,
				},
			},
		},
	}

	select {
	case <-e.stopCh:
		return
	case <-ctx.Done():
		return
	default:
		e.senderMu.Lock()
		_, err := e.metricExporter.Export(ctx, &colmetricpb.ExportMetricsServiceRequest{
			ResourceMetrics: rm,
		})
		e.senderMu.Unlock()
		if err != nil {
			select {
			case <-e.stopCh:
				return
			case <-ctx.Done():
				return
			case errCh <- err:
			}
		}
	}
}

func (e *Exporter) ExportSpan(ctx context.Context, sd *tracesdk.SpanData) {
	e.uploadTraces(ctx, []*tracesdk.SpanData{sd})
}

func (e *Exporter) ExportSpans(ctx context.Context, sds []*tracesdk.SpanData) {
	e.uploadTraces(ctx, sds)
}

func (e *Exporter) uploadTraces(ctx context.Context, sdl []*tracesdk.SpanData) {
	select {
	case <-e.stopCh:
		return

	default:
		if !e.connected() {
			return
		}

		protoSpans := transform.SpanData(sdl)
		if len(protoSpans) == 0 {
			return
		}

		e.senderMu.Lock()
		_, err := e.traceExporter.Export(ctx, &coltracepb.ExportTraceServiceRequest{
			ResourceSpans: protoSpans,
		})
		e.senderMu.Unlock()
		if err != nil {
			e.setStateDisconnected(err)
		}
	}
}
