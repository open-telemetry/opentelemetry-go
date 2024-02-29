// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/otlp/otlptrace/otlptracetest/data.go.tmpl

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracetest // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/otlptracetest"

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// SingleReadOnlySpan returns a one-element slice with a read-only span. It
// may be useful for testing driver's trace export.
func SingleReadOnlySpan() []tracesdk.ReadOnlySpan {
	return tracetest.SpanStubs{
		{
			SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    trace.TraceID{2, 3, 4, 5, 6, 7, 8, 9, 2, 3, 4, 5, 6, 7, 8, 9},
				SpanID:     trace.SpanID{3, 4, 5, 6, 7, 8, 9, 0},
				TraceFlags: trace.FlagsSampled,
			}),
			Parent: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    trace.TraceID{2, 3, 4, 5, 6, 7, 8, 9, 2, 3, 4, 5, 6, 7, 8, 9},
				SpanID:     trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
				TraceFlags: trace.FlagsSampled,
			}),
			SpanKind:          trace.SpanKindInternal,
			Name:              "foo",
			StartTime:         time.Date(2020, time.December, 8, 20, 23, 0, 0, time.UTC),
			EndTime:           time.Date(2020, time.December, 0, 20, 24, 0, 0, time.UTC),
			Attributes:        []attribute.KeyValue{},
			Events:            []tracesdk.Event{},
			Links:             []tracesdk.Link{},
			Status:            tracesdk.Status{Code: codes.Ok},
			DroppedAttributes: 0,
			DroppedEvents:     0,
			DroppedLinks:      0,
			ChildSpanCount:    0,
			Resource:          resource.NewSchemaless(attribute.String("a", "b")),
			InstrumentationLibrary: instrumentation.Library{
				Name:    "bar",
				Version: "0.0.0",
			},
		},
	}.Snapshots()
}
