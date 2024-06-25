// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"context"

	"go.opentelemetry.io/otel/log"
	logsdk "go.opentelemetry.io/otel/sdk/log"
)

// Use a processor which sets the minimum log level for the Logs SDK.
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
// processor such that only [logsdk.Record]s with a [log.Severity] greater than or
// equal to minimum will be processed.
func NewMinSeverityProcessor(minimum log.Severity, processor logsdk.Processor) *MinSeverityProcessor {
	return &MinSeverityProcessor{
		Processor: processor,
		minimum:   minimum,
	}
}

// OnEmit passes the context and record to the underlying [logsdk.Processor] s
// decorates if the [log.Severity] of record is greater than or equal to the
// minimum severity s is configured at.
func (s *MinSeverityProcessor) OnEmit(ctx context.Context, record logsdk.Record) error {
	if !s.enabled(record) {
		return nil
	}
	return s.Processor.OnEmit(ctx, record)
}

// Enabled returns true the [log.Severity] of record is greater than or equal
// to the minimum severity s is configured at. It will return false if the
// severity is less than the minimum.
//
// If the record severity is unset, this will return true.
func (s *MinSeverityProcessor) Enabled(_ context.Context, record logsdk.Record) bool {
	return s.enabled(record)
}

func (s *MinSeverityProcessor) enabled(r logsdk.Record) bool {
	severity := r.Severity()
	return severity == log.SeverityUndefined || s.minimum <= severity
}
