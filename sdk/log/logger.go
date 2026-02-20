// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
)

var now = time.Now

// Compile-time check logger implements log.Logger.
var _ log.Logger = (*logger)(nil)

type logger struct {
	embedded.Logger

	provider             *LoggerProvider
	instrumentationScope instrumentation.Scope

	// recCntIncr increments the count of log records created. It will be nil
	// if observability is disabled.
	recCntIncr func(context.Context)
}

func newLogger(p *LoggerProvider, scope instrumentation.Scope) *logger {
	l := &logger{
		provider:             p,
		instrumentationScope: scope,
	}

	var err error
	l.recCntIncr, err = newRecordCounterIncr()
	if err != nil {
		otel.Handle(err)
	}
	return l
}

func (l *logger) Emit(ctx context.Context, r log.Record) {
	newRecord := l.newRecord(ctx, r)
	for _, p := range l.provider.processors {
		if err := p.OnEmit(ctx, &newRecord); err != nil {
			otel.Handle(err)
		}
	}
}

// Enabled returns true if at least one Processor held by the LoggerProvider
// that created the logger will process for the provided context and param.
//
// If it is not possible to definitively determine the record will be
// processed, true will be returned by default. A value of false will only be
// returned if it can be positively verified that no Processor will process.
func (l *logger) Enabled(ctx context.Context, param log.EnabledParameters) bool {
	p := EnabledParameters{
		InstrumentationScope: l.instrumentationScope,
		Severity:             param.Severity,
		EventName:            param.EventName,
	}

	for _, processor := range l.provider.processors {
		if processor.Enabled(ctx, p) {
			// At least one Processor will process the Record.
			return true
		}
	}
	// No Processor will process the record.
	return false
}

func (l *logger) newRecord(ctx context.Context, r log.Record) Record {
	sc := trace.SpanContextFromContext(ctx)

	newRecord := Record{
		eventName:         r.EventName(),
		timestamp:         r.Timestamp(),
		observedTimestamp: r.ObservedTimestamp(),
		severity:          r.Severity(),
		severityText:      r.SeverityText(),

		traceID:    sc.TraceID(),
		spanID:     sc.SpanID(),
		traceFlags: sc.TraceFlags(),

		resource:                  l.provider.resource,
		scope:                     &l.instrumentationScope,
		attributeValueLengthLimit: l.provider.attributeValueLengthLimit,
		attributeCountLimit:       l.provider.attributeCountLimit,
		allowDupKeys:              l.provider.allowDupKeys,
	}
	if l.recCntIncr != nil {
		l.recCntIncr(ctx)
	}

	// This ensures we deduplicate key-value collections in the log body
	newRecord.SetBody(r.Body())

	// This field SHOULD be set once the event is observed by OpenTelemetry.
	if newRecord.observedTimestamp.IsZero() {
		newRecord.observedTimestamp = now()
	}

	r.WalkAttributes(func(kv log.KeyValue) bool {
		newRecord.AddAttributes(kv)
		return true
	})

	addExceptionFromError(&newRecord, r.GetError())

	return newRecord
}

type stackTracer interface {
	StackTrace() fmt.Formatter
}

func addExceptionFromError(r *Record, err error) {
	if r == nil || err == nil {
		return
	}

	var hasType, hasMessage, hasStacktrace bool
	r.WalkAttributes(func(kv log.KeyValue) bool {
		switch kv.Key {
		case string(semconv.ExceptionTypeKey):
			hasType = true
		case string(semconv.ExceptionMessageKey):
			hasMessage = true
		case string(semconv.ExceptionStacktraceKey):
			hasStacktrace = true
		}
		return !hasType || !hasMessage || !hasStacktrace
	})

	attrs := make([]log.KeyValue, 0, 3)
	if !hasType {
		attrs = append(attrs, log.String(string(semconv.ExceptionTypeKey), errorType(err)))
	}
	if !hasMessage {
		attrs = append(attrs, log.String(string(semconv.ExceptionMessageKey), err.Error()))
	}
	if !hasStacktrace {
		if st := errorStackTrace(err); st != "" {
			attrs = append(attrs, log.String(string(semconv.ExceptionStacktraceKey), st))
		}
	}

	if len(attrs) == 0 {
		return
	}
	r.AddAttributes(attrs...)
}

func errorType(err error) string {
	t := reflect.TypeOf(err)
	if t == nil {
		return ""
	}
	if t.PkgPath() == "" && t.Name() == "" {
		// Likely a builtin type.
		return t.String()
	}
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}

func errorStackTrace(err error) string {
	if st, ok := err.(stackTracer); ok {
		return fmt.Sprintf("%+v", st.StackTrace())
	}
	return ""
}
