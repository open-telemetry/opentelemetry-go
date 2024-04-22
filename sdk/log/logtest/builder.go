// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/sdk/log/logtest"

import (
	"context"
	"slices"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

// TODO: comment.
type RecordBuilder struct {
	timestamp         time.Time
	observedTimestamp time.Time
	// severity          log.Severity
	// severityText      string
	// body              log.Value
	attrs []log.KeyValue
	// traceID           trace.TraceID
	// spanID            trace.SpanID
	// traceFlags        trace.TraceFlags

	resource *resource.Resource
	scope    instrumentation.Scope

	dropped int
}

// TODO: comment.
func (b RecordBuilder) Record() sdklog.Record {
	var record sdklog.Record
	p := processor(func(r sdklog.Record) {
		r.SetTimestamp(b.timestamp)
		r.SetObservedTimestamp(b.observedTimestamp)
		r.SetAttributes(b.attrs...)

		// Generate dropped attributes.
		for i := 0; i < b.dropped; i++ {
			r.AddAttributes(log.KeyValue{})
		}

		record = r
	})

	attributeCountLimit := -1
	if b.dropped > 0 {
		// Make sure that we can generate dropped attributes.
		attributeCountLimit = len(b.attrs)
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(b.resource),
		sdklog.WithAttributeCountLimit(attributeCountLimit),
		sdklog.WithAttributeValueLengthLimit(-1),
		sdklog.WithProcessor(p),
	)

	l := provider.Logger(b.scope.Name,
		log.WithInstrumentationVersion(b.scope.Version),
		log.WithSchemaURL(b.scope.SchemaURL),
	)
	l.Emit(context.Background(), log.Record{}) // This executes the processor function.
	return record
}

// TODO: comment.
func (b RecordBuilder) SetTimestamp(t time.Time) RecordBuilder {
	b.timestamp = t
	return b
}

// TODO: comment.
func (b RecordBuilder) SetObservedTimestamp(t time.Time) RecordBuilder {
	b.observedTimestamp = t
	return b
}

// TODO: comment.
func (b RecordBuilder) AddAttributes(attrs ...log.KeyValue) RecordBuilder {
	b.attrs = slices.Clone(b.attrs)
	b.attrs = append(b.attrs, attrs...)
	return b
}

// TODO: comment.
func (b RecordBuilder) SetAttributes(attrs ...log.KeyValue) RecordBuilder {
	b.attrs = slices.Clone(attrs)
	return b
}

// TODO: comment.
func (b RecordBuilder) SetInstrumentationScope(scope instrumentation.Scope) RecordBuilder {
	b.scope = scope
	return b
}

// TODO: comment.
func (b RecordBuilder) SetResource(r *resource.Resource) RecordBuilder {
	b.resource = r
	return b
}

// TODO: comment.
//
// Notice: The returned record is going to have an attribute count limit.
// Therefore, it will not be possible to add additional attributes on the record
// returned by the builder that has dropped attributes set to a value greater than 0
// (they will be dropped).
func (b RecordBuilder) SetDroppedAttributes(n int) RecordBuilder {
	b.dropped = n
	return b
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
