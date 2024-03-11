// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"sync"
)

var _ Processor = (*SimpleProcessor)(nil)

// SimpleProcessor is an processor that synchronously exports log records.
type SimpleProcessor struct {
	exporter Exporter
	pool     sync.Pool
}

// NewSimpleProcessor is a simple Processor adapter.
//
// Any of the Exporter's methods may be called concurrently with itself
// or with other methods. It is the responsibility of the Exporter to manage
// this concurrency.
func NewSimpleProcessor(exporter Exporter) *SimpleProcessor {
	p := &SimpleProcessor{exporter: exporter}
	p.pool = sync.Pool{
		New: func() any {
			b := make([]Record, 0, 1)
			return &b
		},
	}
	return p
}

// OnEmit batches provided log record.
func (s *SimpleProcessor) OnEmit(ctx context.Context, r Record) error {
	records := s.pool.Get().(*[]Record)
	defer func() {
		*records = (*records)[:0]
		s.pool.Put(records)
	}()

	return s.exporter.Export(ctx, *records)
}

// Shutdown shuts down the expoter.
func (s *SimpleProcessor) Shutdown(ctx context.Context) error {
	return s.exporter.Shutdown(ctx)
}

// ForceFlush flushes the exporter.
func (s *SimpleProcessor) ForceFlush(ctx context.Context) error {
	return s.exporter.ForceFlush(ctx)
}
