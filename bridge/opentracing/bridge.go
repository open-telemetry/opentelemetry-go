// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opentracing

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"google.golang.org/grpc/codes"

	ot "github.com/opentracing/opentracing-go"
	otext "github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"

	otelcore "go.opentelemetry.io/otel/api/core"
	oteltrace "go.opentelemetry.io/otel/api/trace"

	"go.opentelemetry.io/otel/bridge/opentracing/migration"
)

type bridgeSpanContext struct {
	// TODO: have a look at the java implementation of the shim to
	// see what do they do with the baggage items
	baggageItems    map[string]string
	otelSpanContext otelcore.SpanContext
}

var _ ot.SpanContext = &bridgeSpanContext{}

func newBridgeSpanContext(otelSpanContext otelcore.SpanContext, parentOtSpanContext ot.SpanContext) *bridgeSpanContext {
	bCtx := &bridgeSpanContext{
		baggageItems:    nil,
		otelSpanContext: otelSpanContext,
	}
	if parentOtSpanContext != nil {
		parentOtSpanContext.ForeachBaggageItem(func(key, value string) bool {
			bCtx.setBaggageItem(key, value)
			return true
		})
	}
	return bCtx
}

func (c *bridgeSpanContext) ForeachBaggageItem(handler func(k, v string) bool) {
	for k, v := range c.baggageItems {
		if !handler(k, v) {
			break
		}
	}
}

func (c *bridgeSpanContext) setBaggageItem(restrictedKey, value string) {
	if c.baggageItems == nil {
		c.baggageItems = make(map[string]string)
	}
	crk := http.CanonicalHeaderKey(restrictedKey)
	c.baggageItems[crk] = value
}

func (c *bridgeSpanContext) baggageItem(restrictedKey string) string {
	crk := http.CanonicalHeaderKey(restrictedKey)
	return c.baggageItems[crk]
}

type bridgeSpan struct {
	otelSpan      oteltrace.Span
	ctx           *bridgeSpanContext
	tracer        *BridgeTracer
	skipDeferHook bool
}

var _ ot.Span = &bridgeSpan{}

func (s *bridgeSpan) Finish() {
	s.otelSpan.End()
}

func (s *bridgeSpan) FinishWithOptions(opts ot.FinishOptions) {
	var otelOpts []oteltrace.EndOption

	if !opts.FinishTime.IsZero() {
		otelOpts = append(otelOpts, oteltrace.WithEndTime(opts.FinishTime))
	}
	for _, record := range opts.LogRecords {
		s.logRecord(record)
	}
	for _, data := range opts.BulkLogData {
		s.logRecord(data.ToLogRecord())
	}
	s.otelSpan.End(otelOpts...)
}

func (s *bridgeSpan) logRecord(record ot.LogRecord) {
	s.otelSpan.AddEventWithTimestamp(context.Background(), record.Timestamp, "", otLogFieldsToOtelCoreKeyValues(record.Fields)...)
}

func (s *bridgeSpan) Context() ot.SpanContext {
	return s.ctx
}

func (s *bridgeSpan) SetOperationName(operationName string) ot.Span {
	s.otelSpan.SetName(operationName)
	return s
}

func (s *bridgeSpan) SetTag(key string, value interface{}) ot.Span {
	switch key {
	case string(otext.SpanKind):
		// TODO: Should we ignore it?
	case string(otext.Error):
		if b, ok := value.(bool); ok {
			status := codes.OK
			if b {
				status = codes.Unknown
			}
			s.otelSpan.SetStatus(status)
		}
	default:
		s.otelSpan.SetAttribute(otTagToOtelCoreKeyValue(key, value))
	}
	return s
}

func (s *bridgeSpan) LogFields(fields ...otlog.Field) {
	s.otelSpan.AddEvent(context.Background(), "", otLogFieldsToOtelCoreKeyValues(fields)...)
}

type bridgeFieldEncoder struct {
	pairs []otelcore.KeyValue
}

var _ otlog.Encoder = &bridgeFieldEncoder{}

func (e *bridgeFieldEncoder) EmitString(key, value string) {
	e.emitCommon(key, value)
}

func (e *bridgeFieldEncoder) EmitBool(key string, value bool) {
	e.emitCommon(key, value)
}

func (e *bridgeFieldEncoder) EmitInt(key string, value int) {
	e.emitCommon(key, value)
}

func (e *bridgeFieldEncoder) EmitInt32(key string, value int32) {
	e.emitCommon(key, value)
}

func (e *bridgeFieldEncoder) EmitInt64(key string, value int64) {
	e.emitCommon(key, value)
}

func (e *bridgeFieldEncoder) EmitUint32(key string, value uint32) {
	e.emitCommon(key, value)
}

func (e *bridgeFieldEncoder) EmitUint64(key string, value uint64) {
	e.emitCommon(key, value)
}

func (e *bridgeFieldEncoder) EmitFloat32(key string, value float32) {
	e.emitCommon(key, value)
}

func (e *bridgeFieldEncoder) EmitFloat64(key string, value float64) {
	e.emitCommon(key, value)
}

func (e *bridgeFieldEncoder) EmitObject(key string, value interface{}) {
	e.emitCommon(key, value)
}

func (e *bridgeFieldEncoder) EmitLazyLogger(value otlog.LazyLogger) {
	value(e)
}

func (e *bridgeFieldEncoder) emitCommon(key string, value interface{}) {
	e.pairs = append(e.pairs, otTagToOtelCoreKeyValue(key, value))
}

func otLogFieldsToOtelCoreKeyValues(fields []otlog.Field) []otelcore.KeyValue {
	encoder := &bridgeFieldEncoder{}
	for _, field := range fields {
		field.Marshal(encoder)
	}
	return encoder.pairs
}

func (s *bridgeSpan) LogKV(alternatingKeyValues ...interface{}) {
	fields, err := otlog.InterleavedKVToFields(alternatingKeyValues...)
	if err != nil {
		return
	}
	s.LogFields(fields...)
}

func (s *bridgeSpan) SetBaggageItem(restrictedKey, value string) ot.Span {
	s.ctx.setBaggageItem(restrictedKey, value)
	return s
}

func (s *bridgeSpan) BaggageItem(restrictedKey string) string {
	return s.ctx.baggageItem(restrictedKey)
}

func (s *bridgeSpan) Tracer() ot.Tracer {
	return s.tracer
}

func (s *bridgeSpan) LogEvent(event string) {
	s.LogEventWithPayload(event, nil)
}

func (s *bridgeSpan) LogEventWithPayload(event string, payload interface{}) {
	data := ot.LogData{
		Event:   event,
		Payload: payload,
	}
	s.Log(data)
}

func (s *bridgeSpan) Log(data ot.LogData) {
	record := data.ToLogRecord()
	s.LogFields(record.Fields...)
}

type bridgeSetTracer struct {
	isSet      bool
	otelTracer oteltrace.Tracer

	warningHandler BridgeWarningHandler
	warnOnce       sync.Once
}

func (s *bridgeSetTracer) tracer() oteltrace.Tracer {
	if !s.isSet {
		s.warnOnce.Do(func() {
			s.warningHandler("The OpenTelemetry tracer is not set, default no-op tracer is used! Call SetOpenTelemetryTracer to set it up.\n")
		})
	}
	return s.otelTracer
}

// BridgeWarningHandler is a type of handler that receives warnings
// from the BridgeTracer.
type BridgeWarningHandler func(msg string)

// BridgeTracer is an implementation of the OpenTracing tracer, which
// translates the calls to the OpenTracing API into OpenTelemetry
// counterparts and calls the underlying OpenTelemetry tracer.
type BridgeTracer struct {
	setTracer bridgeSetTracer

	warningHandler BridgeWarningHandler
	warnOnce       sync.Once
}

var _ ot.Tracer = &BridgeTracer{}
var _ ot.TracerContextWithSpanExtension = &BridgeTracer{}

// NewBridgeTracer creates a new BridgeTracer. The new tracer forwards
// the calls to the OpenTelemetry Noop tracer, so it should be
// overridden with the SetOpenTelemetryTracer function. The warnings
// handler does nothing by default, so to override it use the
// SetWarningHandler function.
func NewBridgeTracer() *BridgeTracer {
	return &BridgeTracer{
		setTracer: bridgeSetTracer{
			otelTracer: oteltrace.NoopTracer{},
		},
		warningHandler: func(msg string) {},
	}
}

// SetWarningHandler overrides the warning handler.
func (t *BridgeTracer) SetWarningHandler(handler BridgeWarningHandler) {
	t.setTracer.warningHandler = handler
	t.warningHandler = handler
}

// SetWarningHandler overrides the underlying OpenTelemetry
// tracer. The passed tracer should know how to operate in the
// environment that uses OpenTracing API.
func (t *BridgeTracer) SetOpenTelemetryTracer(tracer oteltrace.Tracer) {
	t.setTracer.otelTracer = tracer
	t.setTracer.isSet = true
}

// StartSpan is a part of the implementation of the OpenTracing Tracer
// interface.
func (t *BridgeTracer) StartSpan(operationName string, opts ...ot.StartSpanOption) ot.Span {
	sso := ot.StartSpanOptions{}
	for _, opt := range opts {
		opt.Apply(&sso)
	}
	// TODO: handle links, needs SpanData to be in the API first?
	bRelation, _ := otSpanReferencesToBridgeRelationAndLinks(sso.References)
	attributes, kind, hadTrueErrorTag := otTagsToOtelAttributesKindAndError(sso.Tags)
	checkCtx := migration.WithDeferredSetup(context.Background())
	checkCtx2, otelSpan := t.setTracer.tracer().Start(checkCtx, operationName, func(opts *oteltrace.SpanOptions) {
		opts.Attributes = attributes
		opts.StartTime = sso.StartTime
		opts.Relation = bRelation.ToOtelRelation()
		opts.Record = true
		opts.SpanKind = kind
	})
	if checkCtx != checkCtx2 {
		t.warnOnce.Do(func() {
			t.warningHandler("SDK should have deferred the context setup, see the documentation of go.opentelemetry.io/otel/bridge/opentracing/migration\n")
		})
	}
	if hadTrueErrorTag {
		otelSpan.SetStatus(codes.Unknown)
	}
	var otSpanContext ot.SpanContext
	if bRelation.spanContext != nil {
		otSpanContext = bRelation.spanContext
	}
	sctx := newBridgeSpanContext(otelSpan.SpanContext(), otSpanContext)
	span := &bridgeSpan{
		otelSpan: otelSpan,
		ctx:      sctx,
		tracer:   t,
	}

	return span
}

// ContextWithBridgeSpan sets up the context with the passed
// OpenTelemetry span as the active OpenTracing span.
//
// This function should be used by the OpenTelemetry tracers that want
// to be aware how to operate in the environment using OpenTracing
// API.
func (t *BridgeTracer) ContextWithBridgeSpan(ctx context.Context, span oteltrace.Span) context.Context {
	var otSpanContext ot.SpanContext
	if parentSpan := ot.SpanFromContext(ctx); parentSpan != nil {
		otSpanContext = parentSpan.Context()
	}
	bCtx := newBridgeSpanContext(span.SpanContext(), otSpanContext)
	bSpan := &bridgeSpan{
		otelSpan:      span,
		ctx:           bCtx,
		tracer:        t,
		skipDeferHook: true,
	}
	return ot.ContextWithSpan(ctx, bSpan)
}

// ContextWithSpanHook is an implementation of the OpenTracing tracer
// extension interface. It will call the DeferredContextSetupHook
// function on the tracer if it implements the
// DeferredContextSetupTracerExtension interface.
func (t *BridgeTracer) ContextWithSpanHook(ctx context.Context, span ot.Span) context.Context {
	bSpan, ok := span.(*bridgeSpan)
	if !ok || bSpan.skipDeferHook {
		return ctx
	}
	if tracerWithExtension, ok := bSpan.tracer.setTracer.tracer().(migration.DeferredContextSetupTracerExtension); ok {
		ctx = tracerWithExtension.DeferredContextSetupHook(ctx, bSpan.otelSpan)
	}
	return ctx
}

func otTagsToOtelAttributesKindAndError(tags map[string]interface{}) ([]otelcore.KeyValue, oteltrace.SpanKind, bool) {
	kind := oteltrace.SpanKindInternal
	error := false
	var pairs []otelcore.KeyValue
	for k, v := range tags {
		switch k {
		case string(otext.SpanKind):
			if sk, ok := v.(string); ok {
				kind = oteltrace.SpanKind(sk)
			}
		case string(otext.Error):
			if b, ok := v.(bool); ok && b {
				error = true
			}
		default:
			pairs = append(pairs, otTagToOtelCoreKeyValue(k, v))
		}
	}
	return pairs, kind, error
}

func otTagToOtelCoreKeyValue(k string, v interface{}) otelcore.KeyValue {
	key := otTagToOtelCoreKey(k)
	switch val := v.(type) {
	case bool:
		return key.Bool(val)
	case int64:
		return key.Int64(val)
	case uint64:
		return key.Uint64(val)
	case float64:
		return key.Float64(val)
	case int32:
		return key.Int32(val)
	case uint32:
		return key.Uint32(val)
	case float32:
		return key.Float32(val)
	case int:
		return key.Int(val)
	case uint:
		return key.Uint(val)
	case string:
		return key.String(val)
	default:
		return key.String(fmt.Sprint(v))
	}
}

func otTagToOtelCoreKey(k string) otelcore.Key {
	return otelcore.Key(k)
}

type bridgeRelation struct {
	spanContext      *bridgeSpanContext
	relationshipType oteltrace.RelationshipType
}

func (r bridgeRelation) ToOtelRelation() oteltrace.Relation {
	if r.spanContext == nil {
		return oteltrace.Relation{}
	}
	return oteltrace.Relation{
		SpanContext:      r.spanContext.otelSpanContext,
		RelationshipType: r.relationshipType,
	}
}

func otSpanReferencesToBridgeRelationAndLinks(references []ot.SpanReference) (bridgeRelation, []*bridgeSpanContext) {
	if len(references) == 0 {
		return bridgeRelation{}, nil
	}
	first := references[0]
	relation := bridgeRelation{
		spanContext:      mustGetBridgeSpanContext(first.ReferencedContext),
		relationshipType: otSpanReferenceTypeToOtelRelationshipType(first.Type),
	}
	var links []*bridgeSpanContext
	for _, reference := range references[1:] {
		links = append(links, mustGetBridgeSpanContext(reference.ReferencedContext))
	}
	return relation, links
}

func mustGetBridgeSpanContext(ctx ot.SpanContext) *bridgeSpanContext {
	ourCtx, ok := ctx.(*bridgeSpanContext)
	if !ok {
		panic("oops, some foreign span context here")
	}
	return ourCtx
}

func otSpanReferenceTypeToOtelRelationshipType(srt ot.SpanReferenceType) oteltrace.RelationshipType {
	switch srt {
	case ot.ChildOfRef:
		return oteltrace.ChildOfRelationship
	case ot.FollowsFromRef:
		return oteltrace.FollowsFromRelationship
	default:
		panic("fix yer code, it uses bogus opentracing reference type")
	}
}

// TODO: these headers are most likely bogus
var (
	traceIDHeader       = http.CanonicalHeaderKey("x-otelbridge-trace-id")
	spanIDHeader        = http.CanonicalHeaderKey("x-otelbridge-span-id")
	traceFlagsHeader    = http.CanonicalHeaderKey("x-otelbridge-trace-flags")
	baggageHeaderPrefix = http.CanonicalHeaderKey("x-otelbridge-baggage-")
)

// Inject is a part of the implementation of the OpenTracing Tracer
// interface.
//
// Currently only the HTTPHeaders format is kinda sorta supported.
func (t *BridgeTracer) Inject(sm ot.SpanContext, format interface{}, carrier interface{}) error {
	bridgeSC, ok := sm.(*bridgeSpanContext)
	if !ok {
		return ot.ErrInvalidSpanContext
	}
	if !bridgeSC.otelSpanContext.IsValid() {
		return ot.ErrInvalidSpanContext
	}
	if builtinFormat, ok := format.(ot.BuiltinFormat); !ok || builtinFormat != ot.HTTPHeaders {
		return ot.ErrUnsupportedFormat
	}
	hhcarrier, ok := carrier.(ot.HTTPHeadersCarrier)
	if !ok {
		return ot.ErrInvalidCarrier
	}
	hhcarrier.Set(traceIDHeader, bridgeSC.otelSpanContext.TraceIDString())
	hhcarrier.Set(spanIDHeader, bridgeSC.otelSpanContext.SpanIDString())
	hhcarrier.Set(traceFlagsHeader, traceFlagsToString(bridgeSC.otelSpanContext.TraceFlags))
	bridgeSC.ForeachBaggageItem(func(k, v string) bool {
		// we assume that keys are already canonicalized
		hhcarrier.Set(baggageHeaderPrefix+k, v)
		return true
	})
	return nil
}

func traceFlagsToString(opts byte) string {
	var parts []string
	if opts&otelcore.TraceFlagsSampled == otelcore.TraceFlagsSampled {
		parts = append(parts, "sampled")
	}
	return strings.Join(parts, ",")
}

// Extract is a part of the implementation of the OpenTracing Tracer
// interface.
//
// Currently only the HTTPHeaders format is kinda sorta supported.
func (t *BridgeTracer) Extract(format interface{}, carrier interface{}) (ot.SpanContext, error) {
	if builtinFormat, ok := format.(ot.BuiltinFormat); !ok || builtinFormat != ot.HTTPHeaders {
		return nil, ot.ErrUnsupportedFormat
	}
	hhcarrier, ok := carrier.(ot.HTTPHeadersCarrier)
	if !ok {
		return nil, ot.ErrInvalidCarrier
	}
	bridgeSC := &bridgeSpanContext{}
	err := hhcarrier.ForeachKey(func(k, v string) error {
		ck := http.CanonicalHeaderKey(k)
		switch ck {
		case traceIDHeader:
			traceID, err := otelcore.TraceIDFromHex(v)
			if err != nil {
				return err
			}
			bridgeSC.otelSpanContext.TraceID = traceID
		case spanIDHeader:
			spanID, err := otelcore.SpanIDFromHex(v)
			if err != nil {
				return err
			}
			bridgeSC.otelSpanContext.SpanID = spanID
		case traceFlagsHeader:
			bridgeSC.otelSpanContext.TraceFlags = stringToTraceFlags(v)
		default:
			if strings.HasPrefix(ck, baggageHeaderPrefix) {
				bk := strings.TrimPrefix(ck, baggageHeaderPrefix)
				bridgeSC.setBaggageItem(bk, v)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if !bridgeSC.otelSpanContext.IsValid() {
		return nil, ot.ErrSpanContextNotFound
	}
	return bridgeSC, nil
}

func stringToTraceFlags(s string) byte {
	var opts byte
	for _, part := range strings.Split(s, ",") {
		switch part {
		case "sampled":
			opts |= otelcore.TraceFlagsSampled
		}
	}
	return opts
}
