// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package logtest is a testing helper package.
package logtest // import "go.opentelemetry.io/otel/sdk/log/logtest"

import (
	"context"
	"slices"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// RecordFactory is used to facilitate unit testing implementations of
// [go.opentelemetry.io/otel/sdk/log.Exporter]
// and [go.opentelemetry.io/otel/sdk/log.Processor].
//
// Do not use RecordFactory to create records in production code.
type RecordFactory struct {
	Timestamp         time.Time
	ObservedTimestamp time.Time
	Severity          log.Severity
	SeverityText      string
	Body              log.Value
	Attributes        []log.KeyValue
	TraceID           trace.TraceID
	SpanID            trace.SpanID
	TraceFlags        trace.TraceFlags

	Resource             *resource.Resource
	InstrumentationScope instrumentation.Scope

	DroppedAttributes int
}

// NewRecord returns a log record.
func (b RecordFactory) NewRecord() sdklog.Record {
	var record sdklog.Record
	p := processor(func(r sdklog.Record) {
		r.SetTimestamp(b.Timestamp)
		r.SetObservedTimestamp(b.ObservedTimestamp)
		r.SetSeverity(b.Severity)
		r.SetSeverityText(b.SeverityText)
		r.SetBody(b.Body)
		r.SetAttributes(slices.Clone(b.Attributes)...)

		// Generate dropped attributes.
		for i := 0; i < b.DroppedAttributes; i++ {
			r.AddAttributes(log.KeyValue{})
		}

		r.SetTraceID(b.TraceID)
		r.SetSpanID(b.SpanID)
		r.SetTraceFlags(b.TraceFlags)

		record = r
	})

	attributeCountLimit := -1
	if b.DroppedAttributes > 0 {
		// Make sure that we can generate dropped attributes.
		attributeCountLimit = len(b.Attributes)
	}

	res := b.Resource
	if res == nil {
		res = resource.Empty()
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithAttributeCountLimit(attributeCountLimit),
		sdklog.WithAttributeValueLengthLimit(-1),
		sdklog.WithProcessor(p),
	)

	l := provider.Logger(b.InstrumentationScope.Name,
		log.WithInstrumentationVersion(b.InstrumentationScope.Version),
		log.WithSchemaURL(b.InstrumentationScope.SchemaURL),
	)
	l.Emit(context.Background(), log.Record{}) // This executes the processor function.
	return record
}

type processor func(r sdklog.Record)

func (p processor) OnEmit(ctx context.Context, r sdklog.Record) error {
	p(r)
	return nil
}

func (processor) Enabled(context.Context, sdklog.Record) bool {
	return true
}

func (processor) Shutdown(ctx context.Context) error {
	return nil
}

func (processor) ForceFlush(context.Context) error {
	return nil
}
