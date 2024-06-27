// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/log"
	logsdk "go.opentelemetry.io/otel/sdk/log"
)

// Use a processor which sets the minimum log level.
func ExampleProcessor_filtering() {
	// Existing processor that emits telemetry.
	var processor logsdk.Processor = logsdk.NewBatchProcessor(nil)

	// Wrap the processor so that it respects minimum serverity level.
	processor = NewMinSeverityProcessor(log.SeverityInfo, processor)

	// The created processor can then be registered with
	// the OpenTelemetry Logs SDK using the WithProcessor option.
	_ = logsdk.NewLoggerProvider(
		logsdk.WithProcessor(processor),
	)
}

// MinSeverityProcessor is a [logsdk.Processor] that limits processing of log records
// to only those with a [log.Severity] above a configured threshold.
type MinSeverityProcessor struct {
	logsdk.Processor
	minimum log.Severity
}

// NewMinSeverityProcessor returns a [MinSeverityProcessor] that decorates the
// processor such that only [logsdk.Record] with a [log.Severity] greater than or
// equal to minimum is processed.
func NewMinSeverityProcessor(minimum log.Severity, processor logsdk.Processor) *MinSeverityProcessor {
	return &MinSeverityProcessor{
		Processor: processor,
		minimum:   minimum,
	}
}

// OnEmit passes the context and record to the underlying [logsdk.Processor]
// if the [log.Severity] of record is greater than or equal to the
// minimum severity is configured at.
func (p *MinSeverityProcessor) OnEmit(ctx context.Context, record logsdk.Record) error {
	if !p.enabled(record) {
		return nil
	}
	return p.Processor.OnEmit(ctx, record)
}

// Enabled returns true if the [log.Severity] of record is greater than or equal
// to the minimum severity is configured at. It will return false if the
// severity is less than the minimum.
//
// If the record severity is unset, this will return true.
func (p *MinSeverityProcessor) Enabled(ctx context.Context, record logsdk.Record) bool {
	return p.enabled(record) && p.Processor.Enabled(ctx, record)
}

func (p *MinSeverityProcessor) enabled(r logsdk.Record) bool {
	severity := r.Severity()
	return severity == log.SeverityUndefined || p.minimum <= severity
}

// Use a processor which redacts sensitive data from some attributes.
func ExampleProcessor_redact() {
	// Existing processor that emits telemetry.
	var processor logsdk.Processor = logsdk.NewBatchProcessor(nil)

	// Wrap the processor so that it redacts values from token attributes.
	processor = &RedactTokensProcessor{processor}

	// The created processor can then be registered with
	// the OpenTelemetry Logs SDK using the WithProcessor option.
	_ = logsdk.NewLoggerProvider(
		logsdk.WithProcessor(processor),
	)
}

// RedactTokensProcessor is a [logsdk.Processor] decorator that redacts values
// from attributes containing "token" in the key.
type RedactTokensProcessor struct {
	logsdk.Processor
}

// OnEmit redacts values from attributes containing "token" in the key
// by replacing them with a REDACTED value.
func (p *RedactTokensProcessor) OnEmit(ctx context.Context, record logsdk.Record) error {
	record.WalkAttributes(func(kv log.KeyValue) bool {
		if strings.Contains(strings.ToLower(kv.Key), "token") {
			record.AddAttributes(log.String(kv.Key, "REDACTED"))
		}
		return true
	})
	return p.Processor.OnEmit(ctx, record)
}
