// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"

	"go.opentelemetry.io/otel/log"
)

// Compile-time check MinSeverity implements Processor.
var _ Processor = (*MinSeverityProcessor)(nil)

// MinSeverityProcessor is a [Processor] that limits processing of log records
// to only those with a [log.Severity] above a configured threshold.
type MinSeverityProcessor struct {
	Processor
	minimum log.Severity
}

// NewMinSeverityProcessor returns a [MinSeverityProcessor] that decorates the
// processor such that only [Record]s with a [log.Severity] greater than or
// equal to minimum will be processed.
func NewMinSeverityProcessor(minimum log.Severity, processor Processor) *MinSeverityProcessor {
	return &MinSeverityProcessor{
		Processor: processor,
		minimum:   minimum,
	}
}

// OnEmit passes the context and record to the underlying [Processor] s
// decorates if the [log.Severity] of record is greater than or equal to the
// minimum severity s is configured at.
func (s *MinSeverityProcessor) OnEmit(ctx context.Context, record Record) error {
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
func (s *MinSeverityProcessor) Enabled(_ context.Context, record Record) bool {
	return s.enabled(record)
}

func (s *MinSeverityProcessor) enabled(r Record) bool {
	// TODO (#5067): replace with var definition when added.
	const unsetSeverity log.Severity = 0

	severity := r.Severity()
	return severity == unsetSeverity || s.minimum <= severity
}
