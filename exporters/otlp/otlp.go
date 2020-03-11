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

// code in this package is mostly copied from contrib.go.opencensus.io/exporter/ocagent/connection.go
package otlp

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	coltracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/collector/trace/v1"
	tracepb "github.com/open-telemetry/opentelemetry-proto/gen/go/trace/v1"

	export "go.opentelemetry.io/otel/sdk/export/trace"
)

type Exporter struct {
	// mu protects the non-atomic and non-channel variables
	mu sync.RWMutex
	// senderMu protects the concurrent unsafe send on traceExporter client
	senderMu          sync.Mutex
	started           bool
	traceExporter     coltracepb.TraceServiceClient
	grpcClientConn    *grpc.ClientConn
	lastConnectErrPtr unsafe.Pointer

	startOnce      sync.Once
	stopCh         chan bool
	disconnectedCh chan bool

	backgroundConnectionDoneCh chan bool

	c Config
}

var _ export.SpanBatcher = (*Exporter)(nil)

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
	e.c = Config{}
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

func (e *Exporter) ExportSpans(ctx context.Context, sds []*export.SpanData) {
	e.uploadTraces(ctx, sds)
}

func otSpanDataToPbSpans(sdl []*export.SpanData) []*tracepb.ResourceSpans {
	if len(sdl) == 0 {
		return nil
	}
	protoSpans := make([]*tracepb.Span, 0, len(sdl))
	for _, sd := range sdl {
		if sd != nil {
			protoSpans = append(protoSpans, otSpanToProtoSpan(sd))
		}
	}
	return []*tracepb.ResourceSpans{
		{
			Resource: nil,
			Spans:    protoSpans,
		},
	}
}

func (e *Exporter) uploadTraces(ctx context.Context, sdl []*export.SpanData) {
	select {
	case <-e.stopCh:
		return

	default:
		if !e.connected() {
			return
		}

		protoSpans := otSpanDataToPbSpans(sdl)
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
