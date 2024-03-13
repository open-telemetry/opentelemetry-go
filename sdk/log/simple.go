// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
)

var _ Processor = (*SimpleProcessor)(nil)

// SimpleProcessor is an processor that synchronously exports log records.
type SimpleProcessor struct{}

// NewSimpleProcessor is a simple Processor adapter.
//
// Any of the exporter's methods may be called concurrently with itself
// or with other methods. It is the responsibility of the exporter to manage
// this concurrency.
func NewSimpleProcessor(exporter Exporter) *SimpleProcessor {
	// TODO (#5062): Implement.
	return nil
}

// OnEmit batches provided log record.
func (s *SimpleProcessor) OnEmit(ctx context.Context, r Record) error {
	// TODO (#5062): Implement.
	return nil
}

// Shutdown shuts down the expoter.
func (s *SimpleProcessor) Shutdown(ctx context.Context) error {
	// TODO (#5062): Implement.
	return nil
}

// ForceFlush flushes the exporter.
func (s *SimpleProcessor) ForceFlush(ctx context.Context) error {
	// TODO (#5062): Implement.
	return nil
}
