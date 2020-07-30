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

// tracetest is a testing helper package for the SDK. User can configure no-op or in-memory exporters to verify
// different SDK behaviors or custom instrumentation.
package tracetest // import "go.opentelemetry.io/otel/sdk/export/trace/tracetest"

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/sdk/export/trace"
)

var _ trace.SpanBatcher = (*NoopExporter)(nil)
var _ trace.SpanSyncer = (*NoopExporter)(nil)

// NewNoopExporter returns a new no-op exporter.
// It implements both trace.SpanBatcher and trace.SpanSyncer.
func NewNoopExporter() *NoopExporter {
	return new(NoopExporter)
}

// NoopExporter is an exporter that does nothing.
type NoopExporter struct{}

// ExportSpans implements the trace.SpanBatcher interface.
func (nsb *NoopExporter) ExportSpans(context.Context, []*trace.SpanData) {}

// ExportSpan implements the trace.SpanSyncer interface.
func (nsb *NoopExporter) ExportSpan(context.Context, *trace.SpanData) {}

var _ trace.SpanBatcher = (*InMemoryExporter)(nil)
var _ trace.SpanSyncer = (*InMemoryExporter)(nil)

// NewInMemoryExporter returns a new trace.SpanBatcher that stores in-memory all exported spans.
// It implements both trace.SpanBatcher and trace.SpanSyncer.
func NewInMemoryExporter() *InMemoryExporter {
	return new(InMemoryExporter)
}

// InMemoryExporter is an exporter that stores in-memory all exported spans.
type InMemoryExporter struct {
	mu  sync.Mutex
	sds []*trace.SpanData
}

// ExportSpans implements the trace.SpanBatcher interface.
func (imsb *InMemoryExporter) ExportSpans(_ context.Context, sds []*trace.SpanData) {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	imsb.sds = append(imsb.sds, sds...)
}

// ExportSpan implements the trace.SpanSyncer interface.
func (imsb *InMemoryExporter) ExportSpan(_ context.Context, sd *trace.SpanData) {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	imsb.sds = append(imsb.sds, sd)
}

// Reset the current in-memory storage.
func (imsb *InMemoryExporter) Reset() {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	imsb.sds = nil
}

// GetSpans returns the current in-memory stored spans.
func (imsb *InMemoryExporter) GetSpans() []*trace.SpanData {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	ret := make([]*trace.SpanData, len(imsb.sds))
	copy(ret, imsb.sds)
	return ret
}
