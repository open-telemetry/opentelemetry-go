// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	semconv "go.opentelemetry.io/otel/semconv/v1.45.0"
	"go.opentelemetry.io/otel/trace"
)

var now = time.Now

// attributesBufferSize is the number of attributes newRecord can collect
// without allocating.
const attributesBufferSize = 16

const (
	exceptionTypeKey    = semconv.ExceptionTypeKey
	exceptionMessageKey = semconv.ExceptionMessageKey
)

// This is a compile-time check that logger implements log.Logger.
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
	processors := l.provider.processors
	if len(processors) == 0 {
		if l.provider.stopped.Load() {
			return
		}
		// Emit remains observable without processors, but no lifecycle
		// admission or SDK record construction is needed.
		l.recordCreated(ctx)
		return
	}

	for _, err := range l.emit(ctx, r, processors) {
		otel.Handle(err)
	}
}

func (l *logger) emit(ctx context.Context, r log.Record, processors []Processor) []error {
	if !l.provider.beginProcessorOperation() {
		return nil
	}
	defer l.provider.endProcessorOperation()

	l.recordCreated(ctx)
	newRecord := l.newRecord(ctx, r)
	var errs []error
	for _, processor := range processors {
		if err := processor.OnEmit(ctx, &newRecord); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (l *logger) recordCreated(ctx context.Context) {
	if l.recCntIncr != nil {
		l.recCntIncr(ctx)
	}
}

// Enabled returns true if at least one Processor held by the LoggerProvider
// that created the logger will process a record for the provided context and
// param.
//
// Enabled returns false after the LoggerProvider that created l starts shutting
// down.
//
// If it is not possible to definitively determine whether the record will be
// processed, true will be returned by default. A value of false will only be
// returned if it can be positively verified that no Processor will process it.
func (l *logger) Enabled(ctx context.Context, param log.EnabledParameters) bool {
	p := EnabledParameters{
		InstrumentationScope: l.instrumentationScope,
		Severity:             param.Severity,
		EventName:            param.EventName,
	}

	processors := l.provider.processors
	if len(processors) == 0 || !l.provider.beginProcessorOperation() {
		return false
	}
	defer l.provider.endProcessorOperation()

	for _, processor := range processors {
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
		attributeValueDepthLimit:  l.provider.attrValueDepthLimit(),
		attributeCountLimit:       l.provider.attributeCountLimit,
		allowDupKeys:              l.provider.allowDupKeys,
	}
	// This ensures we deduplicate key-value collections in the log body
	newRecord.SetBody(r.Body())

	// This field SHOULD be set once the event is observed by OpenTelemetry.
	if newRecord.observedTimestamp.IsZero() {
		newRecord.observedTimestamp = now()
	}

	// User-provided exception attributes MUST take precedence. Track message
	// and type independently so a supplied value suppresses only its own
	// derivation.
	//
	// Collect the attributes and add them in one call. Adding them one at a
	// time makes de-duplication index every attribute already added, which
	// is quadratic in the number of attributes.
	var hasExceptionMessage, hasExceptionType bool
	if r.AttributesLen() <= attributesBufferSize {
		var stack [attributesBufferSize]attribute.KeyValue
		_, hasExceptionMessage, hasExceptionType = addRecordAttrs(&newRecord, &r, stack[:0])
	} else {
		buf := getUnique()
		*buf, hasExceptionMessage, hasExceptionType = addRecordAttrs(&newRecord, &r, *buf)
		putUnique(buf)
	}

	// Avoid inspecting the error for attributes when the caller has
	// already supplied the attributes.
	if err := r.Err(); err != nil && (!hasExceptionMessage || !hasExceptionType) {
		// Derive missing exception attributes by default, as required by the
		// Logs SDK specification. Attribute limits may constrain generation,
		// but attributes omitted due to those limits are still counted.
		var attrs [2]attribute.KeyValue
		n := 0

		// Derived attributes are buffered until flush, so the current attribute
		// count stays unchanged while missing values are prepared.
		hasLimit := newRecord.hasAttributeCountLimit()
		var remaining int
		if hasLimit {
			remaining = newRecord.attributeCountLimit - newRecord.AttributesLen()
		}
		if !hasExceptionMessage {
			if msg := err.Error(); msg != "" {
				if hasLimit && remaining <= n {
					dropped := 1
					if !hasExceptionType {
						// Every non-nil error has a concrete type, so avoid
						// resolving it when it cannot be retained.
						dropped++
					}
					newRecord.addDropped(dropped)
					goto flush
				}
				attrs[n] = exceptionMessageKey.String(msg)
				n++
			}
		}
		if !hasExceptionType {
			if hasLimit && remaining <= n {
				newRecord.addDropped(1)
				goto flush
			}
			if errType := errorType(err); errType != "" {
				attrs[n] = exceptionTypeKey.String(errType)
				n++
			}
		}

	flush:
		if n > 0 {
			newRecord.addAttrs(attrs[:n])
		}
	}

	return newRecord
}

// addRecordAttrs adds the attributes of src to dst, which must have no
// attributes, using buf as scratch space. It reports whether src has the
// exception message and type attributes. It returns buf, possibly grown, with
// length 0 and its contents cleared so it can be reused.
func addRecordAttrs(
	dst *Record,
	src *log.Record,
	buf []attribute.KeyValue,
) (_ []attribute.KeyValue, hasExceptionMessage, hasExceptionType bool) {
	src.WalkAttributes(func(kv attribute.KeyValue) bool {
		switch kv.Key {
		case exceptionMessageKey:
			hasExceptionMessage = true
		case exceptionTypeKey:
			hasExceptionType = true
		}
		buf = append(buf, kv)
		return true
	})
	dst.addAttrsInPlace(buf)
	clear(buf)
	return buf[:0], hasExceptionMessage, hasExceptionType
}

// addAttrsInPlace adds attrs to r, which must have no attributes. It may
// overwrite attrs.
//
// The result is the same as calling r.AddAttributes for each attribute in
// turn, but takes linear time instead of quadratic: a repeated key replaces
// the value at its first position, and once the attribute count limit is
// reached, every occurrence of a key that was not retained is dropped and
// counted.
func (r *Record) addAttrsInPlace(attrs []attribute.KeyValue) {
	attrs = filterInvalid(attrs)
	if r.allowDupKeys || len(attrs) <= 1 {
		attrs, drop := r.head(attrs)
		r.addDropped(drop)
		r.addAttrs(attrs)
		return
	}

	index := getIndex()
	defer putIndex(index)

	var n, dropped int
	var dup bool
	for _, a := range attrs {
		if idx, found := index[a.Key]; found {
			dup = true
			attrs[idx] = a
			continue
		}
		if r.hasAttributeCountLimit() && n >= r.attributeCountLimit {
			dropped++
			continue
		}
		index[a.Key] = n
		attrs[n] = a
		n++
	}
	if dup {
		logKeyValuePairDropped()
	}
	r.addDropped(dropped)
	r.addAttrs(attrs[:n])
}

func errorType(err error) string {
	if et, ok := err.(interface{ ErrorType() string }); ok {
		if s := et.ErrorType(); s != "" {
			return s
		}
	}

	t := reflect.TypeOf(unwrapError(err))
	if t == nil {
		return ""
	}

	pkg, name := t.PkgPath(), t.Name()
	if pkg != "" && name != "" {
		return pkg + "." + name
	}

	// The type has no package path or name (predeclared, not-defined,
	// or alias for a not-defined type).
	//
	// The type has no package path or name (predeclared, not-defined,
	// or alias for a not-defined type).
	//
	// This is not guaranteed to be unique, but is a best effort.
	return t.String()
}

var (
	fmtWrapErrorType = reflect.TypeOf(fmt.Errorf("wrapped: %w", errors.New("err")))
	errorsJoinType   = reflect.TypeOf(errors.Join(errors.New("err")))
)

func unwrapError(err error) error {
	for {
		switch reflect.TypeOf(err) {
		case fmtWrapErrorType:
			u := errors.Unwrap(err)
			if u == nil {
				return err // When the wrapped error is nil, use the concrete type of the wrapper.
			}
			err = u
		case errorsJoinType:
			if u, ok := err.(interface{ Unwrap() []error }); ok {
				if errs := u.Unwrap(); len(errs) == 1 && errs[0] != nil {
					err = errs[0]
					continue
				}
			}
			return err
		default:
			return err
		}
	}
}
