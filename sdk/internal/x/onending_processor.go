// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/internal/x"

import "go.opentelemetry.io/otel/trace"

// OnEndingSpanProcessor represents span processors that allow mutating spans
// just before they are ended and made immutable.
type OnEndingSpanProcessor interface {
	// OnEnding is called while the span is finished, and while spans are still
	// mutable. It is called synchronously and cannot block.
	OnEnding(trace.Span)
}
