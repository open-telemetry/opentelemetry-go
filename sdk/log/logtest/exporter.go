// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package logtest is a testing helper package for the SDK. User can configure
// in-memory exporters to verify different SDK behaviors or custom
// instrumentation.
package logtest // import "go.opentelemetry.io/otel/sdk/log/logtest"

import (
	"context"
	"sync"

	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewInMemoryExporter returns a new InMemoryExporter.
func NewInMemoryExporter() *InMemoryExporter {
	return new(InMemoryExporter)
}

// InMemoryExporter is an exporter that stores all received log records
// in-memory.
type InMemoryExporter struct {
	mu      sync.Mutex
	records []sdklog.Record
}

// Export handles the export of records by storing them in memory.
func (imsb *InMemoryExporter) Export(ctx context.Context, records []sdklog.Record) error {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()

	imsb.records = append(imsb.records, records...)
	return nil
}

// Shutdown stops the exporter by clearing the records held in memory.
func (imsb *InMemoryExporter) Shutdown(context.Context) error {
	imsb.Reset()
	return nil
}

// ForceFlush is a noop method in the context of InMemoryExporter.
func (imsb *InMemoryExporter) ForceFlush(ctx context.Context) error {
	return nil
}

// Reset the current in-memory storage.
func (imsb *InMemoryExporter) Reset() {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	imsb.records = []sdklog.Record{}
}

// GetRecords returns the current in-memory stored log records.
func (imsb *InMemoryExporter) GetRecords() []sdklog.Record {
	imsb.mu.Lock()
	defer imsb.mu.Unlock()
	ret := make([]sdklog.Record, len(imsb.records))
	copy(ret, imsb.records)
	return ret
}
