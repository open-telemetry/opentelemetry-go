// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"sync"
)

// Compile-time check SimpleProcessor implements Processor.
var _ Processor = (*SimpleProcessor)(nil)

// SimpleProcessor is an processor that synchronously exports log records.
type SimpleProcessor struct {
	exporter Exporter
}

// NewSimpleProcessor is a simple Processor adapter.
//
// This Processor is not recommended for production use. The synchronous
// nature of this Processor make it good for testing, debugging, or
// showing examples of other features, but it can be slow and have a high
// computation resource usage overhead. [NewBatchProcessor] is recommended
// for production use instead.
func NewSimpleProcessor(exporter Exporter, _ ...SimpleProcessorOption) *SimpleProcessor {
	if exporter == nil {
		// Do not panic on nil exporter.
		exporter = defaultNoopExporter
	}
	return &SimpleProcessor{exporter: exporter}
}

var simpleProcRecordsPool = sync.Pool{
	New: func() any {
		records := make([]Record, 1)
		return &records
	},
}

// OnEmit batches provided log record.
func (s *SimpleProcessor) OnEmit(ctx context.Context, r Record) error {
	records := simpleProcRecordsPool.Get().(*[]Record)
	(*records)[0] = r
	defer func() {
		simpleProcRecordsPool.Put(records)
	}()

	return s.exporter.Export(ctx, *records)
}

// Enabled returns true.
func (s *SimpleProcessor) Enabled(context.Context, Record) bool {
	return true
}

// Shutdown shuts down the expoter.
func (s *SimpleProcessor) Shutdown(ctx context.Context) error {
	return s.exporter.Shutdown(ctx)
}

// ForceFlush flushes the exporter.
func (s *SimpleProcessor) ForceFlush(ctx context.Context) error {
	return s.exporter.ForceFlush(ctx)
}

// SimpleProcessorOption applies a configuration to a [SimpleProcessor].
type SimpleProcessorOption interface {
	apply()
}
