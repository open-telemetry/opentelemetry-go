// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"testing"

	apiTrace "go.opentelemetry.io/otel/trace"
)

type reentrantRecordError struct {
	span apiTrace.Span
}

func (e reentrantRecordError) Error() string {
	e.span.AddEvent("from Error")
	return "boom"
}

func TestRecordErrorAllowsReentrantErrorFormatting(t *testing.T) {
	tp := NewTracerProvider()
	_, span := tp.Tracer(t.Name()).Start(t.Context(), "span")
	t.Cleanup(func() { span.End() })

	span.RecordError(reentrantRecordError{span: span})
}
