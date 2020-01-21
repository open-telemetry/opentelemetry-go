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
	otelcorrelation "go.opentelemetry.io/otel/api/correlation"
	otelglobal "go.opentelemetry.io/otel/api/global"
	otelkey "go.opentelemetry.io/otel/api/key"
	otelpropagation "go.opentelemetry.io/otel/api/propagation"
	oteltrace "go.opentelemetry.io/otel/api/trace"
	otelparent "go.opentelemetry.io/otel/internal/trace/parent"

	"go.opentelemetry.io/otel/bridge/opentracing/migration"
)

type bridgeSpanContext struct {
	baggageItems    otelcorrelation.Map
	otelSpanContext otelcore.SpanContext
}

var _ ot.SpanContext = &bridgeSpanContext{}

func newBridgeSpanContext(otelSpanContext otelcore.SpanContext, parentOtSpanContext ot.SpanContext) *bridgeSpanContext {
	bCtx := &bridgeSpanContext{
		baggageItems:    otelcorrelation.NewEmptyMap(),
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
	c.baggageItems.Foreach(func(kv otelcore.KeyValue) bool {
		return handler(string(kv.Key), kv.Value.Emit())
	})
}

func (c *bridgeSpanContext) setBaggageItem(restrictedKey, value string) {
	crk := http.CanonicalHeaderKey(restrictedKey)
	c.baggageItems = c.baggageItems.Apply(otelcorrelation.MapUpdate{SingleKV: otelkey.New(crk).String(value)})
}

func (c *bridgeSpanContext) baggageItem(restrictedKey string) string {
	crk := http.CanonicalHeaderKey(restrictedKey)
	val, _ := c.baggageItems.Value(otelkey.New(crk))
	return val.Emit()
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
		s.otelSpan.SetAttributes(otTagToOtelCoreKeyValue(key, value))
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

	propagators otelpropagation.Propagators
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
		propagators:    nil,
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

func (t *BridgeTracer) SetPropagators(propagators otelpropagation.Propagators) {
	t.propagators = propagators
}

// StartSpan is a part of the implementation of the OpenTracing Tracer
// interface.
func (t *BridgeTracer) StartSpan(operationName string, opts ...ot.StartSpanOption) ot.Span {
	sso := ot.StartSpanOptions{}
	for _, opt := range opts {
		opt.Apply(&sso)
	}
	// TODO: handle links
	bRelation, _ := otSpanReferencesToBridgeRelationAndLinks(sso.References)
	attributes, kind, hadTrueErrorTag := otTagsToOtelAttributesKindAndError(sso.Tags)
	checkCtx := migration.WithDeferredSetup(context.Background())
	checkCtx2, otelSpan := t.setTracer.tracer().Start(checkCtx, operationName, func(opts *oteltrace.StartConfig) {
		opts.Attributes = attributes
		opts.StartTime = sso.StartTime
		opts.Parent = bRelation.ToOtelParent()
		opts.Record = true
		opts.Links = nil // TODO: handle links
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
	err := false
	var pairs []otelcore.KeyValue
	for k, v := range tags {
		switch k {
		case string(otext.SpanKind):
			if s, ok := v.(string); ok {
				switch strings.ToLower(s) {
				case "client":
					kind = oteltrace.SpanKindClient
				case "server":
					kind = oteltrace.SpanKindServer
				case "producer":
					kind = oteltrace.SpanKindProducer
				case "consumer":
					kind = oteltrace.SpanKindConsumer
				}
			}
		case string(otext.Error):
			if b, ok := v.(bool); ok && b {
				err = true
			}
		default:
			pairs = append(pairs, otTagToOtelCoreKeyValue(k, v))
		}
	}
	return pairs, kind, err
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
	relationshipType ot.SpanReferenceType
}

func (r bridgeRelation) ToOtelParent() context.Context {
	ctx := context.Background()
	if r.spanContext == nil {
		return ctx
	}
	ctx = oteltrace.WithRemoteContext(ctx, otelcore.SpanContext{
		TraceID: r.spanContext.otelSpanContext.TraceID,
		SpanID:  r.spanContext.otelSpanContext.SpanID,
		// TODO: Flags
	})
	// TODO: relationshipType
	return ctx
}

func otSpanReferencesToBridgeRelationAndLinks(references []ot.SpanReference) (bridgeRelation, []*bridgeSpanContext) {
	if len(references) == 0 {
		return bridgeRelation{}, nil
	}
	first := references[0]
	relation := bridgeRelation{
		spanContext:      mustGetBridgeSpanContext(first.ReferencedContext),
		relationshipType: first.Type,
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

// fakeSpan is just a holder of span context, nothing more. It's for
// propagators, so they can get the span context from Go context.
type fakeSpan struct {
	oteltrace.NoopSpan

	sc otelcore.SpanContext
}

func (s fakeSpan) SpanContext() otelcore.SpanContext {
	return s.sc
}

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
	header := http.Header(hhcarrier)
	fs := fakeSpan{
		sc: bridgeSC.otelSpanContext,
	}
	ctx := oteltrace.ContextWithSpan(context.Background(), fs)
	ctx = otelcorrelation.WithMap(ctx, bridgeSC.baggageItems)
	otelpropagation.InjectHTTP(ctx, t.getPropagators(), header)
	return nil
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
	header := http.Header(hhcarrier)
	ctx := otelpropagation.ExtractHTTP(context.Background(), t.getPropagators(), header)
	baggage := otelcorrelation.FromContext(ctx)
	_, otelSC, _ := otelparent.GetContext(ctx, nil)
	bridgeSC := &bridgeSpanContext{
		baggageItems:    baggage,
		otelSpanContext: otelSC,
	}
	if !bridgeSC.otelSpanContext.IsValid() {
		return nil, ot.ErrSpanContextNotFound
	}
	return bridgeSC, nil
}

func (t *BridgeTracer) getPropagators() otelpropagation.Propagators {
	if t.propagators != nil {
		return t.propagators
	}
	return otelglobal.Propagators()
}
