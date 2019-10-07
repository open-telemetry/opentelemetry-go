// Copyright 2019, OpenTelemetry Authors
// Copyright 2017, OpenCensus Authors
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

package stackdriver

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	traceclient "cloud.google.com/go/trace/apiv2"
	"github.com/golang/protobuf/proto"
	"google.golang.org/api/support/bundler"
	tracepb "google.golang.org/genproto/googleapis/devtools/cloudtrace/v2"

	"go.opentelemetry.io/sdk/trace"
)

// traceExporter is an imeplementation of trace.Exporter and trace.BatchExporter
// that uploads spans to Stackdriver Trace in batch.
type traceExporter struct {
	o         *options
	projectID string
	bundler   *bundler.Bundler
	// uploadFn defaults in uploadSpans; it can be replaced for tests.
	uploadFn func(spans []*tracepb.Span)
	overflowLogger
	client *traceclient.Client
}

func newTraceExporter(o *options) (*traceExporter, error) {
	ctx := o.Context
	if ctx == nil {
		ctx = context.Background()
	}
	client, err := traceclient.NewClient(ctx, o.TraceClientOptions...)
	if err != nil {
		return nil, fmt.Errorf("Stackdriver: couldn't initiate trace client: %v", err)
	}
	return newTraceExporterWithClient(o, client), nil
}

const defaultBufferedByteLimit = 0 * 1024 * 1024

func newTraceExporterWithClient(o *options, c *traceclient.Client) *traceExporter {
	e := &traceExporter{
		projectID: o.ProjectID,
		client:    c,
		o:         o,
	}
	b := bundler.NewBundler((*tracepb.Span)(nil), func(bundle interface{}) {
		e.uploadFn(bundle.([]*tracepb.Span))
	})
	if o.BundleDelayThreshold > 0 {
		b.DelayThreshold = o.BundleDelayThreshold
	} else {
		b.DelayThreshold = 2 * time.Second
	}
	if o.BundleCountThreshold > 0 {
		b.BundleCountThreshold = o.BundleCountThreshold
	} else {
		b.BundleCountThreshold = 50
	}
	// The measured "bytes" are not really bytes, see exportReceiver.
	b.BundleByteThreshold = b.BundleCountThreshold * 200
	b.BundleByteLimit = b.BundleCountThreshold * 1000
	if o.TraceSpansBufferMaxBytes > 0 {
		b.BundleByteLimit = o.TraceSpansBufferMaxBytes
	} else {
		b.BundleByteLimit = defaultBufferedByteLimit
	}

	e.bundler = b
	e.uploadFn = e.uploadSpans
	return e
}

// ExportSpan exports a SpanData to Stackdriver Trace.
func (e *traceExporter) ExportSpan(s *trace.SpanData) {
	protoSpan := protoFromSpanData(s, e.projectID)
	protoSize := proto.Size(protoSpan)
	err := e.bundler.Add(protoSpan, protoSize)
	switch err {
	case nil:
		return
	case bundler.ErrOversizedItem:
	case bundler.ErrOverflow:
		e.overflowLogger.log()
	default:
		e.o.handleError(err)
	}
}

// ExportSpans exports a slice of SpanData to Stackdriver Trace in batch
func (e *traceExporter) ExportSpans(sds []*trace.SpanData) {

}

// Shutdown waits for exported trace spans to be uploaded.
//
// This is useful if your program is ending and you do not want to lose recent
// spans.
func (e *traceExporter) Shutdown() {
	e.bundler.Flush()
}

// uploadSpans uploads a set of spans to Stackdriver.
func (e *traceExporter) uploadSpans(spans []*tracepb.Span) {
	req := tracepb.BatchWriteSpansRequest{
		Name:  "projects/" + e.projectID,
		Spans: spans,
	}
	// Create a never-sampled span to prevent traces associated with exporter.
	ctx, cancel := newContextWithTimeout(e.o.Context, e.o.Timeout)
	defer cancel()

	// TODO(ymotongpoo): add this part after OTel support NeverSampler
	// for tracer.Start() initialization.
	//
	// tracer := apitrace.Register()
	// ctx, span := tracer.Start(
	// 	ctx,
	// 	"go.opentelemetry.io/exporter/stackdriver.uploadSpans",
	// )
	// defer span.End()
	// span.SetAttribute(key.New("num_spans").Int64(int64(len(spans))))

	err := e.client.BatchWriteSpans(ctx, &req)
	if err != nil {
		// TODO(ymotongpoo): handle detailed error categories
		// span.SetStatus(codes.Unknown)
		e.o.handleError(err)
	}
}

// overflowLogger ensures that at most one overflow error log message is
// written every 5 seconds.
type overflowLogger struct {
	mu    sync.Mutex
	pause bool
	accum int
}

func (o *overflowLogger) delay() {
	o.pause = true
	time.AfterFunc(5*time.Second, func() {
		o.mu.Lock()
		defer o.mu.Unlock()
		switch {
		case o.accum == 0:
			o.pause = false
		case o.accum == 1:
			log.Println("OpenTelemetry Stackdriver exporter: failed to upload span: buffer full")
			o.accum = 0
			o.delay()
		default:
			log.Printf("OpenTelemetry Stackdriver exporter: failed to upload %d spans: buffer full", o.accum)
			o.accum = 0
			o.delay()
		}
	})
}

func (o *overflowLogger) log() {
	o.mu.Lock()
	defer o.mu.Unlock()
	if !o.pause {
		log.Println("OpenTelemetry Stackdriver exporter: failed to upload span: buffer full")
		o.delay()
	} else {
		o.accum++
	}
}
