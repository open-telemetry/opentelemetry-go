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

// Package tracetest is a testing helper package for the SDK. User can
// configure no-op or in-memory exporters to verify different SDK behaviors or
// custom instrumentation.
package logtest // import "go.opentelemetry.io/otel/sdk/trace/tracetest"

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/sdk/log"
)

var _ log.LogRecordExporter = (*NoopExporter)(nil)

// NewNoopExporter returns a new no-op exporter.
func NewNoopExporter() *NoopExporter {
	return new(NoopExporter)
}

// NoopExporter is an exporter that drops all received spans and performs no
// action.
type NoopExporter struct{}

// ExportSpans handles export of spans by dropping them.
func (nsb *NoopExporter) ExportLogRecords(ctx context.Context, spans []log.ReadOnlyLogRecord) error {
	return nil
}

// Shutdown stops the exporter by doing nothing.
func (nsb *NoopExporter) Shutdown(context.Context) error { return nil }

var _ log.LogRecordExporter = (*InMemoryExporter)(nil)

// NewInMemoryExporter returns a new InMemoryExporter.
func NewInMemoryExporter() *InMemoryExporter {
	return new(InMemoryExporter)
}

// InMemoryExporter is an exporter that stores all received spans in-memory.
type InMemoryExporter struct {
	mu sync.Mutex
	ss LogRecordStubs
}

// ExportSpans handles export of spans by storing them in memory.
func (imsb *InMemoryExporter) ExportLogRecords(ctx context.Context, spans []log.ReadOnlyLogRecord) error {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	imsb.ss = append(imsb.ss, LogRecordStubsFromReadOnlyLogRecords(spans)...)
	return nil
}

// Shutdown stops the exporter by clearing spans held in memory.
func (imsb *InMemoryExporter) Shutdown(context.Context) error {
	imsb.Reset()
	return nil
}

// Reset the current in-memory storage.
func (imsb *InMemoryExporter) Reset() {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	imsb.ss = nil
}

// GetSpans returns the current in-memory stored spans.
func (imsb *InMemoryExporter) GetSpans() LogRecordStubs {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	ret := make(LogRecordStubs, len(imsb.ss))
	copy(ret, imsb.ss)
	return ret
}
