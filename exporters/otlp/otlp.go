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
	commonpb "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	metricpb "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
	resourcepb "github.com/open-telemetry/opentelemetry-proto/gen/go/resource/v1"

	"go.opentelemetry.io/otel/exporters/otlp/internal/transform"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
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
	errAlreadyStarted  = errors.New("already started")
	errNotStarted      = errors.New("not started")
	errDisconnected    = errors.New("exporter disconnected")
	errStopped         = errors.New("exporter stopped")
	errContextCanceled = errors.New("context canceled")
	errTransforming    = errors.New("transforming failed")
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

// result is the product of transforming Records into OTLP Metrics.
type result struct {
	Resource resource.Resource
	Library  string
	Metric   *metricpb.Metric
	Err      error
}

// Export implements the "go.opentelemetry.io/otel/sdk/export/metric".Exporter
// interface. It transforms metric Records into OTLP Metrics and transmits them.
func (e *Exporter) Export(ctx context.Context, cps metricsdk.CheckpointSet) error {
	records, errc := e.source(ctx, cps)

	// Start a fixed number of goroutines to transform records.
	transformed := make(chan result)
	var wg sync.WaitGroup
	wg.Add(int(e.c.numWorkers))
	for i := uint(0); i < e.c.numWorkers; i++ {
		go func() {
			defer wg.Done()
			e.transformer(ctx, records, transformed)
		}()
	}
	go func() {
		wg.Wait()
		close(transformed)
	}()

	// Synchronosly collect the transformed records and transmit.
	err := e.sink(ctx, transformed)
	if err != nil {
		return err
	}

	// source is complete, check for any errors.
	if err := <-errc; err != nil {
		return err
	}
	return nil
}

// source starts a goroutine that sends each one of the Records yielded by
// the CheckpointSet on the returned chan. Any error encoutered will be sent
// on the returned error chan after seeding is complete.
func (e *Exporter) source(ctx context.Context, cps metricsdk.CheckpointSet) (<-chan metricsdk.Record, <-chan error) {
	errc := make(chan error, 1)
	out := make(chan metricsdk.Record)
	// Seed records into process.
	go func() {
		defer close(out)
		// No selected needed since errc is buffered.
		errc <- cps.ForEach(func(r metricsdk.Record) error {
			select {
			case <-e.stopCh:
				return errStopped
			case <-ctx.Done():
				return errContextCanceled
			case out <- r:
			}
			return nil
		})
	}()
	return out, errc
}

// transformer transforms records read from the passed in chan into
// OTLP Metrics which are sent on the out chan.
func (e *Exporter) transformer(ctx context.Context, in <-chan metricsdk.Record, out chan<- result) {
	for r := range in {
		m, err := transform.Record(r)
		// Propagate errors, but do not send empty results.
		if err == nil && m == nil {
			continue
		}
		res := result{
			Resource: r.Descriptor().Resource(),
			Library:  r.Descriptor().LibraryName(),
			Metric:   m,
			Err:      err,
		}
		select {
		case <-e.stopCh:
			return
		case <-ctx.Done():
			return
		case out <- res:
		}
	}
}

// sink collects transformed Records, batches them by id, and exports them.
func (e *Exporter) sink(ctx context.Context, in <-chan result) error {
	var errStrings []string

	type resourceBatch struct {
		Resource *resourcepb.Resource
		// Group by instrumentation library name and then the MetricDescriptor.
		InstrumentationLibraryBatches map[string]map[string]*metricpb.Metric
	}

	// group by unique Resource string.
	grouped := make(map[string]resourceBatch)
	for res := range in {
		if res.Err != nil {
			errStrings = append(errStrings, res.Err.Error())
			continue
		}

		rb, ok := grouped[res.Resource.String()]
		if !ok {
			rb = resourceBatch{
				Resource:                      transform.Resource(&res.Resource),
				InstrumentationLibraryBatches: make(map[string]map[string]*metricpb.Metric),
			}
			grouped[res.Resource.String()] = rb
		}

		mb, ok := rb.InstrumentationLibraryBatches[res.Library]
		if !ok {
			mb = make(map[string]*metricpb.Metric)
			rb.InstrumentationLibraryBatches[res.Library] = mb
		}

		m, ok := mb[res.Metric.GetMetricDescriptor().String()]
		if !ok {
			mb[res.Metric.GetMetricDescriptor().String()] = res.Metric
			continue
		}
		if len(res.Metric.Int64DataPoints) > 0 {
			m.Int64DataPoints = append(m.Int64DataPoints, res.Metric.Int64DataPoints...)
		}
		if len(res.Metric.DoubleDataPoints) > 0 {
			m.DoubleDataPoints = append(m.DoubleDataPoints, res.Metric.DoubleDataPoints...)
		}
		if len(res.Metric.HistogramDataPoints) > 0 {
			m.HistogramDataPoints = append(m.HistogramDataPoints, res.Metric.HistogramDataPoints...)
		}
		if len(res.Metric.SummaryDataPoints) > 0 {
			m.SummaryDataPoints = append(m.SummaryDataPoints, res.Metric.SummaryDataPoints...)
		}
	}

	if len(grouped) == 0 {
		return nil
	}
	if !e.connected() {
		return errDisconnected
	}

	var rms []*metricpb.ResourceMetrics
	for _, rb := range grouped {
		rm := &metricpb.ResourceMetrics{Resource: rb.Resource}
		for ilName, mb := range rb.InstrumentationLibraryBatches {
			ilm := &metricpb.InstrumentationLibraryMetrics{
				Metrics: make([]*metricpb.Metric, 0, len(mb)),
			}
			if ilName != "" {
				ilm.InstrumentationLibrary = &commonpb.InstrumentationLibrary{Name: ilName}
			}
			for _, m := range mb {
				ilm.Metrics = append(ilm.Metrics, m)
			}
			rm.InstrumentationLibraryMetrics = append(rm.InstrumentationLibraryMetrics, ilm)
		}
		rms = append(rms, rm)
	}

	select {
	case <-e.stopCh:
		return errStopped
	case <-ctx.Done():
		return errContextCanceled
	default:
		e.senderMu.Lock()
		_, err := e.metricExporter.Export(ctx, &colmetricpb.ExportMetricsServiceRequest{
			ResourceMetrics: rms,
		})
		e.senderMu.Unlock()
		if err != nil {
			return err
		}
	}

	// Report any transformer errors.
	if len(errStrings) > 0 {
		return fmt.Errorf("%w:\n -%s", errTransforming, strings.Join(errStrings, "\n -"))
	}
	return nil
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
