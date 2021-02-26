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

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

// SimpleSpanProcessor is a SpanProcessor that synchronously sends all
// completed Spans to a trace.Exporter immediately.
type SimpleSpanProcessor struct {
	exporterMu sync.RWMutex
	exporter   export.SpanExporter
	stopOnce   sync.Once
}

var _ SpanProcessor = (*SimpleSpanProcessor)(nil)

// NewSimpleSpanProcessor returns a new SimpleSpanProcessor.
func NewSimpleSpanProcessor(exporter export.SpanExporter) *SimpleSpanProcessor {
	ssp := &SimpleSpanProcessor{
		exporter: exporter,
	}
	return ssp
}

// OnStart does nothing.
func (ssp *SimpleSpanProcessor) OnStart(context.Context, ReadWriteSpan) {}

// OnEnd immediately exports a ReadOnlySpan.
func (ssp *SimpleSpanProcessor) OnEnd(s ReadOnlySpan) {
	ssp.exporterMu.RLock()
	defer ssp.exporterMu.RUnlock()

	if ssp.exporter != nil && s.SpanContext().IsSampled() {
		ss := s.Snapshot()
		if err := ssp.exporter.ExportSpans(context.Background(), []*export.SpanSnapshot{ss}); err != nil {
			otel.Handle(err)
		}
	}
}

// Shutdown shuts down the exporter this SimpleSpanProcessor exports to.
func (ssp *SimpleSpanProcessor) Shutdown(ctx context.Context) error {
	ssp.exporterMu.Lock()
	defer ssp.exporterMu.Unlock()

	var err error
	ssp.stopOnce.Do(func() {
		err = ssp.exporter.Shutdown(ctx)
		// Set exporter to nil so subsequent calls to OnEnd are ignored
		// gracefully.
		ssp.exporter = nil
	})
	return err
}

// ForceFlush does nothing as there is no data to flush.
func (ssp *SimpleSpanProcessor) ForceFlush() {
}
