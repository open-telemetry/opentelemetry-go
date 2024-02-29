// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otel2oc // import "go.opentelemetry.io/otel/bridge/opencensus/internal/otel2oc"

import (
	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/trace"
)

func SpanContext(sc trace.SpanContext) octrace.SpanContext {
	var to octrace.TraceOptions
	if sc.IsSampled() {
		// OpenCensus doesn't expose functions to directly set sampled
		to = 0x1
	}
	return octrace.SpanContext{
		TraceID:      octrace.TraceID(sc.TraceID()),
		SpanID:       octrace.SpanID(sc.SpanID()),
		TraceOptions: to,
	}
}
