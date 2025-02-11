// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package xlog_test

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/xlog"
)

// Use a processor that filters out records based on the provided context.
// It also demonstrates the use of experimental [xlog.FilterProcessor].
func ExampleProcessor_filtering() {
	// Existing processor that emits telemetry.
	var processor log.Processor = log.NewBatchProcessor(nil)

	// Wrap the processor so that it ignores processing log records
	// when a context deriving from WithIgnoreLogs is passed
	// to the logging methods.
	processor = &ContextFilterProcessor{Processor: processor}

	// The created processor can then be registered with
	// the OpenTelemetry Logs SDK using the WithProcessor option.
	_ = log.NewLoggerProvider(
		log.WithProcessor(processor),
	)
}

type key struct{}

var ignoreLogsKey key

// WithIgnoreLogs returns a context which is used by [ContextFilterProcessor]
// to filter out log records.
func WithIgnoreLogs(ctx context.Context) context.Context {
	return context.WithValue(ctx, ignoreLogsKey, true)
}

// ContextFilterProcessor filters out logs when a context deriving from
// [WithIgnoreLogs] is passed to its methods.
type ContextFilterProcessor struct {
	log.Processor

	lazyFilter sync.Once
	// Support the experimental FilterProcessor interface for the embedded processor.
	filter xlog.FilterProcessor
}

// Compile time check.
var _ xlog.FilterProcessor = (*ContextFilterProcessor)(nil)

func (p *ContextFilterProcessor) OnEmit(ctx context.Context, record *log.Record) error {
	if ignoreLogs(ctx) {
		return nil
	}
	return p.Processor.OnEmit(ctx, record)
}

func (p *ContextFilterProcessor) Enabled(ctx context.Context, param xlog.EnabledParameters) bool {
	p.lazyFilter.Do(func() {
		if f, ok := p.Processor.(xlog.FilterProcessor); ok {
			p.filter = f
		}
	})
	return !ignoreLogs(ctx) && (p.filter == nil || p.filter.Enabled(ctx, param))
}

func ignoreLogs(ctx context.Context) bool {
	_, ok := ctx.Value(ignoreLogsKey).(bool)
	return ok
}
