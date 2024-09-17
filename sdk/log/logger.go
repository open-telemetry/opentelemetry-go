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
	param := log.EnabledParameters{}
	param.SetSeverity(r.Severity())
	if !l.Enabled(ctx, param) {
		return
	}

	newRecord := l.newRecord(ctx, r)
	for _, p := range l.provider.processors {
		if err := p.OnEmit(ctx, &newRecord); err != nil {
			otel.Handle(err)
		}
	}
}

// Enabled returns false if at least one Filterer held by the LoggerProvider
// that created the logger will return false for the provided context and param.
func (l *logger) Enabled(ctx context.Context, param log.EnabledParameters) bool {
	for _, flt := range l.provider.filterers {
		if !flt.Filter(ctx, param) {
			return false
		}
	}
	return true
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
