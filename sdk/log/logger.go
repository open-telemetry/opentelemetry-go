// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

var now = time.Now

// Compile-time check logger implements log.Logger.
var _ log.Logger = (*logger)(nil)

type logger struct {
	embedded.Logger

	provider             *LoggerProvider
	instrumentationScope instrumentation.Scope
}

func newLogger(p *LoggerProvider, scope instrumentation.Scope) *logger {
	return &logger{
		provider:             p,
		instrumentationScope: scope,
	}
}

func (l *logger) Emit(ctx context.Context, r log.Record) {
	newRecord := l.newRecord(ctx, r)
	for _, p := range l.provider.processors {
		if err := p.OnEmit(ctx, &newRecord); err != nil {
			otel.Handle(err)
		}
	}
}

func (l *logger) Enabled(ctx context.Context, r log.EnabledParameters) bool {
	newParam := l.newEnabledParameters(ctx, r)
	for _, p := range l.provider.processors {
		if enabled := p.Enabled(ctx, newParam); enabled {
			return true
		}
	}
	return false
}

func (l *logger) newRecord(ctx context.Context, r log.Record) Record {
	sc := trace.SpanContextFromContext(ctx)

	newRecord := Record{
		timestamp:         r.Timestamp(),
		observedTimestamp: r.ObservedTimestamp(),
		severity:          r.Severity(),
		severityText:      r.SeverityText(),
		body:              r.Body(),

		traceID:    sc.TraceID(),
		spanID:     sc.SpanID(),
		traceFlags: sc.TraceFlags(),

		resource:                  l.provider.resource,
		scope:                     &l.instrumentationScope,
		attributeValueLengthLimit: l.provider.attributeValueLengthLimit,
		attributeCountLimit:       l.provider.attributeCountLimit,
	}

	// This field SHOULD be set once the event is observed by OpenTelemetry.
	if newRecord.observedTimestamp.IsZero() {
		newRecord.observedTimestamp = now()
	}

	r.WalkAttributes(func(kv log.KeyValue) bool {
		newRecord.AddAttributes(kv)
		return true
	})

	return newRecord
}

func (l *logger) newEnabledParameters(ctx context.Context, param log.EnabledParameters) EnabledParameters {
	sc := trace.SpanContextFromContext(ctx)

	newParam := EnabledParameters{
		traceID:    sc.TraceID(),
		spanID:     sc.SpanID(),
		traceFlags: sc.TraceFlags(),

		resource: l.provider.resource,
		scope:    &l.instrumentationScope,
	}

	if v, ok := param.Severity(); ok {
		newParam.setSeverity(v)
	}

	return newParam
}

// EnabledParameters represents payload for [Logger]'s Enabled method.
type EnabledParameters struct {
	severity    log.Severity
	severitySet bool

	traceID    trace.TraceID
	spanID     trace.SpanID
	traceFlags trace.TraceFlags

	// resource represents the entity that collected the log.
	resource *resource.Resource

	// scope is the Scope that the Logger was created with.
	scope *instrumentation.Scope
}

// Severity returns the [Severity] level value, or [SeverityUndefined] if no value was set.
// The ok result indicates whether the value was set.
func (r *EnabledParameters) Severity() (value log.Severity, ok bool) {
	return r.severity, r.severitySet
}

// setSeverity sets the [Severity] level.
func (r *EnabledParameters) setSeverity(level log.Severity) {
	r.severity = level
	r.severitySet = true
}

// TraceID returns the trace ID or empty array.
func (r *EnabledParameters) TraceID() trace.TraceID {
	return r.traceID
}

// SpanID returns the span ID or empty array.
func (r *EnabledParameters) SpanID() trace.SpanID {
	return r.spanID
}

// TraceFlags returns the trace flags.
func (r *EnabledParameters) TraceFlags() trace.TraceFlags {
	return r.traceFlags
}

// Resource returns the entity that collected the log.
func (r *EnabledParameters) Resource() resource.Resource {
	if r.resource == nil {
		return *resource.Empty()
	}
	return *r.resource
}

// InstrumentationScope returns the scope that the Logger was created with.
func (r *EnabledParameters) InstrumentationScope() instrumentation.Scope {
	if r.scope == nil {
		return instrumentation.Scope{}
	}
	return *r.scope
}
