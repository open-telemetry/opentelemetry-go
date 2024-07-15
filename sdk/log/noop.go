// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import "context"

// Compile-time check NoopProcessor implements Processor.
var _ Processor = (*NoopProcessor)(nil)

var noopProcessorInstance = &NoopProcessor{}

// NoopProcessor is a [Processor] that does nothing.
type NoopProcessor struct{}

// NewNoopProcessor returns a [Processor] that does nothing.
func NewNoopProcessor() *NoopProcessor {
	return noopProcessorInstance
}

// Enabled returns true.
func (p *NoopProcessor) Enabled(context.Context, Record) bool {
	return true
}

// OnEmit does nothing and returns nil.
func (p *NoopProcessor) OnEmit(ctx context.Context, r Record) error {
	return nil
}

// Shutdown does nothing and returns nil.
func (p *NoopProcessor) Shutdown(ctx context.Context) error {
	return nil
}

// ForceFlush does nothing and returns nil.
func (panic *NoopProcessor) ForceFlush(ctx context.Context) error {
	return nil
}
