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

package tracetest // import "go.opentelemetry.io/otel/sdk/export/trace"

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/sdk/export/trace"
)

// NewNoopSpanBatcher returns a new no-op trace.SpanBatcher.
func NewNoopSpanBatcher() trace.SpanBatcher {
	return new(noopSpanBatcher)
}

type noopSpanBatcher struct {
}

// ExportSpans implements the trace.SpanBatcher interface.
func (nsb *noopSpanBatcher) ExportSpans(context.Context, []*trace.SpanData) {}

// NewInMemorySpanBatcher returns a new trace.SpanBatcher that stores in-memory all exported spans.
func NewInMemorySpanBatcher() *InMemorySpanBatcher {
	return new(InMemorySpanBatcher)
}

type InMemorySpanBatcher struct {
	mu  sync.Mutex
	sds []*trace.SpanData
}

// ExportSpans implements the trace.SpanBatcher interface.
func (imsb *InMemorySpanBatcher) ExportSpans(_ context.Context, sds []*trace.SpanData) {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	imsb.sds = append(imsb.sds, sds...)
}

// Reset the current in-memory storage.
func (imsb *InMemorySpanBatcher) Reset() {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	imsb.sds = nil
}

// GetSpans returns the current in-memory stored spans.
func (imsb *InMemorySpanBatcher) GetSpans() []*trace.SpanData {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	ret := make([]*trace.SpanData, len(imsb.sds))
	copy(ret, imsb.sds)
	return ret
}
