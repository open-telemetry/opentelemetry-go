// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package logtest is a testing helper package.
package logtest // import "go.opentelemetry.io/otel/sdk/log/logtest"

import (
	"reflect"
	"time"
	"unsafe"

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
	EventName         string
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
	InstrumentationScope *instrumentation.Scope

	DroppedAttributes         int
	AttributeValueLengthLimit int
	AttributeCountLimit       int
}

// NewRecord returns a [sdklog.Record] configured from the values of f.
func (f RecordFactory) NewRecord() sdklog.Record {
	// r needs to be addressable for set() below.
	r := new(sdklog.Record)

	// Set to unlimited so attributes are set exactly.
	set(r, "attributeCountLimit", -1)
	set(r, "attributeValueLengthLimit", -1)

	r.SetEventName(f.EventName)
	r.SetTimestamp(f.Timestamp)
	r.SetObservedTimestamp(f.ObservedTimestamp)
	r.SetSeverity(f.Severity)
	r.SetSeverityText(f.SeverityText)
	r.SetBody(f.Body)
	r.SetAttributes(f.Attributes...)
	r.SetTraceID(f.TraceID)
	r.SetSpanID(f.SpanID)
	r.SetTraceFlags(f.TraceFlags)

	set(r, "resource", f.Resource)
	set(r, "scope", f.InstrumentationScope)
	set(r, "dropped", f.DroppedAttributes)
	set(r, "attributeCountLimit", f.AttributeCountLimit)
	set(r, "attributeValueLengthLimit", f.AttributeValueLengthLimit)

	return *r
}

func set(r *sdklog.Record, name string, value any) {
	rVal := reflect.ValueOf(r).Elem()
	rf := rVal.FieldByName(name)
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).
		Elem()
		// nolint: gosec  // conversion of uintptr -> unsafe.Pointer.
	rf.Set(reflect.ValueOf(value))
}
